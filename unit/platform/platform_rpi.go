// +build rpi
// +build !darwin

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package platform

import (
	"github.com/djthorpe/gopi/v2"
	rpi "github.com/djthorpe/gopi/v2/sys/rpi"
)

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.Platform

func (this *platform) Init() error {
	// Initialise
	if err := rpi.BCMHostInit(); err != nil {
		return err
	} else {
		return nil
	}
}

func (this *platform) Close() error {
	// host terminate
	if err := rpi.BCMHostTerminate(); err != nil {
		return err
	}
	return nil
}

func (this *platform) Type() gopi.PlatformType {
	return gopi.PLATFORM_RPI | gopi.PLATFORM_LINUX
}

// Return serial number
func (this *platform) SerialNumber() string {
	// TODO
	return "TO BE DONE"
}
