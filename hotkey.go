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
// example that corporates the golang.design/x/mainthread package:
//
// 	package main
//
// 	import (
// 		"context"
//
// 		"golang.design/x/hotkey"
// 		"golang.design/x/mainthread"
// 	)
//
// 	// initialize mainthread facility so that hotkey can be
// 	// properly registered to the system and handled by the
// 	// application.
// 	func main() { mainthread.Init(fn) }
// 	func fn() { // Use fn as the actual main function.
// 		var (
// 			mods = []hotkey.Modifier{hotkey.ModCtrl}
// 			k    = hotkey.KeyS
// 		)
//
// 		// Register a desired hotkey.
// 		hk, err := hotkey.Register(mods, k)
// 		if err != nil {
// 			panic("hotkey registration failed")
// 		}
//
// 		// Start listen hotkey event whenever you feel it is ready.
// 		triggered := hk.Listen(context.Background())
// 		for range triggered {
// 			println("hotkey ctrl+s is triggered")
// 		}
// 	}
package hotkey

import (
	"context"
	"runtime"

	"golang.design/x/mainthread"
)

// Event represents a hotkey event
type Event struct{}

// Hotkey is a combination of modifiers and key to trigger an event
type Hotkey struct {
	mods []Modifier
	key  Key

	in  chan<- Event
	out <-chan Event
}

// Register registers a combination of hotkeys. If the hotkey has
// registered. This function will invalidates the old registration
// and overwrites its callback.
func Register(mods []Modifier, key Key) (*Hotkey, error) {
	in, out := newEventChan()
	hk := &Hotkey{mods, key, in, out}

	var err error
	mainthread.Call(func() { err = hk.register() })
	if err != nil {
		return nil, err
	}

	runtime.SetFinalizer(hk, func(hk *Hotkey) {
		hk.unregister()
	})

	return hk, nil
}

// Listen handles a hotkey event and triggers a call to fn.
// The hotkey listen hook terminates when the context is canceled.
func (hk *Hotkey) Listen(ctx context.Context) <-chan Event {
	mainthread.Go(func() { hk.handle(ctx) })
	return hk.out
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
