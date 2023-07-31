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
	"crypto/tls"
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

	var err error
	certs := make([]tls.Certificate, 1)
	certs[0], err = tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		log.Fatal(err)
	}

	tlsConfig := &tls.Config{
		Certificates: certs,
		ClientAuth: tls.RequireAnyClientCert,
		InsecureSkipVerify: true,
	} // TODO: Actually verify the peer cert on the server...

	tlsConfig = http3.ConfigureTLSConfig(tlsConfig)

	server := http3.Server {
		TLSConfig: tlsConfig,
		Port: portNumber,
		Handler: mux,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}

}