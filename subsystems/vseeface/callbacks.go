package vseeface

import (
	"github.com/sirupsen/logrus"
	"github.com/zyebytevt/starburst-go/lib"
)

func setExpressionCallback(button *lib.Button) error {
	// This is a very dirty fix to see if this expression is already set.
	if button.GetHighlight() != lib.HighlightNone {
		return nil
	}

	button.SetHighlight(lib.HighlightInProgress)

	keyName := button.Config.Parameters["key"].(string)

	if err := sendKeyShortcut(keyName); err != nil {
		logrus.WithError(err).Error("Failed to send key to VSeeFace!")
	} else {
		for _, button := range expressionButtons {
			button.SetHighlight(lib.HighlightNone)
		}

		if button, exists := expressionButtons[keyName]; exists {
			button.SetHighlight(lib.HighlightActive)
		}
	}

	return nil
}
