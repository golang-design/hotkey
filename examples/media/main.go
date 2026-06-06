// Copyright 2026 The golang.design Initiative Authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.

// Command media demonstrates registering media keys (play/pause, next,
// previous, volume) and reporting their press/release events.
//
// Run it:
//
//	go run golang.design/x/hotkey/examples/media
//
// then press the media keys on your keyboard. Press Ctrl+C to quit.
//
// Platform notes:
//
//   - Linux (X11) and Windows: media keys work with no special permission.
//
//   - macOS: media keys are captured with a CGEventTap, which requires
//     Accessibility / Input Monitoring permission. Because `go run` builds a
//     throwaway binary whose path changes every run, the permission grant will
//     not stick; build a stable binary instead and grant *that* (or grant your
//     terminal) in System Settings → Privacy & Security → Accessibility:
//
//     go build -o mediatest golang.design/x/hotkey/examples/media
//     ./mediatest   # then approve the Accessibility prompt and re-run
//
//     If permission is missing, Register returns an error (it does not crash).
//     macOS has no media-stop key, so KeyMediaStop returns an error there.
package main

import (
	"log"

	"golang.design/x/hotkey"
	"golang.design/x/hotkey/mainthread"
)

func main() { mainthread.Init(fn) } // Required on macOS for the event loop.

func fn() {
	// Media keys are global and take no modifiers.
	media := []struct {
		name string
		key  hotkey.Key
	}{
		{"play/pause", hotkey.KeyMediaPlayPause},
		{"next", hotkey.KeyMediaNext},
		{"previous", hotkey.KeyMediaPrev},
		{"volume up", hotkey.KeyVolumeUp},
		{"volume down", hotkey.KeyVolumeDown},
		{"mute", hotkey.KeyVolumeMute},
	}

	registered := 0
	for _, m := range media {
		hk := hotkey.New(nil, m.key)
		if err := hk.Register(); err != nil {
			log.Printf("skip %-12s: %v", m.name, err)
			continue
		}
		registered++
		log.Printf("registered: %s", m.name)
		go listen(m.name, hk)
	}

	if registered == 0 {
		log.Fatal("no media hotkey could be registered (on macOS, grant Accessibility permission)")
	}

	log.Println("press a media key to see events; Ctrl+C to quit")
	select {} // run until interrupted
}

func listen(name string, hk *hotkey.Hotkey) {
	for {
		<-hk.Keydown()
		log.Printf("%-12s down", name)
		<-hk.Keyup()
		log.Printf("%-12s up", name)
	}
}
