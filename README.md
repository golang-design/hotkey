# hotkey [![PkgGoDev](https://pkg.go.dev/badge/golang.design/x/hotkey)](https://pkg.go.dev/golang.design/x/hotkey) ![](https://changkun.de/urlstat?mode=github&repo=golang-design/hotkey) ![hotkey](https://github.com/golang-design/hotkey/workflows/hotkey/badge.svg?branch=main)

cross platform hotkey package in Go

```go
import "golang.design/x/hotkey"
```

## Features

- Cross platform supports: macOS, Linux (X11), and Windows
- Global hotkey registration without focus on a window

## API Usage

Package hotkey provides the primary facility to register a system-level global hotkey shortcut to notify an application if a user triggers the desired hotkey. A hotkey must be a combination of modifiers and a single key.

Note a platform-specific detail on `macOS` due to the OS restriction (other platforms do not have this restriction): hotkey events must be handled on the "main thread". Therefore, to use this package properly, one must call the `(*Hotkey).Register` method on the main thread, and an OS app main event loop must be established. One can use  the provided `golang.design/x/hotkey/mainthread` for self-contained applications. For applications based on other GUI frameworks, one has to use their provided ability to run the `(*Hotkey).Register` on the main thread. See the [examples](./examples) folder for more examples.

## Who is using this package?

The main purpose of building this package is to support the
[midgard](https://changkun.de/s/midgard) project.

To know more projects, check our [wiki](https://github.com/golang-design/hotkey/wiki) page.

## License

MIT | &copy; 2021 The golang.design Initiative Authors, written by [Changkun Ou](https://changkun.de).