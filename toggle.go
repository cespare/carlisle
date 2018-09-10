package main

import (
	"fmt"

	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/ewmh"
)

type Toggle struct{}

func init() { commands["toggle"] = Toggle{} }

func (t Toggle) Help() string {
	return `toggle usage:

    toggle [maximized|fullscreen]

Toggles the maximized or fullscreen state of the active window.`
}

func (t Toggle) Execute(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("toggle: one argument expected but %d given", len(args))
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

	switch args[0] {
	case "maximized":
		return ewmh.WmStateReqExtra(X, active, ewmh.StateToggle,
			"_NET_WM_STATE_MAXIMIZED_VERT", "_NET_WM_STATE_MAXIMIZED_HORZ", 2)
	case "fullscreen":
		return ewmh.WmStateReq(X, active, ewmh.StateToggle, "_NET_WM_STATE_FULLSCREEN")
	}
	return fmt.Errorf("toggle: unrecognized argument %q", args[0])
}
