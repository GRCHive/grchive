package backblaze

import (
	"encoding/base64"
	"fmt"
	"net/http"
)

func SetAuthorizationHeaderWithAppKey(r *http.Request, key B2Key) {
	authStr := fmt.Sprintf("%s:%s", key.Id, key.Key)

	r.Header.Set("Authorization", fmt.Sprintf("Basic %s",
		base64.StdEncoding.EncodeToString([]byte(authStr))))
}
