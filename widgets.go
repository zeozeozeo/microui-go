// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2024 The Ebitengine Authors

package microui

import "image"

func (c *Context) Button(label string) Res {
	return c.ButtonEx(label, 0, OptAlignCenter)
}

func (c *Context) TextBox(buf *string) Res {
	return c.TextBoxEx(buf, 0)
}

func (c *Context) Slider(value *float64, lo, hi float64) Res {
	return c.SliderEx(value, lo, hi, 0, sliderFmt, OptAlignCenter)
}

func (c *Context) Number(value *float64, step float64) Res {
	return c.NumberEx(value, step, sliderFmt, OptAlignCenter)
}

func (c *Context) Header(label string) Res {
	return c.HeaderEx(label, 0)
}

func (c *Context) TreeNode(label string, f func(res Res)) {
	if res := c.beginTreeNode(label); res != 0 {
		defer c.endTreeNode()
		f(res)
	}
}

func (c *Context) beginTreeNode(label string) Res {
	return c.beginTreeNodeEx(label, 0)
}

func (c *Context) Window(title string, rect image.Rectangle, f func(res Res)) {
	if res := c.beginWindow(title, rect); res != 0 {
		defer c.endWindow()
		f(res)
	}
}

func (c *Context) beginWindow(title string, rect image.Rectangle) Res {
	return c.beginWindowEx(title, rect, 0)
}

func (c *Context) Panel(name string, f func()) {
	c.beginPanel(name)
	defer c.endPanel()
	f()
}

func (c *Context) beginPanel(name string) {
	c.beginPanelEx(name, 0)
}
