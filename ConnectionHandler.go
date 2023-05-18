// This file contains functions responsible for establishing and maintaining network connections and encryption.

package phomeCore

import (
	"github.com/quic-go/quic-go/http3"
	"encoding/base64"
	"net"
	"log"
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

func BeginHTTP3Listener() {
	err := http3.ListenAndServe() 
	if err != nil {
		log.Fatal("Failed to start HTTP3 server!")
	}
}

func DialHttp3() {
	//stub
}

