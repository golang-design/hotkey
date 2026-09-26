// Copyright 2021 The golang.design Initiative Authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.

//go:build (linux || openbsd) && cgo

// Package x11 posts key events for the tests of the parent package. Test
// files cannot use cgo, so the X11 call lives here.
package x11

/*
#cgo linux LDFLAGS: -lX11
#cgo openbsd CFLAGS: -I/usr/X11R6/include
#cgo openbsd LDFLAGS: -L/usr/X11R6/lib -lX11
#include <X11/Xlib.h>

static int sendRootKey(int keysym, unsigned int state) {
  Display *d = XOpenDisplay(NULL);
  if (d == NULL) {
    return 0;
  }
  XKeyEvent ev;
  ev.type = KeyPress;
  ev.display = d;
  ev.root = DefaultRootWindow(d);
  ev.window = DefaultRootWindow(d);
  ev.subwindow = None;
  ev.time = CurrentTime;
  ev.x = ev.y = ev.x_root = ev.y_root = 0;
  ev.same_screen = True;
  ev.keycode = XKeysymToKeycode(d, keysym);
  ev.state = state;
  int ok = XSendEvent(d, DefaultRootWindow(d), False, KeyPressMask, (XEvent *)&ev);
  XFlush(d);
  XCloseDisplay(d);
  return ok;
}
*/
import "C"

// X11 keysyms and modifier masks (see X11/keysymdef.h and X11/X.h).
const (
	KeysymA        = 0x0061 // XK_a
	KeysymControlL = 0xffe3 // XK_Control_L
	KeysymShiftL   = 0xffe1 // XK_Shift_L
	KeysymF9       = 0xffc6 // XK_F9

	ShiftMask   = 1 << 0
	ControlMask = 1 << 2
)

// SendRootKey posts a key press to the root window, the window a hotkey's
// event loop selects KeyPressMask on. Keys the X server routes or replays to
// the root -- because no other client consumed them, or because another client
// released a keyboard grab -- reach a registered hotkey the very same way.
func SendRootKey(keysym int, state uint) bool {
	return C.sendRootKey(C.int(keysym), C.uint(state)) != 0
}
