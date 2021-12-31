//=============================================================================
// File:     sftp.go
// Contents: Transfer files using SFTP
//=============================================================================

package main

import (
	"log"

	"github.com/pkg/sftp"
)

func writeHello(host string) Exitcode {
	clientConn, exitCode := connectViaSSH(host)
	if exitCode != SUCCESS {
		return exitCode
	}
	defer clientConn.Close()

	// open an SFTP session over an existing ssh connection.
	sftpClient, err := sftp.NewClient(clientConn)
	if err != nil {
		log.Fatal(err)
	}
	defer sftpClient.Close()

	// walk a directory
	w := sftpClient.Walk("/home/fuji1021")
	for w.Step() {
		if w.Err() != nil {
			continue
		}
		log.Println(w.Path())
	}

	// leave your mark
	f, err := sftpClient.Create("hello.txt") // in /root/hello.txt
	if err != nil {
		log.Fatal(err)
	}
	if _, err := f.Write([]byte("Hello world!")); err != nil {
		log.Fatal(err)
	}
	f.Close()

	// check it's there
	fi, err := sftpClient.Lstat("hello.txt")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(fi)

	return SUCCESS
}
