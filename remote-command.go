//=============================================================================
// File:     remote-command.go
// Contents: CLI command handlers
//=============================================================================

package main

import (
	"bytes"
	"fmt"

	"golang.org/x/crypto/ssh"
)

// Contact the remote host via SSH and run the specified neph CLI command
// The nephCommand is something like "neph command options"
// Returns the final exitCode of the remote neph CLI command
func remoteNephCommand(remoteHost string, nephCommand string) uint {

	clientConn, exitCode := connectViaSSH(remoteHost)
	if exitCode != SUCCESS {
		return exitCode
	}
	defer clientConn.Close()

	session, err := clientConn.NewSession()
	if err != nil {
		fmt.Printf("failed to create session: %v\n", err)
		return SSH_SESSION_FAILURE
	}
	defer session.Close()

	fmt.Printf("--- Begin remote neph command on %s ---\n", remoteHost)
	defer fmt.Printf("--- End remote neph command on %s ---\n", remoteHost)

	var b bytes.Buffer
	session.Stdout = &b

	err = session.Run(nephCommand)
	if err != nil {
		if err.Error() == "Process exited with status 127" {
			fmt.Printf("'%s' can't be executed until you setup the remote host with 'neph init %s'\n", nephCommand, remoteHost)
		} else {
			fmt.Printf("failed to run '%s' on '%s': %v\n", nephCommand, remoteHost, err)
		}
		ee := err.(*ssh.ExitError)
		rc := uint(ee.Waitmsg.ExitStatus())
		return rc
	}

	fmt.Printf("%s", b.String())
	return SUCCESS
}
