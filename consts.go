//=============================================================================
// File:     consts.go
// Contents: global constants
//=============================================================================

package main

var NEPH_VERSION = "0.0.2"

const (
	SUCCESS                          uint = iota // 0 = everything worked
	BASH_SCRIPT_FAILED                           // 1 = typical failure code coming from a Bash script or executable, never explicilty used by Neph
	SSH_LOCAL_CONFIGURATION_FAILURE              // 2 = SSH on the localhost not configured to be used by Neph
	SSH_REMOTE_CONFIGURATION_FAILURE             // 3 = SSH on the remote host not setup to acept connections
	SSH_CONNECTION_FAILURE                       // 4 = Connecting to remote host via SSH didn't succeed
	SSH_SESSION_FAILURE                          // 5 = Connecting to remote host via SSH didn't succeed
	NEPH_SCRIPT_MISSING                          // 6 = A script in /var/neph/scripts doesn't exist
	NEPH_SCRIPT_NOT_EXECUTABLE                   // 7 = A script in /var/neph/scripts isn't executable
	NEPH_LOGIC_ERROR                             // 8 = Seemingly impossible to happen
	CLI_BAD_ARGUMENTS                            // 9 = Arguments to the Neph CLI rejected
)
