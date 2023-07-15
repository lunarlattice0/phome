// This file contains functions responsible for establishing and maintaining network connections.
// We will use http3 with a cert generated in encryption.go
/*
client requests uuid from server
client verifies TLS cert of HTTP with uuid/tls pairing
client sends signed base64 json data package of uuid + data
server verifies data package with pem
server responds with 200 OK with empty body or 200OK with signed packet in body.
*/

package phomeCore

import (
	"encoding/base64"
	"github.com/quic-go/quic-go/http3"
	"log"
	"net/http"
	"strconv"
)

func EncodeB64(in string) string { // Output can be used in a external program, like a QR generator.
	return base64.URLEncoding.EncodeToString([]byte(in))
}

func DecodeB64(in string) string {
	data, err := base64.URLEncoding.DecodeString(in)
	if err != nil {
		log.Fatal(err)
	}
	return string(data)
}

func handshake(w http.ResponseWriter, r *http.Request) {
	//TODO: Verify peer identity on client and server (2 way)
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