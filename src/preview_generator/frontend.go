package main

import (
	"encoding/json"
	"gitlab.com/b3h47pte/audit-stuff/backblaze_api"
	"gitlab.com/b3h47pte/audit-stuff/core"
	"gitlab.com/b3h47pte/audit-stuff/database"
	"gitlab.com/b3h47pte/audit-stuff/vault_api"
	"gitlab.com/b3h47pte/audit-stuff/webcore"
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

	b2Auth, err := backblaze.B2Auth(core.EnvConfig.Backblaze.Key)
	if err != nil {
		return &webcore.RabbitMQError{err, true}
	}

	// Download file from B2.
	encryptedBytes, err := backblaze.DownloadFile(b2Auth, backblaze.B2File{
		BucketId: msg.File.BucketId,
		FileId:   msg.File.StorageId,
	})
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
	cmd := exec.Command("./src/preview_generator/linux_amd64_stripped/generator",
		"-input", tmpfile.Name(),
		"-outputDir", tempdir)

	err = cmd.Run()
	if err != nil {
		core.Warning("Failed to convert: " + err.Error())
		err = database.MarkPreviewUnavailable(msg.File, core.ServerRole)
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

		// Get the relevant organization.
		org, err := database.FindOrganizationFromId(msg.File.OrgId)
		if err != nil {
			return &webcore.RabbitMQError{err, true}
		}

		tx := database.CreateTx()

		// Create file preview in database and then encrypt/upload to B2.
		previewFile := core.ControlDocumentationFile{
			StorageName:  "PREVIEW" + msg.File.StorageName,
			RelevantTime: msg.File.RelevantTime,
			UploadTime:   msg.File.UploadTime,
			CategoryId:   msg.File.CategoryId,
			OrgId:        msg.File.OrgId,
			AltName:      "PREVIEW" + msg.File.AltName,
			Description:  "PREVIEW",
			UploadUserId: msg.File.UploadUserId,
		}
		err = webcore.UploadNewFileWithTx(&previewFile, previewBytes, core.ServerRole, org, b2Auth, tx)
		if err != nil {
			tx.Rollback()
			return &webcore.RabbitMQError{err, true}
		}

		// Connect file preview to original file.
		err = database.LinkFileWithPreviewWithTx(msg.File, previewFile, core.ServerRole, tx)
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
		Url:   core.EnvConfig.Vault.Url,
		Token: core.EnvConfig.Vault.Token,
	})

	webcore.DefaultRabbitMQ.Connect(*core.EnvConfig.RabbitMQ)
	defer webcore.DefaultRabbitMQ.Cleanup()

	forever := make(chan bool)

	webcore.DefaultRabbitMQ.ReceiveMessages(webcore.FILE_PREVIEW_QUEUE, generatePreview)

	<-forever
}
