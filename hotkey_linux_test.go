// Copyright 2021 The golang.design Initiative Authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.
//
// Written by Changkun Ou <changkun.de>

//go:build linux

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
// Ctrl+Alt+A (Ctrl+Mod2+Mod4+A on Linux)
func TestHotkey(t *testing.T) {
	tt := time.Second * 5

	ctx, cancel := context.WithTimeout(context.Background(), tt)
	defer cancel()

	hk, err := hotkey.Register([]hotkey.Modifier{
		hotkey.ModCtrl, hotkey.Mod2, hotkey.Mod4}, hotkey.KeyA)
	if err != nil {
		t.Errorf("failed to register hotkey: %v", err)
		return
	}

	trigger := hk.Listen(ctx)
	for {
		select {
		case <-ctx.Done():
			return
		case <-trigger:
			fmt.Println("triggered")
		}
	}
}
