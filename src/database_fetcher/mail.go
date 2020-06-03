package main

import (
	"fmt"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	"gitlab.com/grchive/grchive/mail_api"
	"html/template"
)

func sendDbSchemaChangeEmails(db *core.Database) error {
	settings, err := database.GetDatabaseSettings(db.Id)
	if err != nil {
		return err
	}

	org, err := database.FindOrganizationFromId(db.OrgId)
	if err != nil {
		return err
	}

	emailTemplate, err := template.ParseFiles("src/webserver/templates/email/schemaChange.tmpl")
	if err != nil {
		return err
	}

	for _, u := range settings.OnSchemaChangeNotifyUsers {
		msg, err := core.HtmlTemplateToString(emailTemplate, map[string]string{
			"recipient":    u.FullName(),
			"databaseName": db.Name,
			// ??? Not sure how we do this since we don't have the router here.
			"databaseUrl": fmt.Sprintf("%s/dashboard/org/%s/it/databases/%d",
				core.EnvConfig.SelfUri,
				org.OktaGroupName,
				db.Id,
			),
		})

		if err != nil {
			return err
		}

		mailPayload := mail.MailPayload{
			From: core.EnvConfig.Mail.VeriEmailFrom,
			To: mail.Email{
				Name:  u.FullName(),
				Email: u.Email,
			},
			Subject: fmt.Sprintf("GRCHive - Schema Change Detected For %s", db.Name),
			Message: msg,
		}
		mail.MailProvider.SendMail(mailPayload)
	}

	return nil
}
