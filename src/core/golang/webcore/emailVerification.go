package webcore

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	"gitlab.com/grchive/grchive/mail_api"
	"html/template"
	"strconv"
)

func CreateEmailVerificationSubject() string {
	const emailVerificationSubjectBase string = "Welcome to %s! Please Verify your Email."
	return fmt.Sprintf(emailVerificationSubjectBase, core.EnvConfig.Company.CompanyName)
}

const emailVerificationTemplateFname string = "src/webserver/templates/email/verification.tmpl"

func SendEmailVerification(user *core.User) error {
	tx := database.CreateTx()
	err := SendEmailVerificationWithTx(user, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func SendEmailVerificationWithTx(user *core.User, tx *sqlx.Tx) error {
	veri := core.CreateNewEmailVerification(user, core.DefaultUuidGen, core.DefaultClock)

	veriLink, err := core.CreateUrlWithParams(MustGetRouteUrlAbsolute(EmailVerifyRouteName), map[string]string{
		"code": veri.Code,
		"user": strconv.FormatInt(user.Id, 10),
	})

	if err != nil {
		return err
	}

	err = database.StoreEmailVerificationWithTx(veri, tx)
	if err != nil {
		return err
	}

	emailVerificationTemplate, err := template.ParseFiles(emailVerificationTemplateFname)
	if err != nil {
		return err
	}

	message, err := core.HtmlTemplateToString(emailVerificationTemplate, map[string]string{
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
