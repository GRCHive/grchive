package backblaze

type B2File struct {
	BucketId string `json:"bucketId"`
	FileId   string `json:"fileId"`
}
