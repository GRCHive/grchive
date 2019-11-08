package webcore

import (
	"fmt"
	"gitlab.com/b3h47pte/audit-stuff/core"
	"gitlab.com/b3h47pte/audit-stuff/database"
	"gitlab.com/b3h47pte/audit-stuff/mail_api"
	"html/template"
	"strconv"
)

func CreateEmailVerificationSubject() string {
	const emailVerificationSubjectBase string = "Welcome to %s! Please Verify your Email."
	return fmt.Sprintf(emailVerificationSubjectBase, core.EnvConfig.Company.CompanyName)
}

const emailVerificationTemplateFname string = "src/webserver/templates/email/verification.tmpl"

var emailVerificationTemplate = template.Must(template.ParseFiles(emailVerificationTemplateFname))

func SendEmailVerification(user *core.User) error {
	veri := core.CreateNewEmailVerification(user)

	veriLink, err := core.CreateUrlWithParams(MustGetRouteUrlAbsolute(EmailVerifyRouteName), map[string]string{
		"code": veri.Code,
		"user": strconv.FormatInt(user.Id, 10),
	})

	if err != nil {
		return err
	}

	err = database.StoreEmailVerification(veri)
	if err != nil {
		return err
	}

	message, err := core.TemplateToString(emailVerificationTemplate, map[string]string{
		"userFullName":     user.FullName(),
		"productName":      core.EnvConfig.Company.CompanyName,
		"verificationLink": veriLink,
	})
	if err != nil {
		return err
	}

	mailPayload := mail.MailPayload{
		From: core.EnvConfig.Mail.VeriEmailFrom,
		To: mail.Email{
			Name:  user.FullName(),
			Email: user.Email,
		},
		Subject: CreateEmailVerificationSubject(),
		Message: message,
	}

	err = mail.MailProvider.SendMail(mailPayload)
	if err != nil {
		return err
	}

	return nil
}

func CheckEmailVerification(code string, userId int64) bool {
	_, err := database.FindUserVerification(code, userId)
	if err != nil {
		core.Warning("Can't find user verification: " + err.Error())
		return false
	}

	err = database.AcceptUserVerification(code, userId)
	if err != nil {
		core.Warning("Can't accept user verification: " + err.Error())
		return false
	}

	return true
}
