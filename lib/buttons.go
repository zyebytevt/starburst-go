package lib

import (
	"image/color"
	"time"

	"github.com/magicmonkey/go-streamdeck"
	"github.com/magicmonkey/go-streamdeck/buttons"
	"github.com/sirupsen/logrus"

	sddecorators "github.com/magicmonkey/go-streamdeck/decorators"
)

var notifySuccessDecorator = sddecorators.NewBorder(10, color.RGBA{R: 0, G: 255, B: 0, A: 255})
var notifyErrorDecorator = sddecorators.NewBorder(10, color.RGBA{R: 255, G: 0, B: 0, A: 255})

var ActiveStateDecorator = sddecorators.NewBorder(8, color.RGBA{R: 255, G: 255, B: 255, A: 255})
var InProgressStateDecorator = sddecorators.NewBorder(6, color.RGBA{R: 128, G: 128, B: 128, A: 255})

type ActionCallback func(button *Button) error

type NotifyType int

type Button struct {
	streamDeck       *streamdeck.StreamDeck
	index            int
	onPressed        ActionCallback
	currentDecorator *sddecorators.Border

	Config   *ButtonConfig
	UserData map[string]any
}

const (
	NotifySuccess NotifyType = iota
	NotifyError
)

func NewButton(streamDeck *streamdeck.StreamDeck, index int, imagePath string, onPressed ActionCallback, config *ButtonConfig) (*Button, error) {
	internalButton, err := buttons.NewImageFileButton(imagePath)

	if err != nil {
		return nil, err
	}

	button := &Button{
		streamDeck: streamDeck,
		index:      index,
		onPressed:  onPressed,
		Config:     config,
		UserData:   make(map[string]any),
	}

	internalButton.SetActionHandler(button)
	streamDeck.AddButton(index, internalButton)

	return button, nil
}

func (btn *Button) Pressed(b streamdeck.Button) {
	if btn.onPressed != nil {
		if err := btn.onPressed(btn); err != nil {
			go btn.FlashNotify(NotifyError)
			logrus.WithError(err).Error("Button callback failed with error.")
		}
	}
}

func (btn *Button) SetDecorator(decorator *sddecorators.Border) {
	if decorator == nil {
		btn.streamDeck.UnsetDecorator(btn.index)
	} else {
		btn.streamDeck.SetDecorator(btn.index, decorator)
	}

	btn.currentDecorator = decorator
}

func (btn *Button) GetDecorator() *sddecorators.Border {
	return btn.currentDecorator
}

func (btn *Button) FlashNotify(notifyType NotifyType) {
	var decorator *sddecorators.Border
	var count int

	if notifyType == NotifyError {
		decorator = notifyErrorDecorator
		count = 5
	} else if notifyType == NotifySuccess {
		decorator = notifySuccessDecorator
		count = 3
	}

	for i := 0; i < count; i++ {
		btn.streamDeck.UnsetDecorator(btn.index)
		time.Sleep(time.Millisecond * 150)
		btn.streamDeck.SetDecorator(btn.index, decorator)
		time.Sleep(time.Millisecond * 150)
	}

	if btn.currentDecorator == nil {
		btn.streamDeck.UnsetDecorator(btn.index)
	} else {
		btn.streamDeck.SetDecorator(btn.index, btn.currentDecorator)
	}
}
