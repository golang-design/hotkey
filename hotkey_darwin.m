// Copyright 2021 The golang.design Initiative Authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.
//
// Written by Changkun Ou <changkun.de>

//go:build darwin

#include <stdint.h>
#import <Cocoa/Cocoa.h>
#import <Carbon/Carbon.h>

extern void hotkeyCallback(uintptr_t handle);

static OSStatus
eventHandler(EventHandlerCallRef nextHandler, EventRef theEvent, void *userData) {
	EventHotKeyID k;
	GetEventParameter(theEvent, kEventParamDirectObject, typeEventHotKeyID, NULL, sizeof(k), NULL, &k);
	hotkeyCallback((uintptr_t)k.id); // use id as handle
	return noErr;
}

// registerHotkeyWithCallback registers a global system hotkey for callbacks.
int registerHotKey(int mod, int key, uintptr_t handle) {
	EventTypeSpec eventType;
	eventType.eventClass = kEventClassKeyboard;
	eventType.eventKind = kEventHotKeyPressed;
	InstallApplicationEventHandler(
		&eventHandler, 1, &eventType, NULL, NULL
	);

	EventHotKeyID hkid = {.id = handle};
	EventHotKeyRef ref;
	OSStatus s = RegisterEventHotKey(
		key, mod, hkid, GetApplicationEventTarget(), 0, &ref
	);
	if (s != noErr) {
		return -1;
	}
	return 0;
}


// The following three lines of code must run on the main thread.
// It must handle it using golang.design/x/mainthread.
//
// inspired from here: https://github.com/cehoffman/dotfiles/blob/4be8e893517e970d40746a9bdc67fe5832dd1c33/os/mac/iTerm2HotKey.m
void runApp() {
	[NSApplication sharedApplication];
	[NSApp disableRelaunchOnLogin];
	[NSApp run];
}

void stopApp() {
	[NSApp stop: nil];
}