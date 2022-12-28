# hotkey [![PkgGoDev](https://pkg.go.dev/badge/golang.design/x/hotkey)](https://pkg.go.dev/golang.design/x/hotkey) ![](https://changkun.de/urlstat?mode=github&repo=golang-design/hotkey) ![hotkey](https://github.com/golang-design/hotkey/workflows/hotkey/badge.svg?branch=main)

cross platform hotkey package in Go

```go
import "golang.design/x/hotkey"
```

## Features

- Cross platform supports: macOS, Linux (X11), and Windows
- Global hotkey registration without focus on a window

## API Usage

Package hotkey provides the basic facility to register a system-level
global hotkey shortcut so that an application can be notified if a user
triggers the desired hotkey. A hotkey must be a combination of modifiers
and a single key.

```go
package main

import (
	"log"

	"golang.design/x/hotkey"
	"golang.design/x/hotkey/mainthread"
)

func main() { mainthread.Init(fn) } // Not necessary when use in Fyne, Ebiten or Gio.
func fn() {
	hk := hotkey.New([]hotkey.Modifier{hotkey.ModCtrl, hotkey.ModShift}, hotkey.KeyS)
	err := hk.Register()
	if err != nil {
		log.Fatalf("hotkey: failed to register hotkey: %v", err)
		return
	}

	log.Printf("hotkey: %v is registered\n", hk)
	<-hk.Keydown()
	log.Printf("hotkey: %v is down\n", hk)
	<-hk.Keyup()
	log.Printf("hotkey: %v is up\n", hk)
	hk.Unregister()
	log.Printf("hotkey: %v is unregistered\n", hk)
}
```

Note platform specific details:

- On macOS, due to the OS restriction (other platforms does not have this
  restriction), hotkey events must be handled on the "main thread".
  Therefore, in order to use this package properly, one must start an OS
  main event loop on the main thread, For self-contained applications,
  using [golang.design/x/hotkey/mainthread](https://pkg.go.dev/golang.design/x/hotkey/mainthread)
  is possible. It is uncessary or applications based on other GUI frameworks,
  such as fyne, ebiten, or Gio. See the "[./examples](./examples)" folder
  for more examples.
- On Linux (X11), when AutoRepeat is enabled in the X server, the Keyup
  is triggered automatically and continuously as Keydown continues.
- On Linux (X11), some keys may be mapped to multiple Mod keys. To
  correctly register the key combination, one must use the correct
  underlying keycode combination. For example, a regular Ctrl+Alt+S
  might be registered as: Ctrl+Mod2+Mod4+S.
- If this package did not include a desired key, one can always provide
  the keycode to the API. For example, if a key code is 0x15, then the
  corresponding key is `hotkey.Key(0x15)`.

## Examples

| Description | Folder |
|:------------|:------:|
| A minimum example | [minimum](./examples/minimum/main.go) |
| Register multiple hotkeys | [multiple](./examples/multiple/main.go) |
| A example to use in GLFW | [glfw](./examples/glfw/main.go) |
| A example to use in Fyne | [fyne](./examples/fyne/main.go) |
| A example to use in Ebiten | [ebiten](./examples/ebiten/main.go) |
| A example to use in Gio | [gio](./examples/gio/main.go) |

## Who is using this package?

The main purpose of building this package is to support the
[midgard](https://changkun.de/s/midgard) project.

To know more projects, check our [wiki](https://github.com/golang-design/hotkey/wiki) page.

## License

MIT | &copy; 2021 The golang.design Initiative Authors, written by [Changkun Ou](https://changkun.de).