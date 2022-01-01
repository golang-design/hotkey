// Copyright 2022 The golang.design Initiative Authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.
//
// Written by Changkun Ou <changkun.de>

//go:build !darwin

package mainthread

import (
	"fmt"
	"runtime"
	"sync"
)

func init() {
	runtime.LockOSThread()
}

var (
	funcQ    = make(chan funcData, runtime.GOMAXPROCS(0))
	donePool = sync.Pool{New: func() interface{} {
		return make(chan error)
	}}
)

type funcData struct {
	fn   func()
	done chan error
}

// Call runs the function on the main thread.
func Call(f func()) {
	done := donePool.Get().(chan error)
	defer donePool.Put(done)

	data := funcData{fn: f, done: done}
	funcQ <- data
	if err := <-done; err != nil {
		panic(err)
	}
}

// Run runs the main thread.
func Run(main func()) {
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
					}
				}()
				f.fn()
			}()
		case <-done:
			return
		}
	}
}
