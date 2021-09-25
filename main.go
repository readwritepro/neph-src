//=============================================================================
// File:     main.go
// Contents: Executable entry point handling argv
//=============================================================================

package main

import (
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
)

func main() {

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	// determine which command line pattern to follow
	var pattern string
	var argdump string
	for _, argv := range os.Args[1:] {
		if isCommand(argv) {
			pattern += "command "
			argdump += "command:" + argv + "\n"
		} else if isSubCommand(argv) {
			pattern += "subcommand "
			argdump += "subcommand:" + argv + "\n"
		} else if isMetaCommand(argv) {
			pattern += "metacommand "
			argdump += "metacommand:" + argv + "\n"
		} else if isLocalhost(argv) {
			pattern += "localhost "
			argdump += "localhost:" + argv + "\n"
		} else if isRemotehost(argv) {
			pattern += "host "
			argdump += "host:" + argv + "\n"
		} else if isScript(argv) {
			pattern += "script "
			argdump += "script:" + argv + "\n"
		} else if isOption(argv) {
			pattern += "option "
			argdump += "option:" + argv + "\n"
		} else {
			pattern += "unrecognized"
			argdump += "unrecognized:" + argv + "\n"
		}
	}

	var exitCode uint

	if strings.HasPrefix(pattern, "command host") ||
		strings.HasPrefix(pattern, "command localhost") {
		exitCode = executeCommand(os.Args[1], os.Args[2], os.Args[3:])

	} else if strings.HasPrefix(pattern, "command subcommand host") ||
		strings.HasPrefix(pattern, "command subcommand localhost") {
		exitCode = executeSubCommand(os.Args[1], os.Args[2], os.Args[3], os.Args[4:])

	} else if strings.HasPrefix(pattern, "metacommand") {
		exitCode = executeMetaCommand(os.Args[1], os.Args[2:])

	} else {
		fmt.Printf("Unable to figure out what to do with argument pattern '%s'\n", pattern)
		fmt.Printf("%s", argdump)
		fmt.Printf("Try neph help\n")
		exitCode = CLI_BAD_ARGUMENTS
	}
	ec := int(exitCode)
	os.Exit(ec)
}

func isCommand(argv string) bool {
	switch argv {
	case "init", "push", "pull", "scrub", "info", "apply", "examine", "exec":
		return true
	default:
		return false
	}
}

func isSubCommand(argv string) bool {
	switch argv {
	case "configs", "scripts":
		return true
	default:
		return false
	}
}

func isMetaCommand(argv string) bool {
	switch argv {
	case "help", "version":
		return true
	default:
		return false
	}
}

func isOption(argv string) bool {
	switch argv {
	case "--privileged", "--force":
		return true
	default:
		return false
	}
}

func isScript(argv string) bool {
	scriptPath := filepath.Join("/var/neph/scripts", argv)
	if _, err := os.Stat(scriptPath); errors.Is(err, os.ErrNotExist) {
		return false
	} else {
		return true
	}
}

func isLocalhost(argv string) bool {
	if argv == "localhost" {
		return true
	}
	if argv == "127.0.0.1" {
		return true
	}
	if hostname, _ := os.Hostname(); hostname == argv {
		return true
	}
	return false
}

func isRemotehost(argv string) bool {
	iprecords, err := net.LookupIP(argv)
	if err != nil {
		return false
	}
	for _, ip := range iprecords {
		if string(ip) == "127.0.0.1" {
			return false
		}
	}
	return true
}

// execute a neph command
// returns an exitcode where 0 is success, anything else is a failure
func executeCommand(command string, host string, options []string) uint {
	switch command {
	case "init":
		return commandInit(host, options)

	case "push":
		return commandPush(host, options)

	case "pull":
		return commandPull(host, options)

	case "scrub":
		return commandScrub(host, options)

	case "apply":
		return commandApply(host, options)

	case "examine":
		return commandExamine(host, options)

	case "exec":
		return commandExecScript(host, options)
	}

	fmt.Printf("Unhandled command '%s'\n", command)
	fmt.Printf("Try neph help\n")
	return CLI_BAD_ARGUMENTS
}

// execute a neph command
// returns an exitcode where 0 is success, anything else is a failure
func executeSubCommand(command string, subcommand string, host string, options []string) uint {
	cmd := fmt.Sprintf("%s %s", command, subcommand)

	switch cmd {
	case "info configs":
		return commandInfoConfigs(host, options)

	case "info scripts":
		return commandInfoScripts(host, options)
	}

	fmt.Printf("Unhandled command '%s %s'\n", command, subcommand)
	fmt.Printf("Try neph help\n")
	return CLI_BAD_ARGUMENTS
}

// help and version metacommands
func executeMetaCommand(metaCommand string, options []string) uint {
	switch metaCommand {
	case "version":
		fmt.Printf("neph version %s\n", NEPH_VERSION)
		return SUCCESS

	case "help":
		printUsage()
		return SUCCESS

	default:
		fmt.Printf("Unknown meta command %s\n", metaCommand)
		return CLI_BAD_ARGUMENTS
	}
}
