package backblaze

const DownloadFileEndpoint string = "/b2api/v2/b2_delete_file_version"

func DownloadFile(auth *B2AuthToken, file B2File) ([]byte, error) {
	return sendBackblazeDownload(auth, file)
}
