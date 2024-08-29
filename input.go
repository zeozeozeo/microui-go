// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2024 The Ebitengine Authors

package microui

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

func (ctx *Context) inputMouseMove(x, y int) {
	ctx.mousePos = image.Pt(x, y)
}

func mouseButtonToInt(btn ebiten.MouseButton) int {
	switch btn {
	case ebiten.MouseButtonLeft:
		return mouseLeft
	case ebiten.MouseButtonRight:
		return mouseRight
	case ebiten.MouseButtonMiddle:
		return mouseMiddle
	}
	return 0
}

func (ctx *Context) inputMouseDown(x, y int, btn ebiten.MouseButton) {
	ctx.inputMouseMove(x, y)
	ctx.mouseDown |= mouseButtonToInt(btn)
	ctx.mousePressed |= mouseButtonToInt(btn)
}

func (ctx *Context) inputMouseUp(x, y int, btn ebiten.MouseButton) {
	ctx.inputMouseMove(x, y)
	ctx.mouseDown &= ^mouseButtonToInt(btn)
}

func (ctx *Context) inputScroll(x, y int) {
	ctx.scrollDelta.X += x
	ctx.scrollDelta.Y += y
}

func keyToInt(key ebiten.Key) int {
	switch key {
	case ebiten.KeyShift:
		return keyShift
	case ebiten.KeyControl:
		return keyControl
	case ebiten.KeyAlt:
		return keyAlt
	case ebiten.KeyBackspace:
		return keyBackspace
	case ebiten.KeyEnter:
		return keyReturn
	}
	return 0
}

func (ctx *Context) inputKeyDown(key ebiten.Key) {
	ctx.keyPressed |= keyToInt(key)
	ctx.keyDown |= keyToInt(key)
}

func (ctx *Context) inputKeyUp(key ebiten.Key) {
	ctx.keyDown &= ^keyToInt(key)
}

func (ctx *Context) inputText(text []rune) {
	ctx.textInput = text
}
