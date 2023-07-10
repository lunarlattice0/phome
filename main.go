//This is a reference CLI wrapper for phomeCore

package main

import (
	"fmt"
	pc "github.com/Thelolguy1/phome/phomeCore"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

func usage() {
	fmt.Fprintln(os.Stderr, "usage: phome [client | server | showpair | newpair | regenerate]")
	fmt.Fprintln(os.Stderr, "       phome [server] [port number]")
	fmt.Fprintln(os.Stderr, "       phome [newpair] [pairing code of other device]")
	fmt.Fprintln(os.Stderr, "hint: you can run the client and server in as separate processes simultaneously")
	os.Exit(1)
}

func ensureCertsExist(dirs *Directories) {
	// We will use cert.pem as our pairing pubkey and for TLS
	// This will be shared out-of-band in b64 form for pairing to ensure authenticity.
	_, err := os.Stat(filepath.Join(dirs.Certificates, "cert.pem"))

	// Private key
	_, err2 := os.Stat(filepath.Join(dirs.Certificates, "key.pem"))

	// Own UUID
	_, err3 := os.Stat(filepath.Join(dirs.Certificates, "uuid"))
	if err != nil || err2 != nil || err3 != nil {
		pc.GenCerts(dirs.Certificates)
	}
}

func beginListener(dirs *Directories, port int) {
	ensureCertsExist(dirs)

	cert := filepath.Join(dirs.Certificates, "cert.pem")
	key := filepath.Join(dirs.Certificates, "key.pem")

	pc.BeginHTTP3(cert, key, port)
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

		pubkeyFile := filepath.Join(dirs.Certificates, "cert.pem")
		pubKeyData, err := os.ReadFile(pubkeyFile)
		if err != nil {
			log.Fatal(err)
		}

		uuidFile := filepath.Join(dirs.Certificates, "uuid")
		uuidFileData, err := os.ReadFile(uuidFile)
		if err != nil {
			log.Fatal(err)
		}

		newPairingJSON := pc.PairingJSON{PubKey: string(pubKeyData), UUID: string(uuidFileData)}

		pairingJSON := pc.GeneratePairingJSON(&newPairingJSON)
		pairingJSONB64 := pc.EncodeB64(pairingJSON)
		fmt.Println(pairingJSONB64)
	case "newpair":
		if len(os.Args) < 3 {
			usage()
		}

		ensureCertsExist(&dirs)

		uuidFile := filepath.Join(dirs.Certificates, "uuid")
		uuidFileData, err := os.ReadFile(uuidFile)
		if err != nil {
			log.Fatal(err)
		}

		peerPairingStr := pc.DecodeB64(os.Args[2])
		newPeerPairing := pc.DecodePairingJson(peerPairingStr)

		//We don't care about matching certs because the probability is so low.
		if newPeerPairing.UUID == string(uuidFileData) {
			fmt.Fprintln(os.Stderr, "phome: peer has same UUID as this computer, please regenerate certificates and UUIDs on either device.")
			os.Exit(-1)
		}

		err = os.MkdirAll(dirs.PairedDevices, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}

		peerCertDir := filepath.Join(dirs.PairedDevices, newPeerPairing.UUID)
		peerCertFile := filepath.Join(dirs.PairedDevices, newPeerPairing.UUID, "cert.pem")

		// check if peer directory already exists
		if _, err = os.Stat(peerCertFile); !os.IsNotExist(err) {
			fmt.Fprintln(os.Stderr, "phome: peer already paired!")
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
		//fmt.Println("UUID:" + newPeerPairing.UUID)
		//fmt.Println("Cert:\n" + newPeerPairing.PubKey)
	case "server":
		if len(os.Args) < 3 {
			usage()
		}
		pN, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatal(err) // deferring the error to strconv, so that I don't have to use reflect package
		}

		log.Println("Starting HTTP3 listener on port " + os.Args[2])
		beginListener(&dirs, pN)
	default:
		usage()
	}
}
