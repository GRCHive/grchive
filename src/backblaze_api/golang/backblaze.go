package backblaze

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// TODO: Collapse the common parts of the 3 send functions

type B2ApiResponse map[string]*json.RawMessage

func SetAuthorizationHeaderWithAppKey(r *http.Request, key B2Key) {
	authStr := fmt.Sprintf("%s:%s", key.Id, key.Key)

	r.Header.Set("Authorization", fmt.Sprintf("Basic %s",
		base64.StdEncoding.EncodeToString([]byte(authStr))))
}

func SetAuthorizationHeaderWithToken(r *http.Request, token string) {
	r.Header.Set("Authorization", token)
}

func handleBackblazeError(resp *http.Response, obj B2ApiResponse) error {
	if resp.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf(
			"Failed to auth with Backblaze: %d\n\tMessage: %s - %s",
			resp.StatusCode,
			string(*obj["message"]),
			string(*obj["code"]),
		))
	}
	return nil
}

func sendBackblazeApiEndpoint(auth *B2AuthToken, method string, endpoint string, data interface{}, outInt interface{}) error {
	body := &bytes.Buffer{}
	if data != nil {
		jsonBuffer, err := json.Marshal(data)
		if err != nil {
			return err
		}
		body = bytes.NewBuffer(jsonBuffer)
	}

	req, err := http.NewRequest(method, auth.ApiUrl+endpoint, body)
	if err != nil {
		return err
	}
	SetAuthorizationHeaderWithToken(req, auth.Token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBodyData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	rootObj := B2ApiResponse{}
	err = json.Unmarshal(respBodyData, &rootObj)
	if err != nil {
		return err
	}

	if err = handleBackblazeError(resp, rootObj); err != nil {
		return err
	}

	if outInt != nil {
		err = json.Unmarshal(respBodyData, outInt)
		if err != nil {
			return err
		}
	}

	return nil
}

func sendBackblazeDownload(auth *B2AuthToken, file B2File) ([]byte, error) {
	rawBody, err := json.Marshal(map[string]string{
		"fileId": file.FileId,
	})
	if err != nil {
		return nil, err
	}
	body := bytes.NewBuffer(rawBody)

	const downloadEndpoint = "/b2api/v2/b2_download_file_by_id"
	req, err := http.NewRequest("POST", auth.DownloadUrl+downloadEndpoint, body)
	if err != nil {
		return nil, err
	}
	SetAuthorizationHeaderWithToken(req, auth.Token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBodyData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		rootObj := B2ApiResponse{}
		err = json.Unmarshal(respBodyData, &rootObj)
		if err != nil {
			return nil, err
		}

		if err = handleBackblazeError(resp, rootObj); err != nil {
			return nil, err
		}
	}

	return respBodyData, nil
}

func sendBackblazeUpload(uploadToken UploadFileToken, filename string, data []byte, outInt interface{}) error {
	body := bytes.NewBuffer(data)

	req, err := http.NewRequest("POST", uploadToken.Url, body)
	if err != nil {
		return err
	}

	checkSum := sha1.Sum(data)

	SetAuthorizationHeaderWithToken(req, uploadToken.Token)
	req.Header.Set("X-Bz-File-Name", url.PathEscape(filename))
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("Content-Length", fmt.Sprintf("%d", body.Len()))
	req.Header.Set("X-Bz-Content-Sha1", hex.EncodeToString(checkSum[:]))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBodyData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	rootObj := B2ApiResponse{}
	err = json.Unmarshal(respBodyData, &rootObj)
	if err != nil {
		return err
	}

	if err = handleBackblazeError(resp, rootObj); err != nil {
		return err
	}

	err = json.Unmarshal(respBodyData, outInt)
	if err != nil {
		return err
	}
	return nil
}
