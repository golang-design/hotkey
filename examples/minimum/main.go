// Copyright 2021 The golang.design Initiative Authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.

package main

import (
	"golang.design/x/hotkey"
	"golang.design/x/hotkey/mainthread"
)

func main() { mainthread.Init(fn) }
func fn() {
	hk := hotkey.New([]hotkey.Modifier{hotkey.ModCtrl}, hotkey.KeyS)
	var err error
	if mainthread.Call(func() { err = hk.Register() }); err != nil {
		panic(err)
	}

	for range hk.Keydown() {
		println("hotkey ctrl+s is triggered, exit the application, good bye!")
		hk.Unregister()
	}
}
