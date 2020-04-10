package gitea

import (
	"encoding/base64"
	"fmt"
)

type GiteaConfig struct {
	Protocol string
	Host     string
	Port     int32
	Token    string
}

func (c GiteaConfig) apiUrl() string {
	return fmt.Sprintf("%s://%s:%d/api/v1", c.Protocol, c.Host, c.Port)
}

func (c GiteaConfig) apiUrlUserAuth(user GiteaUser) string {
	return fmt.Sprintf("%s://%s:%s@%s:%d/api/v1", c.Protocol, user.Username, user.Password, c.Host, c.Port)
}

type GiteaUserlike interface {
	GetUsername() string
}

type GiteaUser struct {
	Username string
	Password string
	Email    string
	FullName string
}

func (u GiteaUser) GetUsername() string {
	return u.Username
}

type GiteaOrganization struct {
	Username string
	FullName string
}

func (o GiteaOrganization) GetUsername() string {
	return o.Username
}

type GiteaRepository struct {
	Name  string
	Owner string
}

type GiteaToken struct {
	Name  string
	Token string
}

type GiteaCreateFileOptions struct {
	Content string
	Message string
}

func (o GiteaCreateFileOptions) PrepareApiBody() map[string]interface{} {
	base64Data := base64.StdEncoding.EncodeToString([]byte(o.Content))
	return map[string]interface{}{
		"content": base64Data,
		"message": o.Message,
	}
}

type GiteaDeleteFileOptions struct {
	Message string `json:"message"`
	Sha     string `json:"sha"`
}
