package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
)

func handleHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodConnect {
		handleConnect(w, r)
	} else {
		fmt.Fprint(w, "Tunnel by HTTP CONNECT")
	}
}

func handleConnect(w http.ResponseWriter, r *http.Request) {
	remoteConn, err := net.Dial("tcp", r.Host)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to dial destination host: %s", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, fmt.Sprintf("failed to get hijacker: %s", err), http.StatusInternalServerError)
		return
	}

	localConn, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to hijack conn: %s", err), http.StatusInternalServerError)
	}

	go func() {
		io.Copy(localConn, remoteConn)
		remoteConn.Close()
	}()
	go func() {
		io.Copy(remoteConn, localConn)
		localConn.Close()
	}()

}

func main() {
	var listenAddr string
	flag.StringVar(&listenAddr, "addr", ":8080", "listening address")
	flag.Parse()

	s := &http.Server{
		Addr:    listenAddr,
		Handler: http.HandlerFunc(handleHTTP),
	}
	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
