// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2024 The Ebitengine Authors

package microui

import "image"

/*============================================================================
** layout
**============================================================================*/

func (ctx *Context) PushLayout(body image.Rectangle, scroll image.Point) {
	layout := Layout{}
	layout.Body = body.Sub(scroll)
	layout.Max = image.Pt(-0x1000000, -0x1000000)

	// push()
	ctx.LayoutStack = append(ctx.LayoutStack, layout)

	ctx.LayoutRow(1, []int{0}, 0)
}

func (ctx *Context) LayoutBeginColumn() {
	ctx.PushLayout(ctx.LayoutNext(), image.Pt(0, 0))
}

func (ctx *Context) LayoutEndColumn() {
	b := ctx.GetLayout()
	// pop()
	expect(len(ctx.LayoutStack) > 0)
	ctx.LayoutStack = ctx.LayoutStack[:len(ctx.LayoutStack)-1]
	// inherit position/next_row/max from child layout if they are greater
	a := ctx.GetLayout()
	a.Position.X = mu_max(a.Position.X, b.Position.X+b.Body.Min.X-a.Body.Min.X)
	a.NextRow = mu_max(a.NextRow, b.NextRow+b.Body.Min.Y-a.Body.Min.Y)
	a.Max.X = mu_max(a.Max.X, b.Max.X)
	a.Max.Y = mu_max(a.Max.Y, b.Max.Y)
}

func (ctx *Context) LayoutRow(items int, widths []int, height int) {
	layout := ctx.GetLayout()

	expect(len(widths) <= maxWidths)
	copy(layout.Widths[:], widths)

	layout.Items = items
	layout.Position = image.Pt(layout.Indent, layout.NextRow)
	layout.Size.Y = height
	layout.ItemIndex = 0
}

// sets layout size.x
func (ctx *Context) LayoutWidth(width int) {
	ctx.GetLayout().Size.X = width
}

// sets layout size.y
func (ctx *Context) LayoutHeight(height int) {
	ctx.GetLayout().Size.Y = height
}

func (ctx *Context) LayoutSetNext(r image.Rectangle, relative bool) {
	layout := ctx.GetLayout()
	layout.Next = r
	if relative {
		layout.NextType = Relative
	} else {
		layout.NextType = Absolute
	}
}

func (ctx *Context) LayoutNext() image.Rectangle {
	layout := ctx.GetLayout()
	style := ctx.Style
	var res muRect

	if layout.NextType != 0 {
		// handle rect set by `mu_layout_set_next`
		next_type := layout.NextType
		layout.NextType = 0
		res = rectFromRectangle(layout.Next)

		if next_type == Absolute {
			ctx.LastRect = res.rectangle()
			return ctx.LastRect
		}
	} else {
		// handle next row
		if layout.ItemIndex == layout.Items {
			ctx.LayoutRow(layout.Items, nil, layout.Size.Y)
		}

		// position
		res.X = layout.Position.X
		res.Y = layout.Position.Y

		// size
		if layout.Items > 0 {
			res.W = layout.Widths[layout.ItemIndex]
		} else {
			res.W = layout.Size.X
		}
		res.H = layout.Size.Y
		if res.W == 0 {
			res.W = style.Size.X + style.Padding*2
		}
		if res.H == 0 {
			res.H = style.Size.Y + style.Padding*2
		}
		if res.W < 0 {
			res.W += layout.Body.Dx() - res.X + 1
		}
		if res.H < 0 {
			res.H += layout.Body.Dy() - res.Y + 1
		}

		layout.ItemIndex++
	}

	// update position
	layout.Position.X += res.W + style.Spacing
	layout.NextRow = mu_max(layout.NextRow, res.Y+res.H+style.Spacing)

	// apply body offset
	res.X += layout.Body.Min.X
	res.Y += layout.Body.Min.Y

	// update max position
	layout.Max.X = mu_max(layout.Max.X, res.X+res.W)
	layout.Max.Y = mu_max(layout.Max.Y, res.Y+res.H)

	ctx.LastRect = res.rectangle()
	return ctx.LastRect
}
