// Copyright 2021 The golang.design Initiative Authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.
//
// Written by Changkun Ou <changkun.de>

//go:build (linux || darwin) && !cgo

package hotkey_test

import (
	"strings"
	"testing"

	"golang.design/x/hotkey"
)

// TestHotkey: without cgo there is no hotkey to register, and Register says
// so with an error rather than panicking, which crashed a program that could
// have run on without its hotkey (#52).
func TestHotkey(t *testing.T) {
	hk := hotkey.New([]hotkey.Modifier{}, hotkey.Key(0))
	err := hk.Register()
	if err == nil || !strings.Contains(err.Error(), "cgo") {
		t.Fatalf("Register without cgo: %v, want an error naming cgo", err)
	}
	if err := hk.Unregister(); err == nil {
		t.Fatal("Unregister without cgo succeeded")
	}
}
