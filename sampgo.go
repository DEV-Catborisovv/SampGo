package sampgo

import "C"
import (
	"fmt"
	"unsafe"
)

type EventType int

const (
	Repeat EventType = iota
	OnceOnly
)

type event struct {
	Handler interface{}
	Type    EventType
}

var events = make(map[string]event)
var mainEvent func() = nil

//export onTick
func onTick() {
	evt, ok := events["tick"]
	if !ok {
		return
	}

	fn, ok := evt.Handler.(func())
	if !ok {
		return
	}
	fn()

	if evt.Type == OnceOnly {
		// If this event was registered only once, then reassign this event to a new blank struct.
		// PoC code. Still needs to be implemented completely.
	}

	return
}

// On registers an event with a handler.
func On(eventName string, handler interface{}) error {
	_, ok := events[eventName]
	if ok {
		return fmt.Errorf("this handler already exists")
	}

	events[eventName] = event{Handler: handler, Type: Repeat}
	_ = Print(fmt.Sprintf("Registered %s event", eventName))

	return nil
}

// Once registers an event with a handler one time only.
func Once(eventName string, handler interface{}) error {
	_, ok := events[eventName]
	if ok {
		return fmt.Errorf("this handler already exists")
	}

	events[eventName] = event{Handler: handler, Type: OnceOnly}
	_ = Print(fmt.Sprintf("Registered %s event", eventName))

	return nil
}

// Print allows you to print to the SAMP console.
func Print(msg string) error {
	cstr := C.CString(msg)
	defer C.free(unsafe.Pointer(cstr))
	C.goLogprintf(cstr)

	return nil
}
