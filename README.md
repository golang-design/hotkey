# hotkey [![PkgGoDev](https://pkg.go.dev/badge/golang.design/x/hotkey)](https://pkg.go.dev/golang.design/x/hotkey) ![](https://changkun.de/urlstat?mode=github&repo=golang-design/hotkey) ![hotkey](https://github.com/golang-design/hotkey/workflows/hotkey/badge.svg?branch=main)

cross platform hotkey package in Go

```go
import "golang.design/x/hotkey"
```

## Features

- Cross platform supports: macOS, Linux (X11), and Windows
- Global hotkey registration without focus on a window

## API Usage

Package `hotkey` provides the basic facility to register a system-level
hotkey so that the application can be notified if a user triggers the
desired hotkey. By definition, a hotkey is a combination of modifiers
and a single key, and thus register a hotkey that contains multiple
keys is not supported at the moment. Furthermore, because of OS
restriction, hotkey events must be handled on the main thread.

Therefore, in order to use this package properly, here is a complete
example that corporates the mainthread:
package:

```go
package main

import (
	"context"

	"golang.design/x/hotkey"
	"golang.design/x/hotkey/mainthread"
)

// initialize mainthread facility so that hotkey can be
// properly registered to the system and handled by the
// application.
func main() { mainthread.Run(fn) }
func fn() { // Use fn as the actual main function.
	var (
		mods = []hotkey.Modifier{hotkey.ModCtrl, hotkey.ModShift}
		k    = hotkey.KeyS
	)

	// Register a desired hotkey.
	hk := hotkey.New(mods, k)
	if err := hk.Register() err != nil {
		panic("hotkey registration failed")
	}

	// Start listen hotkey event whenever you feel it is ready.
	for range hk.Listen() {
		println("hotkey ctrl+shift+s is triggered")
	}
}
```

## Who is using this package?

The main purpose of building this package is to support the
[midgard](https://changkun.de/s/midgard) project.

To know more projects, check our [wiki](https://github.com/golang-design/clipboard/wiki) page.

## License

MIT | &copy; 2021 The golang.design Initiative Authors, written by [Changkun Ou](https://changkun.de).