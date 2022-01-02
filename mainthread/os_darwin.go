// Copyright 2022 The golang.design Initiative Authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.
//
// Written by Changkun Ou <changkun.de>

//go:build darwin

package mainthread

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa
#import <Cocoa/Cocoa.h>
#import <Dispatch/Dispatch.h>

extern void os_main(void);
extern void wakeupMainThread(void);
static bool isMainThread() {
	return [NSThread isMainThread];
}
*/
import "C"
import (
	"os"
	"runtime"
)

func init() {
	runtime.LockOSThread()
}

// Call calls f on the main thread and blocks until f finishes.
func Call(f func()) {
	if C.isMainThread() {
		f()
		return
	}
	go func() {
		mainFuncs <- f
		C.wakeupMainThread()
	}()
}

// Init initializes the functionality of running arbitrary subsequent functions be called on the main system thread.
//
// Init must be called in the main.main function.
func Init(f func()) {
	go func() {
		f()
		os.Exit(0)
	}()

	C.os_main()
}

var mainFuncs = make(chan func(), 1)

//export dispatchMainFuncs
func dispatchMainFuncs() {
	for {
		select {
		case f := <-mainFuncs:
			f()
		default:
			return
		}
	}
}
