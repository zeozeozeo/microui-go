package microui

import (
	"sort"
)

func expect(x bool) {
	if !x {
		panic("expect() failed")
	}
}

func mu_min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func mu_max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func mu_min_real(a, b Mu_Real) Mu_Real {
	if a < b {
		return a
	}
	return b
}

func mu_max_real(a, b Mu_Real) Mu_Real {
	if a > b {
		return a
	}
	return b
}

func mu_clamp(x, a, b int) int {
	return mu_min(b, mu_max(a, x))
}

func mu_clamp_real(x, a, b Mu_Real) Mu_Real {
	return mu_min_real(b, mu_max_real(a, x))
}

func NewVec2(x, y int) Vec2 {
	return Vec2{x, y}
}

func NewRect(x, y, w, h int) Rect {
	return Rect{x, y, w, h}
}

func NewColor(r, g, b, a uint8) Color {
	return Color{r, g, b, a}
}

func expand_rect(rect Rect, n int) Rect {
	return NewRect(rect.X-n, rect.Y-n, rect.W+n*2, rect.H+n*2)
}

func intersect_rects(r1, r2 Rect) Rect {
	var x1 int = mu_max(r1.X, r2.X)
	var y1 int = mu_max(r1.Y, r2.Y)
	var x2 int = mu_min(r1.X+r1.W, r2.X+r2.W)
	var y2 int = mu_min(r1.Y+r1.H, r2.Y+r2.H)
	if x2 < x1 {
		x2 = x1
	}
	if y2 < y1 {
		y2 = y1
	}
	return NewRect(x1, y1, x2-x1, y2-y1)
}

func rect_overlaps_vec2(r Rect, p Vec2) bool {
	return p.X >= r.X && p.X < r.X+r.W && p.Y >= r.Y && p.Y < r.Y+r.H
}

func hash(hash *mu_Id, data []byte) {
	for i := 0; i < len(data); i++ {
		*hash = (*hash ^ mu_Id(data[i])) * 16777619
	}
}

func (ctx *Context) GetID(data []byte) mu_Id {
	idx := len(ctx.IdStack)
	var res mu_Id
	if idx > 0 {
		res = ctx.IdStack[len(ctx.IdStack)-1]
	} else {
		res = HASH_INITIAL
	}
	hash(&res, data)
	ctx.LastID = res
	return res
}

func (ctx *Context) PushID(data []byte) {
	// push()
	ctx.IdStack = append(ctx.IdStack, ctx.GetID(data))
}

func (ctx *Context) PopID() {
	expect(len(ctx.IdStack) > 0)
	ctx.IdStack = ctx.IdStack[:len(ctx.IdStack)-1]
}

func (ctx *Context) PushClipRect(rect Rect) {
	last := ctx.GetClipRect()
	// push()
	ctx.ClipStack = append(ctx.ClipStack, intersect_rects(rect, last))
}

func (ctx *Context) PopClipRect() {
	expect(len(ctx.ClipStack) > 0)
	ctx.ClipStack = ctx.ClipStack[:len(ctx.ClipStack)-1]
}

func (ctx *Context) GetClipRect() Rect {
	expect(len(ctx.ClipStack) > 0)
	return ctx.ClipStack[len(ctx.ClipStack)-1]
}

func (ctx *Context) CheckClip(r Rect) int {
	cr := ctx.GetClipRect()
	if r.X > cr.X+cr.W || r.X+r.W < cr.X ||
		r.Y > cr.Y+cr.H || r.Y+r.H < cr.Y {
		return MU_CLIP_ALL
	}
	if r.X >= cr.X && r.X+r.W <= cr.X+cr.W &&
		r.Y >= cr.Y && r.Y+r.H <= cr.Y+cr.H {
		return 0
	}
	return MU_CLIP_PART
}

func (ctx *Context) GetLayout() *Layout {
	expect(len(ctx.LayoutStack) > 0)
	return &ctx.LayoutStack[len(ctx.LayoutStack)-1]
}

func (ctx *Context) PopContainer() {
	cnt := ctx.GetCurrentContainer()
	layout := ctx.GetLayout()
	cnt.ContentSize.X = layout.Max.X - layout.Body.X
	cnt.ContentSize.Y = layout.Max.Y - layout.Body.Y
	// pop container, layout and id
	// pop()
	expect(len(ctx.ContainerStack) > 0) // TODO: no expect in original impl
	ctx.ContainerStack = ctx.ContainerStack[:len(ctx.ContainerStack)-1]
	// pop()
	expect(len(ctx.LayoutStack) > 0) // TODO: no expect in original impl
	ctx.LayoutStack = ctx.LayoutStack[:len(ctx.LayoutStack)-1]
	ctx.PopID()
}

func (ctx *Context) GetCurrentContainer() *Container {
	expect(len(ctx.ContainerStack) > 0)
	return ctx.ContainerStack[len(ctx.ContainerStack)-1]
}

func (ctx *Context) getContainer(id mu_Id, opt int) *Container {
	var cnt *Container
	// try to get existing container from pool
	idx := ctx.PoolGet(ctx.ContainerPool[:], id)
	if idx >= 0 {
		if ctx.Containers[idx].Open || (^opt&MU_OPT_CLOSED) != 0 {
			ctx.PoolUpdate(ctx.ContainerPool[:], idx)
		}
		return &ctx.Containers[idx]
	}
	if (opt & MU_OPT_CLOSED) != 0 {
		return nil
	}
	// container not found in pool: init new container
	idx = ctx.PoolInit(ctx.ContainerPool[:], id)
	cnt = &ctx.Containers[idx]
	cnt.Clear()
	cnt.HeadIdx = -1
	cnt.TailIdx = -1
	cnt.Open = true
	ctx.BringToFront(cnt)
	return cnt
}

func (ctx *Context) GetContainer(name string) *Container {
	id := ctx.GetID([]byte(name))
	return ctx.getContainer(id, 0)
}

func (ctx *Context) BringToFront(cnt *Container) {
	ctx.LastZindex++
	cnt.Zindex = ctx.LastZindex
}

func (ctx *Context) SetFocus(id mu_Id) {
	ctx.Focus = id
	ctx.UpdatedFocus = true
}

func (ctx *Context) Begin() {
	expect(ctx.TextWidth != nil && ctx.TextHeight != nil)
	ctx.CommandList = ctx.CommandList[:0]
	ctx.RootList = ctx.RootList[:0]
	ctx.ScrollTarget = nil
	ctx.HoverRoot = ctx.NextHoverRoot
	ctx.NextHoverRoot = nil
	ctx.MouseDelta.X = ctx.MousePos.X - ctx.lastMousePos.X
	ctx.MouseDelta.Y = ctx.MousePos.Y - ctx.lastMousePos.Y
	ctx.Frame++
}

func (ctx *Context) End() {
	// check stacks
	expect(len(ctx.ContainerStack) == 0)
	expect(len(ctx.ClipStack) == 0)
	expect(len(ctx.IdStack) == 0)
	expect(len(ctx.LayoutStack) == 0)

	// handle scroll input
	if ctx.ScrollTarget != nil {
		ctx.ScrollTarget.Scroll.X += ctx.ScrollDelta.X
		ctx.ScrollTarget.Scroll.Y += ctx.ScrollDelta.Y
	}

	// unset focus if focus id was not touched this frame
	if !ctx.UpdatedFocus {
		ctx.Focus = 0
	}
	ctx.UpdatedFocus = false

	// bring hover root to front if mouse was pressed
	if ctx.MousePressed != 0 && ctx.NextHoverRoot != nil &&
		ctx.NextHoverRoot.Zindex < ctx.LastZindex &&
		ctx.NextHoverRoot.Zindex >= 0 {
		ctx.BringToFront(ctx.NextHoverRoot)
	}

	// reset input state
	ctx.KeyPressed = 0
	ctx.TextInput = nil
	ctx.MousePressed = 0
	ctx.ScrollDelta = NewVec2(0, 0)
	ctx.lastMousePos = ctx.MousePos

	// sort root containers by zindex
	// TODO (port): i'm not sure if this works
	sort.SliceStable(ctx.RootList, func(i, j int) bool {
		return ctx.RootList[i].Zindex < ctx.RootList[j].Zindex
	})

	// set root container jump commands
	for i := 0; i < len(ctx.RootList); i++ {
		cnt := ctx.RootList[i]
		// if this is the first container then make the first command jump to it.
		// otherwise set the previous container's tail to jump to this one
		if i == 0 {
			cmd := ctx.CommandList[0]
			expect(cmd.Type == MU_COMMAND_JUMP)
			cmd.Jump.DstIdx = cnt.HeadIdx + 1
			expect(cmd.Jump.DstIdx < MU_COMMANDLIST_SIZE)
		} else {
			prev := ctx.RootList[i-1]
			ctx.CommandList[prev.TailIdx].Jump.DstIdx = cnt.HeadIdx + 1
		}
		// make the last container's tail jump to the end of command list
		if i == len(ctx.RootList)-1 {
			ctx.CommandList[cnt.TailIdx].Jump.DstIdx = len(ctx.CommandList)
		}
	}
}
