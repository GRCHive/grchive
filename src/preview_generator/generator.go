package main

import (
	"flag"
	"github.com/h2non/filetype"
	"gitlab.com/grchive/grchive/core"
	"os"
	"os/exec"
)

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

	file, err := os.Open(*inputFilename)
	if err != nil {
		core.Error("Failed to open file: " + err.Error())
	}

	header := make([]byte, 261)
	_, err = file.Read(header)
	if err != nil {
		core.Error("Failed to read file: " + err.Error())
	}

	if !filetype.IsImage(header) &&
		!filetype.IsDocument(header) {
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
