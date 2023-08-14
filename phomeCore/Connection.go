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
	"errors"
	"encoding/pem"
	"bytes"
)

func handshake(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		log.Println("Received position update request, validating...")
	default:
		log.Println("Dropped invalid request")
		//Silent Drop
		http.Error(w, "", http.StatusBadRequest)
	}
}

// Compare the certificate received from the server with the saved certificate of this UUID.
// An MITM should not be possible since the server must be hosted with the same cert that signed the JSON.
func PCVerifyConnection (rawCerts [][]byte, knownCerts map[string]string) (error) {
	pubKeyBlock := &pem.Block {
		Type: "CERTIFICATE",
		Bytes: rawCerts[0],
	}
	pubKeyPEM := string(pem.EncodeToMemory(pubKeyBlock))
	block, _ := pem.Decode([]byte(pubKeyPEM))

	if block == nil {
		return errors.New("The received certificate did not contain a public key.")
	}


	peerCert, err := x509.ParseCertificate([]byte(block.Bytes))
	if err != nil {
		return errors.New("The received certificate is not a valid X.509 certificate.")
	}

	peerUuid := peerCert.DNSNames[0]
	cachedPeerPEM := knownCerts[peerUuid]

	// No idea if this is the most efficient way, but this is likely safe.

	if len(cachedPeerPEM) == len(pubKeyPEM) {
		for i := range cachedPeerPEM {
			if cachedPeerPEM[i] != pubKeyPEM[i] {
				return errors.New("The received and stored certificates do not match.")
			}
		}
	} else if len(cachedPeerPEM) == 0 {
		return errors.New("The certificate received is from an unpaired device and cannot be used.")
	} else if len(pubKeyPEM) == 0 {
		return errors.New("The peer sent a blank certificate.")
	} else {
		return errors.New("The received and stored certificates are of different lengths and do not match.")
	}
	return nil
}

func BeginClientPeer(certFile string, keyFile string, addr string, knownUuids map[string]string) {
	//generate client TLS config
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
		VerifyPeerCertificate: func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) (error) {
			return PCVerifyConnection(rawCerts, knownUuids)
		},
	}

	//generate http3 roundtripper config
	roundTripper := &http3.RoundTripper{
		TLSClientConfig: tlsConfig,
	}

	client := &http.Client{
		Transport: roundTripper,
	}
	// TEST
	testJSONStruct := JSONBundle{
		Test: "SNEEDFEEDSEED",
	}

	bodyReader := bytes.NewReader([]byte(testJSONStruct.GenerateJSON()))
	resp, err := client.Post(addr, "application/json", bodyReader)
	log.Println(resp) // debug
	if err != nil {
		log.Fatal(err)
	}
	// TEST
}

func BeginHTTP(certFile string, keyFile string, addr string, knownUuids map[string]string) {
	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(handshake))

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
		VerifyPeerCertificate: func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
			return PCVerifyConnection(rawCerts, knownUuids)
		},
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
