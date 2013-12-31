package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/BurntSushi/xgbutil/xinerama"
	"github.com/BurntSushi/xgbutil/xrect"
	"github.com/BurntSushi/xgbutil/xwindow"
)

type MoveResizeState struct {
	W, H             int
	ScreenW, ScreenH int
}

type MoveResize struct {
	X, Y, W, H Arith
}

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
	fmt.Printf("\033[01;34m>>>> m: %v\x1B[m\n", m)
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
	root := xwindow.New(X, X.RootWin())
	rgeom, err := root.Geometry()
	if err != nil {
		return err
	}

	// Locate all the displays
	var heads xinerama.Heads
	if X.ExtInitialized("XINERAMA") {
		heads, err = xinerama.PhysicalHeads(X)
		if err != nil {
			return err
		}
	} else {
		heads = xinerama.Heads{rgeom}
	}

	if len(heads) != 1 {
		return errors.New(">1 heads are not handled for now.")
	}

	// Find the struts so we know what area we have to work with
	clients, err := ewmh.ClientListGet(X)
	if err != nil {
		return err
	}
	for _, client := range clients {
		strut, err := ewmh.WmStrutPartialGet(X, client)
		if err != nil {
			continue
		}

		xrect.ApplyStrut(heads, uint(rgeom.Width()), uint(rgeom.Height()),
			strut.Left, strut.Right, strut.Top, strut.Bottom,
			strut.LeftStartY, strut.LeftEndY, strut.RightStartY, strut.RightEndY,
			strut.TopStartX, strut.TopEndX, strut.BottomStartX, strut.BottomEndX,
		)
	}

	state := &MoveResizeState{
		ScreenW: heads[0].Width(),
		ScreenH: heads[0].Height(),
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
		x = int(xf)
	}
	y := dgeom.Y()
	if m.Y != nil {
		yf, err := m.Y.Eval(state)
		if err != nil {
			return err
		}
		y = int(yf)
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
