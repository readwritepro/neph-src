//=============================================================================
// File:     exec-script.go
// Contents: Execute local or remote Bash scripts located in /var/neph/scripts
//=============================================================================

package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"golang.org/x/crypto/ssh"
)

// Handle the "neph exec" CLI
func commandExecScript(host string, options []string) Exitcode {

	if isLocalhost(host) {
		localScript := options[0]
		return executeLocalScript(localScript, options[1:])
	}

	if isRemotehost(host) {
		remoteScript := options[0]
		return executeRemoteScript(host, remoteScript, options[1:])
	}

	return NEPH_LOGIC_ERROR
}

// Execute a neph script on the localhost
func executeLocalScript(localScript string, options []string) Exitcode {
	scriptPath := filepath.Join("/var/neph/scripts", localScript)
	if _, err := os.Stat(scriptPath); errors.Is(err, os.ErrNotExist) {
		fmt.Printf("local script %s does not exist", scriptPath)
		return NEPH_SCRIPT_MISSING
	}

	cmd := exec.Command(scriptPath)

	fmt.Printf("\n--- Begin script %s ---\n", localScript)
	defer fmt.Printf("\n--- End script %s ---\n", localScript)

	out, err := cmd.Output()
	fmt.Printf("%s", string(out))
	if err != nil {
		fmt.Printf("%v", err)
		return Exitcode(cmd.ProcessState.ExitCode())
	}
	return SUCCESS
}

// execute a neph script on aremote  host
func executeRemoteScript(host string, remoteScript string, options []string) Exitcode {

	clientConn, exitCode := connectViaSSH(host)
	if exitCode != SUCCESS {
		return exitCode
	}
	defer clientConn.Close()

	if !remoteScriptExists(clientConn, remoteScript) {
		return NEPH_SCRIPT_MISSING
	}
	if !isRemoteScriptExecutable(clientConn, remoteScript) {
		return NEPH_SCRIPT_NOT_EXECUTABLE
	}

	return doRemoteScript(clientConn, remoteScript)
}

// Check to see if the script exists on the remote host
// returns true if script exists
func remoteScriptExists(clientConn *ssh.Client, remoteScript string) bool {
	session, err := clientConn.NewSession()
	if err != nil {
		fmt.Printf("failed to create session: %v\n", err)
		return false
	}
	defer session.Close()

	scriptPath := filepath.Join("/var/neph/scripts", remoteScript)
	testCommand := fmt.Sprintf("test -f %s", scriptPath)
	err = session.Run(testCommand)
	if err != nil {
		fmt.Printf("remote script %s does not exist\n", scriptPath)
		return false
	} else {
		return true
	}
}

// Check to see if the script on the remote host is executable
// returns true if script is executable
func isRemoteScriptExecutable(clientConn *ssh.Client, remoteScript string) bool {
	session, err := clientConn.NewSession()
	if err != nil {
		fmt.Printf("failed to create session: %v\n", err)
		return false
	}
	defer session.Close()

	scriptPath := filepath.Join("/var/neph/scripts", remoteScript)
	testCommand := fmt.Sprintf("test -x %s", scriptPath)
	err = session.Run(testCommand)
	if err != nil {
		fmt.Printf("remote script is not executable. Try 'chmod +x %s'\n", scriptPath)
		return false
	} else {
		return true
	}
}

// Run the remote script, capture its output, return its exitCode
// Returns true if the script was executed and returned 0
func doRemoteScript(clientConn *ssh.Client, remoteScript string) Exitcode {
	session, err := clientConn.NewSession()
	if err != nil {
		fmt.Printf("failed to create session: %v\n", err)
		return SSH_SESSION_FAILURE
	}
	defer session.Close()

	fmt.Printf("--- Begin remote script %s ---\n", remoteScript)
	defer fmt.Printf("--- End remote script %s ---\n", remoteScript)

	var b bytes.Buffer
	session.Stdout = &b
	scriptPath := filepath.Join("/var/neph/scripts", remoteScript)
	err = session.Run(scriptPath)
	fmt.Printf("%s", b.String())
	if err != nil {
		fmt.Printf("%s didn't exit cleanly: %v\n", scriptPath, err)
		ee := err.(*ssh.ExitError)
		rc := Exitcode(ee.Waitmsg.ExitStatus())
		return rc
	}
	return SUCCESS
}
