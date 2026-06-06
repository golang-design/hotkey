// Copyright 2021 The golang.design Initiative Authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.
//
// Written by Changkun Ou <changkun.de>

//go:build darwin

package hotkey

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa -framework CoreGraphics -framework ApplicationServices
#include <stdint.h>
#import <Cocoa/Cocoa.h>

extern void keydownCallback(uintptr_t handle);
extern void keyupCallback(uintptr_t handle);

// All hotkeys are delivered through a CGEventTap (regular keys and media
// keys alike), so a single mechanism handles both. registerTap returns NULL
// when the process lacks Accessibility (Input Monitoring) permission.
void* registerTap(uintptr_t handle, int isMedia, int code, uint64_t flags);
void unregisterTap(void* tap);
int isAXTrusted();
*/
import "C"
import (
	"errors"
	"runtime/cgo"
	"sync"
	"unsafe"
)

// Hotkey is a combination of modifiers and key to trigger an event.
//
// On macOS every hotkey is served by a CGEventTap, which requires the
// application to be trusted for Accessibility (Input Monitoring). Register
// returns an error when that permission is missing.
type platformHotkey struct {
	mu         sync.Mutex
	registered bool
	tap        unsafe.Pointer
	handle     cgo.Handle
}

// CGEventFlags modifier masks (see CGEventTypes.h), mapped from the package's
// Carbon-style Modifier values.
const (
	cgFlagShift   = 0x20000
	cgFlagControl = 0x40000
	cgFlagOption  = 0x80000
	cgFlagCommand = 0x100000
)

func (hk *Hotkey) register() error {
	hk.mu.Lock()
	defer hk.mu.Unlock()
	if hk.registered {
		return errAlreadyRegistered
	}

	var (
		isMedia C.int
		code    C.int
		flags   C.uint64_t
	)
	if hk.key&mediaKeyBit != 0 {
		if hk.key == KeyMediaStop {
			// There is no NX_KEYTYPE for media stop on macOS.
			return errors.New("hotkey: media stop is not available on macOS")
		}
		isMedia = 1
		code = C.int(hk.key &^ mediaKeyBit)
	} else {
		code = C.int(hk.key)
		var f C.uint64_t
		for _, m := range hk.mods {
			switch m {
			case ModCmd:
				f |= cgFlagCommand
			case ModShift:
				f |= cgFlagShift
			case ModOption:
				f |= cgFlagOption
			case ModCtrl:
				f |= cgFlagControl
			}
		}
		flags = f
	}

	h := cgo.NewHandle(hk)
	tap := C.registerTap(C.uintptr_t(h), isMedia, code, flags)
	if tap == nil {
		h.Delete()
		return errors.New("hotkey: failed to register, grant the application Accessibility (Input Monitoring) permission")
	}
	hk.tap = tap
	hk.handle = h
	hk.registered = true
	return nil
}

func (hk *Hotkey) unregister() error {
	hk.mu.Lock()
	defer hk.mu.Unlock()
	if !hk.registered {
		return errNotRegistered
	}
	C.unregisterTap(hk.tap)
	hk.tap = nil
	hk.handle.Delete()
	hk.registered = false
	return nil
}

// axTrusted reports whether the process is trusted for Accessibility (Input
// Monitoring). It is used by tests to skip when the environment cannot grant
// permission (e.g. CI runners).
func axTrusted() bool { return C.isAXTrusted() != 0 }

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
type Key uint32

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

// mediaKeyBit marks a Key as a macOS media key, encoded as an NX_KEYTYPE_*
// code (see <IOKit/hidsystem/ev_keymap.h>) rather than a Carbon virtual
// keycode. register() detects the bit and routes such keys to a CGEventTap.
const mediaKeyBit Key = 1 << 16

// Media keys. The low bits are NX_KEYTYPE_* codes. KeyMediaStop has no NX
// equivalent on macOS; registering it returns an error.
const (
	KeyMediaPlayPause Key = mediaKeyBit | 16     // NX_KEYTYPE_PLAY
	KeyMediaNext      Key = mediaKeyBit | 17     // NX_KEYTYPE_NEXT
	KeyMediaPrev      Key = mediaKeyBit | 18     // NX_KEYTYPE_PREVIOUS
	KeyMediaStop      Key = mediaKeyBit | 0xffff // unsupported on macOS
	KeyVolumeUp       Key = mediaKeyBit | 0      // NX_KEYTYPE_SOUND_UP
	KeyVolumeDown     Key = mediaKeyBit | 1      // NX_KEYTYPE_SOUND_DOWN
	KeyVolumeMute     Key = mediaKeyBit | 7      // NX_KEYTYPE_MUTE
)
