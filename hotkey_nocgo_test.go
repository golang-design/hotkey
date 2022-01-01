// Copyright 2021 The golang.design Initiative Authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.
//
// Written by Changkun Ou <changkun.de>

//go:build (linux || darwin) && !cgo

package hotkey_test

import (
	"testing"

	"golang.design/x/hotkey"
)

// TestHotkey should always run success.
// This is a test to run and for manually testing, registered combination:
// Ctrl+Alt+A (Ctrl+Mod2+Mod4+A on Linux)
func TestHotkey(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			return
		}
		t.Fatalf("expect to fail when CGO_ENABLED=0")
	}()

	hk := hotkey.New([]hotkey.Modifier{}, hotkey.Key(0))
	err := hk.Register()
	if err != nil {
		t.Fatal(err)
	}
	hk.Unregister()
}
