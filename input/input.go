/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md

	This package provides input mechanisms, including the touchscreen
	interface for the official Raspberry Pi LED.
*/


// The input package provides a generic input device (mouse, touchscreen)
// for Linux-based input devices
package input

import (
	"io"
	"os"
	"time"
	"image"
	"errors"
	"syscall"
	"encoding/binary"
)

import (
	"fmt"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

// Store state for the non-abstract input driver
type Device struct {
	driver Driver
	poll   int
	event  syscall.EpollEvent
	position image.Point
	last_position image.Point
	slot uint32
	slots []*TouchEvent
}

// Abstract configuration which is used to open and return the
// concrete driver
type Config interface {
	// Opens the driver from configuration, or returns error
	Open() (Driver, error)
}

// Abstract driver interface
type Driver interface {
	// Return name of device
	GetName() string

	// Return the type of device
	GetType() DeviceType

	// Return the file descriptor
	GetFd() *os.File

	// Get number of touch slots (for touch screens)
	GetSlots() uint

	// Close closes the driver and frees the underlying resources
	Close() error
}

// Event
type Event struct {
	Timestamp  time.Duration
	Type       EventType
}

// TouchEvent
type TouchEvent struct {
	Event
	Identifier int
	Slot       uint32
	Point      image.Point
	LastPoint  image.Point
}

// Event callback
type TouchEventCallback func(*Device, *TouchEvent)

// Type of input device
type DeviceType int
type EventType int

// Non-exported raw event data structure sent over the wire
type rawEvent struct {
	Second      uint32
	Microsecond uint32
	Type        uint16
	Code        uint16
	Value       uint32
}

// Callback definition to process an event
type processEventsCallback func(syscall.EpollEvent) error

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

// Internal constants
const (
	MAX_POLL_EVENTS      = 32
	MAX_EVENT_SIZE_BYTES = 1024
)

// Constants which define the type of input device. At the moment, only
// touchscreen & mouse
const (
	TYPE_TOUCHSCREEN DeviceType = iota
	TYPE_MOUSE
)

// Constants which define the type of event.
const (
	EVENT_UNKNOWN EventType = iota
	EVENT_BTN_PRESS
	EVENT_BTN_RELEASE
	EVENT_MOVE
	EVENT_SLOT_PRESS
	EVENT_SLOT_RELEASE
	EVENT_SLOT_MOVE
)

// Event types
// See https://www.kernel.org/doc/Documentation/input/event-codes.txt for
// more information
const (
	EV_SYN uint16 = 0x0000 // SYNC
	EV_KEY uint16 = 0x0001 // KEY PRESS
	EV_ABS uint16 = 0x0003 // ABSOLUTE AXIS VALUE CHANGE
)

// Button touch, press, release (when EV_KEY code)
// See https://www.kernel.org/doc/Documentation/input/event-codes.txt for
// more information
const (
	BTN_TOUCH         uint16 = 0x014A
	BTN_TOUCH_RELEASE uint32 = 0x00000000
	BTN_TOUCH_PRESS   uint32 = 0x00000001
)

// Multi-Touch Types
// See https://www.kernel.org/doc/Documentation/input/multi-touch-protocol.txt
// for more information
const (
	ABS_X              uint16 = 0x0000
	ABS_Y              uint16 = 0x0001
	ABS_MT_SLOT        uint16 = 0x002F // 47 MT slot being modified
	ABS_MT_POSITION_X  uint16 = 0x0035 // 53 Center X of multi touch position
	ABS_MT_POSITION_Y  uint16 = 0x0036 // 54 Center Y of multi touch position
	ABS_MT_TRACKING_ID uint16 = 0x0039 // 57 Unique ID of initiated contact
)

////////////////////////////////////////////////////////////////////////////////
// Opener methods

// Open opens a connection to the touchscreen with the given driver.
func Open(config Config) (*Device, error) {
	driver, err := config.Open()
	if err != nil {
		return nil, err
	}

	device := new(Device)
	device.driver = driver
	device.poll, err = syscall.EpollCreate1(0)
	if err != nil {
		driver.Close()
		return nil, err
	}

	// register the poll with the device
	device.event.Events = syscall.EPOLLIN
	device.event.Fd = int32(driver.GetFd().Fd())
	if err = syscall.EpollCtl(device.poll, syscall.EPOLL_CTL_ADD, int(device.event.Fd), &device.event); err != nil {
		syscall.Close(device.poll)
		driver.Close()
		return nil, err
	}

	// Set position
	device.position = image.ZP
	device.last_position = image.Point{ -1, -1 }

	// GetSlots will return positive non-zero value where the device is a slot
	// based multitouch device, for example, where you can use more than one
	// finger on a touchscreen
	num_slots := driver.GetSlots()
	if num_slots > 0 {
		// set slot to zero, create the slots, set slot values
		device.slot = 0
		device.slots = make([]*TouchEvent, driver.GetSlots())
		for i, _ := range device.slots {
			device.slots[i] = &TouchEvent{ Slot: uint32(i) }
		}
	}

	// success - return device
	return device, nil
}

////////////////////////////////////////////////////////////////////////////////
// Public Device methods

// Closes the device and frees the resources
func (device *Device) Close() error {
	return device.driver.Close()
}

// Gets the name of the input device
func (device *Device) GetName() string {
	return device.driver.GetName()
}

// Gets the type of the input device
func (device *Device) GetType() DeviceType {
	return device.driver.GetType()
}

// Gets the mouse or touchscreen position
func (device *Device) GetPosition() image.Point {
	return device.position
}

// Processes touch events for touch devices, blocks when there are no
// events to process. On error, returns
func (device *Device) ProcessEvents(callback TouchEventCallback) error {
	if device.GetType() != TYPE_TOUCHSCREEN {
		return errors.New("Invalid device type in call to ProcessTouchEvents")
	}
	err := device.waitForEvents(func(event syscall.EpollEvent) error {
		for {
			var event rawEvent
			err := binary.Read(device.driver.GetFd(), binary.LittleEndian, &event)
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}
			err = device.processRawTouchEvent(&event,callback)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

////////////////////////////////////////////////////////////////////////////////
// Private Device methods

// Waits for new raw events, and then executes the callback
func (device *Device) waitForEvents(callback processEventsCallback) error {
	events := make([]syscall.EpollEvent, MAX_POLL_EVENTS)
	for {
		n, err := syscall.EpollWait(device.poll, events, -1)
		if err != nil {
			return err
		}
		if n <= 0 {
			continue
		}
		for _, event := range events[:n] {
			if event.Fd != int32(device.driver.GetFd().Fd()) {
				continue
			}
			callback(event)
		}
	}
	return nil
}

func (device *Device) processRawTouchEvent(event *rawEvent,callback TouchEventCallback) error {
	// Calculate timestamp
	ts := time.Duration(time.Duration(event.Second) * time.Second + time.Duration(event.Microsecond) * time.Microsecond)

	// Parse raw event data
	switch {
	case event.Type == EV_SYN:

		// Fire EVENT_MOVE
		if device.position.Eq(device.last_position) == false {
			callback(device,&TouchEvent{ Event{ Timestamp: ts, Type: EVENT_MOVE }, 0, 0, device.position, device.last_position })
		}
		device.last_position = device.position

		// Don't do slot-based events if there aren't any slots
		if device.slots == nil {
			return nil
		}

		e := device.slots[device.slot]
		e.Timestamp = ts

		// If type of event is not release, then set to press
		if e.Type == EVENT_SLOT_PRESS {
			e.Type = EVENT_SLOT_MOVE
		} else if e.Type != EVENT_SLOT_RELEASE && e.Type != EVENT_SLOT_MOVE {
			e.Type = EVENT_SLOT_PRESS
		}
		callback(device,e)

		// Set slot state back to unknown
		if e.Type == EVENT_SLOT_RELEASE {
			e.Type = EVENT_UNKNOWN
		}

		// Set last position to current one
		e.LastPoint = e.Point

		return nil
	case event.Type == EV_KEY && event.Code == BTN_TOUCH && event.Value == BTN_TOUCH_PRESS:
		callback(device,&TouchEvent{ Event{ Timestamp: ts, Type: EVENT_BTN_PRESS }, 0, 0, device.position, image.ZP })
		return nil
	case event.Type == EV_KEY && event.Code == BTN_TOUCH && event.Value == BTN_TOUCH_RELEASE:
		callback(device,&TouchEvent{ Event{ Timestamp: ts, Type: EVENT_BTN_RELEASE }, 0, 0, device.position, image.ZP })
		return nil
	case event.Type == EV_ABS && event.Code == ABS_X:
		device.position.X = int(event.Value)
		return nil
	case event.Type == EV_ABS && event.Code == ABS_Y:
		device.position.Y = int(event.Value)
		return nil
	case event.Type == EV_ABS && event.Code == ABS_MT_POSITION_X:
		device.slots[device.slot].Point.X = int(event.Value)
		return nil
	case event.Type == EV_ABS && event.Code == ABS_MT_POSITION_Y:
		device.slots[device.slot].Point.Y = int(event.Value)
		return nil
	case event.Type == EV_ABS && event.Code == ABS_MT_SLOT:
		if event.Value >= uint32(len(device.slots)) {
			return errors.New("Invalid slot value")
		}
		device.slot = event.Value
		return nil
	case event.Type == EV_ABS && event.Code == ABS_MT_TRACKING_ID:
		// Identifier is a 16 bit value which we turn into an int
		id := int(int16(event.Value))
		if id == -1 {
			device.slots[device.slot].Type = EVENT_SLOT_RELEASE
		} else {
			device.slots[device.slot].Identifier = id
		}
		return nil
	}
	return errors.New(fmt.Sprintf("Invalid event with type: %v", event.Type))
}





