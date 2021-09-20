// Copyright 2020-2021 The golang.design Initiative Authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.
//
// Written by Changkun Ou <changkun.de>

// Package mainthread offers facilities to schedule functions
// on the main thread. To use this package properly, one must
// call `mainthread.Init` from the main package. For example:
//
// 	package main
//
// 	import "golang.design/x/mainthread"
//
// 	func main() { mainthread.Init(fn) }
//
// 	// fn is the actual main function
// 	func fn() {
// 		// ... do stuff ...
//
// 		// mainthread.Call returns when f1 returns. Note that if f1 blocks
// 		// it will also block the execution of any subsequent calls on the
// 		// main thread.
// 		mainthread.Call(f1)
//
// 		// ... do stuff ...
//
//
// 		// mainthread.Go returns immediately and f2 is scheduled to be
// 		// executed in the future.
// 		mainthread.Go(f2)
//
// 		// ... do stuff ...
// 	}
//
// 	func f1() { ... }
// 	func f2() { ... }
//
// If the given function triggers a panic, and called via `mainthread.Call`,
// then the panic will be propagated to the same goroutine. One can capture
// that panic, when possible:
//
// 	defer func() {
// 		if r := recover(); r != nil {
// 			println(r)
// 		}
// 	}()
//
// 	mainthread.Call(func() { ... }) // if panic
//
// If the given function triggers a panic, and called via `mainthread.Go`,
// then the panic will be cached internally, until a call to the `Error()` method:
//
// 	mainthread.Go(func() { ... }) // if panics
//
// 	// ... do stuff ...
//
// 	if err := mainthread.Error(); err != nil { // can be captured here.
// 		println(err)
// 	}
//
// Note that a panic happens before `mainthread.Error()` returning the
// panicked error. If one needs to guarantee `mainthread.Error()` indeed
// captured the panic, a dummy function can be used as synchornization:
//
// 	mainthread.Go(func() { panic("die") })	// if panics
// 	mainthread.Call(func() {}) 				// for execution synchronization
// 	err := mainthread.Error()				// err must be non-nil
//
// It is possible to cache up to a maximum of 42 panicked errors.
// More errors are ignored.
package mainthread // import "golang.design/x/mainthread"

import (
	"fmt"
	"runtime"
	"sync"
)

func init() {
	runtime.LockOSThread()
}

// Init initializes the functionality of running arbitrary subsequent
// functions be called on the main system thread.
//
// Init must be called in the main.main function.
func Init(main func()) {
	done := donePool.Get().(chan error)
	defer donePool.Put(done)

	go func() {
		defer func() {
			done <- nil
		}()
		main()
	}()

	for {
		select {
		case f := <-funcQ:
			func() {
				defer func() {
					r := recover()
					if f.done != nil {
						if r != nil {
							f.done <- fmt.Errorf("%v", r)
						} else {
							f.done <- nil
						}
					} else {
						if r != nil {
							select {
							case erroQ <- fmt.Errorf("%v", r):
							default:
							}
						}
					}
				}()
				f.fn()
			}()
		case <-done:
			return
		}
	}
}

// Call calls f on the main thread and blocks until f finishes.
func Call(f func()) {
	done := donePool.Get().(chan error)
	defer donePool.Put(done)

	data := funcData{fn: f, done: done}
	funcQ <- data
	if err := <-done; err != nil {
		panic(err)
	}
}

// Go schedules f to be called on the main thread.
func Go(f func()) {
	funcQ <- funcData{fn: f}
}

// Error returns an error that is captured if there are any panics
// happened on the mainthread.
//
// It is possible to cache up to a maximum of 42 panicked errors.
// More errors are ignored.
func Error() error {
	select {
	case err := <-erroQ:
		return err
	default:
		return nil
	}
}

var (
	funcQ    = make(chan funcData, runtime.GOMAXPROCS(0))
	erroQ    = make(chan error, 42)
	donePool = sync.Pool{New: func() interface{} {
		return make(chan error)
	}}
)

type funcData struct {
	fn   func()
	done chan error
}
