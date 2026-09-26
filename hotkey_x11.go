// Copyright 2021 The golang.design Initiative Authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.
//
// Written by Changkun Ou <changkun.de>

//go:build linux || openbsd

package hotkey

/*
#cgo LDFLAGS: -lX11
#cgo openbsd CFLAGS: -I/usr/X11R6/include
#cgo openbsd LDFLAGS: -L/usr/X11R6/lib -lX11

#include <stdint.h>
#include <X11/Xlib.h>

Display *openDisplay();
Window createInvisWindow(Display *d);
void sendCancel(Display *d, Window window);
void cleanupConnection(Display *d, Window window);
int grabHotkey(Display *d, unsigned int* mods, int nmods, int key);
int hotkeyKeycode(Display *d, int key);
void waitHotkey(uintptr_t hkhandle, Display *d, int keycode, unsigned int mods);
*/
import "C"
import (
	"context"
	"errors"
	"runtime"
	"runtime/cgo"
	"sync"
)

func init() {
	// The X11 display is touched from multiple threads (register opens it and
	// grabs, the per-hotkey event loop blocks in XNextEvent, and unregister
	// sends the cancel event), so Xlib must be made thread-safe before the
	// first Xlib call.
	//
	// Whether there is a display to talk to is left to register, which
	// reports it as an error. Checking here panicked every program that
	// linked this package wherever DISPLAY was unset — on a server, over
	// SSH, under Wayland without XWayland — even on runs that never
	// registered a hotkey (#50).
	C.XInitThreads()
}

type platformHotkey struct {
	mu         sync.Mutex
	registered bool
	ctx        context.Context
	cancel     context.CancelFunc
	canceled   chan struct{}
	display    *C.Display
	window     C.Window

	// keycode and mask identify the registered combination so the event
	// loop can tell it apart from other keys.
	keycode C.int
	mask    C.uint
}

// errNoDisplay is what register reports without an X server to talk to.
var errNoDisplay = errors.New("hotkey: cannot open the X11 display; is DISPLAY set? " +
	"Under Wayland, global hotkeys need XWayland")

// grabMu serializes the grab in register across hotkeys, because the C side
// uses a process-global error slot and swaps the process-global X error
// handler while probing for a conflict.
var grabMu sync.Mutex

func (hk *Hotkey) register() error {
	hk.mu.Lock()
	defer hk.mu.Unlock()
	if hk.registered {
		return errAlreadyRegistered
	}

	var mod Modifier
	for _, m := range hk.mods {
		mod = mod | m
	}
	// Grab the hotkey once per NumLock/CapsLock state so it fires regardless
	// of those locks (see lockVariants).
	variants := lockVariants(mod)
	cmods := make([]C.uint, len(variants))
	for i, v := range variants {
		cmods[i] = C.uint(v)
	}

	display := C.openDisplay()
	if display == nil {
		return errNoDisplay
	}
	window := C.createInvisWindow(display)

	// Grab synchronously so a conflict surfaces here as an error instead of
	// crashing the program later via Xlib's default error handler.
	grabMu.Lock()
	rc := C.grabHotkey(display, &cmods[0], C.int(len(cmods)), C.int(hk.key))
	grabMu.Unlock()
	if rc != 0 {
		C.cleanupConnection(display, window)
		return errRegisterFailed
	}

	hk.display = display
	hk.window = window
	hk.keycode = C.hotkeyKeycode(display, C.int(hk.key))
	hk.mask = C.uint(mod)
	hk.registered = true
	hk.ctx, hk.cancel = context.WithCancel(context.Background())
	hk.canceled = make(chan struct{})

	go hk.handle()
	return nil
}

func (hk *Hotkey) unregister() error {
	hk.mu.Lock()
	defer hk.mu.Unlock()
	if !hk.registered {
		return errNotRegistered
	}
	hk.cancel()
	C.sendCancel(hk.display, hk.window)
	<-hk.canceled
	// Tear down the connection only after the event loop has returned, so the
	// XDestroyWindow/XCloseDisplay cannot race the sendCancel above (which
	// would destroy the window out from under the in-flight XSendEvent).
	C.cleanupConnection(hk.display, hk.window)
	hk.display = nil
	hk.registered = false
	return nil
}

// handle delivers events for an already-registered (and grabbed) hotkey until
// it is unregistered.
func (hk *Hotkey) handle() {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	h := cgo.NewHandle(hk)
	defer h.Delete()

	for {
		select {
		case <-hk.ctx.Done():
			close(hk.canceled)
			return
		default:
			C.waitHotkey(C.uintptr_t(h), hk.display, hk.keycode, hk.mask)
		}
	}
}

// X11 lock modifier masks (see /usr/include/X11/X.h).
const (
	x11LockMask Modifier = 1 << 1 // CapsLock
	x11Mod2Mask Modifier = 1 << 4 // usually NumLock
)

// lockVariants returns mod combined with every on/off combination of the
// CapsLock and NumLock masks, deduplicated. An XGrabKey uses an exact
// modifier mask, so without these variants a hotkey would stop firing
// whenever NumLock or CapsLock is toggled on. Duplicates are removed so we
// never grab the same key+mask twice (which itself raises BadAccess).
func lockVariants(mod Modifier) []uint32 {
	extra := []Modifier{0, x11LockMask, x11Mod2Mask, x11LockMask | x11Mod2Mask}
	seen := make(map[uint32]bool, len(extra))
	out := make([]uint32, 0, len(extra))
	for _, e := range extra {
		v := uint32(mod | e)
		if !seen[v] {
			seen[v] = true
			out = append(out, v)
		}
	}
	return out
}

//export hotkeyDown
func hotkeyDown(h uintptr) {
	hk := cgo.Handle(h).Value().(*Hotkey)
	hk.keydownIn <- Event{}
}

//export hotkeyUp
func hotkeyUp(h uintptr) {
	hk := cgo.Handle(h).Value().(*Hotkey)
	hk.keyupIn <- Event{}
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
type Key uint32

// All kinds of keys
const (
	KeySpace Key = 0x0020
	Key1     Key = 0x0030
	Key2     Key = 0x0031
	Key3     Key = 0x0032
	Key4     Key = 0x0033
	Key5     Key = 0x0034
	Key6     Key = 0x0035
	Key7     Key = 0x0036
	Key8     Key = 0x0037
	Key9     Key = 0x0038
	Key0     Key = 0x0039
	KeyA     Key = 0x0061
	KeyB     Key = 0x0062
	KeyC     Key = 0x0063
	KeyD     Key = 0x0064
	KeyE     Key = 0x0065
	KeyF     Key = 0x0066
	KeyG     Key = 0x0067
	KeyH     Key = 0x0068
	KeyI     Key = 0x0069
	KeyJ     Key = 0x006a
	KeyK     Key = 0x006b
	KeyL     Key = 0x006c
	KeyM     Key = 0x006d
	KeyN     Key = 0x006e
	KeyO     Key = 0x006f
	KeyP     Key = 0x0070
	KeyQ     Key = 0x0071
	KeyR     Key = 0x0072
	KeyS     Key = 0x0073
	KeyT     Key = 0x0074
	KeyU     Key = 0x0075
	KeyV     Key = 0x0076
	KeyW     Key = 0x0077
	KeyX     Key = 0x0078
	KeyY     Key = 0x0079
	KeyZ     Key = 0x007a

	KeyReturn Key = 0xff0d
	KeyEscape Key = 0xff1b
	KeyDelete Key = 0xffff
	KeyTab    Key = 0xff09

	KeyLeft  Key = 0xff51
	KeyRight Key = 0xff53
	KeyUp    Key = 0xff52
	KeyDown  Key = 0xff54

	KeyF1  Key = 0xffbe
	KeyF2  Key = 0xffbf
	KeyF3  Key = 0xffc0
	KeyF4  Key = 0xffc1
	KeyF5  Key = 0xffc2
	KeyF6  Key = 0xffc3
	KeyF7  Key = 0xffc4
	KeyF8  Key = 0xffc5
	KeyF9  Key = 0xffc6
	KeyF10 Key = 0xffc7
	KeyF11 Key = 0xffc8
	KeyF12 Key = 0xffc9
	KeyF13 Key = 0xffca
	KeyF14 Key = 0xffcb
	KeyF15 Key = 0xffcc
	KeyF16 Key = 0xffcd
	KeyF17 Key = 0xffce
	KeyF18 Key = 0xffcf
	KeyF19 Key = 0xffd0
	KeyF20 Key = 0xffd1
)

// Media keys. These are X11 XF86 keysyms (see
// /usr/include/X11/XF86keysym.h) and are grabbed through the same path as
// any other key.
const (
	KeyMediaPlayPause Key = 0x1008FF14 // XF86AudioPlay
	KeyMediaNext      Key = 0x1008FF17 // XF86AudioNext
	KeyMediaPrev      Key = 0x1008FF16 // XF86AudioPrev
	KeyMediaStop      Key = 0x1008FF15 // XF86AudioStop
	KeyVolumeUp       Key = 0x1008FF13 // XF86AudioRaiseVolume
	KeyVolumeDown     Key = 0x1008FF11 // XF86AudioLowerVolume
	KeyVolumeMute     Key = 0x1008FF12 // XF86AudioMute
)
