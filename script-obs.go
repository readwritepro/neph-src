package main

import (
	"errors"
	"os"
	"path/filepath"

	eh "github.com/readwritepro/error-handler"
	"github.com/readwritepro/figtree"
)

// The Script type is a set of CLI commands that may be used in connection with figtree parameters
type Script struct {
	scriptName string          // The script name is the basename of the file that contains the script
	scriptPath string          // like /var/neph/scripts/scriptName
	figtree    *figtree.Branch // script root
}

// Script files are created by a text editor and placed in the /var/neph/scripts directory
// Script files are written using figtree syntax
// The LoadScript function reads the script file with the given name
// Returns a Script object with the figtree field pointing to the in-memory representation of the script
// Returns an error if the script file does not exist, or is not in figtree syntax
func LoadScript(scriptName string) (*Script, error) {

	scriptPath := filepath.Join("/var/neph/scripts", scriptName)
	if _, err := os.Stat(scriptPath); errors.Is(err, os.ErrNotExist) {
		return nil, err
	}

	script := &Script{
		scriptName: scriptName,
		scriptPath: scriptPath,
	}

	root, err := figtree.ReadConfig(scriptPath)
	if eh.Invalid(err, scriptPath) {
		return nil, err
	} else {
		script.figtree = root
	}

	return script, nil
}
