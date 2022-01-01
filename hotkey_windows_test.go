// Copyright 2021 The golang.design Initiative Authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.
//
// Written by Changkun Ou <changkun.de>

//go:build windows

package hotkey_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"golang.design/x/hotkey"
)

// TestHotkey should always run success.
// This is a test to run and for manually testing, registered combination:
// Ctrl+Shift+S
func TestHotkey(t *testing.T) {
	tt := time.Second * 5

	ctx, cancel := context.WithTimeout(context.Background(), tt)
	defer cancel()

	hk := hotkey.New([]hotkey.Modifier{hotkey.ModCtrl, hotkey.ModShift}, hotkey.KeyS)
	if err := hk.Register(); err != nil {
		t.Errorf("failed to register hotkey: %v", err)
		return
	}
	for {
		select {
		case <-ctx.Done():
			return
		case <-hk.Listen():
			fmt.Println("triggered")
		}
	}
}
