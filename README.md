# Carlisle

X window management utility.

Status: Alpha. Useful for me. Liable to change.

## Installation

    $ go get github.com/cespare/carlisle

## Usage

    $ carlisle COMMAND arg1 arg2 ...

### moveresize

Set the position and size of the active window.

    $ carlisle moveresize x=v1 y=v2 w=v3 h=v4

Any parameters not provided will remain unchanged from their current values. Each value is an s-expression
which may involve constants; the functions `+`, `-`, `*`, `/`, `min`, and `max`; and the following variables:

* `x` - window's current x-value
* `y` - window's current y-value
* `w` - window's current width
* `h` - window's current height
* `sw` - screen width
* `sh` - screen height

One detail to be aware of is that the coordinates you use are relative to the usable portion of the display.
If you have, say, a panel at the top or bottom of the screen taking, the window coordinates and screen
height/width are all relative to the rectangle excluding the panel(s).

**Example:** Here's how you can use moveresize to make the active window occupy the right half of your screen:

    $ carlisle moveresize x=0 y=0 'w=(* 0.5 sw)' h=sh

### focus

Activate (raise and focus) a window by matching its title.

    $ carlisle focus match=substr

This just looks through the window stack from top to bottom and finds the first window whose title contains
the substring. The match is case-insensitive.

**Example:** This is how you might quickly focus a gvim window:

    $ carlisle focus match=gvim

## Examples

Note that you'll probably want single quotes around your `moveresize` argument strings or your shell will
split them up.

`moveresize` is enough to implement most kinds of positioning commands. For ideas, here is an excerpt from my
`.xbindkeysrc`:

```
# Move 100px right, don't move off-screen
"carlisle moveresize 'x=(min (+ x 100) (- sw w))'"
  mod4 + l

# Move 100px left, don't move off-screen
"carlisle moveresize 'x=(max (- x 100) 0)'"
  mod4 + h

# Move 100px down
"carlisle moveresize 'y=(+ y 100)'"
  mod4 + j

# Move 100px up; not above the top of the screen
"carlisle moveresize 'y=(max (- y 100) 0)'"
  mod4 + k

# Expand 100px to the right
"carlisle moveresize 'w=(+ w 100)'"
  mod4 + shift + l

# Contract 100px from right
"carlisle moveresize 'w=(max (- w 100) 0)'"
  mod4 + shift + h

# Expand 100px down
"carlisle moveresize 'h=(+ h 100)'"
  mod4 + shift + j

# Contract 100px up
"carlisle moveresize 'h=(max (- h 100) 0)'"
  mod4 + shift + k

# Right half of screen
"carlisle moveresize x=0 y=0 'w=(* 0.5 sw)' h=sh"
  mod4 + shift + comma

# Left half of screen
"carlisle moveresize 'x=(* 0.5 sw)' y=0 'w=(* 0.5 sw)' h=sh"
  mod4 + shift + period
```

## To Do

* Multi-head support. (Not difficult, just need to hook up another monitor...)
* Test out with some various WMs (I use XFCE; BurntSushi indicates that KWin, at least, requires some hacks)
* Other commands (see below)


## Ideas

* Carlisle will take a single command and arg strings: `./carlisle command arg1 arg2 ...`
  - Each argument looks like an assignment `a=b`. Spaces are allowed (but you have to surround the argument
    string with quotes: `'a = b'`).
  - Generally, all parameters are optional and have reasonable defaults if not provided
* Move+resize
  - `moveresize 'x=0 y=0 w=(* 0.5 sw) h=sh'` -- put the current window on the left half of the current
    screen. `sw` is screen width, `sh` is screen height.
  - `moveresize 'w=(+ w 100)'` - increase the screen width by 100px (other parameters stay the same)
  - Math: `+`, `-`, `*`, `-`, `min`, `max`
* Focus
  - `focus 'match=foobar'` -- focus the first window with title matching "foobar".
* Move window to a different display
  - `movedisplay 'dir=left'`
* Move desktops
  - `movedesktop 'dir=right'`
* Pick a window by hitting a key or two
  - `startpicker` (no arguments)
  - Inspiration from vimium/slate -- display the window name/icon on each (and maybe shade the desktop). The
    user just needs to type an unambiguous prefix of the shortcut displayed.
* Full-screen, minimize

## Similar tools

* [wmctrl](http://tomas.styblo.name/wmctrl/), a similar tool (doesn't take expressions for movement).
* [Slate](https://github.com/jigish/slate) for Mac OS X.

## License

MIT (see LICENSE.txt)
