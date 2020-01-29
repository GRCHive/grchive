package storage

import (
	"cloud.google.com/go/storage"
	"context"
	"google.golang.org/api/option"
	"io/ioutil"
)

type GCloudStorageApi interface {
	Init(option.ClientOption)
	Upload(bucket string, filename string, data []byte) error
	Download(bucket string, filename string) ([]byte, error)
	Delete(bucket string, filename string) error
}

type RealGCloudStorageApi struct {
	client *storage.Client
}

func (s *RealGCloudStorageApi) Init(opt option.ClientOption) {
	var err error

	ctx := context.Background()
	s.client, err = storage.NewClient(ctx, opt)
	if err != nil {
		panic("Failed to initialize GCloud Storage: " + err.Error())
	}
}

func (s *RealGCloudStorageApi) Upload(bucket string, filename string, data []byte) error {
	obj := s.client.Bucket(bucket).Object(filename)
	wr := obj.NewWriter(context.Background())
	_, err := wr.Write(data)
	if err != nil {
		return err
	}
	return wr.Close()
}

func (s *RealGCloudStorageApi) Download(bucket string, filename string) ([]byte, error) {
	obj := s.client.Bucket(bucket).Object(filename)
	reader, err := obj.NewReader(context.Background())
	if err != nil {
		return nil, err
	}

	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return bytes, reader.Close()
}

func (s *RealGCloudStorageApi) Delete(bucket string, filename string) error {
	obj := s.client.Bucket(bucket).Object(filename)
	return obj.Delete(context.Background())
}
