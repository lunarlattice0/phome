// This file connects to the GeoClue dbus service to acquire the user's location.

package main

import (
	"github.com/godbus/dbus/v5"
	"log"
)

var existingDbusConnection *dbus.Conn

func beginDbusConnectionBus() (dbus.Conn, error) {
	if existingDbusConnection == nil {
		conn, err := dbus.ConnectSystemBus()
		if err != nil {
			return dbus.Conn{}, err
		}
		existingDbusConnection = conn
		return *conn, nil
	} else {
		return *existingDbusConnection, nil
	}
}

func createGeoClueClient (conn dbus.Conn) (string) {
	obj := conn.Object("org.freedesktop.GeoClue2", "/org/freedesktop/GeoClue2/Manager")
	var s string
	err := obj.Call("org.freedesktop.GeoClue2.Manager.CreateClient", 0).Store(&s)
	if err != nil {
		log.Fatal(err)
	}
	return s
}

func deleteGeoClueClient (conn dbus.Conn, destroyClient string) { // This may fail?
	obj := conn.Object("org.freedesktop.GeoClue2", dbus.ObjectPath("/org/freedesktop/GeoClue2/Manager"))
	obj.Call("org.freedesktop.GeoClue2.Manager.DeleteClient", 0, destroyClient)
}

func getLocation () (lat string, long string, err error) {
	conn, err := beginDbusConnectionBus()
	if err != nil {
		return "", "", err
	}
	clientName := createGeoClueClient(conn)
	defer deleteGeoClueClient(conn, clientName)
	client := conn.Object("org.freedesktop.GeoClue2" , dbus.ObjectPath(clientName))
	// We need to set DesktopId
	err = client.Call("org.freedesktop.DBus.Properties.Set", 0, "('DesktopId', <'io.github.Thelolguy1.phome'>)").Err
	if err != nil {
		log.Fatal("SNEED")
	}
	// yuck
	client.Call("org.freedesktop.GeoClue2.Start", 0, "")

	//get location client
	var locationObjectPath string
	client.Call("org.freedesktop.DBus.Properties.Get", 0, "('org.freedesktop.GeoClue2.Client', 'Location'").Store(&locationObjectPath)

	locationObject := conn.Object("org.freedesktop.GeoClue2", dbus.ObjectPath(locationObjectPath))

	locationObject.Call("org.freedesktop.GeoClue2.Location", 0, "Latitude").Store(&lat)
	locationObject.Call("org.freedesktop.GeoClue2.Location", 0, "Longitude").Store(&long)

	return
}





