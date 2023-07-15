// This file is for handling Peers, including pairing and storage of existing paired devices.
package phomeCore

import (
	"encoding/json"
	"github.com/google/uuid"
	"log"
)

type PairingJSON struct { // Do we want to add a common name?
	UUID   string
	PubKey string
}

// Note: GenCerts also generates the localhost UUID.
func GenerateUUID() string {
	id := uuid.New()
	return id.String()
}

// This function generates the initial pairing JSON.
// It is recommended to convert the string output to base64.
func GeneratePairingJSON(newPairingJSON *PairingJSON) string {
	jsonStr, err := json.Marshal(newPairingJSON)
	if err != nil {
		log.Fatal(err)
	}

	return string(jsonStr)
}

// This function unmarshals a peer's pairing JSON and returns a pairing JSON struct.
func DecodePairingJson(pairingJSONstr string) *PairingJSON {
	newPeerPairing := PairingJSON{}
	err := json.Unmarshal([]byte(pairingJSONstr), &newPeerPairing)
	if err != nil {
		log.Fatal(err)
	}
	return &newPeerPairing
}