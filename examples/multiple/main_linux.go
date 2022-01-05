// Copyright 2022 The golang.design Initiative Authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"sync"

	"golang.design/x/hotkey"
	"golang.design/x/hotkey/mainthread"
)

func main() { mainthread.Init(fn) }
func fn() {
	wg := sync.WaitGroup{}
	wg.Add(2)

	hotkey.XSetErrorHandler(func(e hotkey.XErrorEvent) {
		fmt.Printf("X error: %v\n", e)
	})

	go func() {
		defer wg.Done()

		err := listenHotkey(hotkey.KeyS, hotkey.ModCtrl, hotkey.ModShift)
		if err != nil {
			log.Println(err)
		}
	}()
	go func() {
		defer wg.Done()

		err := listenHotkey(hotkey.KeyA, hotkey.ModCtrl, hotkey.ModShift)
		if err != nil {
			log.Println(err)
		}
	}()
	go func() {
		defer wg.Done()

		err := listenHotkey(hotkey.KeyA, hotkey.ModCtrl, hotkey.ModShift)
		if err != nil {
			log.Println(err)
		}
	}()
	wg.Wait()

	hotkey.XRestoreErrorHandler()

	// Will now crash with BadAccess
	wg.Add(2)
	go func() {
		defer wg.Done()

		err := listenHotkey(hotkey.KeyA, hotkey.ModCtrl, hotkey.ModShift)
		if err != nil {
			log.Println(err)
		}
	}()
	go func() {
		defer wg.Done()

		err := listenHotkey(hotkey.KeyA, hotkey.ModCtrl, hotkey.ModShift)
		if err != nil {
			log.Println(err)
		}
	}()
	wg.Wait()
}

func listenHotkey(key hotkey.Key, mods ...hotkey.Modifier) (err error) {
	ms := []hotkey.Modifier{}
	ms = append(ms, mods...)
	hk := hotkey.New(ms, key)

	err = hk.Register()
	if err != nil {
		return
	}

	// Blocks until the hokey is triggered.
	<-hk.Keydown()
	log.Printf("hotkey: %v is down\n", hk)
	<-hk.Keyup()
	log.Printf("hotkey: %v is up\n", hk)
	hk.Unregister()
	return
}
