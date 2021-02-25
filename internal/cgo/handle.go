// Copyright 2021 The golang.design Initiative Authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.
//
// Written by Changkun Ou <changkun.de>

// Package cgo is an implementation of golang.org/issue/37033.
package cgo

import (
	"sync"
	"sync/atomic"
)

// Handle provides a safe representation of Go values to pass between
// C and Go back and forth. The zero value of a handle is not a valid
// handle, and thus safe to use as a sentinel in C APIs.
//
// The underlying type of Handle may change, but the value is guaranteed
// to fit in an integer type that is large enough to hold the bit pattern
// of any pointer. For instance, on the Go side:
//
// 	package main
//
// 	/*
// 	extern void MyGoPrint(unsigned long long handle);
// 	void myprint(unsigned long long handle);
// 	*/
// 	import "C"
// 	import "runtime/cgo"
//
// 	//export MyGoPrint
// 	func MyGoPrint(handle C.ulonglong) {
// 		h := cgo.Handle(handle)
// 		val := h.Value().(int)
// 		println(val)
// 		h.Delete()
// 	}
//
// 	func main() {
// 		val := 42
// 		C.myprint(C.ulonglong(cgo.NewHandle(val)))
// 		// Output: 42
// 	}
//
// and on the C side:
//
// 	// A Go function
// 	extern void MyGoPrint(unsigned long long handle);
//
// 	// A C function
// 	void myprint(unsigned long long handle) {
// 	    MyGoPrint(handle);
// 	}
type Handle uintptr

// NewHandle returns a handle for a given value.
//
// The handle is valid until the program calls Delete on it. The handle
// uses resources, and this package assumes that C code may hold on to
// the handle, so a program must explicitly call Delete when the handle
// is no longer needed.
//
// The intended use is to pass the returned handle to C code, which
// passes it back to Go, which calls Value. See an example in the
// comments of the Handle definition.
func NewHandle(v interface{}) Handle {
	h := atomic.AddUintptr(&idx, 1)
	if h == 0 {
		panic("runtime/cgo: ran out of handle space")
	}

	m.Store(h, v)
	return Handle(h)
}

// Value returns the associated Go value for a valid handle.
//
// The method panics if the handle is invalid already.
func (h Handle) Value() interface{} {
	v, ok := m.Load(uintptr(h))
	if !ok {
		panic("runtime/cgo: misuse of an invalid Handle")
	}
	return v
}

// Delete invalidates a handle. This method must be called when C code no
// longer has a copy of the handle, and the program no longer needs the
// Go value that associated with the handle.
//
// The method panics if the handle is invalid already.
func (h Handle) Delete() {
	_, ok := m.LoadAndDelete(uintptr(h))
	if !ok {
		panic("runtime/cgo: misuse of an invalid Handle")
	}
}

var (
	m   = &sync.Map{} // map[Handle]interface{}
	idx uintptr       // atomic
)
