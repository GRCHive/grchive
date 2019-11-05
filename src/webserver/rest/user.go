package rest

import (
	"encoding/json"
	"gitlab.com/b3h47pte/audit-stuff/core"
	"gitlab.com/b3h47pte/audit-stuff/database"
	"gitlab.com/b3h47pte/audit-stuff/webcore"
	"net/http"
)

type UpdateUserProfileInputs struct {
	UserId    int64  `webcore:"userId"`
	FirstName string `webcore:"firstName"`
	LastName  string `webcore:"lastName"`
}

func updateUserProfile(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	if err := r.ParseForm(); err != nil || len(r.PostForm) == 0 {
		core.Warning("Failed to parse form data: " + core.ErrorString(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	inputs := UpdateUserProfileInputs{}
	err := webcore.UnmarshalRequestForm(r, &inputs)
	if err != nil {
		core.Warning("Can't parse inputs: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user := core.User{
		Id:        inputs.UserId,
		FirstName: inputs.FirstName,
		LastName:  inputs.LastName,
	}

	err = database.UpdateUser(&user)
	if err != nil {
		core.Warning("Can't update user: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(user)
}
