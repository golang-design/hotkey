// Copyright 2021 The golang.design Initiative Authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.
//
// Written by Changkun Ou <changkun.de>

//go:build !windows && !cgo

package hotkey

type platformHotkey struct{}

// Modifier represents a modifier
type Modifier uint32

// Key represents a key.
type Key uint8

func (hk *Hotkey) register() error {
	panic("hotkey: cannot use when CGO_ENABLED=0")
}

// unregister deregisteres a system hotkey.
func (hk *Hotkey) unregister() error {
	panic("hotkey: cannot use when CGO_ENABLED=0")
}
