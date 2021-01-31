package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

func handleHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	if r.Method == "CONNECT" {
		handleConnect(w, r)
	} else {
		fmt.Fprint(w, "Tunnel by HTTP CONNECT")
	}
}

func handleConnect(w http.ResponseWriter, r *http.Request) {

}

func main() {
	var listenAddr string
	flag.StringVar(&listenAddr, "addr", ":8080", "listening address")
	flag.Parse()

	http.HandleFunc("/", handleHTTP)

	if err := http.ListenAndServe(listenAddr, nil); err != nil {
		log.Fatal(err)
	}
}
