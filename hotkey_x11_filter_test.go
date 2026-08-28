// Copyright 2021 The golang.design Initiative Authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.

//go:build (linux || openbsd) && cgo

package hotkey

import (
	"testing"
	"time"

	"golang.design/x/hotkey/internal/x11"
)

func sendRootKey(t *testing.T, keysym int, state uint) {
	t.Helper()
	if !x11.SendRootKey(keysym, state) {
		t.Fatal("cannot post a key event to the root window")
	}
}

// TestOnlyRegisteredCombinationIsReported verifies that the event loop reports
// the registered combination and nothing else. The loop selects key events on
// the root window, so keys belonging to other applications arrive there too:
// reporting them fires the hotkey on foreign shortcuts, and reporting every
// key of the combination fires it once per keystroke instead of once.
func TestOnlyRegisteredCombinationIsReported(t *testing.T) {
	hk := New([]Modifier{ModCtrl, ModShift}, KeyF9)
	if err := hk.Register(); err != nil {
		t.Skipf("cannot register the hotkey on this display: %v", err)
	}
	defer hk.Unregister()

	// A key that belongs to somebody else must not be reported.
	sendRootKey(t, x11.KeysymA, 0)
	select {
	case <-hk.Keydown():
		t.Fatal("a key that is not the registered combination was reported")
	case <-time.After(300 * time.Millisecond):
	}

	// Neither must the modifiers of the combination, which are key presses of
	// their own: only the combination itself counts, exactly once.
	sendRootKey(t, x11.KeysymControlL, 0)
	sendRootKey(t, x11.KeysymShiftL, x11.ControlMask)
	select {
	case <-hk.Keydown():
		t.Fatal("a modifier of the combination was reported on its own")
	case <-time.After(300 * time.Millisecond):
	}

	sendRootKey(t, x11.KeysymF9, x11.ControlMask|x11.ShiftMask)
	select {
	case <-hk.Keydown():
	case <-time.After(2 * time.Second):
		t.Fatal("the registered combination was not reported")
	}
	select {
	case <-hk.Keydown():
		t.Fatal("the registered combination was reported more than once")
	case <-time.After(300 * time.Millisecond):
	}
}
