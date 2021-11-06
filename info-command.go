//=============================================================================
// File:     info-command.go
// Contents: Obtain info about /etc/neph/conf/* or /var/neph/scripts
//=============================================================================

package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
)

// Handle "neph info hosts [remoteHost]"
// Get the hostnames and the IP addresses configured for neph usage
func commandInfoHosts(host string, options []string) uint {
	if host == "localhost" {
		configuredHosts, exitCode := GetAllHostnames()
		if exitCode == SUCCESS {
			for hostname, ipAddress := range configuredHosts {
				fmt.Printf("%s %s\n", hostname, ipAddress)
			}
		}
		return exitCode
	} else {
		return remoteNephCommand(host, "neph info hosts")
	}
}

// Handle "neph info configs [remoteHost]"
// List all of the config files in /etc/neph/conf/
func commandInfoConfigs(host string, options []string) uint {
	if host == "localhost" {
		return walkDir("/etc/neph/conf")
	} else {
		return remoteNephCommand(host, "neph info configs")
	}
}

// Handle "neph info scripts [remoteHost]"
// List all of the scripts in /var/neph/scripts/
func commandInfoScripts(host string, options []string) uint {
	if host == "localhost" {
		return walkDir("/var/neph/scripts")
	} else {
		return remoteNephCommand(host, "neph info scripts")
	}
}

// walk the directory, printing the filenames found
func walkDir(dir string) uint {
	dirEntries, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Printf("Unable to list info for %s: %v\n", dir, err)
		return FS_FAILURE
	}

	for _, entry := range dirEntries {
		if !isHiddenFile(entry.Name()) {
			fullFilename := filepath.Join(dir, entry.Name())
			if entry.IsDir() {
				exitCode := walkDir(fullFilename)
				if exitCode != SUCCESS {
					return exitCode
				}
			} else {
				fmt.Printf("%s\n", fullFilename)
			}
		}
	}
	return SUCCESS
}

// Returns true if the filename begins with a dot
func isHiddenFile(filename string) bool {
	if filename[0:1] == "." {
		return true
	} else {
		return false
	}
}
