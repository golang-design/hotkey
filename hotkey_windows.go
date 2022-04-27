// Copyright 2021 The golang.design Initiative Authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.
//
// Written by Changkun Ou <changkun.de>

//go:build windows

package hotkey

import (
	"errors"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"golang.design/x/hotkey/internal/win"
)

type platformHotkey struct {
	mu         sync.Mutex
	hotkeyId   uint64
	registered bool
	funcs      chan func()
	canceled   chan struct{}
}

var hotkeyId uint64 // atomic

// register registers a system hotkey. It returns an error if
// the registration is failed. This could be that the hotkey is
// conflict with other hotkeys.
func (hk *Hotkey) register() error {
	hk.mu.Lock()
	if hk.registered {
		hk.mu.Unlock()
		return errors.New("hotkey already registered")
	}

	mod := uint8(0)
	for _, m := range hk.mods {
		mod = mod | uint8(m)
	}

	hk.hotkeyId = atomic.AddUint64(&hotkeyId, 1)
	hk.funcs = make(chan func())
	hk.canceled = make(chan struct{})
	go hk.handle()

	var (
		ok   bool
		err  error
		done = make(chan struct{})
	)
	hk.funcs <- func() {
		ok, err = win.RegisterHotKey(0, uintptr(hk.hotkeyId), uintptr(mod), uintptr(hk.key))
		done <- struct{}{}
	}
	<-done
	if !ok {
		close(hk.canceled)
		hk.mu.Unlock()
		return err
	}
	hk.registered = true
	hk.mu.Unlock()
	return nil
}

// unregister deregisteres a system hotkey.
func (hk *Hotkey) unregister() error {
	hk.mu.Lock()
	defer hk.mu.Unlock()
	if !hk.registered {
		return errors.New("hotkey is not registered")
	}

	done := make(chan struct{})
	hk.funcs <- func() {
		win.UnregisterHotKey(0, uintptr(hk.hotkeyId))
		done <- struct{}{}
		close(hk.canceled)
	}
	<-done

	<-hk.canceled
	hk.registered = false
	return nil
}

const (
	// wmHotkey represents hotkey message
	wmHotkey uint32 = 0x0312
	wmQuit   uint32 = 0x0012
)

// handle handles the hotkey event loop.
func (hk *Hotkey) handle() {
	// We could optimize this. So far each hotkey is served in an
	// individual thread. If we have too many hotkeys, then a program
	// have to create too many threads to serve them.
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	tk := time.NewTicker(time.Second / 100)
	for range tk.C {
		msg := win.MSG{}
		if !win.PeekMessage(&msg, 0, 0, 0) {
			select {
			case f := <-hk.funcs:
				f()
			case <-hk.canceled:
				return
			default:
			}
			continue
		}
		if !win.GetMessage(&msg, 0, 0, 0) {
			return
		}

		switch msg.Message {
		case wmHotkey:
			hk.keydownIn <- Event{}

			tk := time.NewTicker(time.Second / 100)
			for range tk.C {
				if win.GetAsyncKeyState(int(hk.key)) == 0 {
					hk.keyupIn <- Event{}
					break
				}
			}
		case wmQuit:
			return
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
type Key uint16

// All kinds of Keys
const (
	KeySpace Key = 0x20
	Key0     Key = 0x30
	Key1     Key = 0x31
	Key2     Key = 0x32
	Key3     Key = 0x33
	Key4     Key = 0x34
	Key5     Key = 0x35
	Key6     Key = 0x36
	Key7     Key = 0x37
	Key8     Key = 0x38
	Key9     Key = 0x39
	KeyA     Key = 0x41
	KeyB     Key = 0x42
	KeyC     Key = 0x43
	KeyD     Key = 0x44
	KeyE     Key = 0x45
	KeyF     Key = 0x46
	KeyG     Key = 0x47
	KeyH     Key = 0x48
	KeyI     Key = 0x49
	KeyJ     Key = 0x4A
	KeyK     Key = 0x4B
	KeyL     Key = 0x4C
	KeyM     Key = 0x4D
	KeyN     Key = 0x4E
	KeyO     Key = 0x4F
	KeyP     Key = 0x50
	KeyQ     Key = 0x51
	KeyR     Key = 0x52
	KeyS     Key = 0x53
	KeyT     Key = 0x54
	KeyU     Key = 0x55
	KeyV     Key = 0x56
	KeyW     Key = 0x57
	KeyX     Key = 0x58
	KeyY     Key = 0x59
	KeyZ     Key = 0x5A

	KeyReturn Key = 0x0D
	KeyEscape Key = 0x1B
	KeyDelete Key = 0x2E
	KeyTab    Key = 0x09

	KeyLeft  Key = 0x25
	KeyRight Key = 0x27
	KeyUp    Key = 0x26
	KeyDown  Key = 0x28

	KeyF1  Key = 0x70
	KeyF2  Key = 0x71
	KeyF3  Key = 0x72
	KeyF4  Key = 0x73
	KeyF5  Key = 0x74
	KeyF6  Key = 0x75
	KeyF7  Key = 0x76
	KeyF8  Key = 0x77
	KeyF9  Key = 0x78
	KeyF10 Key = 0x79
	KeyF11 Key = 0x7A
	KeyF12 Key = 0x7B
	KeyF13 Key = 0x7C
	KeyF14 Key = 0x7D
	KeyF15 Key = 0x7E
	KeyF16 Key = 0x7F
	KeyF17 Key = 0x80
	KeyF18 Key = 0x81
	KeyF19 Key = 0x82
	KeyF20 Key = 0x83
)
