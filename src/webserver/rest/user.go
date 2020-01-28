package rest

import (
	"encoding/json"
	"github.com/google/go-querystring/query"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	"gitlab.com/grchive/grchive/webcore"
	"net/http"
	"time"
)

type UpdateUserProfileInputs struct {
	FirstName string `webcore:"firstName"`
	LastName  string `webcore:"lastName"`
}

type VerifyEmailInputs struct {
	Code   string `webcore:"code"`
	UserId int64  `webcore:"user"`
}

type RequestResendVerificationEmailInputs struct {
	UserId int64 `webcore:"userId"`
}

type SendInviteInputs struct {
	FromUserId int64    `webcore:"fromUserId"`
	FromOrgId  int32    `webcore:"fromOrgId"`
	ToEmails   []string `webcore:"toEmails"`
	RoleId     int64    `webcore:"roleId"`
}

type AcceptInviteInputs struct {
	Email      string `webcore:"email"`
	InviteCode string `webcore:"inviteCode"`
}

func getAllOrganizationsForUser(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	userId, err := webcore.GetUserIdFromRequestUrl(r)
	if err != nil {
		core.Warning("Can't find user id: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	apiKey, err := webcore.GetAPIKeyFromRequest(r)
	if apiKey == nil || err != nil {
		core.Warning("No API Key: " + core.ErrorString(err))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if userId != apiKey.UserId {
		core.Warning("Unauthorized access")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	orgs, err := database.FindAccessibleOrganizationsForUser(userId)
	if err != nil {
		core.Warning("Can't find orgs: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(orgs)
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

	userId, err := webcore.GetUserIdFromRequestUrl(r)
	if err != nil {
		core.Warning("Can't find user id: " + err.Error())
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
	if err != nil || apiUser.Id != userId {
		core.Warning("Can't verify API user access: " + core.ErrorString(err))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	user := core.User{
		FirstName: inputs.FirstName,
		LastName:  inputs.LastName,
		Email:     apiUser.Email,
	}

	err = database.UpdateUserFromEmail(&user)
	if err != nil {
		core.Warning("Can't update user: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(user)
}

func verifyUserEmail(w http.ResponseWriter, r *http.Request) {
	inputs := VerifyEmailInputs{}
	err := webcore.UnmarshalRequestForm(r, &inputs)
	if err != nil {
		core.Warning("Can't parse inputs: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// SHOULD WE DO SOMETHING BETTER ON ERROR?
	ok := webcore.CheckEmailVerification(inputs.Code, inputs.UserId)
	if ok {
		http.Redirect(w, r, webcore.MustGetRouteUrl(webcore.DashboardHomeRouteName), http.StatusTemporaryRedirect)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func requestResendUserVerificationEmail(w http.ResponseWriter, r *http.Request) {
	inputs := RequestResendVerificationEmailInputs{}
	err := webcore.UnmarshalRequestForm(r, &inputs)
	if err != nil {
		core.Warning("Can't parse inputs: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	apiKey, err := webcore.GetAPIKeyFromRequest(r)
	if apiKey == nil || err != nil {
		core.Warning("No API Key: " + core.ErrorString(err))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	user, err := database.FindUserFromId(inputs.UserId)
	if err != nil {
		core.Warning("Can't find user: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if apiKey.UserId != user.Id {
		core.Warning("Unauthorized requesting request.")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	err = webcore.SendEmailVerification(user)
	if err != nil {
		core.Warning("Can't send verification: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func sendInviteToOrganization(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := SendInviteInputs{}
	err := webcore.UnmarshalRequestForm(r, &inputs)
	if err != nil {
		core.Warning("Can't parse inputs: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	role, err := webcore.GetCurrentRequestRole(r, inputs.FromOrgId)
	if err != nil {
		core.Warning("Bad access: " + err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	pendingInvites := []*core.InviteCode{}
	for _, email := range inputs.ToEmails {
		invite := core.InviteCode{
			FromUserId: inputs.FromUserId,
			FromOrgId:  inputs.FromOrgId,
			ToEmail:    email,
			SentTime:   core.CreateNullTime(time.Now().UTC()),
			RoleId:     inputs.RoleId,
		}
		pendingInvites = append(pendingInvites, &invite)
	}

	failureEmail, err := webcore.SendBatchInviteCodes(pendingInvites, role)
	if err != nil {
		core.Warning("Failed to send invites: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		jsonWriter.Encode(struct {
			FailedEmail string
		}{
			FailedEmail: failureEmail,
		})
		return
	}
}

func acceptInviteToOrganization(w http.ResponseWriter, r *http.Request) {
	inputs := AcceptInviteInputs{}
	err := webcore.UnmarshalRequestForm(r, &inputs)
	if err != nil {
		core.Warning("Can't parse inputs: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// If the user isn't registered, redirect them to the registration page with the email
	// and invite code prefilled.
	user, err := database.FindUserFromEmail(inputs.Email)
	if err != nil {
		registerV, _ := query.Values(struct {
			InviteCode string `url:"inviteCode"`
			Email      string `url:"email"`
		}{
			InviteCode: inputs.InviteCode,
			Email:      inputs.Email,
		})

		http.Redirect(w, r,
			webcore.MustGetRouteUrl(webcore.RegisterRouteName)+"?"+registerV.Encode(),
			http.StatusTemporaryRedirect)
		return
	}

	// If the key isn't valid just error out.
	// TODO: Better error?
	invite, err := database.FindInviteCodeFromHash(inputs.InviteCode, inputs.Email, core.ServerRole)
	if err != nil {
		core.Warning("Invalid invite code: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Process the invite and associate it with the user.
	err = webcore.ProcessInviteCodeForUser(invite, user)
	if err != nil {
		core.Warning("Failed to process invite code: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// If the user is logged in, direct them to the dashboard.
	// If the user is not logged in, direct them to the login page.
	_, err = webcore.FindSessionInContext(r.Context())
	if err != nil {
		http.Redirect(w, r,
			webcore.MustGetRouteUrl(webcore.LoginRouteName),
			http.StatusTemporaryRedirect)
	} else {
		http.Redirect(w, r,
			webcore.MustGetRouteUrl(webcore.DashboardHomeRouteName),
			http.StatusTemporaryRedirect)
	}
}
