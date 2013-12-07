package main

import (
	"fmt"
	"log"

	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/BurntSushi/xgbutil/xinerama"
	"github.com/BurntSushi/xgbutil/xrect"
	"github.com/BurntSushi/xgbutil/xwindow"
)

// Errors are all check
func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	X, err := xgbutil.NewConn()
	check(err)

	root := xwindow.New(X, X.RootWin())
	rgeom, err := root.Geometry()
	check(err)

	var heads xinerama.Heads
	if X.ExtInitialized("XINERAMA") {
		heads, err = xinerama.PhysicalHeads(X)
		check(err)
	} else {
		heads = xinerama.Heads{rgeom}
	}

	if len(heads) != 1 {
		log.Fatal(">1 heads are not handled for now.")
	}

	for i, head := range heads {
		fmt.Printf("Head #%d: %s\n", i+1, head)
	}
	fmt.Println("---------")

	clients, err := ewmh.ClientListGet(X)
	check(err)
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
	for i, head := range heads {
		fmt.Printf("Head #%d: %s\n", i+1, head)
	}

	current, err := ewmh.ActiveWindowGet(X)
	check(err)
	fmt.Printf("x=%d, y=%d, w=%d, h=%d\n", heads[0].X(), heads[0].Y(), heads[0].Width()/2, heads[0].Height())
	check(ewmh.MoveresizeWindow(X, current, heads[0].X(), heads[0].Y(), heads[0].Width()/2, heads[0].Height()))

	dgeom, err := xwindow.New(X, current).DecorGeometry()
	check(err)
	fmt.Printf("\033[01;34m>>>> dgeom: %v\x1B[m\n", dgeom)
}
