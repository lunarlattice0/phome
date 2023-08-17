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

func getLocation (conn dbus.Conn) (lat string, long string, err error) {
	dbusConnection, err := beginDbusConnectionBus()
	if err != nil {
		return "", "", err
	}
	clientName := createGeoClueClient(dbusConnection)
	defer deleteGeoClueClient(dbusConnection, clientName)

	client := conn.Object("org.freedesktop.GeoClue2" , dbus.ObjectPath(clientName))
	// We need to set DesktopId
	client.Call("org.freedesktop.DBus.Properties.Set", 0, "('org.freedesktop.GeoClue2.Client', 'DesktopId', <'io.github.Thelolguy1.phome'>)")
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





