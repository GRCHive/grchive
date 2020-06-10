package main

import (
	"github.com/txn2/txeh"
	"gitlab.com/grchive/grchive/core"
	"sync"
)

type HostManager struct {
	hosts *txeh.Hosts

	currentDest string
	mux         sync.Mutex
}

func (h *HostManager) AddOverride(dest string, alias string) error {
	h.mux.Lock()
	h.hosts.AddHost(dest, alias)
	h.currentDest = dest
	return h.hosts.Save()
}

func (h *HostManager) RemoveOverride(dest string) error {
	defer h.mux.Unlock()
	h.hosts.RemoveAddress(dest)
	return h.hosts.Save()
}

func createHostManager() HostManager {
	hsts, err := txeh.NewHostsDefault()
	if err != nil {
		core.Error(err.Error())
	}

	return HostManager{
		hosts: hsts,
	}
}

var GlobalHostManager = createHostManager()
