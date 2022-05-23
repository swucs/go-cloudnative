package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"go-cloudnative/configuration"
	"go-cloudnative/dblayer"
	"go-cloudnative/handler"
	"go-cloudnative/repository"
	"log"
	"net/http"
)

func main() {
	confPath := flag.String("conf", `.\configuration\config.json`, "flag to set the path to the configuration json file")
	flag.Parse()
	//extract configuration
	config, _ := configuration.ExtractConfiguration(*confPath)
	fmt.Println("connecting to database")
	dbhandler, _ := dblayer.NewPersistenceLayer(config.Databasetype, config.DBConnection)
	//reseful API start
	log.Fatal(ServeAPI(config.RestfulEndpoint, dbhandler))
}

func ServeAPI(endpoint string, dbHandler repository.DatabaseHandler) error {
	handler := handler.NewEventHandler(dbHandler)
	r := mux.NewRouter()
	eventsrouter := r.PathPrefix("/events").Subrouter()

	eventsrouter.Methods("GET").Path("/{SearchCriteria}/{search}").HandlerFunc(handler.FindEventHandler)
	eventsrouter.Methods("GET").Path("").HandlerFunc(handler.AllEventHandler)
	eventsrouter.Methods("POST").Path("").HandlerFunc(handler.NewEventHandler)

	return http.ListenAndServe(endpoint, r)
}
