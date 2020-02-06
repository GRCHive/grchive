package main

import (
	"bytes"
	"flag"
	"gitlab.com/grchive/grchive/core"
	"os"
	"os/exec"
	"strings"
)

var supportedApplicationMIME []string = []string{
	"application/msword",
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
	"application/vnd.oasis.opendocument.presentation",
	"application/vnd.oasis.opendocument.spreadsheet",
	"application/vnd.oasis.opendocument.text",
	"application/pdf",
	"application/vnd.ms-powerpoint",
	"application/vnd.openxmlformats-officedocument.presentationml.presentation",
	"application/x-sh",
	"application/xhtml+xml",
	"application/vnd.ms-excel",
	"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
}

func isFileSupported(filename string) bool {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("./src/preview_generator/mime_helper",
		"--filename", filename)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		core.Error("Failed to get MIME: " + err.Error() + " " + stderr.String())
	}

	mime := stdout.String()
	if strings.HasPrefix(mime, "text/") || strings.HasPrefix(mime, "image/") {
		return true
	}

	for _, m := range supportedApplicationMIME {
		if strings.TrimSpace(mime) == m {
			return true
		}
	}

	return false
}

func main() {
	inputFilename := flag.String("input", "", "Filename for the file to generate a preview for.")
	outputDirectory := flag.String("outputDir", "", "Directory to output the preview file in.")
	flag.Parse()

	_, err := os.Stat(*inputFilename)
	if *inputFilename == "" || os.IsNotExist(err) {
		os.Exit(1)
	}

	_, err = os.Stat(*outputDirectory)
	if *outputDirectory == "" || os.IsNotExist(err) {
		os.Exit(1)
	}

	if !isFileSupported(*inputFilename) {
		core.Error("Unsupported filetype.")
	}

	cmd := exec.Command("libreoffice",
		"--headless",
		"--convert-to", "pdf",
		"--outdir", *outputDirectory,
		*inputFilename)

	err = cmd.Run()
	if err != nil {
		core.Error("Failed to run: " + err.Error())
	}
}
