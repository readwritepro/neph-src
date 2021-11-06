//=============================================================================
// File:     help.go
// Contents: Usage text
//=============================================================================

package main

import "fmt"

func printUsage() {
	fmt.Print(`
The neph command installs, configures, and executes cloud setup software on a remote device
using passwordless SSH with root privileges.

Usage 1) neph [init|push|pull|scrub] host
Usage 2) neph info [configs|scripts|hosts] [host|localhost]
Usage 3) neph apply [host|localhost] configfile dtbfile
Usage 4) neph examine [host|localhost] configfile
Usage 5) neph exec [host|localhost] script
Usage 6) neph [version|help]

    init          copy the neph executable, scripts, and figtree files from the local host to the remote host
                  neph init host [--privileged]

    push          copy missing files, update older files, delete obsolete files to the remote host
                  neph push host [--force]

    pull          copy missing files, update older files, delete obsolete files from remote host
                  neph pull host [--force]

    scrub         remove figtree files (from this device) that were used by a former remote host
                  neph scrub host	

    info hosts    list the hostnames and IP addresses known by the specified host
                  neph info hosts [host]

    info configs  list configurations in /etc/neph/conf
                  neph info configs [host]

    info scripts  list scripts in /var/neph/scripts
                  neph info scripts [host]

    apply         apply a DTB (delimited text block) to a config file
                  neph apply host configfile dtbfile

    examine       examine a config file and print the contents of its DTB (delimited text block)
                  neph examine host configfile

    exec          execute the specified script on the local or remote host
                  neph localhost script-file
                  neph remotehost script-file

Options:
    --force      copy, update, and delete scripts and configurations without checking timestamps 
    --privileged elevates the target host to be a privileged device by sending it the private ssh key

File Locations:
    /usr/bin/neph                    CLI executable (chmod 700)
    /etc/neph/conf                   figtree configuration files (chmod 600)
    /root/.ssh/neph-rsa-private-key  PEM formatted SSH key (chmod 600)
    /var/neph/scripts                script files (chmod 700)

`)
}
