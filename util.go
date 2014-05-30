package main

import (
	"fmt"

	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/BurntSushi/xgbutil/xinerama"
	"github.com/BurntSushi/xgbutil/xrect"
	"github.com/BurntSushi/xgbutil/xwindow"
)

// findHeads returns all the display heads.
func findHeads(X *xgbutil.XUtil) (xinerama.Heads, error) {
	root := xwindow.New(X, X.RootWin())
	rgeom, err := root.Geometry()
	if err != nil {
		return nil, err
	}

	// Locate all the displays
	if X.ExtInitialized("XINERAMA") {
		return xinerama.PhysicalHeads(X)
	}
	return xinerama.Heads{rgeom}, nil
}

// applyStruts finds and applies all struts to the given heads.
func applyStruts(X *xgbutil.XUtil, heads xinerama.Heads) error {
	clients, err := ewmh.ClientListGet(X)
	if err != nil {
		return err
	}

	rgeom := heads[0] // TODO: This definitely is wrong. Test on a xinerama machine

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
	return nil
}

// findAssociatedHead determines which head is associated with a window by using a simple heuristic: pick the
// first head containing the centerpoint of the window.
func findAssociatedHead(window xrect.Rect, heads xinerama.Heads) (index int, err error) {
	if len(heads) == 1 {
		return 0, nil
	}

	centerX := window.X() + (window.Width() / 2)
	centerY := window.Y() + (window.Height() / 2)
	for i, h := range heads {
		if h.X() <= centerX && (h.X()+h.Width()) >= centerX &&
			h.Y() <= centerY && (h.Y()+h.Height()) >= centerY {
			return i, nil
		}
	}
	return 0, fmt.Errorf("cannot locate head for window %s", window)
}
