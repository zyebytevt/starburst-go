package vseeface

import (
	"github.com/sirupsen/logrus"
	"github.com/zyebytevt/streaming-backend/lib"
)

func setExpressionCallback(config *lib.ButtonConfig) {
	keyName := config.Parameters["key"].(string)

	if err := sendKeyShortcut(keyName); err != nil {
		logrus.WithError(err).Error("Failed to send key to VSeeFace!")
	} else {
		for _, button := range expressionButtons {
			button.SetActive(false)
		}

		if button, exists := expressionButtons[keyName]; exists {
			button.SetActive(true)
		}
	}
}
