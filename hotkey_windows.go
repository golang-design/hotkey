// Copyright 2021 The golang.design Initiative Authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.
//
// Written by Changkun Ou <changkun.de>

//go:build windows

package hotkey

import (
	"context"

	"golang.design/x/hotkey/internal/win"
)

// register registers a system hotkey. It returns an error if
// the registration is failed. This could be that the hotkey is
// conflict with other hotkeys.
func (hk *Hotkey) register() error {
	mod := uint8(0)
	for _, m := range hk.mods {
		mod = mod | uint8(m)
	}
	ok, err := win.RegisterHotKey(0, 1, uintptr(mod), uintptr(hk.key))
	if !ok {
		return err
	}
	return nil
}

// unregister deregisteres a system hotkey.
func (hk *Hotkey) unregister() {
	// FIXME: unregister hotkey
}

// msgHotkey represents hotkey message
const msgHotkey uint32 = 0x0312

// handle handles the hotkey event loop.
func (hk *Hotkey) handle(ctx context.Context) {
	msg := win.MSG{}
	// KNOWN ISSUE: This message loop need to press an additional
	// hotkey to eventually return if ctx.Done() is ready.
	for win.GetMessage(&msg, 0, 0, 0) {
		switch msg.Message {
		case msgHotkey:
			select {
			case <-ctx.Done():
				return
			default:
				hk.in <- Event{}
			}
		}
	}
}

// Modifier represents a modifier.
// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-registerhotkey
type Modifier uint8

// All kinds of Modifiers
const (
	ModAlt   Modifier = 0x1
	ModCtrl  Modifier = 0x2
	ModShift Modifier = 0x4
	ModWin   Modifier = 0x8
)

// Key represents a key.
// https://docs.microsoft.com/en-us/windows/win32/inputdev/virtual-key-codes
type Key uint8

// All kinds of Keys
const (
	Key0 Key = 0x30
	Key1 Key = 0x31
	Key2 Key = 0x32
	Key3 Key = 0x33
	Key4 Key = 0x34
	Key5 Key = 0x35
	Key6 Key = 0x36
	Key7 Key = 0x37
	Key8 Key = 0x38
	Key9 Key = 0x39
	KeyA Key = 0x41
	KeyB Key = 0x42
	KeyC Key = 0x43
	KeyD Key = 0x44
	KeyE Key = 0x45
	KeyF Key = 0x46
	KeyG Key = 0x47
	KeyH Key = 0x48
	KeyI Key = 0x49
	KeyJ Key = 0x4A
	KeyK Key = 0x4B
	KeyL Key = 0x4C
	KeyM Key = 0x4D
	KeyN Key = 0x4E
	KeyO Key = 0x4F
	KeyP Key = 0x50
	KeyQ Key = 0x51
	KeyR Key = 0x52
	KeyS Key = 0x53
	KeyT Key = 0x54
	KeyU Key = 0x55
	KeyV Key = 0x56
	KeyW Key = 0x57
	KeyX Key = 0x58
	KeyY Key = 0x59
	KeyZ Key = 0x5A
)
