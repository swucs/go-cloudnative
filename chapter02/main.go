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
	httpErrChan, httptlsErrChan := ServeAPI(config.RestfulEndpoint, config.RestfulTLSEndpoint, dbhandler)
	select {
	case err := <-httpErrChan:
		log.Fatal("HTTP Error: ", err)
	case err := <-httptlsErrChan:
		log.Fatal("HTTPS Error: ", err)
	}

}

func ServeAPI(endpoint string, tlsendpoint string, dbHandler repository.DatabaseHandler) (chan error, chan error) {
	handler := handler.NewEventHandler(dbHandler)
	r := mux.NewRouter()
	eventsrouter := r.PathPrefix("/events").Subrouter()

	eventsrouter.Methods("GET").Path("/{SearchCriteria}/{search}").HandlerFunc(handler.FindEventHandler)
	eventsrouter.Methods("GET").Path("").HandlerFunc(handler.AllEventHandler)
	eventsrouter.Methods("POST").Path("").HandlerFunc(handler.NewEventHandler)

	httpErrChan := make(chan error)
	httptlsErrChan := make(chan error)

	go func() {
		httptlsErrChan <- http.ListenAndServeTLS(tlsendpoint, "cert.pem", "key.pem", r)
	}()

	go func() {
		httpErrChan <- http.ListenAndServeTLS(endpoint, "cert.pem", "key.pem", r)
	}()

	return httpErrChan, httptlsErrChan
}
