//=============================================================================
// File:     delimited-block.go
// Contents: Read and replace HEPH delimited blocks of text from a config file
//=============================================================================

package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

// Scan the given file looking for a NEPH delimited block.
// Returns the text within the delimited block or an empty string if no such block exists
func getDelimitedBlock(targetFile string) (string, error) {

	if _, err := os.Stat(targetFile); errors.Is(err, os.ErrNotExist) {
		fmt.Printf("getDelimitedBlock: no such file %s\n", targetFile)
		return "", err
	}

	inFile, err := os.Open(targetFile)
	if err != nil {
		fmt.Printf("getDelimitedBlock: can't open file %s\n", targetFile)
		return "", err
	}
	defer inFile.Close()

	// create a scanner that uses the "ScanLines" splitter
	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)

	var blockText = ""
	var bInsideBlock = false
	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "#-----BEGIN NEPH-----") {
			bInsideBlock = true
		} else if strings.HasPrefix(line, "#-----END NEPH-----") {
			bInsideBlock = false
		} else if bInsideBlock {
			blockText += line + "\n"
		}
	}
	return blockText, nil
}

// Replace the given file's existing NEPH delimited block with the provided blockText
// If the file doesn't have a NEPH delimited block, append a new block at the end of the file
func replaceDelimitedBlock(targetFile string, blockText string) error {
	if _, err := os.Stat(targetFile); errors.Is(err, os.ErrNotExist) {
		fmt.Printf("replaceDelimitedBlock: no such file %s\n", targetFile)
		return err
	}

	inFile, err := os.Open(targetFile)
	if err != nil {
		fmt.Printf("replaceDelimitedBlock: can't open file for reading %s\n", targetFile)
		return err
	}
	defer inFile.Close()

	tempFile := targetFile + ".tmp"
	outFile, err := os.Create(tempFile)
	if err != nil {
		fmt.Printf("replaceDelimitedBlock: can't open tmp file for writing %s\n", tempFile)
		return err
	}
	defer outFile.Close()

	scanner := bufio.NewScanner(inFile)
	writer := bufio.NewWriter(outFile)

	var bInsideBlock = false
	var bBlockWritten = false
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "#-----BEGIN NEPH-----") {
			bInsideBlock = true
			writer.WriteString("#-----BEGIN NEPH-----\n")
			writer.WriteString(blockText)
			writer.WriteString("#-----END NEPH-----\n")
			bBlockWritten = true
		} else if strings.HasPrefix(line, "#-----END NEPH-----") {
			bInsideBlock = false
		} else if !bInsideBlock {
			writer.WriteString(line + "\n")
		}
	}

	if !bBlockWritten {
		writer.WriteString("#-----BEGIN NEPH-----\n")
		writer.WriteString(blockText)
		writer.WriteString("#-----END NEPH-----\n")
		bBlockWritten = true
	}
	writer.Flush()

	saveFile := targetFile + ".bak"
	if _, err = os.Stat(saveFile); err == nil {
		err = os.Remove(saveFile)
		if err != nil {
			fmt.Printf("replaceDelimitedBlock: unable to remove previous backup %s\n", saveFile)
		}
	}
	err = os.Rename(targetFile, saveFile)
	if err != nil {
		fmt.Printf("replaceDelimitedBlock: unable to save %s to %s\n", targetFile, saveFile)
	}
	err = os.Rename(tempFile, targetFile)
	if err != nil {
		fmt.Printf("replaceDelimitedBlock: unable to save %s to %s\n", tempFile, targetFile)
	}

	return nil
}
