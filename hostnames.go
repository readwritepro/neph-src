//=============================================================================
// File:     hostnames.go
// Contents: Read the /etc/neph/conf/hostnames figtree
//=============================================================================

/*
hostnames {
    nk024       165.227.3.8
    nk025       165.227.11.3
    nk026       138.68.26.133
    nk027       178.128.74.100
    nk028       167.99.98.215
}
*/

package main

import (
	"fmt"

	"github.com/readwritepro/figtree"
)

// Get the IPv4 address of the given hostname by reading it from /etc/neph/conf/hostnames
func GetHostname(hostname string) (string, Exitcode) {
	root, err := figtree.ReadConfig(HOSTNAMES_CONF)
	if err != nil {
		fmt.Printf("unable to read hostnames configuration from %s\n", HOSTNAMES_CONF)
		return "", NEPH_CONFIG_MISSING
	}

	key := fmt.Sprintf("/hostnames/%s", hostname)
	ipAddress, err := root.GetValue(key)
	if err != nil {
		fmt.Printf("hostname '%s' not listed in %s configuration\n", hostname, HOSTNAMES_CONF)
		return "", NEPH_CONFIG_MISSING
	}

	return ipAddress, SUCCESS
}

// Get a map of all configured hosts => IP address
func GetAllHostnames() (map[string]string, Exitcode) {
	configuredHosts := make(map[string]string)

	root, err := figtree.ReadConfig(HOSTNAMES_CONF)
	if err != nil {
		fmt.Printf("unable to read hostnames configuration from %s\n", HOSTNAMES_CONF)
		return configuredHosts, NEPH_CONFIG_MISSING
	}

	hostnamesItem, err := root.QueryOne("/hostnames")
	if err == figtree.ErrNotFound {
		fmt.Printf("%s is missing the all important 'hostnames' section\n", HOSTNAMES_CONF)
		return configuredHosts, NEPH_CONFIG_ERROR
	}

	hostnamesBranch, _ := hostnamesItem.Branch()
	for _, item := range hostnamesBranch.Items {
		hostname := item.Key()
		ipAddress, err := item.Value()
		if err == nil {
			configuredHosts[hostname] = ipAddress
		}
	}

	return configuredHosts, SUCCESS
}
