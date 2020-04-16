package core

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
)

var (
	ErrNoFilesFound   error = errors.New("No files.")
	ErrDirectoryFound       = errors.New("Directory found when file expected.")
	ErrSymLinkFound         = errors.New("Symbolic link found.")
)

func FindFirstFileInDirectory(dir string) (string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return "", err
	}

	if len(files) == 0 {
		return "", ErrNoFilesFound
	}

	return filepath.Join(dir, files[0].Name()), nil
}

func CopyFile(src string, dst string) error {
	// We can ignore the case where the two files are the same.
	// The only thing to be wary of is when src is a symbolic link.
	if src == dst {
		return nil
	}

	// If an error pops up here it probably means
	// that the input file doesn't exist.
	info, err := os.Stat(src)
	if err != nil {
		return err
	}

	if info.Mode()&os.ModeSymlink != 0 {
		return ErrSymLinkFound
	}

	if info.IsDir() {
		return ErrDirectoryFound
	}

	inData, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(dst, inData, info.Mode())
	if err != nil {
		return err
	}

	return nil
}

func FindAbsolutePathThroughSymbolicLink(path string) (string, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return path, err
	}

	if info.Mode()&os.ModeSymlink != 0 {
		path, err = os.Readlink(path)
		if err != nil {
			return path, err
		}
	}

	return path, nil
}
