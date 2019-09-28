package rest

import (
	"encoding/json"
	"gitlab.com/b3h47pte/audit-stuff/core"
	"gitlab.com/b3h47pte/audit-stuff/database"
	"net/http"
	"strings"
)

type tGettingStartedInterest struct {
	Name  string `form:"name" binding:"required"`
	Email string `form:"email" binding:"required"`
}

func postGettingStartedInterest(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	// Retrieve the client's name and email from the input form.
	if err := r.ParseForm(); err != nil || len(r.PostForm) == 0 {
		core.Warning("Failed to parse form data.")
		w.WriteHeader(http.StatusBadRequest)
		jsonWriter.Encode(struct{}{})
		return
	}

	nameData := r.PostForm["name"]
	emailData := r.PostForm["email"]

	if len(nameData) == 0 || len(emailData) == 0 {
		core.Warning("Empty name or email.")
		w.WriteHeader(http.StatusBadRequest)
		jsonWriter.Encode(struct{}{})
		return
	}

	var data tGettingStartedInterest
	data.Name = strings.TrimSpace(nameData[0])
	data.Email = strings.TrimSpace(emailData[0])

	// Save name and email to the database.
	isDuplicate, err := database.AddNewGettingStartedInterest(data.Name, data.Email)

	// If the error is related to having a duplicate then we should let the user know.
	// Otherwise, our service probably failed somewhere which hopefully got logged.
	if err != nil {
		if isDuplicate {
			core.Warning("Detected duplicate entry.")
			w.WriteHeader(http.StatusBadRequest)
			jsonWriter.Encode(struct {
				IsDuplicate bool
			}{
				IsDuplicate: true,
			})
		} else {
			core.Warning("Failed to add getting started interest.")
			w.WriteHeader(http.StatusInternalServerError)
			jsonWriter.Encode(struct{}{})
		}
		return
	}

	jsonWriter.Encode(struct{}{})
}
