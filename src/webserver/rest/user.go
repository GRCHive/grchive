package rest

import (
	"encoding/json"
	"gitlab.com/b3h47pte/audit-stuff/core"
	"gitlab.com/b3h47pte/audit-stuff/database"
	"gitlab.com/b3h47pte/audit-stuff/webcore"
	"net/http"
)

type UpdateUserProfileInputs struct {
	FirstName string `webcore:"firstName"`
	LastName  string `webcore:"lastName"`
}

func updateUserProfile(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := UpdateUserProfileInputs{}
	err := webcore.UnmarshalRequestForm(r, &inputs)
	if err != nil {
		core.Warning("Can't parse inputs: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	email, err := webcore.GetUserEmailFromRequestUrl(r)
	if err != nil {
		core.Warning("Can't find user email: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	apiKey, err := webcore.GetAPIKeyFromRequest(r)
	if apiKey == nil || err != nil {
		core.Warning("No API Key: " + core.ErrorString(err))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	apiUser, err := database.FindUserFromId(apiKey.UserId)
	if err != nil || apiUser.Email != email {
		core.Warning("Can't verify API user access: " + core.ErrorString(err))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	user := core.User{
		FirstName: inputs.FirstName,
		LastName:  inputs.LastName,
		Email:     email,
	}

	err = database.UpdateUserFromEmail(&user)
	if err != nil {
		core.Warning("Can't update user: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(user)
}
