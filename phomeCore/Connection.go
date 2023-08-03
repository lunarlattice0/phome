// This file contains functions responsible for establishing and maintaining network connections.
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
	"crypto/x509"
)

func handshake(w http.ResponseWriter, r *http.Request) {
	//TODO: Verify peer identity on client and server (2 way)
	// requireanyclientcert
	// verifyconnection
	switch r.Method {
	case http.MethodPost:
		log.Println("Received position update request, validating...")
	default:
		log.Println("Dropped invalid request")
		//Silent Drop
		http.Error(w, "", http.StatusBadRequest)
	}
}

func verifyPeer(rawCerts [][]byte, _ [][]*x509.Certificate) (error) {
	// We have no chain of trust, so we cannot establish verified chains.
	currentCert := rawCerts[0]
	log.Println(currentCert)
	return nil
}

func BeginHTTP(certFile string, keyFile string, addr string) {
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
		ClientAuth: tls.RequireAndVerifyClientCert, // changeme
		InsecureSkipVerify: true,
		//VerifyPeerCertificate: verifyPeer,
	}

	h3Server := &http3.Server{
		Addr: addr,
		TLSConfig: tlsConfig,
		Handler: mux,
	}

	httpServer := &http.Server {
		Addr: addr,
		TLSConfig: h3Server.TLSConfig,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h3Server.SetQuicHeaders(w.Header())
			mux.ServeHTTP(w, r)
		}),
	}

	go func() {
		httpServer.ListenAndServeTLS("", "") // provided by tlsconfig

	}()
	if err := h3Server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}