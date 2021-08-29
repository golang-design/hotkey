// Copyright 2021 The golang.design Initiative Authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.
//
// Written by Changkun Ou <changkun.de>

//go:build darwin

package hotkey_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"golang.design/x/hotkey"
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
		hk, err := hotkey.Register([]hotkey.Modifier{hotkey.ModCtrl, hotkey.ModShift}, hotkey.KeyS)
		if err != nil {
			t.Errorf("failed to register hotkey: %v", err)
			return
		}

		trigger := hk.Listen(ctx)
		for {
			select {
			case <-ctx.Done():
				cancel()
				done <- struct{}{}
				return
			case <-trigger:
				fmt.Println("triggered 1")
			}
		}
	}()

	go func() {
		hk, err := hotkey.Register([]hotkey.Modifier{hotkey.ModCtrl, hotkey.ModOption}, hotkey.KeyS)
		if err != nil {
			t.Errorf("failed to register hotkey: %v", err)
			return
		}

		trigger := hk.Listen(ctx)
		for {
			select {
			case <-ctx.Done():
				cancel()
				done <- struct{}{}
				return
			case <-trigger:
				fmt.Println("triggered 2")
			}
		}
	}()

	<-done
	<-done
}
