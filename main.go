package main

import (
	"github.com/gorilla/mux"
	"go-cloudnative/handler"
	"net/http"
)

func main() {

}

func ServeAPI(endpoint string) error {
	handler := handler.EventServiceHandler{}
	r := mux.NewRouter()
	eventsrouter := r.PathPrefix("/events").Subrouter()

	eventsrouter.Methods("GET").Path("/{SearchCriteria}/{search}").HandlerFunc(handler.FindEventHandler)
	eventsrouter.Methods("GET").Path("").HandlerFunc(handler.AllEventHandler)
	eventsrouter.Methods("POST").Path("").HandlerFunc(handler.NewEventHandler)

	return http.ListenAndServe(endpoint, r)
}
