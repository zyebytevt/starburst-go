package lib

import (
	"image/color"
	"time"

	"github.com/magicmonkey/go-streamdeck"
	"github.com/magicmonkey/go-streamdeck/buttons"
	"github.com/sirupsen/logrus"

	sddecorators "github.com/magicmonkey/go-streamdeck/decorators"
)

type HighlightType int

const (
	HighlightNone HighlightType = iota
	HighlightInProgress
	HighlightActive
)

var errorButtonBorder = sddecorators.NewBorder(10, color.RGBA{R: 255, G: 0, B: 0, A: 255})

var activeButtonBorder = sddecorators.NewBorder(8, color.RGBA{R: 0, G: 255, B: 0, A: 255})
var inProgressButtonBorder = sddecorators.NewBorder(6, color.RGBA{R: 0, G: 255, B: 150, A: 255})

type ActionCallback func(button *Button) error

type Button struct {
	streamDeck       *streamdeck.StreamDeck
	index            int
	onPressed        ActionCallback
	currentHighlight HighlightType

	Config *ButtonConfig
}

func NewButton(streamDeck *streamdeck.StreamDeck, index int, imagePath string, onPressed ActionCallback, config *ButtonConfig) (*Button, error) {
	internalButton, err := buttons.NewImageFileButton(imagePath)

	if err != nil {
		return nil, err
	}

	button := &Button{
		streamDeck:       streamDeck,
		index:            index,
		onPressed:        onPressed,
		currentHighlight: HighlightNone,
		Config:           config,
	}

	internalButton.SetActionHandler(button)
	streamDeck.AddButton(index, internalButton)

	return button, nil
}

func (btn *Button) Pressed(b streamdeck.Button) {
	if btn.onPressed != nil {
		if err := btn.onPressed(btn); err != nil {
			go btn.FlashError(5)
			logrus.WithError(err).Error("Button callback failed with error.")
		}
	}
}

func (btn *Button) SetHighlight(highlight HighlightType) {
	if highlight == HighlightNone {
		btn.streamDeck.UnsetDecorator(btn.index)
	} else if highlight == HighlightInProgress {
		btn.streamDeck.SetDecorator(btn.index, inProgressButtonBorder)
	} else if highlight == HighlightActive {
		btn.streamDeck.SetDecorator(btn.index, activeButtonBorder)
	}

	btn.currentHighlight = highlight
}

func (btn *Button) GetHighlight() HighlightType {
	return btn.currentHighlight
}

func (btn *Button) FlashError(count int) {
	for i := 0; i < count; i++ {
		btn.streamDeck.UnsetDecorator(btn.index)
		time.Sleep(time.Millisecond * 150)
		btn.streamDeck.SetDecorator(btn.index, errorButtonBorder)
		time.Sleep(time.Millisecond * 150)
	}

	btn.SetHighlight(btn.currentHighlight)
}
