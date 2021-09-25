//=============================================================================
// File:     exec-script.go
// Contents: Execute local and remote Bash scripts located in /var/neph/scripts
//=============================================================================

package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"golang.org/x/crypto/ssh"
)

// execute a neph script on the localhost
func executeLocalScript(localScript string, options []string) uint {
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
		return uint(cmd.ProcessState.ExitCode())
	}
	return SUCCESS
}

// execute a neph script on a host
func executeRemoteScript(host string, remoteScript string, options []string) uint {

	// The server can be contacted, without authentication, using the ssh-keyscan utility, which will
	// retreive its public host key, which is stored on the server at /etc/ssh/ssh_host_rsa_key.pub
	// The -t flag specifies the key type: rsa, ecdsa, ed25519
	// The return value is a string with three parts: hostname, key type, hostPublicKey, like "nk024 ssh-rsa AAAAB3Nz...CjV+IgUn"
	sshKeyscanPath, err := findOrInstall("ssh-keyscan", "openssh-clients")
	if err != nil {
		fmt.Printf("ssh-keyscan utility not found and not installed: %v\n", err)
		return SSH_LOCAL_CONFIGURATION_FAILURE
	}
	knownHostsEntry, err := exec.Command(sshKeyscanPath, "-t", "rsa", host).Output()
	if err != nil {
		fmt.Printf("ssh-keyscan was not able to obtain RSA type public host key, from %s: %v\n", host, err)
		fmt.Print("check the remote host's configuration at /etc/ssh/sshd_config and make sure the entry 'HostKey /etc/ssh/ssh_host_rsa_key' references an OPENSSH PRIVATE KEY\n")
		return SSH_REMOTE_CONFIGURATION_FAILURE
	}

	// A host key is a cryptographic key used for authenticating computers in the SSH protocol.
	// Host keys are key pairs, typically using the RSA, DSA, or ECDSA algorithms.
	// Public host keys are stored on and/or distributed to SSH clients,
	// and private keys are stored on SSH servers.
	_, _, hostPublicKey, _, _, err := ssh.ParseKnownHosts([]byte(knownHostsEntry))
	if err != nil {
		fmt.Printf("failed to parse knownHostsEntry %s: %v\n", knownHostsEntry, err)
		return SSH_REMOTE_CONFIGURATION_FAILURE
	}

	// Read the contents of the root user's PEM encoded private key file on the local host
	userPrivateKey, err := ioutil.ReadFile("/root/.ssh/neph-rsa-private-key")
	if err != nil {
		fmt.Printf("unable to read PEM encoded private key %s: %v\n", "/root/.ssh/neph-rsa-private-key", err)
		return SSH_LOCAL_CONFIGURATION_FAILURE
	}
	// Parse the PEM encoded file to get the signer
	signer, err := ssh.ParsePrivateKey(userPrivateKey)
	if err != nil {
		fmt.Printf("unable to parse PEM encoded private key %s: %v\n", userPrivateKey, err)
		return SSH_LOCAL_CONFIGURATION_FAILURE
	}

	// On the remote server, the public key must be copied to a file within the user's home directory at /root/. ssh/authorized_keys.
	// (With Digital Ocean, this is done during droplet provisioning.)
	// The authorized_keys file contains a list of public keys, one-per-line, that are authorized to log into this account.
	clientConfig := &ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.FixedHostKey(hostPublicKey),
	}
	hostWithPort := host + ":22"
	clientConn, err := ssh.Dial("tcp", hostWithPort, clientConfig)
	if err != nil {
		fmt.Printf("failed to dial %s: %v", hostWithPort, err)
		return SSH_CONNECTION_FAILURE
	}

	if !remoteCommandExists(clientConn, remoteScript) {
		return NEPH_SCRIPT_MISSING
	}
	if !remoteCommandExecutable(clientConn, remoteScript) {
		return NEPH_SCRIPT_NOT_EXECUTABLE
	}

	return executeRemoteCommand(clientConn, remoteScript)
}

// Check to see if the script exists on the remote host
// returns true if script exists
func remoteCommandExists(clientConn *ssh.Client, remoteScript string) bool {
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
func remoteCommandExecutable(clientConn *ssh.Client, remoteScript string) bool {
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

// Run a single command, using an SSH Session, on an established SSH clientConnetion
// returns true if the command was executed and returned 0
func executeRemoteCommand(clientConn *ssh.Client, remoteScript string) uint {
	session, err := clientConn.NewSession()
	if err != nil {
		fmt.Printf("failed to create session: %v\n", err)
		return SSH_SESSION_FAILURE
	}
	defer session.Close()

	fmt.Printf("\n--- Begin remote script %s ---\n", remoteScript)
	defer fmt.Printf("\n--- End remote script %s ---\n", remoteScript)

	var b bytes.Buffer
	session.Stdout = &b
	scriptPath := filepath.Join("/var/neph/scripts", remoteScript)
	if err := session.Run(scriptPath); err != nil {
		fmt.Printf("failed to run %s: %v\n", scriptPath, err)
		ee := err.(*ssh.ExitError)
		rc := uint(ee.Waitmsg.ExitStatus())
		return rc
	}
	fmt.Printf("%s", b.String())
	return SUCCESS
}

// Get the path to the specified tool. If it is not found, attempt to install it.
// cliTool is the executable file name
// distributionPackage is the DNF package that installs the executable file
// returns the path to the executable
// returns an error if it can't be found and wasn't installed
func findOrInstall(cliTool string, distributionPackage string) (string, error) {
	cliToolPath, err := exec.LookPath(cliTool)
	if err != nil {
		fmt.Printf("%s\nAttempting to install %s...\n", err, cliTool)

		cmd := exec.Command("dnf", "install", "-y", distributionPackage)
		stdout, err := cmd.Output()
		fmt.Print(string(stdout))
		if err != nil {
			fmt.Printf("installation of %s failed\n%v", distributionPackage, err)
			return "", err
		}

		cliToolPath, err = exec.LookPath(cliTool)
		if err != nil {
			fmt.Printf("%s not found\n %v", cliTool, err)
			return "", err
		}
	}
	return cliToolPath, nil
}
