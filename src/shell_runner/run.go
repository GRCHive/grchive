package main

import (
	"gitlab.com/grchive/grchive/core"
)

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
		return runPowershellScript(data, server, t)
	}

	return nil
}
