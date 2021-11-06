//=============================================================================
// File:     ssh.go
// Contents: Connect to remote via SSH
//=============================================================================

package main

import (
	"fmt"
	"io/ioutil"
	"os/exec"

	"golang.org/x/crypto/ssh"
)

// connect to remote host via SSH
// returns a clientConnection and an exitCode
// The caller must Close the ssh.Clinet connection when finished using it
func connectViaSSH(host string) (*ssh.Client, uint) {

	// The server can be contacted, without authentication, using the ssh-keyscan utility, which will
	// retreive its public host key, which is stored on the server at /etc/ssh/ssh_host_rsa_key.pub
	// The -t flag specifies the key type: rsa, ecdsa, ed25519
	// The return value is a string with three parts: hostname, key type, hostPublicKey, like "nk024 ssh-rsa AAAAB3Nz...CjV+IgUn"
	sshKeyscanPath, err := findOrInstall("ssh-keyscan", "openssh-clients")
	if err != nil {
		fmt.Printf("ssh-keyscan utility not found and not installed: %v\n", err)
		return nil, SSH_LOCAL_CONFIGURATION_FAILURE
	}
	knownHostsEntry, err := exec.Command(sshKeyscanPath, "-t", "rsa", host).Output()
	if err != nil {
		fmt.Printf("ssh-keyscan was not able to obtain RSA type public host key, from %s: %v\n", host, err)
		fmt.Print("check the remote host's configuration at /etc/ssh/sshd_config and make sure the entry 'HostKey /etc/ssh/ssh_host_rsa_key' references an OPENSSH PRIVATE KEY\n")
		return nil, SSH_REMOTE_CONFIGURATION_FAILURE
	}

	// A host key is a cryptographic key used for authenticating computers in the SSH protocol.
	// Host keys are key pairs, typically using the RSA, DSA, or ECDSA algorithms.
	// Public host keys are stored on and/or distributed to SSH clients,
	// and private keys are stored on SSH servers.
	_, _, hostPublicKey, _, _, err := ssh.ParseKnownHosts([]byte(knownHostsEntry))
	if err != nil {
		fmt.Printf("failed to parse knownHostsEntry %s: %v\n", knownHostsEntry, err)
		return nil, SSH_REMOTE_CONFIGURATION_FAILURE
	}

	// Read the contents of the root user's PEM encoded private key file on the local host
	userPrivateKey, err := ioutil.ReadFile("/root/.ssh/neph-rsa-private-key")
	if err != nil {
		fmt.Printf("unable to read PEM encoded private key %s: %v\n", "/root/.ssh/neph-rsa-private-key", err)
		return nil, SSH_LOCAL_CONFIGURATION_FAILURE
	}
	// Parse the PEM encoded file to get the signer
	signer, err := ssh.ParsePrivateKey(userPrivateKey)
	if err != nil {
		fmt.Printf("unable to parse PEM encoded private key %s: %v\n", userPrivateKey, err)
		return nil, SSH_LOCAL_CONFIGURATION_FAILURE
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
		fmt.Printf("failed to dial %s: %v\n", hostWithPort, err)
		return nil, SSH_CONNECTION_FAILURE
	}

	return clientConn, SUCCESS
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
