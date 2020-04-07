package webcore

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	"gitlab.com/grchive/grchive/gitea_api"
	"io"
	"io/ioutil"
	"strings"
)

// Better way to do this than hard coding?
const giteaTemplateTarPath string = "devops/gitea/gitea-project-template.tar.gz"

var giteaTemplateRawData *bytes.Buffer = nil
var giteaTemplateSHA256 string = ""

func loadGiteaTemplate() {
	data, err := ioutil.ReadFile(giteaTemplateTarPath)
	if err != nil {
		core.Error("Failed to read gitea template: " + err.Error())
	}
	giteaTemplateRawData = bytes.NewBuffer(data)

	rawSha256 := sha256.Sum256(giteaTemplateRawData.Bytes())
	giteaTemplateSHA256 = hex.EncodeToString(rawSha256[:])
}

func UpdateGiteaRepositoryTemplate(orgId int32) error {
	// If the current hash doesn't exist or if it doesn't match up with what
	// the latest version is, then we want to regenerate the template.
	orgHash, err := database.GetGiteaTemplateHashForOrg(orgId)
	if err != nil {
		return err
	}

	if orgHash == giteaTemplateSHA256 {
		return nil
	}

	gzf, err := gzip.NewReader(giteaTemplateRawData)
	if err != nil {
		return err
	}

	linked, err := database.GetLinkedGiteaRepository(orgId)
	if err != nil {
		return err
	}

	repo := gitea.GiteaRepository{
		Name:  linked.GiteaRepo,
		Owner: linked.GiteaOrg,
	}

	// Every file with a *.tmpl extension should be run through Golang template
	// generation while every other file should just be copied as is.
	tr := tar.NewReader(gzf)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		// We don't need to put any non-files (e.g. directories) into
		// the Git repository.
		if header.Typeflag != tar.TypeReg {
			continue
		}

		// This should just read until we hit the EOF for the file.
		fileData, err := ioutil.ReadAll(tr)
		if err != nil {
			return err
		}

		if strings.HasSuffix(header.Name, ".tmpl") {

		}

		// There isn't a universal commit file sort of function via the
		// Gitea API so make two HTTP requests for each file:
		// 1) Try to create a new file.
		// 2) If that fails for whatever reason, try to update the file.
		// 3) Only fail if both of those operations fail.
		gitPath := strings.TrimPrefix(header.Name, "./")

		_, err = gitea.GlobalGiteaApi.RepositoryCreateFile(
			repo,
			gitPath,
			string(fileData),
		)

		if err != nil {
			_, sha, err := gitea.GlobalGiteaApi.RepositoryGetFile(repo, gitPath)
			if err != nil {
				return err
			}

			_, err = gitea.GlobalGiteaApi.RepositoryUpdateFile(
				repo,
				gitPath,
				string(fileData),
				sha,
			)

			if err != nil {
				return err
			}
		}
	}

	return database.StoreGiteaTemplateHashForOrg(orgId, giteaTemplateSHA256)
}
