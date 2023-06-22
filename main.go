//This is a reference CLI wrapper for phomeCore

package main

import (
	pc "github.com/Thelolguy1/phome/phomeCore"
	"log"
	"os"
	"path/filepath"
)

func main() {
	log.Println("Phome Started")
	dirs := GetDirectories()

	_, err := os.Stat(filepath.Join(dirs.Certificates, "cert.pem"))
	_, err2 := os.Stat(filepath.Join(dirs.Certificates, "key.pem"))

	if (err != nil && err2 != nil) {
		log.Println("Generating HTTP certificates for the first time...")
		pc.GenCerts(dirs.Certificates)
	}
}
