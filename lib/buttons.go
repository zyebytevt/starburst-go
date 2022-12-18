package lib

import (
	"image/color"

	"github.com/magicmonkey/go-streamdeck"
	"github.com/magicmonkey/go-streamdeck/buttons"

	sddecorators "github.com/magicmonkey/go-streamdeck/decorators"
)

var activeButtonBorder = sddecorators.NewBorder(5, color.RGBA{R: 255, G: 0, B: 0, A: 255})

type Button struct {
	streamDeck *streamdeck.StreamDeck
	index      int
	onPressed  func(userData any)
	userData   any
}

func NewButton(streamDeck *streamdeck.StreamDeck, index int, name string, onPressed func(userData any), userData any) *Button {
	internalButton := buttons.NewTextButton(name)

	button := &Button{streamDeck: streamDeck, index: index, onPressed: onPressed, userData: userData}

	internalButton.SetActionHandler(button)
	streamDeck.AddButton(index, internalButton)

	return button
}

func (btn *Button) Pressed(b streamdeck.Button) {
	if btn.onPressed != nil {
		btn.onPressed(btn.userData)
	}
}

func (btn *Button) SetActive(active bool) {
	if active {
		btn.streamDeck.SetDecorator(btn.index, activeButtonBorder)
	} else {
		btn.streamDeck.UnsetDecorator(btn.index)
	}
}
