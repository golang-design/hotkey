// Copyright 2021 The golang.design Initiative Authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.
//
// Written by Changkun Ou <changkun.de>

//go:build darwin

package hotkey

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa -framework Carbon
#include <stdint.h>
#import <Cocoa/Cocoa.h>
#import <Carbon/Carbon.h>

extern void hotkeyCallback(uintptr_t handle);

int registerHotKey(int mod, int key, uintptr_t handle);
void runApp();
void stopApp();
*/
import "C"
import (
	"context"
	"errors"
	"runtime/cgo"
)

// handle handles the hotkey event loop.
func (hk *Hotkey) handle(ctx context.Context) {
	// Note: This call never returns.
	C.runApp()
}

func (hk *Hotkey) register() error {
	// Note: we use handle number as hotkey id in the C side.
	// A cgo handle could ran out of space, but since in hotkey purpose
	// we won't have that much number of hotkeys. So this should be fine.

	h := cgo.NewHandle(hk.in)
	var mod Modifier
	for _, m := range hk.mods {
		mod += m
	}

	ret := C.registerHotKey(C.int(mod), C.int(hk.key), C.uintptr_t(h))
	if ret == C.int(-1) {
		return errors.New("register failed")
	}
	return nil
}

func (hk *Hotkey) unregister() {
	// TODO: unregister registered hotkeys.

	// KNOWN ISSUE: This call seems does not terminate the app.
	C.stopApp()
}

//export hotkeyCallback
func hotkeyCallback(h uintptr) {
	ch := cgo.Handle(h).Value().(chan<- Event)
	ch <- Event{}
}

// Modifier represents a modifier.
// See: /Library/Developer/CommandLineTools/SDKs/MacOSX.sdk/System/Library/Frameworks/Carbon.framework/Versions/A/Frameworks/HIToolbox.framework/Versions/A/Headers/Events.h
type Modifier uint32

// All kinds of Modifiers
const (
	ModCtrl   Modifier = 0x1000
	ModShift  Modifier = 0x200
	ModOption Modifier = 0x800
	ModCmd    Modifier = 0x100
)

// Key represents a key.
// See: /Library/Developer/CommandLineTools/SDKs/MacOSX.sdk/System/Library/Frameworks/Carbon.framework/Versions/A/Frameworks/HIToolbox.framework/Versions/A/Headers/Events.h
type Key uint8

// All kinds of keys
const (
	Key1 Key = 18
	Key2 Key = 19
	Key3 Key = 20
	Key4 Key = 21
	Key5 Key = 23
	Key6 Key = 22
	Key7 Key = 26
	Key8 Key = 28
	Key9 Key = 25
	Key0 Key = 29
	KeyA Key = 0
	KeyB Key = 11
	KeyC Key = 8
	KeyD Key = 2
	KeyE Key = 14
	KeyF Key = 3
	KeyG Key = 5
	KeyH Key = 4
	KeyI Key = 34
	KeyJ Key = 38
	KeyK Key = 40
	KeyL Key = 37
	KeyM Key = 46
	KeyN Key = 45
	KeyO Key = 31
	KeyP Key = 35
	KeyQ Key = 12
	KeyR Key = 15
	KeyS Key = 1
	KeyT Key = 17
	KeyU Key = 32
	KeyV Key = 9
	KeyW Key = 13
	KeyX Key = 7
	KeyY Key = 16
	KeyZ Key = 6
)
