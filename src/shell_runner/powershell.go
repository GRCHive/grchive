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

const WSmanConfiguration = "PowerShell.7"

func runPowershellWSman(
	data *RunData,
	server *core.Server,
	t *PerServerTracker,
	userpass *core.ServerSSHPasswordConnection,
) error {
	outputBuffer := bytes.Buffer{}

	wrapper, err := ioutil.TempFile(os.TempDir(), "powershell-wrapper-*.ps1")
	if err != nil {
		return err
	}
	defer os.Remove(wrapper.Name())

	localScript, err := ioutil.TempFile(os.TempDir(), "powershell-script-*.ps1")
	if err != nil {
		return err
	}
	defer os.Remove(localScript.Name())

	_, err = localScript.WriteString(data.scriptText)
	if err != nil {
		return err
	}

	_, err = wrapper.WriteString(fmt.Sprintf(`
$pw = ConvertTo-SecureString "$Env:PASSWORD" -AsPlainText -Force
$cred = New-Object System.Management.Automation.PSCredential ("%s", $pw)
$session = New-PSSession -ComputerName %s -Authentication Basic -Credential $cred -ConfigurationName "%s" -UseSSL -SessionOption (New-PSSessionOption -SkipCACheck -SkipCNCheck)
Invoke-Command -Session $session -FilePath %s
	`, userpass.Username, server.IpAddress, WSmanConfiguration, localScript.Name()))
	if err != nil {
		return err
	}

	cmd := exec.Command("pwsh", "-f", wrapper.Name())
	cmd.Stdout = &outputBuffer
	cmd.Stderr = &outputBuffer
	cmd.Env = []string{
		fmt.Sprintf("PASSWORD=%s", userpass.Password),
	}

	err = cmd.Run()
	if err != nil {
		outputBuffer.WriteString("\nFailed to run: " + err.Error())
	}

	t.MarkSuccessFailure(outputBuffer.String(), err == nil)
	return err
}

func runPowershellScript(
	data *RunData,
	server *core.Server,
	t *PerServerTracker,
) error {
	serverData := data.perServer[server.Id]

	if serverData.ConnectionChoices.SshPassword != nil {
		core.Info("Trying Username/Password...")
		err := runPowershellWSman(data, server, t, serverData.ConnectionChoices.SshPassword)
		if err == nil {
			core.Info("\tSuccess!")
			return nil
		}
		core.Warning("Failed to use the username/pasword connection: " + err.Error())
	}

	return errors.New("No valid connection to server to run Powershell script.")
}
