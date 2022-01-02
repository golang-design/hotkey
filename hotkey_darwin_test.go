// Copyright 2021 The golang.design Initiative Authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.
//
// Written by Changkun Ou <changkun.de>

//go:build darwin && cgo

package hotkey_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"golang.design/x/hotkey"
	"golang.design/x/hotkey/mainthread"
)

// TestHotkey should always run success.
// This is a test to run and for manually testing the registration of multiple
// hotkeys. Registered hotkeys:
// Ctrl+Shift+S
// Ctrl+Option+S
func TestHotkey(t *testing.T) {
	tt := time.Second * 5
	done := make(chan struct{}, 2)
	ctx, cancel := context.WithTimeout(context.Background(), tt)
	go func() {
		hk := hotkey.New([]hotkey.Modifier{hotkey.ModCtrl, hotkey.ModShift}, hotkey.KeyS)

		var err error
		if mainthread.Call(func() { err = hk.Register() }); err != nil {
			t.Errorf("failed to register hotkey: %v", err)
			return
		}
		for {
			select {
			case <-ctx.Done():
				cancel()
				done <- struct{}{}
				return
			case <-hk.Keydown():
				fmt.Println("triggered ctrl+shift+s")
			}
		}
	}()

	go func() {
		hk := hotkey.New([]hotkey.Modifier{hotkey.ModCtrl, hotkey.ModOption}, hotkey.KeyS)

		var err error
		if mainthread.Call(func() { err = hk.Register() }); err != nil {
			t.Errorf("failed to register hotkey: %v", err)
			return
		}

		for {
			select {
			case <-ctx.Done():
				cancel()
				done <- struct{}{}
				return
			case <-hk.Keydown():
				fmt.Println("triggered ctrl+option+s")
			}
		}
	}()

	<-done
	<-done
}
