// This file contains functions responsible for establishing and maintaining network connections.
// We will use http3 with a cert generated in encryption.go
/*
client verifies with tls.config
client sends package
server responds with 200OK and response in body.
*/

package phomeCore

import (
	"github.com/quic-go/quic-go/http3"
	"log"
	"net/http"
	"strconv"
)

func handshake(w http.ResponseWriter, r *http.Request) {
	//TODO: Verify peer identity on client and server (2 way)
	// requireanyclientcert
	// verifyconnection
	switch r.Method {
	case http.MethodPost:
		log.Println("Received position update request, validating...")
		//TODO: handle request
	case http.MethodGet:
		log.Println("Received request to verify UUID")
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func BeginHTTP3(certFile string, keyFile string, portNumber int) {
	mux := http.NewServeMux()
	mux.Handle("/handshake", http.HandlerFunc(handshake))

	hostAdr := "localhost:" + strconv.Itoa(portNumber)

	if err := http3.ListenAndServe(hostAdr, certFile, keyFile, mux); err != nil {
		log.Fatal(err)
	}
}