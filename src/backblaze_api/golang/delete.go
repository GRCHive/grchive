package backblaze

const DeleteFileEndpoint string = "/b2api/v2/b2_delete_file_version"

func DeleteFile(auth *B2AuthToken, filename string, file B2File) error {
	return sendBackblazeApiEndpoint(auth, "POST", DeleteFileEndpoint, map[string]string{
		"fileName": filename,
		"fileId":   file.FileId,
	}, nil)
}
