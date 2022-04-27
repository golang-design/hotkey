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

extern void keydownCallback(uintptr_t handle);
extern void keyupCallback(uintptr_t handle);
int registerHotKey(int mod, int key, uintptr_t handle, EventHotKeyRef* ref);
int unregisterHotKey(EventHotKeyRef ref);
*/
import "C"
import (
	"errors"
	"runtime/cgo"
	"sync"
)

// Hotkey is a combination of modifiers and key to trigger an event
type platformHotkey struct {
	mu         sync.Mutex
	registered bool
	hkref      C.EventHotKeyRef
}

func (hk *Hotkey) register() error {
	hk.mu.Lock()
	defer hk.mu.Unlock()
	if hk.registered {
		return errors.New("hotkey already registered")
	}

	// Note: we use handle number as hotkey id in the C side.
	// A cgo handle could ran out of space, but since in hotkey purpose
	// we won't have that much number of hotkeys. So this should be fine.

	h := cgo.NewHandle(hk)
	var mod Modifier
	for _, m := range hk.mods {
		mod += m
	}

	ret := C.registerHotKey(C.int(mod), C.int(hk.key), C.uintptr_t(h), &hk.hkref)
	if ret == C.int(-1) {
		return errors.New("failed to register the hotkey")
	}

	hk.registered = true
	return nil
}

func (hk *Hotkey) unregister() error {
	hk.mu.Lock()
	defer hk.mu.Unlock()
	if !hk.registered {
		return errors.New("hotkey is not registered")
	}

	ret := C.unregisterHotKey(hk.hkref)
	if ret == C.int(-1) {
		return errors.New("failed to unregister the current hotkey")
	}
	hk.registered = false
	return nil
}

//export keydownCallback
func keydownCallback(h uintptr) {
	hk := cgo.Handle(h).Value().(*Hotkey)
	hk.keydownIn <- Event{}
}

//export keyupCallback
func keyupCallback(h uintptr) {
	hk := cgo.Handle(h).Value().(*Hotkey)
	hk.keyupIn <- Event{}
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
	KeySpace Key = 49
	Key1     Key = 18
	Key2     Key = 19
	Key3     Key = 20
	Key4     Key = 21
	Key5     Key = 23
	Key6     Key = 22
	Key7     Key = 26
	Key8     Key = 28
	Key9     Key = 25
	Key0     Key = 29
	KeyA     Key = 0
	KeyB     Key = 11
	KeyC     Key = 8
	KeyD     Key = 2
	KeyE     Key = 14
	KeyF     Key = 3
	KeyG     Key = 5
	KeyH     Key = 4
	KeyI     Key = 34
	KeyJ     Key = 38
	KeyK     Key = 40
	KeyL     Key = 37
	KeyM     Key = 46
	KeyN     Key = 45
	KeyO     Key = 31
	KeyP     Key = 35
	KeyQ     Key = 12
	KeyR     Key = 15
	KeyS     Key = 1
	KeyT     Key = 17
	KeyU     Key = 32
	KeyV     Key = 9
	KeyW     Key = 13
	KeyX     Key = 7
	KeyY     Key = 16
	KeyZ     Key = 6

	KeyReturn Key = 0x24
	KeyEscape Key = 0x35
	KeyDelete Key = 0x33
	KeyTab    Key = 0x30

	KeyLeft  Key = 0x7B
	KeyRight Key = 0x7C
	KeyUp    Key = 0x7E
	KeyDown  Key = 0x7D

	KeyF1  Key = 0x7A
	KeyF2  Key = 0x78
	KeyF3  Key = 0x63
	KeyF4  Key = 0x76
	KeyF5  Key = 0x60
	KeyF6  Key = 0x61
	KeyF7  Key = 0x62
	KeyF8  Key = 0x64
	KeyF9  Key = 0x65
	KeyF10 Key = 0x6D
	KeyF11 Key = 0x67
	KeyF12 Key = 0x6F
	KeyF13 Key = 0x69
	KeyF14 Key = 0x6B
	KeyF15 Key = 0x71
	KeyF16 Key = 0x6A
	KeyF17 Key = 0x40
	KeyF18 Key = 0x4F
	KeyF19 Key = 0x50
	KeyF20 Key = 0x5A
)
