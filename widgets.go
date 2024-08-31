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

func (c *Context) TreeNode(label string, f func(res Res)) {
	if res := c.beginTreeNode(label); res != 0 {
		defer c.endTreeNode()
		f(res)
	}
}

func (ctx *Context) beginTreeNode(label string) Res {
	return ctx.beginTreeNodeEx(label, 0)
}

func (ctx *Context) Window(title string, rect image.Rectangle, f func(res Res)) {
	if res := ctx.beginWindow(title, rect); res != 0 {
		defer ctx.endWindow()
		f(res)
	}
}

func (ctx *Context) beginWindow(title string, rect image.Rectangle) Res {
	return ctx.beginWindowEx(title, rect, 0)
}

func (c *Context) Panel(name string, f func()) {
	c.beginPanel(name)
	defer c.endPanel()
	f()
}

func (ctx *Context) beginPanel(name string) {
	ctx.beginPanelEx(name, 0)
}
