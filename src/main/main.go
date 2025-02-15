package main

import (
	"flag"
	"fmt"
	"github.com/VMerlin/receipt-processor/src/processor"
	"log"
	"net/http"
)

var port = flag.String("port", "8080", "endpoint port")

func main() {
	flag.Parse()
	service := processor.New()
	err := StartServer(*port, service)
	if err != nil {
		log.Fatal(fmt.Sprint(" Encountered Issue when starting server", err))
	}
}

func StartServer(port string, service processor.Service) error {
	httpServer := &http.Server{
		Addr: ":" + port,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			processor.Handle(service, w, r)
		}),
	}
	fmt.Printf(" Starting server on %s\n", port)
	return httpServer.ListenAndServe()
}
