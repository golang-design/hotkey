// Copyright 2022 The golang.design Initiative Authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.
//
// Written by Changkun Ou <changkun.de>

//go:build windows || linux || openbsd || (darwin && !cgo)

package mainthread

import "golang.design/x/mainthread"

// Call calls f on the main thread and blocks until f finishes.
func Call(f func()) { mainthread.Call(f) }

// Init initializes the functionality of running arbitrary subsequent functions be called on the main system thread.
//
// Init must be called in the main.main function.
func Init(main func()) { mainthread.Init(main) }
