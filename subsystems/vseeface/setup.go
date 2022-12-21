package vseeface

import (
	"errors"
	"os/exec"
	"strings"

	"github.com/magicmonkey/go-streamdeck"
	"github.com/sirupsen/logrus"
	"github.com/zyebytevt/starburst-go/lib"
)

var expressionButtons map[string]*lib.Button = make(map[string]*lib.Button)
var vsfWindowId string

var actionCallbacks map[string]lib.ActionCallback = map[string]lib.ActionCallback{
	"set_expression": setExpressionCallback,
}

func Setup(streamDeck *streamdeck.StreamDeck) error {
	logrus.Info("Initializing VSeeFace addon...")

	cmd := exec.Command("/usr/bin/xdotool", "search", "--name", "VSeeFace.*")
	output, err := cmd.Output()
	if err != nil {
		return errors.New("could not find VSeeFace window")
	}

	vsfWindowId = strings.Trim(string(output), "\n ")

	configs, err := lib.GetConfigsForKey("vseeface.buttons")

	if err != nil {
		return err
	}

	firstExpressionButton := true

	for _, config := range configs {
		btn, _ := lib.CreateButtonFromConfig(streamDeck, config, actionCallbacks)

		if config.ActionName == "set_expression" {
			// TODO: Think about this index thingy fix
			if firstExpressionButton {
				btn.SetDecorator(lib.ActiveStateDecorator)
				firstExpressionButton = false
			}
			expressionButtons[config.Parameters["key"].(string)] = btn
		}
	}

	return nil
}
