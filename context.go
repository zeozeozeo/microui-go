// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2024 The Ebitengine Authors

package microui

func drawFrame(ctx *Context, rect Rect, colorid int) {
	ctx.DrawRect(rect, ctx.Style.Colors[colorid])
	if colorid == ColorScrollBase ||
		colorid == ColorScrollThumb ||
		colorid == ColorTitleBG {
		return
	}

	// draw border
	if ctx.Style.Colors[ColorBorder].A != 0 {
		ctx.DrawBox(expand_rect(rect, 1), ctx.Style.Colors[ColorBorder])
	}
}

func initContext(ctx *Context) {
	ctx.DrawFrame = drawFrame
	ctx._style = default_style
	ctx.Style = &ctx._style
}

func NewContext() *Context {
	ctx := &Context{}
	initContext(ctx)
	return ctx
}
