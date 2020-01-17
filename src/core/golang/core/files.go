package core

import (
	"errors"
	"io/ioutil"
	"path/filepath"
)

func FindFirstFileInDirectory(dir string) (string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return "", err
	}

	if len(files) == 0 {
		return "", errors.New("No files")
	}

	return filepath.Join(dir, files[0].Name()), nil
}
