// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2024 The Ebitengine Authors

package microui

import (
	"image"
	"sort"
	"unsafe"
)

func expect(x bool) {
	if !x {
		panic("expect() failed")
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func minF(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func maxF(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func clamp(x, a, b int) int {
	return min(b, max(a, x))
}

func clampF(x, a, b float64) float64 {
	return minF(b, maxF(a, x))
}

func hash(hash *ID, data []byte) {
	for i := 0; i < len(data); i++ {
		*hash = (*hash ^ ID(data[i])) * 16777619
	}
}

func ptrToBytes(ptr unsafe.Pointer) []byte {
	return unsafe.Slice((*byte)(unsafe.Pointer(&ptr)), unsafe.Sizeof(ptr))
}

// id returns a hash value based on the data and the last ID on the stack.
func (c *Context) id(data []byte) ID {
	const (
		hashInitial = 2166136261 // 32bit fnv-1a hash
	)

	idx := len(c.idStack)
	var res ID
	if idx > 0 {
		res = c.idStack[len(c.idStack)-1]
	} else {
		res = hashInitial
	}
	hash(&res, data)
	c.LastID = res
	return res
}

func (c *Context) pushID(data []byte) {
	// push()
	c.idStack = append(c.idStack, c.id(data))
}

func (c *Context) popID() {
	c.idStack = c.idStack[:len(c.idStack)-1]
}

func (c *Context) PushClipRect(rect image.Rectangle) {
	last := c.GetClipRect()
	// push()
	c.clipStack = append(c.clipStack, rect.Intersect(last))
}

func (c *Context) PopClipRect() {
	c.clipStack = c.clipStack[:len(c.clipStack)-1]
}

func (c *Context) GetClipRect() image.Rectangle {
	return c.clipStack[len(c.clipStack)-1]
}

func (c *Context) CheckClip(r image.Rectangle) int {
	cr := c.GetClipRect()
	if !r.Overlaps(cr) {
		return ClipAll
	}
	if r.In(cr) {
		return 0
	}
	return ClipPart
}

func (c *Context) GetLayout() *Layout {
	return &c.layoutStack[len(c.layoutStack)-1]
}

func (c *Context) popContainer() {
	cnt := c.GetCurrentContainer()
	layout := c.GetLayout()
	cnt.ContentSize.X = layout.Max.X - layout.Body.Min.X
	cnt.ContentSize.Y = layout.Max.Y - layout.Body.Min.Y
	// pop container, layout and id
	// pop()
	c.containerStack = c.containerStack[:len(c.containerStack)-1]
	// pop()
	c.layoutStack = c.layoutStack[:len(c.layoutStack)-1]
	c.popID()
}

func (c *Context) GetCurrentContainer() *Container {
	return c.containerStack[len(c.containerStack)-1]
}

func (c *Context) getContainer(id ID, opt Option) *Container {
	// try to get existing container from pool
	idx := c.poolGet(c.containerPool[:], id)
	if idx >= 0 {
		if c.containers[idx].Open || (^opt&OptClosed) != 0 {
			c.poolUpdate(c.containerPool[:], idx)
		}
		return &c.containers[idx]
	}
	if (opt & OptClosed) != 0 {
		return nil
	}
	// container not found in pool: init new container
	idx = c.poolInit(c.containerPool[:], id)
	cnt := &c.containers[idx]
	*cnt = Container{}
	cnt.HeadIdx = -1
	cnt.TailIdx = -1
	cnt.Open = true
	c.BringToFront(cnt)
	return cnt
}

func (c *Context) GetContainer(name string) *Container {
	id := c.id([]byte(name))
	return c.getContainer(id, 0)
}

func (c *Context) BringToFront(cnt *Container) {
	c.LastZindex++
	cnt.Zindex = c.LastZindex
}

func (c *Context) SetFocus(id ID) {
	c.Focus = id
	c.UpdatedFocus = true
}

func (c *Context) Update(f func()) {
	c.begin()
	defer c.end()
	f()
}

func (c *Context) begin() {
	c.updateInput()

	c.commandList = c.commandList[:0]
	c.rootList = c.rootList[:0]
	c.ScrollTarget = nil
	c.HoverRoot = c.NextHoverRoot
	c.NextHoverRoot = nil
	c.mouseDelta.X = c.mousePos.X - c.lastMousePos.X
	c.mouseDelta.Y = c.mousePos.Y - c.lastMousePos.Y
	c.tick++
}

func (c *Context) end() {
	// check stacks
	expect(len(c.containerStack) == 0)
	expect(len(c.clipStack) == 0)
	expect(len(c.idStack) == 0)
	expect(len(c.layoutStack) == 0)

	// handle scroll input
	if c.ScrollTarget != nil {
		c.ScrollTarget.Scroll.X += c.scrollDelta.X
		c.ScrollTarget.Scroll.Y += c.scrollDelta.Y
	}

	// unset focus if focus id was not touched this frame
	if !c.UpdatedFocus {
		c.Focus = 0
	}
	c.UpdatedFocus = false

	// bring hover root to front if mouse was pressed
	if c.mousePressed != 0 && c.NextHoverRoot != nil &&
		c.NextHoverRoot.Zindex < c.LastZindex &&
		c.NextHoverRoot.Zindex >= 0 {
		c.BringToFront(c.NextHoverRoot)
	}

	// reset input state
	c.keyPressed = 0
	c.textInput = nil
	c.mousePressed = 0
	c.scrollDelta = image.Pt(0, 0)
	c.lastMousePos = c.mousePos

	// sort root containers by zindex
	sort.SliceStable(c.rootList, func(i, j int) bool {
		return c.rootList[i].Zindex < c.rootList[j].Zindex
	})

	// set root container jump commands
	for i := 0; i < len(c.rootList); i++ {
		cnt := c.rootList[i]
		// if this is the first container then make the first command jump to it.
		// otherwise set the previous container's tail to jump to this one
		if i == 0 {
			cmd := c.commandList[0]
			expect(cmd.typ == commandJump)
			cmd.jump.dstIdx = cnt.HeadIdx + 1
			expect(cmd.jump.dstIdx < commandListSize)
		} else {
			prev := c.rootList[i-1]
			c.commandList[prev.TailIdx].jump.dstIdx = cnt.HeadIdx + 1
		}
		// make the last container's tail jump to the end of command list
		if i == len(c.rootList)-1 {
			c.commandList[cnt.TailIdx].jump.dstIdx = len(c.commandList)
		}
	}
}
