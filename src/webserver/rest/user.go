package rest

import (
	"encoding/json"
	"gitlab.com/b3h47pte/audit-stuff/core"
	"gitlab.com/b3h47pte/audit-stuff/database"
	"gitlab.com/b3h47pte/audit-stuff/webcore"
	"io/ioutil"
	"net/http"
	"strings"
)

func updateUserProfile(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil || len(r.PostForm) == 0 {
		core.Warning("Failed to parse form data: " + core.ErrorString(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	parsedUserData, err := webcore.FindSessionParsedDataInContext(r.Context())
	if err != nil {
		core.Warning("Failed to find parsed user data: " + core.ErrorString(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	oktaUserId := parsedUserData.CurrentUser.OktaUserId
	firstName := r.PostForm["firstName"]
	lastName := r.PostForm["lastName"]

	if len(firstName) == 0 || len(lastName) == 0 {
		core.Warning("Empty form data.")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	postBody, err := json.Marshal(map[string]interface{}{
		"profile": map[string]string{
			"firstName": firstName[0],
			"lastName":  lastName[0],
		},
	})
	if err != nil {
		core.Warning("Failed to create body: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	request, err := http.NewRequest("POST", webcore.CreateOktaUserUpdateUrl(oktaUserId), strings.NewReader(string(postBody)))
	if err != nil {
		core.Warning("Failed to create request: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "SSWS "+core.LoadEnvConfig().Okta.ApiKey)

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		core.Warning("Failed to post to okta: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		core.Warning("Failed to read Okta resp: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if resp.StatusCode != http.StatusOK {
		core.Warning("Okta request failed: " + string(body))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// At this point the user might have changed their name so
	// we should force a refresh of their session so that their
	// information gets updated. The client will have to refresh on their own.
	session, err := webcore.FindSessionInContext(r.Context())
	if err == nil {
		oldSessionId, err := webcore.RefreshUserSession(session, r)
		if err != nil {
			core.Warning("Failed refresh session: " + err.Error())
		} else {
			err = database.UpdateUserSession(session, oldSessionId)
			if err != nil {
				core.Warning("Failed update session: " + err.Error())
			}
		}
	} else {
		core.Warning("Failed to find session: " + err.Error())
	}

	w.WriteHeader(http.StatusOK)
}
