package main

import (
	"fmt"

	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/ewmh"
)

type Maximize struct{}

func init() { commands["maximize"] = &Maximize{} }

func (m *Maximize) Help() string {
	return `maximize usage:
    maximize
Toggles the maximized state of the active window.`
}

func (m *Maximize) Execute(args []string) error {
	if len(args) > 0 {
		return fmt.Errorf("maximize: no arguments expected but %d given", len(args))
	}

	// Connect to X
	X, err := xgbutil.NewConn()
	if err != nil {
		return err
	}

	// Find the active window
	active, err := ewmh.ActiveWindowGet(X)
	if err != nil {
		return err
	}

	return ewmh.WmStateReqExtra(X, active, ewmh.StateToggle,
		"_NET_WM_STATE_MAXIMIZED_VERT", "_NET_WM_STATE_MAXIMIZED_HORZ", 2)
}
