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

// TestOnMainFromTheMainThread: Register and Unregister do their work on the
// main thread, and asked from the main thread itself they used to wait for
// it forever, which libdispatch traps: the process crashed. Unlike the
// hotkey tests, this needs no Accessibility permission.
func TestOnMainFromTheMainThread(t *testing.T) {
	for _, tt := range []struct {
		name string
		call func(func())
	}{
		{"from the main thread", mainthread.Call}, // Call does not wait
		{"from another goroutine", func(f func()) { go f() }},
	} {
		t.Run(tt.name, func(t *testing.T) {
			done := make(chan bool, 1)
			tt.call(func() { done <- hotkey.OnMainProbe() })
			select {
			case ran := <-done:
				if !ran {
					t.Fatal("the work did not run")
				}
			case <-time.After(5 * time.Second):
				t.Fatal("the work never ran")
			}
		})
	}
}

// TestHotkey should always run success.
// This is a test to run and for manually testing the registration of multiple
// hotkeys. Registered hotkeys:
// Ctrl+Shift+S
// Ctrl+Option+S
func TestHotkey(t *testing.T) {
	if !hotkey.AXTrusted() {
		t.Skip("skipping: process is not trusted for Accessibility (Input Monitoring); grant permission to run this test")
	}
	tt := time.Second * 5
	done := make(chan struct{}, 2)
	ctx, cancel := context.WithTimeout(context.Background(), tt)
	go func() {
		hk := hotkey.New([]hotkey.Modifier{hotkey.ModCtrl, hotkey.ModShift}, hotkey.KeyS)
		if err := hk.Register(); err != nil {
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
		if err := hk.Register(); err != nil {
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
