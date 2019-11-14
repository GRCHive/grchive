package webcore

import (
	"fmt"
	"github.com/google/go-querystring/query"
	"gitlab.com/b3h47pte/audit-stuff/core"
	"gitlab.com/b3h47pte/audit-stuff/database"
	"gitlab.com/b3h47pte/audit-stuff/mail_api"
	"html/template"
)

func CreateInvitationSubject() string {
	const base string = "Invitation to %s!"
	return fmt.Sprintf(base, core.EnvConfig.Company.CompanyName)
}

var emailInvitationTemplate = template.Must(template.ParseFiles("src/webserver/templates/email/invite.tmpl"))

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

	message, err := core.TemplateToString(emailInvitationTemplate, map[string]string{
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
