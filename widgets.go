// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2024 The Ebitengine Authors

package microui

import "image"

func (ctx *Context) Button(label string) Res {
	return ctx.ButtonEx(label, 0, OptAlignCenter)
}

func (ctx *Context) TextBox(buf *string) Res {
	return ctx.TextBoxEx(buf, 0)
}

func (ctx *Context) Slider(value *float64, lo, hi float64) Res {
	return ctx.SliderEx(value, lo, hi, 0, sliderFmt, OptAlignCenter)
}

func (ctx *Context) Number(value *float64, step float64) Res {
	return ctx.NumberEx(value, step, sliderFmt, OptAlignCenter)
}

func (ctx *Context) Header(label string) Res {
	return ctx.HeaderEx(label, 0)
}

func (ctx *Context) BeginTreeNode(label string) Res {
	return ctx.BeginTreeNodeEx(label, 0)
}

func (ctx *Context) BeginWindow(title string, rect image.Rectangle) Res {
	return ctx.BeginWindowEx(title, rect, 0)
}

func (ctx *Context) BeginPanel(name string) {
	ctx.BeginPanelEx(name, 0)
}
