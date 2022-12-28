// Copyright 2022 The golang.design Initiative Authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.
package main

import (
	"log"

	"golang.design/x/hotkey"
	"golang.design/x/hotkey/mainthread"
)

func main() { mainthread.Init(fn) } // Not necessary when use in Fyne, Ebiten or Gio.
func fn() {
	hk := hotkey.New([]hotkey.Modifier{hotkey.ModCtrl, hotkey.ModShift}, hotkey.KeyS)
	err := hk.Register()
	if err != nil {
		log.Fatalf("hotkey: failed to register hotkey: %v", err)
		return
	}

	log.Printf("hotkey: %v is registered\n", hk)
	<-hk.Keydown()
	log.Printf("hotkey: %v is down\n", hk)
	<-hk.Keyup()
	log.Printf("hotkey: %v is up\n", hk)
	hk.Unregister()
	log.Printf("hotkey: %v is unregistered\n", hk)
}
