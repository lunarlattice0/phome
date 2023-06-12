// This file contains functions responsible for establishing and maintaining network connections.

package phomeCore

import (
	"net/http"
	"encoding/base64"
	"golang.org/x/net/websocket"
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

func ClientWS(hostIP string){
	origin := "http://localhost"
	url := "ws://" + hostIP + ":56000/ws"
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Println("Failed to connect to host")
	} else {
		log.Println("Connected to " + hostIP)
	}

	if _, err := ws.Write([]byte("hello, world!\n")); err != nil {
		log.Println("Failed to message " + hostIP)
	}

	var retmsg = make([]byte, 512)
	var n int

	if n, err = ws.Read(retmsg); err != nil {
		log.Println("Couldn't read returned message")
	}
	log.Println(retmsg[:n])
}

func HostWS() {
	es := func(ws *websocket.Conn){ //stub
		log.Println(ws)
	}

	http.Handle("/ws", websocket.Handler(es))
	err := http.ListenAndServe(":5600", nil)
	if err != nil {
		panic("HTTPServer says:" + err.Error())
	}
}
