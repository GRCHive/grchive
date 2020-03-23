package main

import (
	"fmt"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/gcloud_api"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

type LibraryStorage struct {
	LocalPath string
	Checksum  string
}

// Key: version
var libraryMap map[string]LibraryStorage = map[string]LibraryStorage{}
var libraryMutex sync.RWMutex = sync.RWMutex{}

func mustCreateTempLibStorageDir(dir string, err error) string {
	if err != nil {
		core.Error("Failed to create temp dir: " + err.Error())
	}

	// Otherwise the user in the docker container can't read it.
	err = os.Chmod(dir, os.FileMode(0755))
	if err != nil {
		core.Error("Failed to change directory permissions: " + err.Error())
	}

	core.Info("LIBRARY STORAGE: ", dir)
	return dir
}

var libraryStorageDir string = mustCreateTempLibStorageDir(ioutil.TempDir(os.TempDir(), "grchive-library"))

func pullCoreLibraryFromVersion(version string) (string, error) {
	libraryMutex.RLock()
	val, ok := libraryMap[version]
	libraryMutex.RUnlock()

	jarFilename := fmt.Sprintf("grchive_public_core.v%s.jar", version)
	sha256Filename := fmt.Sprintf("grchive_public_core.v%s.sha256", version)
	storage := gcloud.DefaultGCloudApi.GetStorageApi()

	// Always download latest SHA256 so we can check against what our previously downloaded SHA256 is.
	// If they don't match, we want to force ourselves to re-pull.
	latestSha256, err := storage.Download(
		core.EnvConfig.Kotlin.LibraryBucket,
		sha256Filename,
		[]byte{},
	)

	if err != nil {
		return "", err
	}

	forceRepull := (!ok || val.Checksum != string(latestSha256))
	if forceRepull {
		jarData, err := storage.Download(
			core.EnvConfig.Kotlin.LibraryBucket,
			jarFilename,
			[]byte{},
		)

		if err != nil {
			return "", err
		}

		filename := filepath.Join(libraryStorageDir, jarFilename)
		val = LibraryStorage{
			LocalPath: filename,
			Checksum:  string(latestSha256),
		}

		// The mutex should capture the write to disk as well
		// to make sure we don't have multiple threads trying to
		// write to the same file.
		libraryMutex.Lock()
		defer libraryMutex.Unlock()

		err = ioutil.WriteFile(filename, jarData, os.FileMode(0755))
		if err != nil {
			return "", err
		}

		libraryMap[version] = val
	}

	return val.LocalPath, nil
}
