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
	Upload(bucket string, filename string, data []byte, key []byte) (int64, error)
	Download(bucket string, filename string, key []byte) ([]byte, error)
	DownloadVersioned(bucket string, filename string, generation int64, key []byte) ([]byte, error)
	Delete(bucket string, filename string) error
	DeleteVersioned(bucket string, filename string, generation int64) error
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

func (s *RealGCloudStorageApi) Upload(bucket string, filename string, data []byte, key []byte) (int64, error) {
	obj := s.client.Bucket(bucket).Object(filename)
	wr := obj.NewWriter(context.Background())

	wr.SendCRC32C = true

	var uploadData []byte
	var err error

	if len(key) > 0 {
		uploadData, err = appendHMACSHA512(data, key)
		if err != nil {
			return -1, err
		}
	} else {
		uploadData = data
	}

	wr.ObjectAttrs.CRC32C = crc32.Checksum(uploadData, crc32cTable)

	_, err = wr.Write(uploadData)
	if err != nil {
		return -1, err
	}

	err = wr.Close()
	if err != nil {
		return -1, err
	}

	attrs, err := obj.Attrs(context.Background())
	if err != nil {
		return -1, err
	}

	return attrs.Generation, nil
}

func (s *RealGCloudStorageApi) Download(bucket string, filename string, key []byte) ([]byte, error) {
	return s.DownloadVersioned(bucket, filename, -1, key)
}

func (s *RealGCloudStorageApi) DownloadVersioned(bucket string, filename string, generation int64, key []byte) ([]byte, error) {
	obj := s.client.Bucket(bucket).Object(filename)

	if generation != -1 {
		obj = obj.Generation(generation)
	}

	reader, err := obj.NewReader(context.Background())
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	if len(key) == 0 {
		return bytes, nil
	}

	rawData, err := verifyHMACSHA512(bytes, key)
	if err != nil {
		if err == ErrNoHMAC {
			// Chances are that this is an older file so reupload the file
			// with an HMAC attached.
			_, err = s.Upload(bucket, filename, rawData, key)
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

func (s *RealGCloudStorageApi) DeleteVersioned(bucket string, filename string, generation int64) error {
	obj := s.client.Bucket(bucket).Object(filename).Generation(generation)
	return obj.Delete(context.Background())
}
