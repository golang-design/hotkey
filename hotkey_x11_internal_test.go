// Copyright 2021 The golang.design Initiative Authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.
//
// Written by Changkun Ou <changkun.de>

//go:build (linux || openbsd) && cgo

package hotkey

import (
	"reflect"
	"testing"
)

// TestLockVariants verifies that a hotkey is grabbed for every NumLock/CapsLock
// state. Before this fix only the exact modifier mask was grabbed, so a hotkey
// stopped firing whenever NumLock or CapsLock was toggled on (issue #25).
func TestLockVariants(t *testing.T) {
	const (
		lock = 1 << 1 // CapsLock
		num  = 1 << 4 // NumLock
	)
	for _, tt := range []struct {
		name string
		mod  Modifier
		want []uint32
	}{
		{
			name: "ctrl+shift grabs all four lock combinations",
			mod:  ModCtrl | ModShift, // 0b101 = 5
			want: []uint32{5, 5 | lock, 5 | num, 5 | lock | num},
		},
		{
			name: "mod already containing NumLock is deduplicated",
			mod:  ModCtrl | Mod2, // Mod2 == NumLock mask
			want: []uint32{uint32(ModCtrl | Mod2), uint32(ModCtrl|Mod2) | lock},
		},
	} {
		if got := lockVariants(tt.mod); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%s: lockVariants(%#b) = %v, want %v", tt.name, tt.mod, got, tt.want)
		}
	}
}
