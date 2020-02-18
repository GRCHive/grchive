package main

import (
	"bytes"
	"flag"
	"gitlab.com/grchive/grchive/core"
	"os"
	"os/exec"
	"path/filepath"
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

type FileSupport struct {
	Supported bool
	IsPdf     bool
}

func isFileSupported(filename string) FileSupport {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("dependencies/python/python-3.8.1/bin/bin/python3",
		"src/preview_generator/mime_helper",
		"--filename", filename)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Env = append(os.Environ(),
		"PYTHONHOME=./dependencies/python/python-3.8.1/bin/lib/python3.8",
		"PYTHONPATH=./dependencies/python/python-3.8.1/bin/lib/python3.8")
	err := cmd.Run()
	if err != nil {
		core.Error("Failed to get MIME: " + err.Error() + " " + stderr.String())
	}

	mime := stdout.String()
	if strings.HasPrefix(mime, "text/") || strings.HasPrefix(mime, "image/") {
		return FileSupport{false, false}
	}

	for _, m := range supportedApplicationMIME {
		if strings.TrimSpace(mime) == m {
			return FileSupport{true, m == "application/pdf"}
		}
	}

	return FileSupport{false, false}
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

	support := isFileSupported(*inputFilename)
	if !support.Supported {
		core.Error("Unsupported filetype.")
	}

	var cmd *exec.Cmd
	if support.IsPdf {
		outputFilename := *outputDirectory + "/" + filepath.Base(*inputFilename)
		cmd = exec.Command("gs",
			"-sDEVICE=pdfwrite",
			"-dCompatibilityLevel=1.4",
			"-dPDFSettings=/ebook",
			"-dNOPAUSE",
			"-dQUIET",
			"-dBATCH",
			"-sOutputFile="+outputFilename,
			*inputFilename)
	} else {
		cmd = exec.Command("libreoffice",
			"--headless",
			"--convert-to", "pdf",
			"--outdir", *outputDirectory,
			*inputFilename)
	}
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		core.Warning("STDOUT: " + stdout.String())
		core.Warning("STDERR: " + stderr.String())
		core.Error("Failed to run: " + err.Error())
	}
}
