// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2024 The Ebitengine Authors

package microui

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

/*============================================================================
** input handlers
**============================================================================*/

func (ctx *Context) InputMouseMove(x, y int) {
	ctx.MousePos = image.Pt(x, y)
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

func (ctx *Context) InputMouseDown(x, y int, btn ebiten.MouseButton) {
	ctx.InputMouseMove(x, y)
	ctx.MouseDown |= mouseButtonToInt(btn)
	ctx.MousePressed |= mouseButtonToInt(btn)
}

func (ctx *Context) InputMouseUp(x, y int, btn ebiten.MouseButton) {
	ctx.InputMouseMove(x, y)
	ctx.MouseDown &= ^mouseButtonToInt(btn)
}

func (ctx *Context) InputScroll(x, y int) {
	ctx.ScrollDelta.X += x
	ctx.ScrollDelta.Y += y
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

func (ctx *Context) InputKeyDown(key ebiten.Key) {
	ctx.KeyPressed |= keyToInt(key)
	ctx.KeyDown |= keyToInt(key)
}

func (ctx *Context) InputKeyUp(key ebiten.Key) {
	ctx.KeyDown &= ^keyToInt(key)
}

func (ctx *Context) InputText(text []rune) {
	ctx.TextInput = text
}
