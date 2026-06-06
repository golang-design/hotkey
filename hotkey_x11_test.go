// Copyright 2021 The golang.design Initiative Authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.
//
// Written by Changkun Ou <changkun.de>

//go:build (linux || openbsd) && cgo

package hotkey_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"golang.design/x/hotkey"
)

// TestKeycodesAreDistinct guards against duplicated keysym constants.
// KeyTab used to alias KeyEscape (both 0xff1b), which made Tab hotkeys
// register Escape instead.
func TestKeycodesAreDistinct(t *testing.T) {
	keys := map[string]hotkey.Key{
		"KeySpace": hotkey.KeySpace, "KeyReturn": hotkey.KeyReturn,
		"KeyEscape": hotkey.KeyEscape, "KeyDelete": hotkey.KeyDelete,
		"KeyTab":  hotkey.KeyTab,
		"KeyLeft": hotkey.KeyLeft, "KeyRight": hotkey.KeyRight,
		"KeyUp": hotkey.KeyUp, "KeyDown": hotkey.KeyDown,
	}
	seen := map[hotkey.Key]string{}
	for name, code := range keys {
		if prev, ok := seen[code]; ok {
			t.Errorf("key %s and %s share the same keycode 0x%x", prev, name, code)
		}
		seen[code] = name
	}
}

// TestHotkey should always run success.
// This is a test to run and for manually testing, registered combination:
// Ctrl+Alt+A (Ctrl+Mod2+Mod4+A on Linux)
func TestHotkey(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	hk := hotkey.New([]hotkey.Modifier{
		hotkey.ModCtrl, hotkey.Mod2, hotkey.Mod4}, hotkey.KeyA)
	if err := hk.Register(); err != nil {
		t.Errorf("failed to register hotkey: %v", err)
		return
	}
	for {
		select {
		case <-ctx.Done():
			return
		case <-hk.Keydown():
			fmt.Println("triggered")
		}
	}
}
