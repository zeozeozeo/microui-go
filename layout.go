// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2024 The Ebitengine Authors

package microui

import "image"

func (c *Context) pushLayout(body image.Rectangle, scroll image.Point) {
	layout := layout{}
	layout.body = body.Sub(scroll)
	layout.max = image.Pt(-0x1000000, -0x1000000)

	// push()
	c.layoutStack = append(c.layoutStack, layout)

	c.LayoutRow(1, []int{0}, 0)
}

func (c *Context) LayoutColumn(f func()) {
	c.layoutBeginColumn()
	defer c.layoutEndColumn()
	f()
}

func (c *Context) layoutBeginColumn() {
	c.pushLayout(c.LayoutNext(), image.Pt(0, 0))
}

func (c *Context) layoutEndColumn() {
	b := c.layout()
	// pop()
	c.layoutStack = c.layoutStack[:len(c.layoutStack)-1]
	// inherit position/next_row/max from child layout if they are greater
	a := c.layout()
	a.position.X = max(a.position.X, b.position.X+b.body.Min.X-a.body.Min.X)
	a.nextRow = max(a.nextRow, b.nextRow+b.body.Min.Y-a.body.Min.Y)
	a.max.X = max(a.max.X, b.max.X)
	a.max.Y = max(a.max.Y, b.max.Y)
}

func (c *Context) LayoutRow(items int, widths []int, height int) {
	layout := c.layout()

	expect(len(widths) <= maxWidths)
	copy(layout.widths[:], widths)

	layout.items = items
	layout.position = image.Pt(layout.indent, layout.nextRow)
	layout.size.Y = height
	layout.itemIndex = 0
}

// LayoutWidth sets layout size.x
func (c *Context) LayoutWidth(width int) {
	c.layout().size.X = width
}

// LayoutHeight sets layout size.y
func (c *Context) LayoutHeight(height int) {
	c.layout().size.Y = height
}

func (c *Context) LayoutNext() image.Rectangle {
	layout := c.layout()
	style := c.Style
	var res image.Rectangle

	// handle next row
	if layout.itemIndex == layout.items {
		c.LayoutRow(layout.items, nil, layout.size.Y)
	}

	// position
	res = image.Rect(layout.position.X, layout.position.Y, layout.position.X+res.Dx(), layout.position.Y+res.Dy())

	// size
	if layout.items > 0 {
		res.Max.X = res.Min.X + layout.widths[layout.itemIndex]
	} else {
		res.Max.X = res.Min.X + layout.size.X
	}
	res.Max.Y = res.Min.Y + layout.size.Y
	if res.Dx() == 0 {
		res.Max.X = res.Min.X + style.Size.X + style.Padding*2
	}
	if res.Dy() == 0 {
		res.Max.Y = res.Min.Y + style.Size.Y + style.Padding*2
	}
	if res.Dx() < 0 {
		res.Max.X += layout.body.Dx() - res.Min.X + 1
	}
	if res.Dy() < 0 {
		res.Max.Y += layout.body.Dy() - res.Min.Y + 1
	}

	layout.itemIndex++

	// update position
	layout.position.X += res.Dx() + style.Spacing
	layout.nextRow = max(layout.nextRow, res.Max.Y+style.Spacing)

	// apply body offset
	res = res.Add(layout.body.Min)

	// update max position
	layout.max.X = max(layout.max.X, res.Max.X)
	layout.max.Y = max(layout.max.Y, res.Max.Y)

	c.LastRect = res
	return c.LastRect
}
