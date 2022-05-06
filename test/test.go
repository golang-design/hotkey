package main

import (
	"fmt"
	h "golang.design/x/hotkey"
	"time"
)

func main() {
	hk := h.New([]h.Modifier{h.ModCtrl}, h.KeyM)
	hk.Register()
	var cancel chan bool = make(chan bool)
	go func() {
		time.Sleep(time.Second * 1)
		fmt.Println("Hotkey will be unregistered")
		hk.Unregister()
		fmt.Println("Hotkey unregistered")
		hk.Register()
		fmt.Println("Registered again")
		cancel <- true
	}()
	select {
	case <-hk.Keydown():
	case <-cancel:
	}
}
