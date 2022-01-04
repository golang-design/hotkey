// Copyright 2022 The golang.design Initiative Authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.

package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"golang.design/x/hotkey"
)

func main() {
	w := app.New().NewWindow("golang.design/x/hotkey")
	label := widget.NewLabel("Hello golang.design!")
	button := widget.NewButton("Hi!", func() { label.SetText("Welcome :)") })
	w.SetContent(container.NewVBox(label, button))

	go func() {
		// Register a desired hotkey.
		hk := hotkey.New([]hotkey.Modifier{hotkey.ModCtrl, hotkey.ModShift}, hotkey.KeyS)
		if err := hk.Register(); err != nil {
			panic("hotkey registration failed")
		}
		// Start listen hotkey event whenever it is ready.
		for range hk.Keydown() {
			button.Tapped(&fyne.PointEvent{})
		}
	}()

	w.ShowAndRun()
}
