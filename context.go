// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2024 The Ebitengine Authors

package microui

import (
	"image"
)

func (c *Context) drawFrame(rect image.Rectangle, colorid int) {
	c.DrawRect(rect, c.Style.Colors[colorid])
	if colorid == ColorScrollBase ||
		colorid == ColorScrollThumb ||
		colorid == ColorTitleBG {
		return
	}

	// draw border
	if c.Style.Colors[ColorBorder].A != 0 {
		c.DrawBox(rect.Inset(-1), c.Style.Colors[ColorBorder])
	}
}

func initContext(ctx *Context) {
	ctx.Style = &defaultStyle
}

func NewContext() *Context {
	ctx := &Context{}
	initContext(ctx)
	return ctx
}
