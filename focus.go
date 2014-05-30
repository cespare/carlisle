package main

import (
	"fmt"
	"strings"

	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/BurntSushi/xgbutil/icccm"
)

type Focus struct{}

func init() { commands["focus"] = &Focus{} }

func (f *Focus) Help() string {
	return `focus usage:
    focus match=string
The string is matched case-insensitively against all window titles, starting
from the top of the stack.`
}

func (f *Focus) parseArgs(args []string) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("focus: only one argument expected but %d given", len(args))
	}
	parts := strings.SplitN(args[0], "=", 2)
	if len(parts) != 2 {
		return "", fmt.Errorf("focus: bad parameter %q", args[0])
	}
	if parts[0] != "match" {
		return "", fmt.Errorf("focus: unrecognized parameter %q", parts[0])
	}
	match := strings.TrimSpace(parts[1])
	if match == "" {
		return "", fmt.Errorf("focus: empty match")
	}
	return match, nil
}

func (f *Focus) Execute(args []string) error {
	match, err := f.parseArgs(args)
	if err != nil {
		return err
	}

	X, err := xgbutil.NewConn()
	if err != nil {
		return err
	}
	// This returns an array from bottom -> top
	clients, err := ewmh.ClientListStackingGet(X)
	if err != nil {
		return err
	}
	// Go from the top to the bottom; find the first match
	for i := len(clients) - 1; i >= 0; i-- {
		client := clients[i]
		name, err := ewmh.WmNameGet(X, client)
		if err != nil || name == "" {
			name, err = icccm.WmNameGet(X, client)
			if err != nil {
				continue
			}
		}
		if strings.Contains(strings.ToLower(name), match) {
			fmt.Printf("focus: found match %q\n", name)
			if err := ewmh.ActiveWindowReq(X, client); err != nil {
				return err
			}
			return nil
		}
	}
	return fmt.Errorf("focus: no window matched %q", match)
}
