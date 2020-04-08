package webcore

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	"gitlab.com/grchive/grchive/gitea_api"
	"html/template"
	"io"
	"io/ioutil"
	"strconv"
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

func loadTemplateParamsForOrg(org *core.Organization) (map[string]string, string, error) {
	// We also need to check the hash on the template parameters since if that
	// changes, we also want to regenerate.
	templateParams := map[string]string{
		// TODO: Let user choose this
		"GRCHIVE_ORG_IDENTIFIER": org.OktaGroupName,
		// TODO: Do we need to actually update teh version in the pom?
		"GRCHIVE_CLIENT_LIB_VERSION": "0.1",
		// TODO: ????
		"GRCHIVE_ORG_URL":  "",
		"ARTIFACTORY_HOST": core.EnvConfig.Artifactory.Host,
		"ARTIFACTORY_PORT": strconv.FormatInt(int64(core.EnvConfig.Artifactory.Port), 10),
	}

	templateJsonRaw, err := json.Marshal(templateParams)
	if err != nil {
		return nil, "", err
	}

	rawSha256 := sha256.Sum256(templateJsonRaw)
	sha256hex := hex.EncodeToString(rawSha256[:])
	return templateParams, sha256hex, nil
}

func UpdateGiteaRepositoryTemplate(orgId int32) error {
	org, err := database.FindOrganizationFromId(orgId)
	if err != nil {
		return err
	}

	// If the current hash doesn't exist or if it doesn't match up with what
	// the latest version is, then we want to regenerate the template.
	orgHash, templateHash, err := database.GetGiteaTemplateHashForOrg(orgId)
	if err != nil {
		return err
	}

	templateParams, testTemplateHash, err := loadTemplateParamsForOrg(org)
	if err != nil {
		return err
	}

	if orgHash == giteaTemplateSHA256 && templateHash == testTemplateHash {
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

		strData := string(fileData)

		useFilename := header.Name
		if strings.HasSuffix(header.Name, ".tmpl") {
			useFilename = strings.TrimSuffix(useFilename, ".tmpl")

			tmpl, err := template.New(header.Name).Parse(strData)
			if err != nil {
				return err
			}

			strData, err = core.TemplateToString(tmpl, templateParams)

			if err != nil {
				return err
			}
		}

		// There isn't a universal commit file sort of function via the
		// Gitea API so make two HTTP requests for each file:
		// 1) Try to create a new file.
		// 2) If that fails for whatever reason, try to update the file.
		// 3) Only fail if both of those operations fail.
		gitPath := strings.TrimPrefix(useFilename, "./")

		_, _, err = gitea.GlobalGiteaApi.RepositoryCreateFile(
			repo,
			gitPath,
			strData,
		)

		if err != nil {
			_, sha, err := gitea.GlobalGiteaApi.RepositoryGetFile(repo, gitPath, "master")
			if err != nil {
				return err
			}

			_, _, err = gitea.GlobalGiteaApi.RepositoryUpdateFile(
				repo,
				gitPath,
				strData,
				sha,
			)

			if err != nil {
				return err
			}
		}
	}

	return database.StoreGiteaTemplateHashForOrg(orgId, giteaTemplateSHA256, testTemplateHash)
}
