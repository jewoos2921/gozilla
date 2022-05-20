package rest

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"gozilla/persistence"
	"net/http"
	"strings"
)

type EventServiceHandler struct {
	dbhandler persistence.DatabaseHandler
}

func (eh *EventServiceHandler) findEventHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	criteria, ok := vars["SearchCriteria"]
	if !ok {
		w.WriteHeader(400)
		fmt.Fprint(w, `{error: No search criteria found, you can either search by id via /id/4 to search by name via /name/coldplayconcert}`)
		return
	}
	searchkey, ok := vars["search"]
	if !ok {
		w.WriteHeader(400)
		fmt.Fprint(w, `{error: No search criteria found, you can either search by id via /id/4 to search by name via /name/coldplayconcert}`)
		return
	}
	var event persistence.Event
	var err error
	switch strings.ToLower(criteria) {
	case "name":
		event, err = eh.dbhandler.FindEventByName(searchkey)
	case "id":
		id, err := hex.DecodeString(searchkey)
		if err == nil {
			event, err = eh.dbhandler.FindEvent(id)
		}
	}
	if err != nil {
		fmt.Fprintf(w, "{error %s}", err)
		return
	}
	w.Header().Set("Content-Type", "application/json;charset=utf8")
	json.NewEncoder(w).Encode(&event)
}
func (eh *EventServiceHandler) allEventHandler(w http.ResponseWriter, r *http.Request) {
	events, err := eh.dbhandler.FindAllAvailableEvents()
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "{error: Error occured while trying to find all available events %s}", err)
		return
	}
	w.Header().Set("Content-Type", "application/json;charset=utf8")
	err = json.NewEncoder(w).Encode(&events)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "{error: Error occured while trying encode events to JSON %s}", err)
		return
	}
}
func (eh *EventServiceHandler) newEventHandler(w http.ResponseWriter, r *http.Request) {
	event := persistence.Event{}
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "{error: Error occured while decoding event data %s}", err)
		return
	}
	id, err := eh.dbhandler.AddEvent(event)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "{error: Error occured while persisting event %d %s}", id, err)
		return
	}
}

func newEventHandler(databasehandler persistence.DatabaseHandler) *EventServiceHandler {
	return &EventServiceHandler{
		dbhandler: databasehandler,
	}
}
func ServeAPI(endpoint, tlsendpoint string, dbhandler persistence.DatabaseHandler) (chan error, chan error) {
	handler := newEventHandler(dbhandler)
	r := mux.NewRouter()
	eventsrouter := r.PathPrefix("/events").Subrouter()

	eventsrouter.Methods("GET").Path("/{SearchCriteria}/{search}").HandlerFunc(handler.findEventHandler)
	eventsrouter.Methods("GET").Path("").HandlerFunc(handler.allEventHandler)
	eventsrouter.Methods("POST").Path("").HandlerFunc(handler.newEventHandler)

	httpErrChan := make(chan error)
	httpLsErrChan := make(chan error)
	go func() { httpLsErrChan <- http.ListenAndServeTLS(tlsendpoint, "cert.pem", "key.perm", r) }()
	go func() { httpErrChan <- http.ListenAndServe(endpoint, r) }()
	return httpErrChan, httpLsErrChan
}
