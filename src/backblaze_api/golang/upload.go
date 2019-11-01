package backblaze

import (
	"errors"
)

func UploadFile(auth B2AuthToken, bucketId string, data []byte) (B2File, error) {
	file := B2File{}
	return file, errors.New("BOOP")
}
