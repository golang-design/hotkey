// Copyright 2021 The golang.design Initiative Authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.
//
// Written by Changkun Ou <changkun.de>

//go:build linux

#include <stdint.h>
#include <stdio.h>
#include <X11/Xlib.h>
#include <X11/Xutil.h>

extern void hotkeyDown(uintptr_t hkhandle);
extern void hotkeyUp(uintptr_t hkhandle);
extern void onError(int typ, unsigned long serial, unsigned char error_code, unsigned char request_code, unsigned char minor_code);

int displayTest() {
	Display* d = NULL;
	for (int i = 0; i < 42; i++) {
		d = XOpenDisplay(0);
		if (d == NULL) continue;
		break;
	}
	if (d == NULL) {
		return -1;
	}
	return 0;
}

static int handleErrors(Display* dpy, XErrorEvent* pErr)
{
	onError(pErr->type, pErr->serial, pErr->error_code, pErr->request_code, pErr->minor_code );
	return 0;
}

static int (*defaultErrHandler)(Display*, XErrorEvent*);

int setErrorHandler() {
	if (!defaultErrHandler) {
		defaultErrHandler = XSetErrorHandler(handleErrors);
	}
}

int restoreErrorHandler() {
	XSetErrorHandler(defaultErrHandler);
}

// waitHotkey blocks until the hotkey is triggered.
// this function crashes the program if the hotkey already grabbed by others.
int waitHotkey(uintptr_t hkhandle, unsigned int mod, int key) {
	Display* d = NULL;
	for (int i = 0; i < 42; i++) {
		d = XOpenDisplay(0);
		if (d == NULL) continue;
		break;
	}
	if (d == NULL) {
		return -1;
	}
	int keycode = XKeysymToKeycode(d, key);
	XGrabKey(d, keycode, mod, DefaultRootWindow(d), False, GrabModeAsync, GrabModeAsync);
	XSelectInput(d, DefaultRootWindow(d), KeyPressMask);
	XEvent ev;
	while(1) {
		XNextEvent(d, &ev);
		switch(ev.type) {
		case KeyPress:
			hotkeyDown(hkhandle);
			continue;
		case KeyRelease:
			hotkeyUp(hkhandle);
			XUngrabKey(d, keycode, mod, DefaultRootWindow(d));
			XCloseDisplay(d);
			return 0;
		}
	}
}
