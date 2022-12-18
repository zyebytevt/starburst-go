package vseeface

import (
	"os/exec"
	"strings"
	"time"

	"github.com/magicmonkey/go-streamdeck"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/zyebytevt/streaming-backend/lib"
)

type ExpressionConfig struct {
	KeyName     string `mapstructure:"key"`
	ButtonText  string `mapstructure:"button_text"`
	ButtonIndex int    `mapstructure:"button_index"`
}

var expressionButtons map[string]*lib.Button = make(map[string]*lib.Button)

func Setup(streamDeck *streamdeck.StreamDeck) {
	logrus.Info("Initializing VSeeFace addon...")

	configs := make([]*ExpressionConfig, 0)
	viper.UnmarshalKey("vseeface.expressions", &configs)

	for i, config := range configs {
		btn := lib.NewButton(streamDeck, config.ButtonIndex, config.ButtonText, func(userData any) {
			funcConfig := userData.(*ExpressionConfig)
			if err := sendKeyShortcut(funcConfig.KeyName); err != nil {
				logrus.WithError(err).Error("Failed to send key to VSeeFace!")
			} else {
				for _, button := range expressionButtons {
					button.SetActive(false)
				}

				if button, exists := expressionButtons[funcConfig.KeyName]; exists {
					button.SetActive(true)
				}
			}
		}, config)

		if i == 0 {
			btn.SetActive(true)
		}
		expressionButtons[config.KeyName] = btn
	}
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
