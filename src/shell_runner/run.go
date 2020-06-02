package main

import (
	"bytes"
	"errors"
	"fmt"
	"gitlab.com/grchive/grchive/core"
	"io/ioutil"
	"os"
	"os/exec"
)

func sharedSshArgs(data *RunData, server *core.Server, username string) []string {
	args := []string{
		"-o",
		"StrictHostKeyChecking=no",
		fmt.Sprintf("%s@%s", username, server.IpAddress),
	}
	return args
}

func genericRunSshStdin(cmd *exec.Cmd, script string, t *PerServerTracker) error {
	outputBuffer := bytes.Buffer{}
	inputBuffer := bytes.NewBufferString(script)

	cmd.Stdin = inputBuffer
	cmd.Stdout = &outputBuffer
	cmd.Stderr = &outputBuffer

	err := cmd.Run()

	if err != nil {
		outputBuffer.WriteString("\nFailed to run: " + err.Error())
	}

	t.MarkSuccessFailure(outputBuffer.String(), err == nil)
	return err
}

func runBashScriptSshPassword(
	data *RunData,
	server *core.Server,
	t *PerServerTracker,
	sshPassword *core.ServerSSHPasswordConnection,
) error {
	args := append(
		[]string{
			"-p",
			sshPassword.Password,
			"ssh",
		},
		sharedSshArgs(data, server, sshPassword.Username)...,
	)

	cmd := exec.Command("sshpass", args...)
	return genericRunSshStdin(cmd, data.scriptText, t)
}

func runBashScriptSshKey(
	data *RunData,
	server *core.Server,
	t *PerServerTracker,
	sshKey *core.ServerSSHKeyConnection,
) error {
	keyFile, err := ioutil.TempFile(os.TempDir(), "privkey")
	if err != nil {
		return err
	}
	defer os.Remove(keyFile.Name())

	_, err = keyFile.WriteString(sshKey.PrivateKey)
	if err != nil {
		return err
	}

	args := append(
		[]string{
			"-i",
			keyFile.Name(),
		},
		sharedSshArgs(data, server, sshKey.Username)...,
	)

	cmd := exec.Command("ssh", args...)
	return genericRunSshStdin(cmd, data.scriptText, t)
}

func runBashScript(
	data *RunData,
	server *core.Server,
	t *PerServerTracker,
) error {
	serverData := data.perServer[server.Id]

	// The user isn't given a choice of connection preferences.
	// We use the first one that works.
	if serverData.ConnectionChoices.SshKey != nil {
		core.Info("Trying SSH Key...")
		err := runBashScriptSshKey(data, server, t, serverData.ConnectionChoices.SshKey)
		if err == nil {
			core.Info("\tSuccess!")
			return nil
		}
		core.Warning("Failed to use the SSH Key connection: " + err.Error())
	}

	if serverData.ConnectionChoices.SshPassword != nil {
		core.Info("Trying SSH Password...")
		err := runBashScriptSshPassword(data, server, t, serverData.ConnectionChoices.SshPassword)
		if err == nil {
			core.Info("\tSuccess!")
			return nil
		}
		core.Warning("Failed to use the SSH pasword connection: " + err.Error())
	}

	return errors.New("No valid connection to server to run Bash script.")
}

func runScriptOnServer(
	data *RunData,
	server *core.Server,
	t *PerServerTracker,
) error {
	defer t.Finish()
	err := t.MarkStart()
	if err != nil {
		return err
	}

	switch data.shell.TypeId {
	case core.BashShellId:
		return runBashScript(data, server, t)
	case core.PowerShellId:
		return errors.New("Powershell Currently Unsupported")
	}

	return nil
}
