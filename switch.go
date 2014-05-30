package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/BurntSushi/xgbutil/xwindow"
)

type Switch struct{}

func init() { commands["switch"] = &Switch{} }

// TODO: handle non-horizontal multi-monitor layouts.

func (s *Switch) Help() string {
	return `switch usage:
    switch [head=n|dir=[left|right]]
Switches the active window to be on a different head. An absolute, 0-indexed head may be provided, or else a
direction (left or right, with wrapping).`
}

const (
	dirNone = iota
	dirLeft
	dirRight
)

func (s *Switch) Execute(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("switch: expected one argument but %d given", len(args))
	}
	arg := args[0]

	headIndex := 0
	relativeDir := dirNone

	parts := strings.SplitN(arg, "=", 2)
	if len(parts) != 2 {
		return fmt.Errorf("bad parameter: %s", arg)
	}
	name := strings.TrimSpace(parts[0])
	switch name {
	case "head":
		var err error
		headIndex, err = strconv.Atoi(strings.TrimSpace(parts[1]))
		if err != nil {
			return fmt.Errorf("cannot parse head number: %s", parts[1])
		}
	case "dir":
		switch parts[1] {
		case "left":
			relativeDir = dirLeft
		case "right":
			relativeDir = dirRight
		default:
			return fmt.Errorf("bad direction: %s", parts[1])
		}
	default:
		return fmt.Errorf("bad parameter: %s", name)
	}

	// Connect to X
	X, err := xgbutil.NewConn()
	if err != nil {
		return err
	}

	heads, err := findHeads(X)
	if err != nil {
		return err
	}
	if len(heads) < 2 {
		return nil // Nothing to do
	}

	// Find the active window
	active, err := ewmh.ActiveWindowGet(X)
	if err != nil {
		return err
	}
	dgeom, err := xwindow.New(X, active).DecorGeometry()
	if err != nil {
		return err
	}

	i, err := findAssociatedHead(dgeom, heads)
	if err != nil {
		return err
	}

	if relativeDir == dirNone {
		if headIndex < 0 || headIndex >= len(heads) {
			return fmt.Errorf("bad head index %d", headIndex)
		}
	} else {
		switch relativeDir {
		case dirLeft:
			headIndex = (i + len(heads) - 1) % len(heads)
		case dirRight:
			headIndex = (i + 1) % len(heads)
		}
	}

	if headIndex == i {
		return nil // Nothing to do
	}

	newHead := heads[headIndex]
	x := newHead.X() + (newHead.Width() / 2) - (dgeom.Width() / 2)
	y := newHead.Y() + (newHead.Height() / 2) - (dgeom.Height() / 2)
	if err := ewmh.MoveresizeWindow(X, active, x, y, dgeom.Width(), dgeom.Height()); err != nil {
		return err
	}
	return nil
}
