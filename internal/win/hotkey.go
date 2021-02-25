// Copyright 2021 The golang.design Initiative Authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.
//
// Written by Changkun Ou <changkun.de>

//go:build windows
// +build windows

package win

import (
	"syscall"
	"unsafe"
)

var (
	user32           = syscall.NewLazyDLL("user32")
	registerHotkey   = user32.NewProc("RegisterHotKey")
	unregisterHotkey = user32.NewProc("UnregisterHotKey")
	getMessage       = user32.NewProc("GetMessageW")
	sendMessage      = user32.NewProc("SendMessageW")
)

// RegisterHotKey defines a system-wide hot key.
// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-registerhotkey
func RegisterHotKey(hwnd uintptr, id int, mod uintptr, k uintptr) (bool, error) {
	ret, _, err := registerHotkey.Call(
		hwnd, uintptr(id), mod, k,
	)
	return ret != 0, err
}

// UnregisterHotKey frees a hot key previously registered by the calling
// thread.
// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-unregisterhotkey
func UnregisterHotKey(hwnd uintptr, id int) (bool, error) {
	ret, _, err := unregisterHotkey.Call(hwnd, uintptr(id))
	return ret != 0, err
}

// MSG contains message information from a thread's message queue.
//
// https://docs.microsoft.com/en-us/windows/win32/api/winuser/ns-winuser-msg
type MSG struct {
	HWnd    uintptr
	Message uint32
	WParam  uintptr
	LParam  uintptr
	Time    uint32
	Pt      struct { //POINT
		x, y int32
	}
}

// SendMessage sends the specified message to a window or windows.
// The SendMessage function calls the window procedure for the specified
// window and does not return until the window procedure has processed
// the message.
// The return value specifies the result of the message processing;
// it depends on the message sent.
//
// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-sendmessage
func SendMessage(hwnd uintptr, msg uint32, wParam, lParam uintptr) uintptr {
	ret, _, _ := sendMessage.Call(
		hwnd,
		uintptr(msg),
		wParam,
		lParam,
	)

	return ret
}

// GetMessage retrieves a message from the calling thread's message
// queue. The function dispatches incoming sent messages until a posted
// message is available for retrieval.
//
// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-getmessage
func GetMessage(msg *MSG, hWnd uintptr, msgFilterMin, msgFilterMax uint32) bool {
	ret, _, _ := getMessage.Call(
		uintptr(unsafe.Pointer(msg)),
		hWnd,
		uintptr(msgFilterMin),
		uintptr(msgFilterMax),
	)

	return ret != 0
}
