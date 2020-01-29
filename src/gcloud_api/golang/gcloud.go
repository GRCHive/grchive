package gcloud

import (
	"gitlab.com/grchive/grchive/gcloud_api/storage"
	"google.golang.org/api/option"
)

type GCloudApi interface {
	// Initialization
	InitFromJson(filename string)

	// Sub-API
	GetStorageApi() storage.GCloudStorageApi
}

type RealGCloudApi struct {
	Storage storage.GCloudStorageApi
}

func (a *RealGCloudApi) InitFromJson(filename string) {
	clientOptions := option.WithCredentialsFile(filename)

	a.Storage = &storage.RealGCloudStorageApi{}
	a.Storage.Init(clientOptions)
}

func (a *RealGCloudApi) GetStorageApi() storage.GCloudStorageApi {
	return a.Storage
}

var DefaultGCloudApi GCloudApi = &RealGCloudApi{}
