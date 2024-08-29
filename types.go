// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2024 The Ebitengine Authors

package microui

import (
	"image"
	"image/color"
)

type ID uintptr

type poolItem struct {
	id         ID
	lastUpdate int
}

type baseCommand struct {
	typ int
}

type jumpCommand struct {
	dstIdx int
}

type clipCommand struct {
	rect image.Rectangle
}

type rectCommand struct {
	rect  image.Rectangle
	color color.Color
}

type Font interface{} // Font is interface{}, microui does not manage fonts

type textCommand struct {
	font  Font
	pos   image.Point
	color color.Color
	str   string
}

type iconCommand struct {
	rect  image.Rectangle
	icon  Icon
	color color.Color
}

type Layout struct {
	Body      image.Rectangle
	Next      image.Rectangle
	Position  image.Point
	Size      image.Point
	Max       image.Point
	Widths    [maxWidths]int
	Items     int
	ItemIndex int
	NextRow   int
	NextType  int
	Indent    int
}

type command struct {
	typ  int
	idx  int
	base baseCommand // type 0 (TODO)
	jump jumpCommand // type 1
	clip clipCommand // type 2
	rect rectCommand // type 3
	text textCommand // type 4
	icon iconCommand // type 5
}

type Container struct {
	HeadIdx     int
	TailIdx     int
	Rect        image.Rectangle
	Body        image.Rectangle
	ContentSize image.Point
	Scroll      image.Point
	Zindex      int
	Open        bool
}

type Style struct {
	Font          Font
	Size          image.Point
	Padding       int
	Spacing       int
	Indent        int
	TitleHeight   int
	ScrollbarSize int
	ThumbSize     int
	Colors        [ColorMax + 1]color.RGBA
}

type Context struct {
	// callbacks

	TextWidth  func(font Font, str string) int
	TextHeight func(font Font) int
	DrawFrame  func(ctx *Context, rect image.Rectangle, colorid int)

	// core state

	Style         *Style
	Hover         ID
	Focus         ID
	LastID        ID
	LastRect      image.Rectangle
	LastZindex    int
	UpdatedFocus  bool
	tick          int
	HoverRoot     *Container
	NextHoverRoot *Container
	ScrollTarget  *Container
	NumberEditBuf string
	NumberEdit    ID

	// stacks

	commandList    []*command
	rootList       []*Container
	containerStack []*Container
	clipStack      []image.Rectangle
	idStack        []ID
	layoutStack    []Layout

	// retained state pools

	containerPool [containerPoolSize]poolItem
	containers    [containerPoolSize]Container
	treeNodePool  [treeNodePoolSize]poolItem

	// input state

	mousePos     image.Point
	lastMousePos image.Point
	mouseDelta   image.Point
	scrollDelta  image.Point
	mouseDown    int
	mousePressed int
	keyDown      int
	keyPressed   int
	textInput    []rune
}
