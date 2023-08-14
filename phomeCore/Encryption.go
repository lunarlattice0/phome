// This file handles encryption and decryption of TLS
package phomeCore

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"
)

// Remember to check that these paths are valid in your implementation!
type SelfIDs struct {
	CertPath string
	KeyPath  string
}

// GenCerts generates certificates for server TLS and client verification.
func (ids *SelfIDs) GenCerts() error {

	certPemFile := ids.CertPath
	keyFile := ids.KeyPath

	uuidStr := GenerateUUID()

	//Modified from https://go.dev/src/crypto/tls/generate_cert.go
	//Please see SUBLICENSE for licensing of the bottom code.
	_, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return (err)
	}

	keyUsage := x509.KeyUsageDigitalSignature
	//notBefore := time.Now()
	//notAfter := notBefore.Add(99*365*24*time.Hour) // 99 year expiry

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return (err)
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"phomeCoreCert"},
		},
		//	NotBefore: notBefore,
		//	NotAfter:  notAfter,
		KeyUsage:              keyUsage,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA:                  true,
		DNSNames:              []string{uuidStr},
	}
	template.KeyUsage |= x509.KeyUsageCertSign
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, priv.Public().(ed25519.PublicKey), priv)
	if err != nil {
		return (err)
	}

	certOut, err := os.Create(certPemFile)
	if err != nil {
		return (err)
	}

	if err := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}); err != nil {
		return (err)
	}
	if err := certOut.Close(); err != nil {
		return (err)
	}

	keyOut, err := os.OpenFile(keyFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)

	if err != nil {
		return (err)
	}

	privBytes, err := x509.MarshalPKCS8PrivateKey(priv)

	if err != nil {
		return (err)
	}

	if err := pem.Encode(keyOut, &pem.Block{Type: "PRIVATE KEY", Bytes: privBytes}); err != nil {
		return (err)
	}

	if err := keyOut.Close(); err != nil {
		return (err)
	}

	//Note: Neither Firefox nor Chrome(ium) support ED25519, so phomeCore is not web-browser accessible.
	//This should be fine, as we do not plan for browsers to be able to manually send in location reports.
	return nil
}
