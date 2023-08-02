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
	"net"
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

	// Adapted from quic-go ListenAndServe (??? implementation)
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

	quicServer := &http3.Server{
		TLSConfig: tlsConfig,
		Handler: mux,
	}

	httpServer := &http.Server {
		//TLSConfig: tlsConfig,
		//Addr: address,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			quicServer.SetQuicHeaders(w.Header())
			mux.ServeHTTP(w, r)
		}),
	}

	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		log.Fatal(err)
	}
	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer udpConn.Close()

	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	tcpConn, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer tcpConn.Close()

	tlsConn := tls.NewListener(tcpConn, tlsConfig)
	defer tlsConn.Close()

	hErr := make(chan error)
	qErr := make(chan error)
	go func() {
		hErr <- httpServer.Serve(tlsConn)
	}()
	go func() {
		qErr <- quicServer.Serve(udpConn)
	}()

	select {
	case err := <-hErr:
		quicServer.Close()
		log.Fatal(err)
	case err := <-qErr:
		log.Fatal(err)
	}
}