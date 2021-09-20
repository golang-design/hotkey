# mainthread [![PkgGoDev](https://pkg.go.dev/badge/golang.design/x/mainthread)](https://pkg.go.dev/golang.design/x/mainthread) ![mainthread](https://github.com/golang-design/mainthread/workflows/mainthread/badge.svg?branch=main) ![](https://changkun.de/urlstat?mode=github&repo=golang-design/mainthread)

schedule functions to run on the main thread

```go
import "golang.design/x/mainthread"
```

## Features

- Main thread scheduling
- Schedule functions without memory allocation

## API Usage

Package mainthread offers facilities to schedule functions
on the main thread. To use this package properly, one must
call `mainthread.Init` from the main package. For example:

```go
package main

import "golang.design/x/mainthread"

func main() { mainthread.Init(fn) }

// fn is the actual main function
func fn() {
	// ... do stuff ...

	// mainthread.Call returns when f1 returns. Note that if f1 blocks
	// it will also block the execution of any subsequent calls on the
	// main thread.
	mainthread.Call(f1)

	// ... do stuff ...


	// mainthread.Go returns immediately and f2 is scheduled to be
	// executed in the future.
	mainthread.Go(f2)

	// ... do stuff ...
}

func f1() { ... }
func f2() { ... }
```

If the given function triggers a panic, and called via `mainthread.Call`,
then the panic will be propagated to the same goroutine. One can capture
that panic, when possible:

```go
defer func() {
	if r := recover(); r != nil {
		println(r)
	}
}()

mainthread.Call(func() { ... }) // if panic
```

If the given function triggers a panic, and called via `mainthread.Go`,
then the panic will be cached internally, until a call to the `Error()` method:

```go
mainthread.Go(func() { ... }) // if panics

// ... do stuff ...

if err := mainthread.Error(); err != nil { // can be captured here.
	println(err)
}
```

Note that a panic happens before `mainthread.Error()` returning the
panicked error. If one needs to guarantee `mainthread.Error()` indeed
captured the panic, a dummy function can be used as synchornization:

```go
mainthread.Go(func() { panic("die") })	// if panics
mainthread.Call(func() {}) 				// for execution synchronization
err := mainthread.Error()				// err must be non-nil
```


It is possible to cache up to a maximum of 42 panicked errors.
More errors are ignored.

## When do you need this package?

Read this to learn more about the design purpose of this package:
https://golang.design/research/zero-alloc-call-sched/

## Who is using this package?

The initial purpose of building this package is to support writing
graphical applications in Go. To know projects that are using this
package, check our [wiki](https://github.com/golang-design/mainthread/wiki)
page.


## License

MIT | &copy; 2021 The golang.design Initiative Authors, written by [Changkun Ou](https://changkun.de).