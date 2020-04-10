package webcore

import (
	"errors"
	"fmt"
	"github.com/google/go-querystring/query"
	"github.com/jmoiron/sqlx"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	"gitlab.com/grchive/grchive/mail_api"
	"html/template"
)

func CreateInvitationSubject() string {
	const base string = "Invitation to %s!"
	return fmt.Sprintf(base, core.EnvConfig.Company.CompanyName)
}

func SendInviteCodeEmailCode(invite *core.InviteCode, code string) error {
	user, err := database.FindUserFromId(invite.FromUserId)
	if err != nil {
		return err
	}

	org, err := database.FindOrganizationFromId(invite.FromOrgId)
	if err != nil {
		return err
	}

	params := struct {
		InviteCode string `url:"inviteCode"`
		Email      string `url:"email"`
	}{
		InviteCode: code,
		Email:      invite.ToEmail,
	}

	v, _ := query.Values(params)
	inviteLink := MustGetRouteUrlAbsolute(AcceptInviteRouteName) + "?" + v.Encode()
	registerLink := MustGetRouteUrlAbsolute(RegisterRouteName) + "?" + v.Encode()

	emailInvitationTemplate, err := template.ParseFiles("src/webserver/templates/email/invite.tmpl")
	if err != nil {
		return err
	}

	message, err := core.HtmlTemplateToString(emailInvitationTemplate, map[string]string{
		"userFullName":     user.FullName(),
		"productName":      core.EnvConfig.Company.CompanyName,
		"inviteLink":       inviteLink,
		"registerLink":     registerLink,
		"inviteCode":       code,
		"organizationName": org.Name,
	})
	if err != nil {
		return err
	}

	mailPayload := mail.MailPayload{
		From: core.EnvConfig.Mail.VeriEmailFrom,
		To: mail.Email{
			Name:  invite.ToEmail,
			Email: invite.ToEmail,
		},
		Subject: CreateInvitationSubject(),
		Message: message,
	}

	err = mail.MailProvider.SendMail(mailPayload)
	if err != nil {
		return err
	}

	return nil
}

func SendSingleInviteCode(invite *core.InviteCode, role *core.Role) error {
	// Make sure the user isn't already a part of the organization.
	exists, err := database.IsUserEmailInOrganization(invite.ToEmail, invite.FromOrgId)
	if err != nil {
		return err
	}

	if exists {
		return errors.New("Can not send invite to a user in the organization.")
	}

	tx := database.CreateTx()

	// Store pending invite into the database so we can
	// keep track of whether or not the email was sent.
	// Later on, if the database entry exists, we know the
	// email was sent successfully, if the database entry
	// does not exist, we know the email failed and was rolled back.
	code, err := database.InsertInviteCodeWithTx(invite, role, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = SendInviteCodeEmailCode(invite, code)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// If an error occurs, returns the invite email that caused the error.
func SendBatchInviteCodes(invites []*core.InviteCode, role *core.Role) (string, error) {
	for _, inv := range invites {
		if err := SendSingleInviteCode(inv, role); err != nil {
			return inv.ToEmail, err
		}
	}
	return "", nil
}

func ProcessInviteCodeForUserWithTx(invite *core.InviteCode, user *core.User, tx *sqlx.Tx) error {
	// Need to do two things here: add the user to the correct organization.
	// Mark the invite as being used.
	org, err := database.FindOrganizationFromId(invite.FromOrgId)
	if err != nil {
		return err
	}

	err = database.AddUserToOrganizationWithTx(user, org, tx)
	if err != nil {
		return err
	}

	err = database.InsertUserRoleForOrgWithTx(user.Id, org.Id, invite.RoleId, core.ServerRole, tx)
	if err != nil {
		return err
	}

	err = database.MarkInviteAsUsedWithTx(invite, tx)
	if err != nil {
		return err
	}

	return nil
}

func ProcessInviteCodeForUser(invite *core.InviteCode, user *core.User) error {
	tx := database.CreateTx()
	err := ProcessInviteCodeForUserWithTx(invite, user, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}
