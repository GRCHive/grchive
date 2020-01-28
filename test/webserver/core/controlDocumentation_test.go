package core_test

import (
	"github.com/stretchr/testify/assert"
	"gitlab.com/grchive/grchive/core"
	"testing"
)

func TestUniqueKey(t *testing.T) {
	for _, ref := range []struct {
		file core.ControlDocumentationFile
		key  string
	}{
		{
			file: core.ControlDocumentationFile{
				Id: 3,
			},
			key: "controlDocFile-3",
		},
		{
			file: core.ControlDocumentationFile{
				Id: 392,
			},
			key: "controlDocFile-392",
		},
	} {
		assert.Equal(t, ref.key, ref.file.UniqueKey())
	}
}

func TestStorageFilename(t *testing.T) {
	for _, ref := range []struct {
		file     core.ControlDocumentationFile
		org      core.Organization
		filename string
	}{
		{
			file: core.ControlDocumentationFile{
				Id: 3,
			},
			org: core.Organization{
				Id:   1,
				Name: "Test",
			},
			filename: "org-1-Test/controlDocFile-3",
		},
		{
			file: core.ControlDocumentationFile{
				Id: 392,
			},
			org: core.Organization{
				Id:   3,
				Name: "Blah Blah",
			},
			filename: "org-3-Blah Blah/controlDocFile-392",
		},
	} {
		assert.Equal(t, ref.filename, ref.file.StorageFilename(&ref.org))
	}

}
