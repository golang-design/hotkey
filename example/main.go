// Copyright 2021 The golang.design Initiative Authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.

package main

import (
	"context"

	"golang.design/x/hotkey"
	"golang.design/x/mainthread"
)

func main() { mainthread.Init(fn) }
func fn() {
	var (
		mods = []hotkey.Modifier{hotkey.ModCtrl}
		k    = hotkey.KeyS
	)

	hk, err := hotkey.Register(mods, k)
	if err != nil {
		panic("hotkey registration failed")
	}

	triggered := hk.Listen(context.Background())
	for range triggered {
		println("hotkey ctrl+s is triggered")
	}
}
