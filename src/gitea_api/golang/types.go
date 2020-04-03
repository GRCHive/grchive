package gitea

import (
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
