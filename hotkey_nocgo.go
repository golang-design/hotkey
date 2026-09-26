// Copyright 2021 The golang.design Initiative Authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.
//
// Written by Changkun Ou <changkun.de>

//go:build !windows && !cgo

package hotkey

import "errors"

type platformHotkey struct{}

// Modifier represents a modifier
type Modifier uint32

// Key represents a key.
type Key uint32

// errNoCgo is what a hotkey reports in a build without cgo, which this
// platform's implementation needs. It used to be a panic, which took down a
// program that could have run on without its hotkey (#52).
var errNoCgo = errors.New("hotkey: unavailable in a build without cgo (CGO_ENABLED=0)")

func (hk *Hotkey) register() error { return errNoCgo }

// unregister deregisteres a system hotkey.
func (hk *Hotkey) unregister() error { return errNoCgo }
