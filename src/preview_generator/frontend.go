package main

import (
	"bytes"
	"encoding/json"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	"gitlab.com/grchive/grchive/gcloud_api"
	"gitlab.com/grchive/grchive/vault_api"
	"gitlab.com/grchive/grchive/webcore"
	"io/ioutil"
	"os"
	"os/exec"
)

func generatePreview(data []byte) *webcore.RabbitMQError {
	msg := webcore.FilePreviewMessage{}
	core.Info(string(data))
	err := json.Unmarshal(data, &msg)
	if err != nil {
		return &webcore.RabbitMQError{err, false}
	}

	storage := gcloud.DefaultGCloudApi.GetStorageApi()

	// Download file from B2.
	encryptedBytes, err := storage.Download(msg.Storage.BucketId, msg.Storage.StorageId)
	if err != nil {
		return &webcore.RabbitMQError{err, true}
	}

	// Unencrypt file.
	decryptedBytes, err := vault.TransitDecrypt(msg.File.UniqueKey(), encryptedBytes)
	if err != nil {
		return &webcore.RabbitMQError{err, true}
	}

	// Write file to disk in temporary file.
	tmpfile, err := ioutil.TempFile("", "og")
	if err != nil {
		return &webcore.RabbitMQError{err, true}
	}
	defer os.Remove(tmpfile.Name())

	_, err = tmpfile.Write(decryptedBytes)
	if err != nil {
		return &webcore.RabbitMQError{err, true}
	}

	// Create temporary output directory.
	tempdir, err := ioutil.TempDir("", "preview")
	defer os.RemoveAll(tempdir)

	// Convert. Any failure to convert means the preview is unavailble.
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command("./src/preview_generator/linux_amd64_stripped/generator",
		"-input", tmpfile.Name(),
		"-outputDir", tempdir)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		core.Warning("Failed to convert: " + err.Error())
		core.Warning("STDOUT: " + stdout.String())
		core.Warning("STDERR: " + stderr.String())
		err = database.MarkPreviewUnavailable(msg.File, msg.Storage, core.ServerRole)
		if err != nil {
			return &webcore.RabbitMQError{err, false}
		}
	} else {
		// Read in the preview file.
		filename, err := core.FindFirstFileInDirectory(tempdir)
		if err != nil {
			return &webcore.RabbitMQError{err, true}
		}

		previewBytes, err := ioutil.ReadFile(filename)
		if err != nil {
			return &webcore.RabbitMQError{err, true}
		}

		tx := database.CreateTx()

		storageFilename := msg.Storage.StorageId + "PREVIEW"

		// Create file preview in database and then encrypt/upload to B2.
		_, previewStorage, err := webcore.UploadNewFileWithTx(
			storage,
			&msg.File,
			msg.Storage.BucketId,
			&storageFilename,
			"PREVIEW"+msg.Storage.StorageName+".pdf",
			previewBytes,
			core.ServerRole,
			msg.Storage.UploadUserId,
			tx,
			nil,
			true,  // useExistingMetadata
			false) // addToFileVersion
		if err != nil {
			tx.Rollback()
			return &webcore.RabbitMQError{err, true}
		}

		// Connect file preview to original file.
		err = database.LinkFileWithPreviewWithTx(msg.File, msg.Storage, *previewStorage, core.ServerRole, tx)
		if err != nil {
			tx.Rollback()
			return &webcore.RabbitMQError{err, true}
		}

		if err = tx.Commit(); err != nil {
			return &webcore.RabbitMQError{err, true}
		}
	}

	return nil
}

func main() {
	core.Init()
	database.Init()
	vault.Initialize(vault.VaultConfig{
		Url:      core.EnvConfig.Vault.Url,
		Username: core.EnvConfig.Vault.Username,
		Password: core.EnvConfig.Vault.Password,
	})
	gcloud.DefaultGCloudApi.InitFromJson(core.EnvConfig.Gcloud.AuthFilename)

	webcore.DefaultRabbitMQ.Connect(*core.EnvConfig.RabbitMQ)
	defer webcore.DefaultRabbitMQ.Cleanup()

	forever := make(chan bool)

	webcore.DefaultRabbitMQ.ReceiveMessages(webcore.FILE_PREVIEW_QUEUE, generatePreview)

	<-forever
}
