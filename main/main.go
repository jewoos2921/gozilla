package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"os"
)

//func main() {
//	confPath := flag.String("conf", `.\configuration\config.json`, "flag to set the path to the configuration json file")
//	flag.Parse()
//
//	config, _ := configuration.ExtractConfiguration(*confPath)
//	fmt.Println("Connecting to database")
//
//	dbhandler, _ := dblayer.NewPersistenceLayer(config.DatabaseType, config.DBConnection)
//
//	httpErrChan, httpLsErrChan := rest.ServeAPI(config.RestfulEndPoint,
//		config.RestfulTLSEndPint, dbhandler)
//	select {
//	case err := <-httpErrChan:
//		log.Fatal("HTTP error: ", err)
//	case err := <-httpLsErrChan:
//		log.Fatal("HTTPS error: ", err)
//	}
//}

func main() {
	amqpURL := os.Getenv("AMQP_URL")
	if amqpURL == "" {
		amqpURL = "amqp://guest:guest@localhost:5672"
	}

	connection, err := amqp.Dial(amqpURL)
	if err != nil {
		panic("could not establish AMQP connection: " + err.Error())
	}

	channel, err := connection.Channel()
	err = channel.ExchangeDeclare("events", "topic", true, false, false, false, nil)

	if err != nil {
		panic(err)
	}

	message := amqp.Publishing{
		Body: []byte("Hello World"),
	}
	err = channel.Publish("events", "some-routing-key", false, false, message)
	if err != nil {
		panic("error while publishing message: " + err.Error())
	}

	_, err = channel.QueueDeclare("my_queue", true, false, false, false, nil)
	if err != nil {
		panic("error while declaring the queue: " + err.Error())
	}

	err = channel.QueueBind("my_queue", "#", "events", false, nil)
	if err != nil {
		panic("error while binding the queue: " + err.Error())
	}
	msgs, err := channel.Consume("my_queue", "", false, false, false, false, nil)
	if err != nil {
		panic("error while consuming the queue: " + err.Error())
	}
	for msg := range msgs {
		fmt.Println("message received: " + string(msg.Body))
		msg.Ack(false)
	}
}
