package core

type Server struct {
	Id              int64  `db:"id"`
	OrgId           int32  `db:"org_id"`
	Name            string `db:"name"`
	Description     string `db:"description"`
	OperatingSystem string `db:"operating_system"`
	Location        string `db:"location"`
	IpAddress       string `db:"ip_address"`
}

type ServerHandle struct {
	Id    int64
	OrgId int32
}

type ServerSSHGenericConnection struct {
	Id       int64
	Username string
}

type ServerSSHPasswordConnection struct {
	Id       int64  `db:"id"`
	ServerId int64  `db:"server_id"`
	OrgId    int32  `db:"org_id"`
	Username string `db:"username"`
	Password string `db:"password"`
}

type ServerSSHKeyConnection struct {
	Id         int64  `db:"id"`
	ServerId   int64  `db:"server_id"`
	OrgId      int32  `db:"org_id"`
	Username   string `db:"username"`
	PrivateKey string `db:"private_key"`
}

type ServerConnectionOptions struct {
	SshPassword *ServerSSHPasswordConnection
	SshKey      *ServerSSHKeyConnection
}

func (c ServerSSHPasswordConnection) Generic() ServerSSHGenericConnection {
	return ServerSSHGenericConnection{
		Id:       c.Id,
		Username: c.Username,
	}
}

func (c ServerSSHKeyConnection) Generic() ServerSSHGenericConnection {
	return ServerSSHGenericConnection{
		Id:       c.Id,
		Username: c.Username,
	}
}
