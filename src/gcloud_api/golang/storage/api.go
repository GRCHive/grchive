package storage

import (
	"cloud.google.com/go/storage"
	"context"
	"google.golang.org/api/option"
	"hash/crc32"
	"io/ioutil"
)

type GCloudStorageApi interface {
	Init(option.ClientOption)
	Upload(bucket string, filename string, data []byte, key []byte) error
	Download(bucket string, filename string, key []byte) ([]byte, error)
	Delete(bucket string, filename string) error
}

type RealGCloudStorageApi struct {
	client *storage.Client
}

var crc32cTable *crc32.Table = crc32.MakeTable(crc32.Castagnoli)

func (s *RealGCloudStorageApi) Init(opt option.ClientOption) {
	var err error

	ctx := context.Background()
	s.client, err = storage.NewClient(ctx, opt)
	if err != nil {
		panic("Failed to initialize GCloud Storage: " + err.Error())
	}
}

func (s *RealGCloudStorageApi) Upload(bucket string, filename string, data []byte, key []byte) error {
	obj := s.client.Bucket(bucket).Object(filename)
	wr := obj.NewWriter(context.Background())

	wr.SendCRC32C = true

	uploadData, err := appendHMACSHA512(data, key)
	if err != nil {
		return err
	}

	wr.ObjectAttrs.CRC32C = crc32.Checksum(uploadData, crc32cTable)

	_, err = wr.Write(uploadData)
	if err != nil {
		return err
	}
	return wr.Close()
}

func (s *RealGCloudStorageApi) Download(bucket string, filename string, key []byte) ([]byte, error) {
	obj := s.client.Bucket(bucket).Object(filename)
	reader, err := obj.NewReader(context.Background())
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	rawData, err := verifyHMACSHA512(bytes, key)
	if err != nil {
		if err == ErrNoHMAC {
			// Chances are that this is an older file so reupload the file
			// with an HMAC attached.
			err = s.Upload(bucket, filename, rawData, key)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	return rawData, nil
}

func (s *RealGCloudStorageApi) Delete(bucket string, filename string) error {
	obj := s.client.Bucket(bucket).Object(filename)
	return obj.Delete(context.Background())
}
