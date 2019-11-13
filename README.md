# Carlisle

X window management utility.

Status: Alpha. Useful for me. Liable to change.

## Installation

    $ go get github.com/cespare/carlisle

## Usage

    $ carlisle COMMAND arg1 arg2 ...

Do `carlisle -h` for help, or `carlisle COMMAND help` for help with a specific
command.

### moveresize

Set the position and size of the active window.

    $ carlisle moveresize x=v1 y=v2 w=v3 h=v4

Any parameters not provided will remain unchanged from their current values.
Each value is an s-expression which may involve constants; the functions `+`,
`-`, `*`, `/`, `min`, and `max`; and the following variables:

* `x` - window's current x-value
* `y` - window's current y-value
* `w` - window's current width
* `h` - window's current height
* `sw` - screen width
* `sh` - screen height

One detail to be aware of is that the coordinates you use are relative to the
usable portion of the display. If you have, say, a panel at the top or bottom of
the screen, the window coordinates and screen height/width are all relative to
the rectangle excluding the panel(s).

**Example:** Here's how you can use moveresize to make the active window occupy
the right half of your screen:

    $ carlisle moveresize 'x=(* 0.5 sw)' y=0 'w=(* 0.5 sw)' h=sh

### focus

Activate (raise and focus) a window by matching its title.

    $ carlisle focus match=substr

This just looks through the window stack from top to bottom and finds the first
window whose title contains the substring. The match is case-insensitive.

**Example:** This is how you might quickly focus a gvim window:

    $ carlisle focus match=gvim

### toggle

Toggles the maximized or fullscreen state of the active window.

    $ carlisle toggle maximized

### switch

Move the active window to a different head. You can provide absolute (`head=N`,
0-indexed) or relative (`dir=left`, `dir=right`) positions. Relative directions
wrap around.

    $ carlisle switch dir=right

## Examples

Note that you'll probably want single quotes around your `moveresize` argument
strings or your shell will split them up.

`moveresize` is enough to implement most kinds of positioning commands. Here are
some ideas:

```
# Move 100px down
$ carlisle moveresize 'y=(+ y 100)'

# Move 100px right, don't move off-screen
$ carlisle moveresize 'x=(min (+ x 100) (- sw w))'

# Left half
$ carlisle moveresize x=0 y=0 'w=(* 0.5 sw)' h=sh

# Expand 100px down without expanding off-screen
$ carlisle moveresize 'h=(min (+ h 100) (- sh y))'
```

I use [xbindkeys](http://www.nongnu.org/xbindkeys) to bind hotkeys to Carlisle
commands, but you can use whatever tool you're most comfortable with (your
desktop environment or window manager probably provides such functionality). If
you'd like more ideas, check out my [`.xbindkeysrc`](https://github.com/cespare/dotfiles/blob/master/.xbindkeysrc).

## Similar tools

* [wmctrl](http://tomas.styblo.name/wmctrl/), a similar tool (doesn't take
  expressions for movement).
* [Slate](https://github.com/jigish/slate) for Mac OS X.

## License

MIT (see LICENSE.txt)
