package configuration

import (
	"encoding/json"
	"fmt"
	"gozilla/dblayer"
	"os"
)

var (
	DBTypeDefault       = dblayer.DBTYPE("mongodb")
	DBConnectionDefault = "mongodb://127.0.0.1"
	RestfulEPDefault    = "localhost:8181"
	RestfulTLSEPDefault = "localhost:9191"
)

type ServiceConfig struct {
	DatabaseType      dblayer.DBTYPE `json:"databaseType"`
	DBConnection      string         `json:"dbconnection"`
	RestfulEndPoint   string         `json:"restfulapi_endpoint"`
	RestfulTLSEndPint string         `json:"restfulapi-tlsendpint"`
}

func ExtractConfiguration(filename string) (ServiceConfig, error) {
	conf := ServiceConfig{
		DBTypeDefault,
		DBConnectionDefault,
		RestfulEPDefault,
		RestfulTLSEPDefault,
	}
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Configuration file not found. Continuing with default values.")
		return conf, err
	}
	err = json.NewDecoder(file).Decode(&conf)
	return conf, err
}
