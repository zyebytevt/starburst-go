package vseeface

import (
	"os/exec"
	"strings"
	"time"

	"github.com/magicmonkey/go-streamdeck"
	"github.com/sirupsen/logrus"
	"github.com/zyebytevt/streaming-backend/lib"
)

var expressionButtons map[string]*lib.Button = make(map[string]*lib.Button)

var actionCallbacks map[string]lib.ActionCallback = map[string]lib.ActionCallback{
	"set_expression": setExpressionCallback,
}

func Setup(streamDeck *streamdeck.StreamDeck) error {
	logrus.Info("Initializing VSeeFace addon...")

	configs, err := lib.GetConfigsForKey("vseeface.buttons")

	if err != nil {
		return err
	}

	for i, config := range configs {
		btn := lib.CreateButtonFromConfig(streamDeck, config, actionCallbacks)

		if config.ActionName == "set_expression" {
			// TODO: Think about this index thingy fix
			if i == 0 {
				btn.SetActive(true)
			}
			expressionButtons[config.Parameters["key"].(string)] = btn
		}
	}

	return nil
}

func sendKeyShortcut(key string) error {
	cmd := exec.Command("/usr/bin/xdotool", "search", "--name", "VSeeFace.*")
	output, err := cmd.Output()
	if err != nil {
		return err
	}

	vsfWindowId := strings.Trim(string(output), "\n ")

	exec.Command("/usr/bin/xdotool", "keydown", "--window", vsfWindowId, "CTRL").Run()
	time.Sleep(time.Millisecond * 100)
	exec.Command("/usr/bin/xdotool", "keydown", "--window", vsfWindowId, "SHIFT").Run()
	time.Sleep(time.Millisecond * 100)
	exec.Command("/usr/bin/xdotool", "keydown", "--window", vsfWindowId, key).Run()
	time.Sleep(time.Millisecond * 100)
	exec.Command("/usr/bin/xdotool", "keyup", "--window", vsfWindowId, key).Run()
	time.Sleep(time.Millisecond * 100)
	exec.Command("/usr/bin/xdotool", "keyup", "--window", vsfWindowId, "SHIFT").Run()
	time.Sleep(time.Millisecond * 100)
	exec.Command("/usr/bin/xdotool", "keyup", "--window", vsfWindowId, "CTRL").Run()

	return nil
}
