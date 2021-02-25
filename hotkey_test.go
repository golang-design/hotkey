// Copyright 2021 The golang.design Initiative Authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.
//
// Written by Changkun Ou <changkun.de>

package hotkey_test

import (
	"os"
	"testing"

	"golang.design/x/mainthread"
)

// The test cannot be run twice since the mainthread loop may not be terminated:
// go test -v -count=1
func TestMain(m *testing.M) {
	mainthread.Init(func() { os.Exit(m.Run()) })
}
