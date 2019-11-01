package backblaze

const GetUrlRetryAmount int = 1
const GetUploadUrlEndpoint string = "/b2api/v2/b2_get_upload_url"
const UploadFileEndpoint string = "/b2api/v2/b2_get_upload_url"

type UploadFileToken struct {
	Url   string `json:"uploadUrl"`
	Token string `json:"authorizationToken"`
}

func GetUploadUrl(auth *B2AuthToken, bucketId string) (UploadFileToken, error) {
	var err error
	token := UploadFileToken{}

	for i := 0; i < GetUrlRetryAmount; i++ {
		err = sendBackblazeApiEndpoint(auth, "POST", GetUploadUrlEndpoint, map[string]string{
			"bucketId": bucketId,
		}, &token)

		if err == nil {
			break
		}
	}

	if err != nil {
		return UploadFileToken{}, err
	}

	return token, nil
}

func UploadFile(auth *B2AuthToken, bucketId string, filename string, data []byte) (B2File, error) {
	uploadToken, err := GetUploadUrl(auth, bucketId)
	if err != nil {
		return B2File{}, err
	}

	file := B2File{}
	err = sendBackblazeUpload(uploadToken, filename, data, &file)
	if err != nil {
		return B2File{}, err
	}

	return file, nil
}
