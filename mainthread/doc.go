// Copyright 2022 The golang.design Initiative Authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.
//
// Written by Changkun Ou <changkun.de>

// Package mainthread wrapps the golang.design/x/mainthread, and
// provides a different implementation for macOS so that it can
// handle main thread events for the NSApplication.
package mainthread
