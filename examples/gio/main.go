// Copyright 2022 The golang.design Initiative Authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.

package main

import (
	"os"

	"gioui.org/app"
	"gioui.org/unit"
	"golang.design/x/hotkey"
)

func main() {
	go fn()
	app.Main()
}

func fn() {
	w := app.NewWindow(app.Title("golang.design/x/hotkey"), app.Size(unit.Dp(200), unit.Dp(200)))

	go reghk()

	for range w.Events() {
	}
}

func reghk() {
	// Register a desired hotkey.
	hk := hotkey.New([]hotkey.Modifier{hotkey.ModCtrl, hotkey.ModShift}, hotkey.KeyS)
	if err := hk.Register(); err != nil {
		panic("hotkey registration failed")
	}

	// Unregister the hotkey when keydown event is triggered
	for range hk.Keydown() {
		hk.Unregister()
	}

	// Exist the app.
	os.Exit(0)
}
