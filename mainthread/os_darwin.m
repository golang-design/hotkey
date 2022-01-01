// Copyright 2022 The golang.design Initiative Authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.
//
// Written by Changkun Ou <changkun.de>

//go:build darwin

#include <stdint.h>
#import <Cocoa/Cocoa.h>

extern void dispatchMainFuncs();

void wakeupMainThread(void) {
	dispatch_async(dispatch_get_main_queue(), ^{
		dispatchMainFuncs();
	});
}

// The following three lines of code must run on the main thread.
// It must handle it using golang.design/x/mainthread.
//
// inspired from here: https://github.com/cehoffman/dotfiles/blob/4be8e893517e970d40746a9bdc67fe5832dd1c33/os/mac/iTerm2HotKey.m
void os_main(void) {
	[NSApplication sharedApplication];
	[NSApp disableRelaunchOnLogin];
	[NSApp run];
}