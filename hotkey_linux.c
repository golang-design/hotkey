// Copyright 2021 The golang.design Initiative Authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.
//
// Written by Changkun Ou <changkun.de>

//go:build linux

#include <stdio.h>
#include <X11/Xlib.h>
#include <X11/Xutil.h>

int displayTest() {
	Display* d = XOpenDisplay(0);
	if (d == NULL) {
		return -1;
	}
	XCloseDisplay(d);
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
//             printf("ERROR: key combination already grabbed by another client.\n");
//             return 0;
//         }
//     }
//     return 0;
// }

// waitHotkey blocks until the hotkey is triggered.
// this function crashes the program if the hotkey already grabbed by others.
int waitHotkey(unsigned int mod, int key) {
	// FIXME: handle registered hotkey properly.
	// XSetErrorHandler(handleErrors);

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
			XUngrabKey(d, keycode, mod, DefaultRootWindow(d));
			XCloseDisplay(d);
			return 0;
		}
	}
}