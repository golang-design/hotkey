// Copyright 2022 The golang.design Initiative Authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.

package main

import (
	"fmt"

	"golang.design/x/hotkey"
	"golang.design/x/hotkey/mainthread"
)

func main() { mainthread.Init(fn) }
func fn() {
	hk := hotkey.New([]hotkey.Modifier{hotkey.ModCtrl, hotkey.ModShift}, hotkey.KeyS)
	err := hk.Register()
	if err != nil {
		return
	}
	fmt.Printf("hotkey: %v is registered\n", hk)
	<-hk.Keydown()
	fmt.Printf("hotkey: %v is down\n", hk)
	<-hk.Keyup()
	fmt.Printf("hotkey: %v is up\n", hk)
	hk.Unregister()
	fmt.Printf("hotkey: %v is unregistered\n", hk)
}
