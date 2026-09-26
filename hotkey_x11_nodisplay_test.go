// Copyright 2026 The golang.design Initiative Authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.

//go:build (linux || openbsd) && cgo

package hotkey_test

import (
	"os"
	"os/exec"
	"strings"
	"testing"

	"golang.design/x/hotkey"
)

// TestNoDisplay runs the package where no X server can be reached. Its init
// used to panic there, which crashed every program that linked the package
// on a server, over SSH, or under Wayland without XWayland, even if it never
// registered a hotkey (#50). The missing display must instead be an error
// from Register.
//
// It re-runs the test binary without DISPLAY, because by the time a test
// runs in this process, init has already been through the display it has.
func TestNoDisplay(t *testing.T) {
	if os.Getenv("HOTKEY_TEST_NO_DISPLAY") == "1" {
		hk := hotkey.New([]hotkey.Modifier{hotkey.ModCtrl}, hotkey.KeyS)
		if err := hk.Register(); err == nil {
			hk.Unregister()
			t.Fatal("Register succeeded without a display")
		} else if !strings.Contains(err.Error(), "X11 display") {
			t.Fatalf("Register without a display: %v, want the reason", err)
		}
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=^TestNoDisplay$", "-test.v")
	for _, kv := range os.Environ() {
		if !strings.HasPrefix(kv, "DISPLAY=") {
			cmd.Env = append(cmd.Env, kv)
		}
	}
	cmd.Env = append(cmd.Env, "HOTKEY_TEST_NO_DISPLAY=1")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("without a display: %v\n%s", err, out)
	}
	if !strings.Contains(string(out), "--- PASS: TestNoDisplay") {
		t.Fatalf("the test did not run without a display:\n%s", out)
	}
}
