package main

import (
	"fmt"
	"regexp"
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

func (f *Focus) parseArgs(args []string) (*regexp.Regexp, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("focus: only one argument expected but %d given", len(args))
	}
	parts := strings.SplitN(args[0], "=", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("focus: bad parameter %q", args[0])
	}
	if parts[0] != "match" {
		return nil, fmt.Errorf("focus: unrecognized parameter %q", parts[0])
	}
	re, err := regexp.Compile(parts[1])
	if err != nil {
		return nil, fmt.Errorf("invalid match regexp: %s", err)
	}
	return re, nil
}

func (f *Focus) Execute(args []string) error {
	re, err := f.parseArgs(args)
	if err != nil {
		return err
	}

	x, err := xgbutil.NewConn()
	if err != nil {
		return err
	}
	// This returns an array from bottom -> top
	clients, err := ewmh.ClientListStackingGet(x)
	if err != nil {
		return err
	}
	// Go from the top to the bottom; find the first match
	for i := len(clients) - 1; i >= 0; i-- {
		client := clients[i]
		name, err := ewmh.WmNameGet(x, client)
		if err != nil || name == "" {
			name, err = icccm.WmNameGet(x, client)
			if err != nil {
				continue
			}
		}
		if re.MatchString(name) {
			fmt.Printf("focus: found match %q\n", name)
			if err := ewmh.ActiveWindowReq(x, client); err != nil {
				return err
			}
			return nil
		}
	}
	return fmt.Errorf("focus: no window matched %q", re)
}
