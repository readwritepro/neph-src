//=============================================================================
// File:     commands.go
// Contents: CLI command handlers
//=============================================================================

package main

func commandInit(host string, options []string) uint {
	return SUCCESS
}

func commandPush(host string, options []string) uint {
	return SUCCESS
}

func commandPull(host string, options []string) uint {
	return SUCCESS
}

func commandScrub(host string, options []string) uint {
	return SUCCESS
}

func commandInfoConfigs(host string, options []string) uint {
	return SUCCESS
}

func commandInfoScripts(host string, options []string) uint {
	return SUCCESS
}

func commandApply(host string, options []string) uint {
	return SUCCESS
}

func commandExamine(host string, options []string) uint {
	return SUCCESS
}

func commandExecScript(host string, options []string) uint {

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
