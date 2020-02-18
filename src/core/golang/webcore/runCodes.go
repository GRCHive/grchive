package webcore

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/google/uuid"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	"gitlab.com/grchive/grchive/mail_api"
	"html/template"
	"strings"
	"time"
)

var RunCodeExpirationTime = time.Hour * 48

func GenerateRandomRunCode(requestId int64, orgId int32) (*core.DbSqlQueryRunCode, string, error) {
	rawCode := strings.ReplaceAll(uuid.New().String(), "-", "")
	salt, err := core.RandomHexString(numBytesForSalt)
	if err != nil {
		return nil, "", err
	}

	base64Code := base64.StdEncoding.EncodeToString([]byte(rawCode))
	saltedCode := base64Code + "." + salt

	hashedCode := sha256.Sum256([]byte(saltedCode))

	runCode := core.DbSqlQueryRunCode{
		RequestId:      requestId,
		OrgId:          orgId,
		ExpirationTime: time.Now().UTC().Add(RunCodeExpirationTime),
		HashedCode:     base64.StdEncoding.EncodeToString(hashedCode[:]),
		Salt:           salt,
	}

	return &runCode, rawCode, nil
}

func SendRunCodeViaEmail(runCode *core.DbSqlQueryRunCode, rawCode string) error {
	request, err := database.GetSqlRequest(runCode.RequestId, runCode.OrgId, core.ServerRole)
	if err != nil {
		return err
	}

	user, err := database.FindUserFromId(request.UploadUserId)
	if err != nil {
		return err
	}

	query, err := database.GetSqlQueryFromId(request.QueryId, request.OrgId, core.ServerRole)
	if err != nil {
		return err
	}

	metadata, err := database.GetSqlMetadataFromId(query.MetadataId, query.OrgId, core.ServerRole)
	if err != nil {
		return err
	}

	emailTemplate, err := template.ParseFiles("src/webserver/templates/email/runCode.tmpl")
	if err != nil {
		return err
	}

	message, err := core.TemplateToString(emailTemplate, map[string]string{
		"queryName": fmt.Sprintf("%s v%d", metadata.Name, query.Version),
		"runCode":   rawCode,
	})
	if err != nil {
		return err
	}

	mailPayload := mail.MailPayload{
		From: core.EnvConfig.Mail.VeriEmailFrom,
		To: mail.Email{
			Name:  fmt.Sprintf("%s %s", user.FirstName, user.LastName),
			Email: user.Email,
		},
		Subject: "SQL Query Request Approval",
		Message: message,
	}

	err = mail.MailProvider.SendMail(mailPayload)
	if err != nil {
		return err
	}

	return nil
}
