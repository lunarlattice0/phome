// This file contains functions responsible for establishing and maintaining network connections and encryption.

package phomeCore

// We will be using chacha20-poly1305 and  key exchange is out of band.

import (
	"crypto/chacha20poly1305"
	"encoding/base64"
	"net"
	"time"
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


