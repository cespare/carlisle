# Carlisle

X window management utility.

Status: pre-pre-alpha. Just experimenting with the APIs right now.

## Notes (to-do list)

* Carlisle will take a single command and arg strings: `./carlisle command arg1 arg2 ...`
  - Each argument looks like an assignment `a=b`. Spaces are allowed (but you have to surround the argument
    string with quotes: `'a = b'`).
  - Generally, all parameters are optional and have reasonable defaults if not provided
* Move+resize
  - `moveresize 'x=0 y=0 w=(* 0.5 sw) h=sh'` -- put the current window on the left half of the current
    screen. `sw` is screen width, `sh` is screen height.
  - `moveresize 'w=(+ w 100)'` - increase the screen width by 100px (other parameters stay the same)
  - Other parameters?
    - Horizontal/vertical gravity when resizing
  - Math: `+`, `-`, `*`, `-`, `min`, `max`
* Focus
  - `focus 'match=foobar'` -- focus the first window with title matching "foobar".
* Move screens
  - `movescreen 'dir=left'`
* Move desktops
  - `movedesktop 'dir=right'`
* Pick a window by hitting a key or two
  - `startpicker` (no arguments)
  - Inspiration from vimium/slate -- display the window name/icon on each (and maybe shade the desktop). The
    user just needs to type an unambiguous prefix of the shortcut displayed.
