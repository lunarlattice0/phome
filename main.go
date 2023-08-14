//This is a reference CLI wrapper for phomeCore
package main

import (
	"fmt"
	pc "github.com/Thelolguy1/phome/phomeCore"
	"log"
	"os"
	"path/filepath"
	"io/fs"
	"crypto/x509"
	"encoding/pem"
)

var selfIDs = pc.SelfIDs{}
var peerMap = make(map[string]string)

func usage() {
	fmt.Fprintln(os.Stderr, "usage: phome [client | server | showpair | newpair | regenerate]")
	fmt.Fprintln(os.Stderr, "       phome [server] [IP:port]")
	fmt.Fprintln(os.Stderr, "       phome [client] [IP:port]")
	fmt.Fprintln(os.Stderr, "       phome [newpair] [pairing code of other device]")
	fmt.Fprintln(os.Stderr, "       hint: use 0.0.0.0 to accept all connections. (e.g. 0.0.0.0:60300)")
	os.Exit(1)
}

func ensureCertsExist(dirs *Directories) {
	err := os.MkdirAll(dirs.Certificates, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	selfIDs.CertPath = filepath.Join(dirs.Certificates, "cert.pem")
	selfIDs.KeyPath = filepath.Join(dirs.Certificates, "key.pem")

	// We will use cert.pem as our pairing pubkey and for TLS
	// This will be shared out-of-band in b64 form for pairing to ensure authenticity.
	_, err = os.Stat(selfIDs.CertPath)

	// Private key
	_, err2 := os.Stat(selfIDs.KeyPath)

	if err != nil || err2 != nil {
		selfIDs.GenCerts()
	}
}

// TODO: Move this to phomeCore for gobind/cgo support.
// Loads known Peer Uuids and certs into a map for the server for fast access.
func loadPeerUuidCerts (dirs *Directories) {
	err := filepath.Walk(dirs.PairedDevices, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			log.Fatal(err)
		}

		if info.IsDir() == true && info.Name() != "PairedDevices" {
			peerUuidFileData, err := os.ReadFile(filepath.Join(path, "cert.pem"))
			if err != nil {
				log.Fatal(err)
			}
			peerMap[info.Name()] = string(peerUuidFileData)
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
}

// While this reference client uses go hashmaps to return client certificates, you may use any other method.
func lookupPeerMap (peerUuid string) (cert string) {
	return peerMap[peerUuid]
}

func main() {
	if len(os.Args) < 2 {
		usage()
	}

	dirs := GetDirectories()
	switch os.Args[1] {
	case "regenerate":
		if err := os.RemoveAll(dirs.Certificates); err != nil {
			log.Fatal(err)
		}
		ensureCertsExist(&dirs)
	case "showpair":
		ensureCertsExist(&dirs)

		pubKeyFile := selfIDs.CertPath
		pubKeyData, err := os.ReadFile(pubKeyFile)
		if err != nil {
			log.Fatal(err)
		}
		newPairingJSON := pc.JSONBundle{PubKey: string(pubKeyData)}

		pairingJSONB64 := pc.EncodeB64(newPairingJSON.GenerateJSON())
		fmt.Println(pairingJSONB64)
	case "newpair":
		if len(os.Args) < 3 {
			usage()
		}

		ensureCertsExist(&dirs)

		peerPairingStr := pc.DecodeB64(os.Args[2])
		newPeerPairing := new(pc.JSONBundle)
		newPeerPairing.DecodeJSON(peerPairingStr)

		// Uuid Decoding order
		// PEM >> PKCS8 (ASN1) >> Certificate.DNSName (uuid)

		// Own uuid
		selfCertPEM, err := os.ReadFile(selfIDs.CertPath)
		if err != nil {
			log.Fatal(err)
		}

		block, _ := pem.Decode(selfCertPEM)
		if block == nil {
			log.Fatal("No public key found in own certificate. Please regenerate!")
		}

		selfCert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			log.Fatal(err)
		}

		selfUuid := selfCert.DNSNames[0]

		// peer uuid
		block, _ = pem.Decode([]byte(newPeerPairing.PubKey))
		if block == nil {
			log.Fatal("No public key found in peer's certificate!")
		}

		peerCert, err := x509.ParseCertificate([]byte(block.Bytes))
		if err != nil {
			log.Fatal(err)
		}
		peerUuid := peerCert.DNSNames[0]

		//We don't care about matching certs because the probability is so low.
		if peerUuid == string(selfUuid) {
			fmt.Fprintln(os.Stderr, "Peer has same Uuid as this computer. Please check your entry or regenerate certificates.")
			os.Exit(-1)
		}

		err = os.MkdirAll(dirs.PairedDevices, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}

		peerCertDir := filepath.Join(dirs.PairedDevices, peerUuid)
		peerCertFile := filepath.Join(dirs.PairedDevices, peerUuid, "cert.pem")

		// check if peer directory already exists
		if _, err = os.Stat(peerCertFile); !os.IsNotExist(err) {
			fmt.Fprintln(os.Stderr, "Peer already paired!")
			os.Exit(-1)
		}

		err = os.MkdirAll(peerCertDir, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}

		peerCertBytes := []byte(newPeerPairing.PubKey)
		if err := os.WriteFile(peerCertFile, peerCertBytes, 0600); err != nil {
			log.Fatal(err)
		}

		fmt.Println("Successfully paired and stored new peer.")
		//fmt.Println("Uuid:" + newPeerPairing.Uuid)
		//fmt.Println("Cert:\n" + newPeerPairing.PubKey)
	case "server":
		if len(os.Args) < 3 {
			usage()
		}

		ensureCertsExist(&dirs)
		loadPeerUuidCerts(&dirs)
		address := string(os.Args[2])

		log.Println("Starting HTTP listener on port " + address)


		cert := selfIDs.CertPath
		key := selfIDs.KeyPath

		pc.BeginHTTP(cert, key, address, lookupPeerMap)
	case "client":
		if len(os.Args) < 3 {
			usage()
		}
		ensureCertsExist(&dirs)
		loadPeerUuidCerts(&dirs)
		//resolve peer address
		//NOTE: We use IP:Port format but http.Client takes "https://IP:Port"
		addr := "https://" + string(os.Args[2])
		pc.BeginClientPeer(
			selfIDs.CertPath,
			selfIDs.KeyPath,
			addr,
			lookupPeerMap,
		)
	default:
		usage()
	}
}
