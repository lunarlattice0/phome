// This file contains functions responsible for establishing and maintaining network connections.
// We will use http3 with a cert generated in encryption.go

package phomeCore

import (
	"github.com/quic-go/quic-go/http3"
	"encoding/base64"
	"log"
	"net/http"
	"strconv"
	)

func EncodeB64(in string) string { // Output can be used in a external program, like a QR generator.
	return base64.URLEncoding.EncodeToString([]byte (in))
}

func DecodeB64(in string) (string, error) {
	data, err := base64.URLEncoding.DecodeString(in)
	if (err != nil) {
		return "", err
	} else {
		return string(data), nil
	}
}

func acceptUpload(w http.ResponseWriter, r *http.Request) {
	log.Println("Received position update request, validating...")

}

func BeginHTTP3(certFile string, keyFile string, portNumber int) {
	mux := http.NewServeMux()
	mux.Handle("/upload", http.HandlerFunc(acceptUpload))

	hostAdr := "localhost:" + strconv.Itoa(portNumber)

	if err := http3.ListenAndServe(hostAdr, certFile, keyFile, mux); err != nil {
		log.Fatal(err)
	}
}