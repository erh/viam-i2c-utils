package main

import (
	"viami2cutils"

	toggleswitch "go.viam.com/rdk/components/switch"
	"go.viam.com/rdk/module"
	"go.viam.com/rdk/resource"
)

func main() {
	module.ModularMain(resource.APIModel{toggleswitch.API, viami2cutils.I2cOnOff})
}
