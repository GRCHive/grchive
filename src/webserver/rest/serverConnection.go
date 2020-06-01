package rest

import (
	"bytes"
	"encoding/json"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	"gitlab.com/grchive/grchive/webcore"
	"io"
	"net/http"
)

func getEncryptedSshKeyFromRequest(r *http.Request) (string, error) {
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		core.Warning("Can't find uploaded file: " + err.Error())
		return "", err
	}
	defer file.Close()

	if fileHeader.Size > webcore.MaxFileSizeBytes {
		core.Warning("File too large.")
		return "", err
	}

	var buffer bytes.Buffer
	_, err = io.Copy(&buffer, file)
	if err != nil {
		core.Warning("Could not read file: " + err.Error())
		return "", err
	}

	encryptedKey, err := webcore.CreateEncryptedPassword(buffer.String())
	if err != nil {
		core.Warning("Failed to find encrypt password: " + err.Error())
		return "", err
	}

	return encryptedKey, nil
}

type NewServerConnectionSSHPasswordInputs struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func newServerConnectionSSHPassword(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := NewServerConnectionSSHPasswordInputs{}
	err := webcore.UnmarshalRequestForm(r, &inputs)
	if err != nil {
		core.Warning("Can't parse inputs: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	server, err := webcore.FindServerInContext(r.Context())
	if err != nil {
		core.Warning("Failed to find server in context: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	org, err := webcore.FindOrganizationInContext(r.Context())
	if err != nil {
		core.Warning("Failed to find org in context: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	encryptedPassword, err := webcore.CreateEncryptedPassword(inputs.Password)
	if err != nil {
		core.Warning("Failed to find encrypt password: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	conn := core.ServerSSHPasswordConnection{
		ServerId: server.Id,
		OrgId:    org.Id,
		Username: inputs.Username,
		Password: encryptedPassword,
	}

	tx := database.CreateTx()
	err = database.WrapTx(tx, func() error {
		return database.NewServerSSHPasswordConnectionWithTx(tx, &conn)
	})

	if err != nil {
		core.Warning("Failed to create new SSH password connection: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(conn.Generic())
}

type NewServerConnectionSSHKeyInputs struct {
	Username string `webcore:"username"`
}

func newServerConnectionSSHKey(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := NewServerConnectionSSHKeyInputs{}
	err := webcore.UnmarshalRequestForm(r, &inputs)
	if err != nil {
		core.Warning("Can't parse inputs: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	server, err := webcore.FindServerInContext(r.Context())
	if err != nil {
		core.Warning("Failed to find server in context: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	org, err := webcore.FindOrganizationInContext(r.Context())
	if err != nil {
		core.Warning("Failed to find org in context: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	conn := core.ServerSSHKeyConnection{
		ServerId: server.Id,
		OrgId:    org.Id,
		Username: inputs.Username,
	}

	conn.PrivateKey, err = getEncryptedSshKeyFromRequest(r)
	if err != nil {
		core.Warning("Failed to get private key: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	tx := database.CreateTx()
	err = database.WrapTx(tx, func() error {
		return database.NewServerSSHKeyConnectionWithTx(tx, &conn)
	})

	if err != nil {
		core.Warning("Failed to create new SSH password connection: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(conn.Generic())
}

func deleteServerConnectionSSHPassword(w http.ResponseWriter, r *http.Request) {
	conn, err := webcore.FindServerSSHPasswordConnectionInContext(r.Context())
	if err != nil {
		core.Warning("Failed to find SSH password connection: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tx := database.CreateTx()
	err = database.WrapTx(tx, func() error {
		return database.DeleteSSHPasswordConnectionWithTx(tx, conn.Id)
	})

	if err != nil {
		core.Warning("Failed to delete SSH password connection: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func deleteServerConnectionSSHKey(w http.ResponseWriter, r *http.Request) {
	conn, err := webcore.FindServerSSHKeyConnectionInContext(r.Context())
	if err != nil {
		core.Warning("Failed to find SSH key connection: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tx := database.CreateTx()
	err = database.WrapTx(tx, func() error {
		return database.DeleteSSHKeyConnectionWithTx(tx, conn.Id)
	})

	if err != nil {
		core.Warning("Failed to delete SSH password connection: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func getServerConnectionSSHPassword(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	conn, err := webcore.FindServerSSHPasswordConnectionInContext(r.Context())
	if err != nil {
		core.Warning("Failed to find SSH password connection: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	conn.Password, err = webcore.DecryptEncryptedPassword(conn.Password)
	if err != nil {
		core.Warning("Failed to decrypt password: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(conn)
}

func getServerConnectionSSHKey(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	conn, err := webcore.FindServerSSHKeyConnectionInContext(r.Context())
	if err != nil {
		core.Warning("Failed to find SSH password connection: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	conn.PrivateKey, err = webcore.DecryptEncryptedPassword(conn.PrivateKey)
	if err != nil {
		core.Warning("Failed to decrypt private key: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(conn)
}

func editServerConnectionSSHPassword(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	conn, err := webcore.FindServerSSHPasswordConnectionInContext(r.Context())
	if err != nil {
		core.Warning("Failed to find SSH password connection: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	inputs := NewServerConnectionSSHPasswordInputs{}
	err = webcore.UnmarshalRequestForm(r, &inputs)
	if err != nil {
		core.Warning("Can't parse inputs: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	encryptedPassword, err := webcore.CreateEncryptedPassword(inputs.Password)
	if err != nil {
		core.Warning("Failed to find encrypt password: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	conn.Username = inputs.Username
	conn.Password = encryptedPassword

	tx := database.CreateTx()
	err = database.WrapTx(tx, func() error {
		return database.EditServerSSHPasswordConnectionWithTx(tx, conn)
	})

	if err != nil {
		core.Warning("Failed to edit server SSH password connection: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(conn.Generic())
}

func editServerConnectionSSHKey(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	conn, err := webcore.FindServerSSHKeyConnectionInContext(r.Context())
	if err != nil {
		core.Warning("Failed to find SSH key connection: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	inputs := NewServerConnectionSSHKeyInputs{}
	err = webcore.UnmarshalRequestForm(r, &inputs)
	if err != nil {
		core.Warning("Can't parse inputs: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	conn.Username = inputs.Username
	conn.PrivateKey, err = getEncryptedSshKeyFromRequest(r)
	if err != nil {
		core.Warning("Can't get encrypted private key: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	tx := database.CreateTx()
	err = database.WrapTx(tx, func() error {
		return database.EditServerSSHKeyConnectionWithTx(tx, conn)
	})

	if err != nil {
		core.Warning("Failed to edit server SSH key connection: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(conn.Generic())
}

func allServerConnections(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	server, err := webcore.FindServerInContext(r.Context())
	if err != nil {
		core.Warning("Failed to find server in context: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	org, err := webcore.FindOrganizationInContext(r.Context())
	if err != nil {
		core.Warning("Failed to find org in context: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response := struct {
		SshPassword *core.ServerSSHGenericConnection
		SshKey      *core.ServerSSHGenericConnection
	}{}

	sshPassword, err := database.GetSSHPasswordConnectionForServer(server.Id, org.Id)
	if err != nil {
		core.Warning("Failed to get SSH password connection: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if sshPassword != nil {
		genSshPassword := sshPassword.Generic()
		response.SshPassword = &genSshPassword
	}

	sshKey, err := database.GetSSHKeyConnectionForServer(server.Id, org.Id)
	if err != nil {
		core.Warning("Failed to get SSH key connection: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if sshKey != nil {
		genSshKey := sshKey.Generic()
		response.SshKey = &genSshKey
	}

	jsonWriter.Encode(response)
}
