// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2024 The Ebitengine Authors

package microui

import "image"

func (ctx *Context) Button(label string) bool {
	return ctx.ButtonEx(label, 0, OptAlignCenter) != 0
}

func (ctx *Context) TextBox(buf *string) int {
	return ctx.TextBoxEx(buf, 0)
}

func (ctx *Context) Slider(value *float64, lo, hi float64) int {
	return ctx.SliderEx(value, lo, hi, 0, sliderFmt, OptAlignCenter)
}

func (ctx *Context) Number(value *float64, step float64) int {
	return ctx.NumberEx(value, step, sliderFmt, OptAlignCenter)
}

func (ctx *Context) Header(label string) bool {
	return ctx.HeaderEx(label, 0) != 0
}

func (ctx *Context) BeginTreeNode(label string) bool {
	return ctx.BeginTreeNodeEx(label, 0) != 0
}

func (ctx *Context) BeginWindow(title string, rect image.Rectangle) bool {
	return ctx.BeginWindowEx(title, rect, 0) != 0
}

func (ctx *Context) BeginPanel(name string) {
	ctx.BeginPanelEx(name, 0)
}
