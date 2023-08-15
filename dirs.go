package main

import (
	"os"
	"path/filepath"
)

type Directories struct {
	// Storage for HTTP3 TLS keys and own UUID.
	Certificates string

	// Per-UUID storage of peer keys.
	PairedDevices string

	//XDG dirs
	Cache  string
	Config string
	Data   string
}

func GetDirectories() (Directories, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return Directories{}, err
	}

	xdgDirs := map[string]string{
		"XDG_CACHE_HOME":  filepath.Join(homeDir, ".cache"),
		"XDG_CONFIG_HOME": filepath.Join(homeDir, ".config"),
		"XDG_DATA_HOME":   filepath.Join(homeDir, ".local", "share"),
	}

	for varName := range xdgDirs {
		if value, ok := os.LookupEnv(varName); ok {
			xdgDirs[varName] = value
		}
	}

	dirs := Directories{
		Cache:  filepath.Join(xdgDirs["XDG_CACHE_HOME"], "phome"),
		Config: filepath.Join(xdgDirs["XDG_CONFIG_HOME"], "phome"),
		Data:   filepath.Join(xdgDirs["XDG_DATA_HOME"], "phome"),
	}

	dirs.Certificates = filepath.Join(dirs.Data, "Certificates")
	dirs.PairedDevices = filepath.Join(dirs.Data, "PairedDevices")
	return dirs, nil
}
