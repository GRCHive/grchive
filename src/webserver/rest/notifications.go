package rest

import (
	"encoding/json"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	"gitlab.com/grchive/grchive/webcore"
	"net/http"
)

type AllNotificationInputs struct {
	UserId int64 `webcore:"userId"`
}

func allNotifications(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := AllNotificationInputs{}
	err := webcore.UnmarshalRequestForm(r, &inputs)
	if err != nil {
		core.Warning("Can't parse inputs: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	key, err := webcore.FindApiKeyInContext(r.Context())
	if err != nil {
		core.Warning("Can't find API key: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if inputs.UserId != key.UserId {
		core.Warning("Invalid access for user's notifications.")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	notifications, err := database.AllNotificationsForUserId(inputs.UserId)
	if err != nil {
		core.Warning("Failed to get user notifications.")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(notifications)
}

type MarkNotificationReadInputs struct {
	UserId          int64   `json:"userId"`
	NotificationIds []int64 `json:"notificationIds"`
	All             bool    `json:"all"`
}

func markNotificationRead(w http.ResponseWriter, r *http.Request) {
	inputs := MarkNotificationReadInputs{}
	err := webcore.UnmarshalRequestForm(r, &inputs)
	if err != nil {
		core.Warning("Can't parse inputs: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	key, err := webcore.FindApiKeyInContext(r.Context())
	if err != nil {
		core.Warning("Can't find API key: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if inputs.UserId != key.UserId {
		core.Warning("Invalid access for user's notifications.")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if inputs.All {
		err = database.MarkAllNotificationsAsRead(inputs.UserId)
	} else {
		err = database.MarkNotificationsAsRead(inputs.UserId, inputs.NotificationIds)
	}

	if err != nil {
		core.Warning("Failed to mark notifications as read: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
