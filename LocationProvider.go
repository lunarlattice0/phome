// This file connects to the GeoClue dbus service to acquire the user's location.

package main

import (
	"github.com/godbus/dbus/v5"
)

var existingDbusConnection *dbus.Conn

func beginDbusConnectionBus() (dbus.Conn, error) {
	if existingDbusConnection == nil {
		conn, err := dbus.ConnectSessionBus()
		if err != nil {
			return dbus.Conn{}, err
		}
		defer conn.Close()
		return *conn, nil
	} else {
		return *existingDbusConnection, nil
	}
}

func createGeoClueClient (conn dbus.Conn) (string) {
	obj := conn.Object("org.freedesktop.GeoClue2", "/org/freedesktop/GeoClue2/Manager")
	obj.Call("org.freedesktop.GeoClue2.Manager.CreateClient", 0, "")
	var s string
	obj.Call("org.freedesktop.GeoClue2.Manager.GetClient", 0, "").Store(&s)
	return s
}

func deleteGeoClueClient (conn dbus.Conn, destroyClient string) { // This may fail?
	obj := conn.Object("org.freedesktop.GeoClue2", "/org/freedesktop/GeoClue2/Manager")
	obj.Call("org.freedesktop.GeoClue2.Manager.DeleteClient", 0, destroyClient)
}

func getLocation (lat string, long string, error) {
	dbusConnection, err := beginDbusConnectionBus()
	if err != nil {
		return "", "", err
	}
	clientName := createGeoClueClient(dbusConnection)
	defer deleteGeoClueClient(dbusConnection, clientName)

	client := conn.Object("org.freedesktop.GeoClue2", clientName)
	// We need to set DesktopId


}





