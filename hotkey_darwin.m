// Copyright 2021 The golang.design Initiative Authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.
//
// Written by Changkun Ou <changkun.de>

//go:build darwin

#include <stdint.h>
#import <Cocoa/Cocoa.h>
#import <Carbon/Carbon.h>
extern void keydownCallback(uintptr_t handle);
extern void keyupCallback(uintptr_t handle);

static OSStatus
keydownHandler(EventHandlerCallRef nextHandler, EventRef theEvent, void *userData) {
	EventHotKeyID k;
	GetEventParameter(theEvent, kEventParamDirectObject, typeEventHotKeyID, NULL, sizeof(k), NULL, &k);
	keydownCallback((uintptr_t)k.id); // use id as handle
	return noErr;
}

static OSStatus
keyupHandler(EventHandlerCallRef nextHandler, EventRef theEvent, void *userData) {
	EventHotKeyID k;
	GetEventParameter(theEvent, kEventParamDirectObject, typeEventHotKeyID, NULL, sizeof(k), NULL, &k);
	keyupCallback((uintptr_t)k.id); // use id as handle
	return noErr;
}

// registerHotkeyWithCallback registers a global system hotkey for callbacks.
int registerHotKey(int mod, int key, uintptr_t handle, EventHotKeyRef* ref) {
	__block OSStatus s;
	dispatch_sync(dispatch_get_main_queue(), ^{
		EventTypeSpec keydownEvent;
		keydownEvent.eventClass = kEventClassKeyboard;
		keydownEvent.eventKind = kEventHotKeyPressed;
		EventTypeSpec keyupEvent;
		keyupEvent.eventClass = kEventClassKeyboard;
		keyupEvent.eventKind = kEventHotKeyReleased;
		InstallApplicationEventHandler(
			&keydownHandler, 1, &keydownEvent, NULL, NULL
		);
		InstallApplicationEventHandler(
			&keyupHandler, 1, &keyupEvent, NULL, NULL
		);

		EventHotKeyID hkid = {.id = handle};
		s = RegisterEventHotKey(
			key, mod, hkid, GetApplicationEventTarget(), 0, ref
		);
	});
	if (s != noErr) {
		return -1;
	}
	return 0;
}

int unregisterHotKey(EventHotKeyRef ref) {
	OSStatus s = UnregisterEventHotKey(ref);
	if (s != noErr) {
		return -1;
	}
	return 0;
}
