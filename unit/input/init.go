/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package input

import (
	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
)

func init() {
	gopi.UnitRegister(gopi.UnitConfig{
		Name:     "gopi/input",
		Type:     gopi.UNIT_INPUT_MANAGER,
		Requires: []string{"gopi/filepoll"},
		New: func(app gopi.App) (gopi.Unit, error) {
			return gopi.New(InputManager{
				FilePoll: app.UnitInstance("gopi/filepoll").(gopi.FilePoll),
			}, app.Log().Clone("gopi/input"))
		},
	})
}