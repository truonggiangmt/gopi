/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/

// This package implements the concrete interface for mouse pointing devices
// for the Raspberry Pi. In order to use it, you need to open the mouse using
// the abstract instance:
//
//   mouse, err := input.Open(mouse.Config{ })
//   if err != nil { handle error }
//   defer mouse.Close()
//
package mouse // import "github.com/djthorpe/gopi/device/mouse"

// System imports
import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// Local imports
import (
	"github.com/djthorpe/gopi/input"
)

////////////////////////////////////////////////////////////////////////////////

// The configuration options for the mouse. Basically, there are no
// configuration options for this device.
type Config struct {
	Name string
}

// The driver state
type Driver struct {
	device string
	name   string
	file   *os.File
}

////////////////////////////////////////////////////////////////////////////////

const (
	PATH_INPUT_DEVICES = "/sys/class/input/event*"
)

////////////////////////////////////////////////////////////////////////////////
// input.Opener interface

// Concrete Open method
func (config Config) Open() (input.Driver, error) {
	var err error

	driver := new(Driver)
	driver.name, driver.device, err = getDeviceNameAndPath(&config)
	if err != nil {
		return nil, err
	}

	// open driver
	driver.file, err = os.Open(driver.device)
	if err != nil {
		return nil, err
	}

	return driver, nil
}

////////////////////////////////////////////////////////////////////////////////
// input.Driver interface

func (this *Driver) Close() error {
	return this.file.Close()
}

func (this *Driver) GetName() string {
	return this.name
}

func (this *Driver) GetType() input.DeviceType {
	return input.TYPE_MOUSE
}

func (this *Driver) GetFd() *os.File {
	return this.file
}

func (this *Driver) GetSlots() uint {
	return 0
}

func (this *Driver) String() string {
	return fmt.Sprintf("<device.Mouse>{ device=%v }", this.device)
}

////////////////////////////////////////////////////////////////////////////////
// Private Methods

func getDeviceNameAndPath(config *Config) (string, string, error) {
	files, err := filepath.Glob(PATH_INPUT_DEVICES)
	if err != nil {
		return "", "", err
	}
	for _, file := range files {
		buf, err := ioutil.ReadFile(path.Join(file, "device", "name"))
		if err != nil {
			continue
		}
		if path.Base(file) == config.Name {
			return strings.TrimSpace(string(buf)), path.Join("/", "dev", "input", path.Base(file)), nil
		}
	}
	return "", "", errors.New("Device not found")
}