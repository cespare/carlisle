package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/BurntSushi/xgbutil/xwindow"
)

type MoveResizeState struct {
	X, Y, W, H       int
	ScreenW, ScreenH int
}

type MoveResize struct {
	X, Y, W, H Arith
}

func init() { commands["moveresize"] = &MoveResize{} }

func (m *MoveResize) Help() string {
	return `moveresize usage:

    moveresize x=x1 y=y1 w=w1 h=h1

The values can be math s-expressions.`
}

func (m *MoveResize) parseArgs(args []string) error {
	for _, arg := range args {
		parts := strings.SplitN(arg, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("Bad parameter: %s", arg)
		}
		var arith Arith
		var err error
		name := strings.TrimSpace(parts[0])

		switch name {
		case "x", "y", "w", "h":
			arith, err = ParseArith([]byte(parts[1]))
			if err != nil {
				return err
			}
		}
		switch name {
		case "x":
			if m.X != nil {
				return errors.New("Duplicate param x")
			}
			m.X = arith
		case "y":
			if m.Y != nil {
				return errors.New("Duplicate param y")
			}
			m.Y = arith
		case "w":
			if m.W != nil {
				return errors.New("Duplicate param w")
			}
			m.W = arith
		case "h":
			if m.H != nil {
				return errors.New("Duplicate param h")
			}
			m.H = arith
		}
	}
	return nil
}

func (m *MoveResize) Execute(args []string) error {
	if err := m.parseArgs(args); err != nil {
		return err
	}
	if m.X == nil && m.Y == nil && m.W == nil && m.H == nil {
		return errors.New("moveresize: nothing to do")
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

	// Find the active window
	current, err := ewmh.ActiveWindowGet(X)
	if err != nil {
		return err
	}
	dgeom, err := xwindow.New(X, current).DecorGeometry()
	if err != nil {
		return err
	}

	i, err := findAssociatedHead(dgeom, heads)
	if err != nil {
		return err
	}
	activeHead := heads[i]

	if err := applyStruts(X, heads); err != nil {
		return err
	}

	state := &MoveResizeState{
		ScreenW: activeHead.Width(),
		ScreenH: activeHead.Height(),
	}

	// We make it appear to the user as though (0,0) is at the top left of the usable space (so if you have a
	// 25px bar at the top of your screen, what the user sees as the origin is actually x = 0, y = 25).
	state.X = dgeom.X() - activeHead.X()
	state.Y = dgeom.Y() - activeHead.Y()
	state.W = dgeom.Width()
	state.H = dgeom.Height()

	extents, err := ewmh.FrameExtentsGet(X, current)
	if err != nil {
		return err
	}

	x := dgeom.X()
	if m.X != nil {
		xf, err := m.X.Eval(state)
		if err != nil {
			return err
		}
		x = int(xf) + activeHead.X()
	}
	y := dgeom.Y()
	if m.Y != nil {
		yf, err := m.Y.Eval(state)
		if err != nil {
			return err
		}
		y = int(yf) + activeHead.Y()
	}
	w := dgeom.Width()
	if m.W != nil {
		wf, err := m.W.Eval(state)
		if err != nil {
			return err
		}
		w = int(wf)
	}
	w -= (extents.Left + extents.Right)
	h := dgeom.Height()
	if m.H != nil {
		hf, err := m.H.Eval(state)
		if err != nil {
			return err
		}
		h = int(hf)
	}
	h -= (extents.Top + extents.Bottom)

	if err := ewmh.MoveresizeWindow(X, current, x, y, w, h); err != nil {
		return err
	}
	return nil
}
