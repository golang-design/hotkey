// Copyright 2021 The golang.design Initiative Authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.
//
// Written by Changkun Ou <changkun.de>

//go:build linux || openbsd

#include <X11/Xlib.h>
#include <X11/Xutil.h>
#include <stdint.h>
#include <string.h> // memset

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

// lastGrabError records the X error code raised during the most recent
// grabHotkey call (0 if none). grabHotkey is serialized on the Go side, so a
// single global slot is sufficient.
static int lastGrabError = 0;

int grabErrorHandler(Display *d, XErrorEvent *e) {
  lastGrabError = e->error_code;
  return 0;
}

// grabHotkey grabs key on display d once per modifier mask in mods (the
// NumLock/CapsLock variants). It installs a temporary error handler and
// XSyncs so that a BadAccess -- the combination is already grabbed by another
// client -- is reported synchronously instead of terminating the program via
// Xlib's default handler. Returns 0 on success, 1 if the combination is
// unavailable. Callers must serialize this (it uses a process-global error
// slot and swaps the process-global X error handler).
int grabHotkey(Display *d, unsigned int *mods, int nmods, int key) {
  lastGrabError = 0;
  int keycode = XKeysymToKeycode(d, key);
  XErrorHandler old = XSetErrorHandler(grabErrorHandler);
  for (int i = 0; i < nmods; i++) {
    XGrabKey(d, keycode, mods[i], DefaultRootWindow(d), False, GrabModeAsync,
             GrabModeAsync);
  }
  XSync(d, False); // force the server to deliver any grab error to our handler
  XSetErrorHandler(old);
  if (lastGrabError == BadAccess) {
    for (int i = 0; i < nmods; i++) {
      XUngrabKey(d, keycode, mods[i], DefaultRootWindow(d));
    }
    return 1;
  }
  XSelectInput(d, DefaultRootWindow(d), KeyPressMask);
  return 0;
}

// waitHotkey delivers key events on display d until a cancel ClientMessage
// (see sendCancel) breaks the loop out of XNextEvent so an unregister can take
// effect without waiting for the next keypress. The grab is established once
// by grabHotkey and held until cleanupConnection.
void waitHotkey(uintptr_t hkhandle, Display *d) {
  XEvent ev;
  while (1) {
    XNextEvent(d, &ev);
    switch (ev.type) {
    case KeyPress:
      hotkeyDown(hkhandle);
      continue;
    case KeyRelease:
      hotkeyUp(hkhandle);
      continue;
    case ClientMessage:
      return;
    }
  }
}
