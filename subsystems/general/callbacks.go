package general

import (
	"os/exec"

	"github.com/zyebytevt/starburst-go/lib"
)

func setBrightnessCallback(button *lib.Button) {
	value := button.Config.Parameters["value"].(int)

	if button.Config.Parameters["absolute"].(bool) {
		brightness = value
	} else {
		brightness += value
	}

	if brightness < 0 {
		brightness = 0
	} else if brightness > 100 {
		brightness = 100
	}

	sd.SetBrightness(brightness)
}

func executeCallback(button *lib.Button) {
	exec.Command(button.Config.Parameters["program"].(string), button.Config.Parameters["cmdline"].(string)).Run()
}
