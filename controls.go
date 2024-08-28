// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2024 The Ebitengine Authors

package microui

import (
	"fmt"
	"image"
	"strconv"
	"unsafe"
)

/*============================================================================
** controls
**============================================================================*/

func (ctx *Context) InHoverRoot() bool {
	for i := len(ctx.ContainerStack) - 1; i >= 0; i-- {
		if ctx.ContainerStack[i] == ctx.HoverRoot {
			return true
		}
		// only root containers have their `head` field set; stop searching if we've
		// reached the current root container
		if ctx.ContainerStack[i].HeadIdx >= 0 {
			break
		}
	}
	return false
}

func (ctx *Context) DrawControlFrame(id mu_Id, rect Rect, colorid int, opt int) {
	if (opt & OptNoFrame) != 0 {
		return
	}
	if ctx.Focus == id {
		colorid += 2
	} else if ctx.Hover == id {
		colorid++
	}
	ctx.DrawFrame(ctx, rect, colorid)
}

func (ctx *Context) DrawControlText(str string, rect Rect, colorid int, opt int) {
	var pos image.Point
	font := ctx.Style.Font
	tw := ctx.TextWidth(font, str)
	ctx.PushClipRect(rect)
	pos.Y = rect.Y + (rect.H-ctx.TextHeight(font))/2
	if (opt & OptAlignCenter) != 0 {
		pos.X = rect.X + (rect.W-tw)/2
	} else if (opt & OptAlignRight) != 0 {
		pos.X = rect.X + rect.W - tw - ctx.Style.Padding
	} else {
		pos.X = rect.X + ctx.Style.Padding
	}
	ctx.DrawText(font, str, pos, ctx.Style.Colors[colorid])
	ctx.PopClipRect()
}

func (ctx *Context) MouseOver(rect Rect) bool {
	return rect_overlaps_vec2(rect, ctx.MousePos) &&
		rect_overlaps_vec2(ctx.GetClipRect(), ctx.MousePos) &&
		ctx.InHoverRoot()
}

func (ctx *Context) UpdateControl(id mu_Id, rect Rect, opt int) {
	mouseover := ctx.MouseOver(rect)

	if ctx.Focus == id {
		ctx.UpdatedFocus = true
	}
	if (opt & OptNoInteract) != 0 {
		return
	}
	if mouseover && ctx.MouseDown == 0 {
		ctx.Hover = id
	}

	if ctx.Focus == id {
		if ctx.MousePressed != 0 && !mouseover {
			ctx.SetFocus(0)
		}
		if ctx.MouseDown == 0 && (^opt&OptHoldFocus) != 0 {
			ctx.SetFocus(0)
		}
	}

	if ctx.Hover == id {
		if ctx.MousePressed != 0 {
			ctx.SetFocus(id)
		} else if !mouseover {
			ctx.Hover = 0
		}
	}
}

func (ctx *Context) Text(text string) {
	var start_idx, end_idx, p int
	font := ctx.Style.Font
	color := ctx.Style.Colors[ColorText]
	ctx.LayoutBeginColumn()
	ctx.LayoutRow(1, []int{-1}, ctx.TextHeight(font))
	for end_idx < len(text) {
		r := ctx.LayoutNext()
		w := 0
		end_idx = p
		start_idx = end_idx
		for end_idx < len(text) && text[end_idx] != '\n' {
			word := p
			for p < len(text) && text[p] != ' ' && text[p] != '\n' {
				p++
			}
			w += ctx.TextWidth(font, text[word:p])
			if w > r.W && end_idx != start_idx {
				break
			}
			if p < len(text) {
				w += ctx.TextWidth(font, string(text[p]))
			}
			end_idx = p
			p++
		}
		ctx.DrawText(font, text[start_idx:end_idx], image.Pt(r.X, r.Y), color)
		p = end_idx + 1
	}
	ctx.LayoutEndColumn()
}

func (ctx *Context) Label(text string) {
	ctx.DrawControlText(text, ctx.LayoutNext(), ColorText, 0)
}

func (ctx *Context) ButtonEx(label string, icon int, opt int) int {
	var res int = 0
	var id mu_Id
	if len(label) > 0 {
		id = ctx.GetID([]byte(label))
	} else {
		iconPtr := &icon
		// TODO: investigate if this okay, if icon represents an icon ID we might need
		// to refer to the value instead of a pointer, like commented below:
		// unsafe.Slice((*byte)(unsafe.Pointer(&icon)), unsafe.Sizeof(icon)))
		id = ctx.GetID(PtrToBytes(unsafe.Pointer(iconPtr)))
	}
	r := ctx.LayoutNext()
	ctx.UpdateControl(id, r, opt)
	// handle click
	if ctx.MousePressed == mouseLeft && ctx.Focus == id {
		res |= ResSubmit
	}
	// draw
	ctx.DrawControlFrame(id, r, ColorButton, opt)
	if len(label) > 0 {
		ctx.DrawControlText(label, r, ColorText, opt)
	}
	if icon != 0 {
		ctx.DrawIcon(icon, r, ctx.Style.Colors[ColorText])
	}
	return res
}

func (ctx *Context) Checkbox(label string, state *bool) int {
	var res int = 0
	id := ctx.GetID(PtrToBytes(unsafe.Pointer(state)))
	r := ctx.LayoutNext()
	box := NewRect(r.X, r.Y, r.H, r.H)
	ctx.UpdateControl(id, r, 0)
	// handle click
	if ctx.MousePressed == mouseLeft && ctx.Focus == id {
		res |= ResChange
		*state = !*state
	}
	// draw
	ctx.DrawControlFrame(id, box, ColorBase, 0)
	if *state {
		ctx.DrawIcon(IconCheck, box, ctx.Style.Colors[ColorText])
	}
	r = NewRect(r.X+box.W, r.Y, r.W-box.W, r.H)
	ctx.DrawControlText(label, r, ColorText, 0)
	return res
}

func (ctx *Context) TextboxRaw(buf *string, id mu_Id, r Rect, opt int) int {
	var res int = 0
	ctx.UpdateControl(id, r, opt|OptHoldFocus)
	buflen := len(*buf)

	if ctx.Focus == id {
		// handle text input
		if len(ctx.TextInput) > 0 {
			*buf += string(ctx.TextInput)
			res |= ResChange
		}
		// handle backspace
		if (ctx.KeyPressed&MU_KEY_BACKSPACE) != 0 && buflen > 0 {
			*buf = (*buf)[:buflen-1]
			res |= ResChange
		}
		// handle return
		if (ctx.KeyPressed & MU_KEY_RETURN) != 0 {
			ctx.SetFocus(0)
			res |= ResSubmit
		}
	}

	// draw
	ctx.DrawControlFrame(id, r, ColorBase, opt)
	if ctx.Focus == id {
		color := ctx.Style.Colors[ColorText]
		font := ctx.Style.Font
		textw := ctx.TextWidth(font, *buf)
		texth := ctx.TextHeight(font)
		ofx := r.W - ctx.Style.Padding - textw - 1
		textx := r.X + mu_min(ofx, ctx.Style.Padding)
		texty := r.Y + (r.H-texth)/2
		ctx.PushClipRect(r)
		ctx.DrawText(font, *buf, image.Pt(textx, texty), color)
		ctx.DrawRect(NewRect(textx+textw, texty, 1, texth), color)
		ctx.PopClipRect()
	} else {
		ctx.DrawControlText(*buf, r, ColorText, opt)
	}

	return res
}

func (ctx *Context) NumberTextBox(value *float32, r Rect, id mu_Id) bool {
	if ctx.MousePressed == mouseLeft && (ctx.KeyDown&MU_KEY_SHIFT) != 0 &&
		ctx.Hover == id {
		ctx.NumberEdit = id
		ctx.NumberEditBuf = fmt.Sprintf(MU_REAL_FMT, *value)
	}
	if ctx.NumberEdit == id {
		res := ctx.TextboxRaw(&ctx.NumberEditBuf, id, r, 0)
		if (res&ResSubmit) != 0 || ctx.Focus != id {
			nval, err := strconv.ParseFloat(ctx.NumberEditBuf, 32)
			if err != nil {
				nval = 0
			}
			*value = float32(nval)
			ctx.NumberEdit = 0
		} else {
			return true
		}
	}
	return false
}

func (ctx *Context) TextBoxEx(buf *string, opt int) int {
	id := ctx.GetID(PtrToBytes(unsafe.Pointer(buf)))
	r := ctx.LayoutNext()
	return ctx.TextboxRaw(buf, id, r, opt)
}

func (ctx *Context) SliderEx(value *float32, low, high, step float32, format string, opt int) int {
	var thumb Rect
	var x, w, res int = 0, 0, 0
	last := *value
	v := last
	id := ctx.GetID(PtrToBytes(unsafe.Pointer(value)))
	base := ctx.LayoutNext()

	// handle text input mode
	if ctx.NumberTextBox(&v, base, id) {
		return res
	}

	// handle normal mode
	ctx.UpdateControl(id, base, opt)

	// handle input
	if ctx.Focus == id && (ctx.MouseDown|ctx.MousePressed) == mouseLeft {
		v = low + float32(ctx.MousePos.X-base.X)*(high-low)/float32(base.W)
		if step != 0 {
			v = ((v + step/2) / step) * step
		}
	}
	// clamp and store value, update res
	*value = mu_clamp_real(v, low, high)
	if last != v {
		res |= ResChange
	}

	// draw base
	ctx.DrawControlFrame(id, base, ColorBase, opt)
	// draw thumb
	w = ctx.Style.ThumbSize
	x = int((v - low) * float32(base.W-w) / (high - low))
	thumb = NewRect(base.X+x, base.Y, w, base.H)
	ctx.DrawControlFrame(id, thumb, ColorButton, opt)
	// draw text
	text := fmt.Sprintf(format, v)
	ctx.DrawControlText(text, base, ColorText, opt)

	return res
}

func (ctx *Context) NumberEx(value *float32, step float32, format string, opt int) int {
	var res int = 0
	id := ctx.GetID(PtrToBytes(unsafe.Pointer(&value)))
	base := ctx.LayoutNext()
	last := *value

	// handle text input mode
	if ctx.NumberTextBox(value, base, id) {
		return res
	}

	// handle normal mode
	ctx.UpdateControl(id, base, opt)

	// handle input
	if ctx.Focus == id && ctx.MouseDown == mouseLeft {
		*value += float32(ctx.MouseDelta.X) * step
	}
	// set flag if value changed
	if *value != last {
		res |= ResChange
	}

	// draw base
	ctx.DrawControlFrame(id, base, ColorBase, opt)
	// draw text
	text := fmt.Sprintf(format, *value)
	ctx.DrawControlText(text, base, ColorText, opt)

	return res
}

func (ctx *Context) MuHeader(label string, istreenode bool, opt int) int {
	var r Rect
	var active, expanded bool
	id := ctx.GetID([]byte(label))
	idx := ctx.PoolGet(ctx.TreeNodePool[:], id)
	ctx.LayoutRow(1, []int{-1}, 0)

	active = idx >= 0
	if (opt & OptExpanded) != 0 {
		expanded = !active
	} else {
		expanded = active
	}
	r = ctx.LayoutNext()
	ctx.UpdateControl(id, r, 0)

	// handle click (TODO (port): check if this is correct)
	clicked := ctx.MousePressed == mouseLeft && ctx.Focus == id
	v1, v2 := 0, 0
	if active {
		v1 = 1
	}
	if clicked {
		v2 = 1
	}
	active = (v1 ^ v2) == 1

	// update pool ref
	if idx >= 0 {
		if active {
			ctx.PoolUpdate(ctx.TreeNodePool[:], idx)
		} else {
			ctx.TreeNodePool[idx] = MuPoolItem{}
		}
	} else if active {
		ctx.PoolInit(ctx.TreeNodePool[:], id)
	}

	// draw
	if istreenode {
		if ctx.Hover == id {
			ctx.DrawFrame(ctx, r, ColorButtonHover)
		}
	} else {
		ctx.DrawControlFrame(id, r, ColorButton, 0)
	}
	var icon_id int
	if expanded {
		icon_id = IconExpanded
	} else {
		icon_id = IconCollapsed
	}
	ctx.DrawIcon(
		icon_id,
		NewRect(r.X, r.Y, r.H, r.H),
		ctx.Style.Colors[ColorText],
	)
	r.X += r.H - ctx.Style.Padding
	r.W -= r.H - ctx.Style.Padding
	ctx.DrawControlText(label, r, ColorText, 0)

	if expanded {
		return ResActive
	}
	return 0
}

func (ctx *Context) HeaderEx(label string, opt int) int {
	return ctx.MuHeader(label, false, opt)
}

func (ctx *Context) BeginTreeNodeEx(label string, opt int) int {
	res := ctx.MuHeader(label, true, opt)
	if (res & ResActive) != 0 {
		ctx.GetLayout().Indent += ctx.Style.Indent
		// push()
		ctx.IdStack = append(ctx.IdStack, ctx.LastID)
	}
	return res
}

func (ctx *Context) EndTreeNode() {
	ctx.GetLayout().Indent -= ctx.Style.Indent
	ctx.PopID()
}

// x = x, y = y, w = w, h = h
func (ctx *Context) scrollbarVertical(cnt *Container, b *Rect, cs image.Point) {
	maxscroll := cs.Y - b.H
	if maxscroll > 0 && b.H > 0 {
		var base, thumb Rect
		id := ctx.GetID([]byte("!scrollbar" + "y"))

		// get sizing / positioning
		base = *b
		base.X = b.X + b.W
		base.W = ctx.Style.ScrollbarSize

		// handle input
		ctx.UpdateControl(id, base, 0)
		if ctx.Focus == id && ctx.MouseDown == mouseLeft {
			cnt.Scroll.Y += ctx.MouseDelta.Y * cs.Y / base.H
		}
		// clamp scroll to limits
		cnt.Scroll.Y = mu_clamp(cnt.Scroll.Y, 0, maxscroll)

		// draw base and thumb
		ctx.DrawFrame(ctx, base, ColorScrollBase)
		thumb = base
		thumb.H = mu_max(ctx.Style.ThumbSize, base.H*b.H/cs.Y)
		thumb.Y += cnt.Scroll.Y * (base.H - thumb.H) / maxscroll
		ctx.DrawFrame(ctx, thumb, ColorScrollThumb)

		// set this as the scroll_target (will get scrolled on mousewheel)
		// if the mouse is over it
		if ctx.MouseOver(*b) {
			ctx.ScrollTarget = cnt
		}
	} else {
		cnt.Scroll.Y = 0
	}
}

// x = y, y = x, w = h, h = w
func (ctx *Context) scrollbarHorizontal(cnt *Container, b *Rect, cs image.Point) {
	maxscroll := cs.X - b.W
	if maxscroll > 0 && b.W > 0 {
		var base, thumb Rect
		id := ctx.GetID([]byte("!scrollbar" + "x"))

		// get sizing / positioning
		base = *b
		base.Y = b.Y + b.H
		base.H = ctx.Style.ScrollbarSize

		// handle input
		ctx.UpdateControl(id, base, 0)
		if ctx.Focus == id && ctx.MouseDown == mouseLeft {
			cnt.Scroll.X += ctx.MouseDelta.X * cs.X / base.W
		}
		// clamp scroll to limits
		cnt.Scroll.X = mu_clamp(cnt.Scroll.X, 0, maxscroll)

		// draw base and thumb
		ctx.DrawFrame(ctx, base, ColorScrollBase)
		thumb = base
		thumb.W = mu_max(ctx.Style.ThumbSize, base.W*b.W/cs.X)
		thumb.X += cnt.Scroll.X * (base.W - thumb.W) / maxscroll
		ctx.DrawFrame(ctx, thumb, ColorScrollThumb)

		// set this as the scroll_target (will get scrolled on mousewheel)
		// if the mouse is over it
		if ctx.MouseOver(*b) {
			ctx.ScrollTarget = cnt
		}
	} else {
		cnt.Scroll.X = 0
	}
}

// if `swap` is true, X = Y, Y = X, W = H, H = W
func (ctx *Context) AddScrollbar(cnt *Container, b *Rect, cs image.Point, swap bool) {
	if swap {
		ctx.scrollbarHorizontal(cnt, b, cs)
	} else {
		ctx.scrollbarVertical(cnt, b, cs)
	}
}

func (ctx *Context) Scrollbars(cnt *Container, body *Rect) {
	sz := ctx.Style.ScrollbarSize
	cs := cnt.ContentSize
	cs.X += ctx.Style.Padding * 2
	cs.Y += ctx.Style.Padding * 2
	ctx.PushClipRect(*body)
	// resize body to make room for scrollbars
	if cs.Y > cnt.Body.H {
		body.W -= sz
	}
	if cs.X > cnt.Body.W {
		body.H -= sz
	}
	// to create a horizontal or vertical scrollbar almost-identical code is
	// used; only the references to `x|y` `w|h` need to be switched
	ctx.AddScrollbar(cnt, body, cs, false)
	ctx.AddScrollbar(cnt, body, cs, true)
	ctx.PopClipRect()
}

func (ctx *Context) PushContainerBody(cnt *Container, body Rect, opt int) {
	if (^opt & OptNoScroll) != 0 {
		ctx.Scrollbars(cnt, &body)
	}
	ctx.PushLayout(expand_rect(body, -ctx.Style.Padding), cnt.Scroll)
	cnt.Body = body
}

func (ctx *Context) BeginRootContainer(cnt *Container) {
	// push()
	ctx.ContainerStack = append(ctx.ContainerStack, cnt)
	// push container to roots list and push head command
	// push()
	ctx.RootList = append(ctx.RootList, cnt)
	cnt.HeadIdx = ctx.PushJump(-1)
	// set as hover root if the mouse is overlapping this container and it has a
	// higher zindex than the current hover root
	if rect_overlaps_vec2(cnt.Rect, ctx.MousePos) &&
		(ctx.NextHoverRoot == nil || cnt.Zindex > ctx.NextHoverRoot.Zindex) {
		ctx.NextHoverRoot = cnt
	}
	// clipping is reset here in case a root-container is made within
	// another root-containers's begin/end block; this prevents the inner
	// root-container being clipped to the outer
	// push()
	ctx.ClipStack = append(ctx.ClipStack, UnclippedRect)
}

func (ctx *Context) EndRootContainer() {
	// push tail 'goto' jump command and set head 'skip' command. the final steps
	// on initing these are done in mu_end()
	cnt := ctx.GetCurrentContainer()
	cnt.TailIdx = ctx.PushJump(-1)
	ctx.CommandList[cnt.HeadIdx].Jump.DstIdx = len(ctx.CommandList) //- 1
	// pop base clip rect and container
	ctx.PopClipRect()
	ctx.PopContainer()
}

func (ctx *Context) BeginWindowEx(title string, rect Rect, opt int) int {
	var body Rect
	id := ctx.GetID([]byte(title))
	cnt := ctx.getContainer(id, opt)
	if cnt == nil || !cnt.Open {
		return 0
	}
	// push()
	ctx.IdStack = append(ctx.IdStack, id)

	if cnt.Rect.W == 0 {
		cnt.Rect = rect
	}
	ctx.BeginRootContainer(cnt)
	body = cnt.Rect
	rect = body

	// draw frame
	if (^opt & OptNoFrame) != 0 {
		ctx.DrawFrame(ctx, rect, ColorWindowBG)
	}

	// do title bar
	if (^opt & OptNoTitle) != 0 {
		tr := rect
		tr.H = ctx.Style.TitleHeight
		ctx.DrawFrame(ctx, tr, ColorTitleBG)

		// do title text
		if (^opt & OptNoTitle) != 0 {
			id := ctx.GetID([]byte("!title"))
			ctx.UpdateControl(id, tr, opt)
			ctx.DrawControlText(title, tr, ColorTitleText, opt)
			if id == ctx.Focus && ctx.MouseDown == mouseLeft {
				cnt.Rect.X += ctx.MouseDelta.X
				cnt.Rect.Y += ctx.MouseDelta.Y
			}
			body.Y += tr.H
			body.H -= tr.H
		}

		// do `close` button
		if (^opt & OptNoClose) != 0 {
			id := ctx.GetID([]byte("!close"))
			r := NewRect(tr.X+tr.W-tr.H, tr.Y, tr.H, tr.H)
			tr.W -= r.W
			ctx.DrawIcon(IconClose, r, ctx.Style.Colors[ColorTitleText])
			ctx.UpdateControl(id, r, opt)
			if ctx.MousePressed == mouseLeft && id == ctx.Focus {
				cnt.Open = false
			}
		}
	}

	ctx.PushContainerBody(cnt, body, opt)

	// do `resize` handle
	if (^opt & OptNoResize) != 0 {
		sz := ctx.Style.TitleHeight
		id := ctx.GetID([]byte("!resize"))
		r := NewRect(rect.X+rect.W-sz, rect.Y+rect.H-sz, sz, sz)
		ctx.UpdateControl(id, r, opt)
		if id == ctx.Focus && ctx.MouseDown == mouseLeft {
			cnt.Rect.W = mu_max(96, cnt.Rect.W+ctx.MouseDelta.X)
			cnt.Rect.H = mu_max(64, cnt.Rect.H+ctx.MouseDelta.Y)
		}
	}

	// resize to content size
	if (opt & OptAutoSize) != 0 {
		r := ctx.GetLayout().Body
		cnt.Rect.W = cnt.ContentSize.X + (cnt.Rect.W - r.W)
		cnt.Rect.H = cnt.ContentSize.Y + (cnt.Rect.H - r.H)
	}

	// close if this is a popup window and elsewhere was clicked
	if (opt&OptPopup) != 0 && ctx.MousePressed != 0 && ctx.HoverRoot != cnt {
		cnt.Open = false
	}

	ctx.PushClipRect(cnt.Body)
	return ResActive
}

func (ctx *Context) EndWindow() {
	ctx.PopClipRect()
	ctx.EndRootContainer()
}

func (ctx *Context) OpenPopup(name string) {
	cnt := ctx.GetContainer(name)
	// set as hover root so popup isn't closed in begin_window_ex()
	ctx.NextHoverRoot = cnt
	ctx.HoverRoot = ctx.NextHoverRoot
	// position at mouse cursor, open and bring-to-front
	cnt.Rect = NewRect(ctx.MousePos.X, ctx.MousePos.Y, 1, 1)
	cnt.Open = true
	ctx.BringToFront(cnt)
}

func (ctx *Context) BeginPopup(name string) int {
	opt := OptPopup | OptAutoSize | OptNoResize |
		OptNoScroll | OptNoTitle | OptClosed
	return ctx.BeginWindowEx(name, NewRect(0, 0, 0, 0), opt)
}

func (ctx *Context) EndPopup() {
	ctx.EndWindow()
}

func (ctx *Context) BeginPanelEx(name string, opt int) {
	var cnt *Container
	ctx.PushID([]byte(name))
	cnt = ctx.getContainer(ctx.LastID, opt)
	cnt.Rect = ctx.LayoutNext()
	if (^opt & OptNoFrame) != 0 {
		ctx.DrawFrame(ctx, cnt.Rect, ColorPanelBG)
	}
	// push()
	ctx.ContainerStack = append(ctx.ContainerStack, cnt)
	ctx.PushContainerBody(cnt, cnt.Rect, opt)
	ctx.PushClipRect(cnt.Body)
}

func (ctx *Context) EndPanel() {
	ctx.PopClipRect()
	ctx.PopContainer()
}
