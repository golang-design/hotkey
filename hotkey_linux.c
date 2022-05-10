// Copyright 2021 The golang.design Initiative Authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.
//
// Written by Changkun Ou <changkun.de>

// go:build linux

#include <X11/Xlib.h>
#include <X11/Xutil.h>
#include <stdint.h>
// Needed for memset()
#include <string.h>

extern void hotkeyDown(uintptr_t hkhandle);
extern void hotkeyUp(uintptr_t hkhandle);

int displayTest() {
  Display *d = NULL;
  for (int i = 0; i < 42; i++) {
    d = XOpenDisplay(0);
    if (d == NULL)
      continue;
    break;
  }
  if (d == NULL) {
    return -1;
  }
  return 0;
}

// FIXME: handle bad access properly.
// int handleErrors( Display* dpy, XErrorEvent* pErr )
// {
//     printf("X Error Handler called, values: %d/%lu/%d/%d/%d\n",
//         pErr->type,
//         pErr->serial,
//         pErr->error_code,
//         pErr->request_code,
//         pErr->minor_code );
//     if( pErr->request_code == 33 ){  // 33 (X_GrabKey)
//         if( pErr->error_code == BadAccess ){
//             printf("ERROR: key combination already grabbed by another
//             client.\n"); return 0;
//         }
//     }
//     return 0;
// }

Display *openDisplay() {
  Display *d = NULL;
  for (int i = 0; i < 42; i++) {
    d = XOpenDisplay(0);
    if (d == NULL)
      continue;
    break;
  }
  return d;
}

// Creates an invisible window, which can receive ClientMessage events. On
// hotkey cancel a ClientMessageEvent is generated on the window. The event is
// catched and the event loop terminates. x: 0 y: 0 w: 1 h: 1 border_width: 1
// depth: 0
// class: InputOnly (window will not be drawn)
// visual: default visual of display
// no attributes will be set (0, &attr)
Window createInvisWindow(Display *d) {
  XSetWindowAttributes attr;
  return XCreateWindow(d, DefaultRootWindow(d), 0, 0, 1, 1, 0, 0, InputOnly,
                       DefaultVisual(d, 0), 0, &attr);
}

// Sends a custom ClientMessage of type (Atom) "go_hotkey_cancel_hotkey"
// Passed value 'True' of XInternAtom creates the Atom, if it does not exist yet
void sendCancel(Display *d, Window window) {
  Atom atom = XInternAtom(d, "golangdesign_hotkey_cancel_hotkey", True);
  XClientMessageEvent clientEvent;
  memset(&clientEvent, 0, sizeof(clientEvent));
  clientEvent.type = ClientMessage;
  clientEvent.send_event = True;
  clientEvent.display = d;
  clientEvent.window = window;
  clientEvent.message_type = atom;
  clientEvent.format = 8;

  XEvent event;
  event.type = ClientMessage;
  event.xclient = clientEvent;
  XSendEvent(d, window, False, 0, &event);
  XFlush(d);
}

// Closes the connection and destroys the invisible 'cancel' window
void cleanupConnection(Display *d, Window w) {
  XDestroyWindow(d, w);
  XCloseDisplay(d);
}

// waitHotkey blocks until the hotkey is triggered.
// this function crashes the program if the hotkey already grabbed by others.
int waitHotkey(uintptr_t hkhandle, unsigned int mod, int key, Display *d,
               Window w) {
  int keycode = XKeysymToKeycode(d, key);
  XGrabKey(d, keycode, mod, DefaultRootWindow(d), False, GrabModeAsync,
           GrabModeAsync);
  XSelectInput(d, DefaultRootWindow(d), KeyPressMask);
  XEvent ev;
  while (1) {
    XNextEvent(d, &ev);
    switch (ev.type) {
    case KeyPress:
      hotkeyDown(hkhandle);
      continue;
    case KeyRelease:
      hotkeyUp(hkhandle);
      XUngrabKey(d, keycode, mod, DefaultRootWindow(d));
      return 0;
    case ClientMessage:
      return 0;
    }
  }
}
