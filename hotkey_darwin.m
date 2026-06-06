// Copyright 2021 The golang.design Initiative Authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.
//
// Written by Changkun Ou <changkun.de>

//go:build darwin

#include <stdint.h>
#import <ApplicationServices/ApplicationServices.h>
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

// --- Media keys via CGEventTap ---
//
// Media keys (play/pause, next, previous, volume) are not Carbon hot keys;
// they are delivered as NSSystemDefined events (type 14, subtype 8). We tap
// the session event stream, decode the embedded NX_KEYTYPE code and the
// up/down state, forward a matching press to Go, and consume the event.

#define NX_SYSDEFINED 14

typedef struct {
	uintptr_t handle;
	int nxKeyCode;
	CFMachPortRef tap;
	CFRunLoopSourceRef source;
} mediaTap;

static CGEventRef mediaTapCallback(CGEventTapProxy proxy, CGEventType type,
                                   CGEventRef event, void *userInfo) {
	mediaTap *mt = (mediaTap *)userInfo;

	// The system disables a tap that is too slow or interrupted; re-enable it.
	if (type == kCGEventTapDisabledByTimeout ||
	    type == kCGEventTapDisabledByUserInput) {
		if (mt->tap != NULL) {
			CGEventTapEnable(mt->tap, true);
		}
		return event;
	}

	NSEvent *e = [NSEvent eventWithCGEvent:event];
	if (e == nil || [e type] != NSEventTypeSystemDefined || [e subtype] != 8) {
		return event;
	}

	int data1 = (int)[e data1];
	int keyCode = (data1 & 0xffff0000) >> 16;
	if (keyCode != mt->nxKeyCode) {
		return event; // not the key this hotkey wants; pass through
	}
	int keyState = (data1 & 0x0000ff00) >> 8; // 0xA: down, 0xB: up
	if (keyState == 0x0a) {
		keydownCallback(mt->handle);
	} else if (keyState == 0x0b) {
		keyupCallback(mt->handle);
	}
	return NULL; // consume the event so the system does not also act on it
}

// registerMediaTap installs a CGEventTap for the given NX_KEYTYPE code.
// Returns NULL on failure, most commonly because the application has not been
// granted Accessibility (Input Monitoring) permission.
void* registerMediaTap(uintptr_t handle, int nxKeyCode) {
	// A keyboard event tap requires Accessibility (Input Monitoring) trust.
	// Check explicitly: CGEventTapCreate can otherwise return a non-NULL but
	// inert tap when untrusted, which would look like success but never fire.
	if (!AXIsProcessTrusted()) {
		return NULL;
	}

	mediaTap *mt = malloc(sizeof(mediaTap));
	mt->handle = handle;
	mt->nxKeyCode = nxKeyCode;
	mt->tap = NULL;
	mt->source = NULL;
	dispatch_sync(dispatch_get_main_queue(), ^{
		CGEventMask mask = CGEventMaskBit(NX_SYSDEFINED);
		CFMachPortRef tap = CGEventTapCreate(
			kCGSessionEventTap, kCGHeadInsertEventTap,
			kCGEventTapOptionDefault, mask, mediaTapCallback, mt);
		if (tap == NULL) {
			return;
		}
		mt->tap = tap;
		mt->source = CFMachPortCreateRunLoopSource(kCFAllocatorDefault, tap, 0);
		CFRunLoopAddSource(CFRunLoopGetMain(), mt->source, kCFRunLoopCommonModes);
		CGEventTapEnable(tap, true);
	});
	if (mt->tap == NULL) {
		free(mt);
		return NULL;
	}
	return mt;
}

void unregisterMediaTap(void* p) {
	if (p == NULL) {
		return;
	}
	mediaTap *mt = (mediaTap *)p;
	dispatch_sync(dispatch_get_main_queue(), ^{
		if (mt->tap != NULL) {
			CGEventTapEnable(mt->tap, false);
		}
		if (mt->source != NULL) {
			CFRunLoopRemoveSource(CFRunLoopGetMain(), mt->source,
			                      kCFRunLoopCommonModes);
			CFRelease(mt->source);
		}
		if (mt->tap != NULL) {
			CFRelease(mt->tap);
		}
	});
	free(mt);
}
