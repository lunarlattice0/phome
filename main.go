//This is a reference CLI wrapper for phomeCore

package main

import (
	pc "github.com/Thelolguy1/phome/phomeCore"
	"log"
	"os"
	"path/filepath"
	"fmt"
	"strconv"
)

func usage() {
		fmt.Fprintln(os.Stderr, "usage: phome [client|server]")
		fmt.Fprintln(os.Stderr, "	phome [server] [port number]")
		fmt.Fprintln(os.Stderr, "hint: you can run the client and server in different terminals simultaneously")
		os.Exit(1)
}

func main() {
	log.Println("phome development build (NOT FOR PRODUCTION USE)")
	if len(os.Args) < 2 {
		usage()
	}

	dirs := GetDirectories()
	switch os.Args[1] {
		case "server":
			if len(os.Args) < 3 {
				usage()
			}
			pN, err := strconv.Atoi(os.Args[2])
			if err != nil {
				log.Fatal(err)
			}

			log.Println("Starting HTTP3 listener on port " + os.Args[2])
			beginListener(&dirs, pN)
		default:
			usage()
	}
}

func beginListener(dirs *Directories, port int) {
	_, err := os.Stat(filepath.Join(dirs.Certificates, "cert.pem"))
	_, err2 := os.Stat(filepath.Join(dirs.Certificates, "key.pem"))

	if (err != nil || err2 != nil) {
		log.Println("Generating HTTP certificates for the first time...")
		pc.GenCerts(dirs.Certificates)
	}

	cert := filepath.Join(dirs.Certificates, "cert.pem")
	key := filepath.Join(dirs.Certificates, "key.pem")

	pc.BeginHTTP3(cert, key, port)
}