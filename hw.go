/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package gopi

import (
	"strings"
	"time"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type (
	PlatformType uint
)

////////////////////////////////////////////////////////////////////////////////
// INTERFACES

type Platform interface {

	// Return platform type
	Type() PlatformType

	// Return serial number
	SerialNumber() string

	// Return uptime
	Uptime() time.Duration

	// Return 1, 5 and 15 minute load averages
	LoadAverages() (float64, float64, float64)

	// Implements gopi.Unit
	Unit
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	PLATFORM_NONE   PlatformType = 0
	PLATFORM_DARWIN PlatformType = (1 << iota) >> 1
	PLATFORM_RPI
	PLATFORM_LINUX
	PLATFORM_MIN = PLATFORM_DARWIN
	PLATFORM_MAX = PLATFORM_LINUX
)

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (p PlatformType) String() string {
	str := ""
	if p == 0 {
		return p.FlagString()
	}
	for v := PLATFORM_MIN; v <= PLATFORM_MAX; v <<= 1 {
		if p&v == v {
			str += "|" + v.FlagString()
		}
	}
	return strings.TrimPrefix(str, "|")
}

func (p PlatformType) FlagString() string {
	switch p {
	case PLATFORM_NONE:
		return "PLATFORM_NONE"
	case PLATFORM_DARWIN:
		return "PLATFORM_DARWIN"
	case PLATFORM_RPI:
		return "PLATFORM_RPI"
	case PLATFORM_LINUX:
		return "PLATFORM_LINUX"
	default:
		return "[?? Invalid PlatformType value]"
	}
}
