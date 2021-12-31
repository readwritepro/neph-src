//=============================================================================
// File:     init-command.go
// Contents: Initialize a remote host by sending it this executable, configs and scripts
//=============================================================================

package main

import "fmt"

// Handle "neph init remoteHost"
func commandInit(host string, options []string) Exitcode {

	fmt.Printf("--- Begin neph init on %s ---\n", host)
	defer fmt.Printf("--- End neph init on %s ---\n", host)

	writeHello(host)
	return SUCCESS
}
