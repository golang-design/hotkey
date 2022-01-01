// Copyright 2021 The golang.design Initiative Authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.
//
// Written by Changkun Ou <changkun.de>

// Package hotkey provides the basic facility to register a system-level
// hotkey so that the application can be notified if a user triggers the
// desired hotkey. By definition, a hotkey is a combination of modifiers
// and a single key, and thus register a hotkey that contains multiple
// keys is not supported at the moment. Furthermore, because of OS
// restriction, hotkey events must be handled on the main thread.
//
// Therefore, in order to use this package properly, here is a complete
// example that corporates the golang.design/x/hotkey/mainthread package:
//
// 	package main
//
// 	import (
// 		"context"
//
// 		"golang.design/x/hotkey"
// 		"golang.design/x/hotkey/mainthread"
// 	)
//
// 	// initialize mainthread facility so that hotkey can be
// 	// properly registered to the system and handled by the
// 	// application.
// 	func main() { mainthread.Run(fn) }
// 	func fn() { // Use fn as the actual main function.
// 		var (
// 			mods = []hotkey.Modifier{hotkey.ModCtrl, hotkey.ModShift}
// 			k    = hotkey.KeyS
// 		)
//
// 		// Register a desired hotkey.
// 		hk := hotkey.New(mods, k)
// 		if err := hk.Register(); err != nil {
// 			panic("hotkey registration failed")
// 		}
//
// 		// Start listen hotkey event whenever you feel it is ready.
// 		for range hk.Listen() {
// 			println("hotkey ctrl+shift+s is triggered")
// 		}
// 	}
package hotkey

// Event represents a hotkey event
type Event struct{}

// Hotkey is a combination of modifiers and key to trigger an event
type Hotkey struct {
	platformHotkey

	mods []Modifier
	key  Key

	in  chan<- Event
	out <-chan Event
}

func New(mods []Modifier, key Key) *Hotkey {
	in, out := newEventChan()
	return &Hotkey{mods: mods, key: key, in: in, out: out}
}

// Register registers a combination of hotkeys. If the hotkey has
// registered. This function will invalidates the old registration
// and overwrites its callback.
func (hk *Hotkey) Register() error {
	return hk.register()
}

// Listen handles a hotkey event and triggers a call to fn.
// The hotkey listen hook terminates when the context is canceled.
func (hk *Hotkey) Listen() <-chan Event { return hk.out }

// Unregister unregisters the hotkey.
func (hk *Hotkey) Unregister() error {
	err := hk.unregister()
	if err != nil {
		return err
	}

	// Setup a new event channel.
	close(hk.in)
	in, out := newEventChan()
	hk.in = in
	hk.out = out
	return nil
}

// newEventChan returns a sender and a receiver of a buffered channel
// with infinite capacity.
func newEventChan() (chan<- Event, <-chan Event) {
	in, out := make(chan Event), make(chan Event)

	go func() {
		var q []Event

		for {
			e, ok := <-in
			if !ok {
				close(out)
				return
			}
			q = append(q, e)
			for len(q) > 0 {
				select {
				case out <- q[0]:
					q[0] = Event{}
					q = q[1:]
				case e, ok := <-in:
					if ok {
						q = append(q, e)
						break
					}
					for _, e := range q {
						out <- e
					}
					close(out)
					return
				}
			}
		}
	}()
	return in, out
}
