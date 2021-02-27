// Copyright 2021 The golang.design Initiative Authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.
//
// Written by Changkun Ou <changkun.de>

//go:build linux
// +build linux

package hotkey

/*
#cgo LDFLAGS: -lX11

int displayTest();
int waitHotkey(unsigned int mod, int key);
*/
import "C"
import "context"

func init() {
	if C.displayTest() != 0 {
		panic("cannot use hotkey package")
	}
}

// Nothing needs to do for register
func (hk *Hotkey) register() error { return nil }

// Nothing needs to do for unregister
func (hk *Hotkey) unregister() {}

// handle registers an application global hotkey to the system,
// and returns a channel that will signal if the hotkey is triggered.
//
// No customization for the hotkey, the hotkey is always: Ctrl+Mod4+s
func (hk *Hotkey) handle(ctx context.Context) {
	// KNOWN ISSUE: if a hotkey is grabbed by others, C side will crash the program

	var mod Modifier
	for _, m := range hk.mods {
		mod = mod | m
	}
	for {
		select {
		case <-ctx.Done():
			return
		default:
			ret := C.waitHotkey(C.uint(mod), C.int(hk.key))
			if ret != 0 {
				continue
			}
			hk.in <- Event{}
		}
	}
}

// Modifier represents a modifier.
type Modifier uint32

// All kinds of Modifiers
// See /usr/include/X11/X.h
const (
	ModCtrl  Modifier = (1 << 2)
	ModShift Modifier = (1 << 0)
	Mod1     Modifier = (1 << 3)
	Mod2     Modifier = (1 << 4)
	Mod3     Modifier = (1 << 5)
	Mod4     Modifier = (1 << 6)
	Mod5     Modifier = (1 << 7)
)

// Key represents a key.
// See /usr/include/X11/keysymdef.h
type Key uint8

// All kinds of keys
const (
	Key1 Key = 0x0030
	Key2 Key = 0x0031
	Key3 Key = 0x0032
	Key4 Key = 0x0033
	Key5 Key = 0x0034
	Key6 Key = 0x0035
	Key7 Key = 0x0036
	Key8 Key = 0x0037
	Key9 Key = 0x0038
	Key0 Key = 0x0039
	KeyA Key = 0x0061
	KeyB Key = 0x0062
	KeyC Key = 0x0063
	KeyD Key = 0x0064
	KeyE Key = 0x0065
	KeyF Key = 0x0066
	KeyG Key = 0x0067
	KeyH Key = 0x0068
	KeyI Key = 0x0069
	KeyJ Key = 0x006a
	KeyK Key = 0x006b
	KeyL Key = 0x006c
	KeyM Key = 0x006d
	KeyN Key = 0x006e
	KeyO Key = 0x006f
	KeyP Key = 0x0070
	KeyQ Key = 0x0071
	KeyR Key = 0x0072
	KeyS Key = 0x0073
	KeyT Key = 0x0074
	KeyU Key = 0x0075
	KeyV Key = 0x0076
	KeyW Key = 0x0077
	KeyX Key = 0x0078
	KeyY Key = 0x0079
	KeyZ Key = 0x007a
)
