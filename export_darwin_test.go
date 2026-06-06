// Copyright 2026 The golang.design Initiative Authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.

//go:build darwin && cgo

package hotkey

// AXTrusted exposes axTrusted to tests so they can skip when the process is
// not trusted for Accessibility (e.g. on CI runners).
var AXTrusted = axTrusted
