// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2024 The Ebitengine Authors

package microui

import "image"

func (c *Context) pushLayout(body image.Rectangle, scroll image.Point) {
	layout := Layout{}
	layout.Body = body.Sub(scroll)
	layout.Max = image.Pt(-0x1000000, -0x1000000)

	// push()
	c.layoutStack = append(c.layoutStack, layout)

	c.LayoutRow(1, []int{0}, 0)
}

func (c *Context) LayoutBeginColumn() {
	c.pushLayout(c.LayoutNext(), image.Pt(0, 0))
}

func (c *Context) LayoutEndColumn() {
	b := c.Layout()
	// pop()
	c.layoutStack = c.layoutStack[:len(c.layoutStack)-1]
	// inherit position/next_row/max from child layout if they are greater
	a := c.Layout()
	a.Position.X = max(a.Position.X, b.Position.X+b.Body.Min.X-a.Body.Min.X)
	a.NextRow = max(a.NextRow, b.NextRow+b.Body.Min.Y-a.Body.Min.Y)
	a.Max.X = max(a.Max.X, b.Max.X)
	a.Max.Y = max(a.Max.Y, b.Max.Y)
}

func (c *Context) LayoutRow(items int, widths []int, height int) {
	layout := c.Layout()

	expect(len(widths) <= maxWidths)
	copy(layout.Widths[:], widths)

	layout.Items = items
	layout.Position = image.Pt(layout.Indent, layout.NextRow)
	layout.Size.Y = height
	layout.ItemIndex = 0
}

// LayoutWidth sets layout size.x
func (c *Context) LayoutWidth(width int) {
	c.Layout().Size.X = width
}

// LayoutHeight sets layout size.y
func (c *Context) LayoutHeight(height int) {
	c.Layout().Size.Y = height
}

func (c *Context) LayoutNext() image.Rectangle {
	layout := c.Layout()
	style := c.Style
	var res image.Rectangle

	// handle next row
	if layout.ItemIndex == layout.Items {
		c.LayoutRow(layout.Items, nil, layout.Size.Y)
	}

	// position
	res = image.Rect(layout.Position.X, layout.Position.Y, layout.Position.X+res.Dx(), layout.Position.Y+res.Dy())

	// size
	if layout.Items > 0 {
		res.Max.X = res.Min.X + layout.Widths[layout.ItemIndex]
	} else {
		res.Max.X = res.Min.X + layout.Size.X
	}
	res.Max.Y = res.Min.Y + layout.Size.Y
	if res.Dx() == 0 {
		res.Max.X = res.Min.X + style.Size.X + style.Padding*2
	}
	if res.Dy() == 0 {
		res.Max.Y = res.Min.Y + style.Size.Y + style.Padding*2
	}
	if res.Dx() < 0 {
		res.Max.X += layout.Body.Dx() - res.Min.X + 1
	}
	if res.Dy() < 0 {
		res.Max.Y += layout.Body.Dy() - res.Min.Y + 1
	}

	layout.ItemIndex++

	// update position
	layout.Position.X += res.Dx() + style.Spacing
	layout.NextRow = max(layout.NextRow, res.Max.Y+style.Spacing)

	// apply body offset
	res = res.Add(layout.Body.Min)

	// update max position
	layout.Max.X = max(layout.Max.X, res.Max.X)
	layout.Max.Y = max(layout.Max.Y, res.Max.Y)

	c.LastRect = res
	return c.LastRect
}
