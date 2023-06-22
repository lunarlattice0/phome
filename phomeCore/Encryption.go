// This file handles encryption and decryption of TLS and JWT

package phomeCore

import (
	"crypto/rand"
	"crypto/ed25519"
	"crypto/x509"
	"crypto/x509/pkix"
	"log"
	"time"
	"math/big"
	"os"
	"encoding/pem"
	"path/filepath"
)

func GenCerts(targetDir string) {
	err := os.MkdirAll(targetDir, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}



	var CertPemFile = filepath.Join(targetDir, "cert.pem")
	var KeyFile = filepath.Join(targetDir, "key.pem")

	//Modified from https://go.dev/src/crypto/tls/generate_cert.go
	//Please see SUBLICENSE for licensing of the bottom code.
	_, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		log.Fatalf("Failed to generate ed25519 key: %v", err)
	}

	keyUsage := x509.KeyUsageDigitalSignature
	notBefore := time.Now()
	notAfter := notBefore.Add(99999)

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		log.Fatalf("Failed to generate serial number: %v", err)
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"phomeCoreCert"},
		},
		NotBefore: notBefore,
		NotAfter:  notAfter,
		KeyUsage:              	keyUsage,
		ExtKeyUsage:           	[]x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: 	true,
		IsCA:			true,
	}
	template.KeyUsage |= x509.KeyUsageCertSign
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, priv.Public().(ed25519.PublicKey), priv)
	if err != nil {
		log.Fatalf("Failed to create certificate: %v", err)
	}

	certOut, err := os.Create(CertPemFile)
	if err != nil {
		log.Fatalf("Failed to open cert.pem for writing: %v", err)
	}

	if err := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}); err != nil {
		log.Fatalf("Failed to write data to cert.pem: %v", err)
	}
	if err := certOut.Close(); err != nil {
		log.Fatalf("Error closing cert.pem: %v", err)
	}
	log.Println("Wrote cert.pem")

	keyOut, err := os.OpenFile(KeyFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)

	if err != nil {
		log.Fatalf("Failed to open key.pem for writing: %v", err)
	}

	privBytes, err := x509.MarshalPKCS8PrivateKey(priv)

	if err != nil {
		log.Fatalf("Unable to marshal private key: %v", err)
	}

	if err := pem.Encode(keyOut, &pem.Block{Type: "PRIVATE KEY", Bytes: privBytes}); err != nil {
		log.Fatalf("Failed to write data to key.pem: %v", err)
	}

	if err := keyOut.Close(); err != nil {
		log.Fatalf("Error closing key.pem: %v", err)
	}

	log.Print("Wrote key.pem\n")
}