// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2024 The Ebitengine Authors

package microui

import (
	"image"
	"image/color"
)

type ID uintptr

type PoolItem struct {
	ID         ID
	LastUpdate int
}

type baseCommand struct {
	Type int
}

type jumpCommand struct {
	Base   baseCommand
	DstIdx int
}

type clipCommand struct {
	Base baseCommand
	Rect image.Rectangle
}

type rectCommand struct {
	Base  baseCommand
	Rect  image.Rectangle
	Color color.Color
}

type Font interface{} // Font is interface{}, microui does not manage fonts

type textCommand struct {
	Base  baseCommand
	Font  Font
	Pos   image.Point
	Color color.Color
	Str   string
}

type iconCommand struct {
	Base  baseCommand
	Rect  image.Rectangle
	ID    int
	Color color.Color
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
	Type int
	Idx  int
	Base baseCommand // type 0 (TODO)
	Jump jumpCommand // type 1
	Clip clipCommand // type 2
	Rect rectCommand // type 3
	Text textCommand // type 4
	Icon iconCommand // type 5
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
	Frame         int
	HoverRoot     *Container
	NextHoverRoot *Container
	ScrollTarget  *Container
	NumberEditBuf string
	NumberEdit    ID

	// stacks

	commandList    []*command
	RootList       []*Container
	ContainerStack []*Container
	ClipStack      []image.Rectangle
	IDStack        []ID
	LayoutStack    []Layout

	// retained state pools

	ContainerPool [containerPoolSize]PoolItem
	Containers    [containerPoolSize]Container
	TreeNodePool  [treeNodePoolSize]PoolItem

	// input state

	MousePos     image.Point
	lastMousePos image.Point
	MouseDelta   image.Point
	ScrollDelta  image.Point
	MouseDown    int
	MousePressed int
	KeyDown      int
	KeyPressed   int
	TextInput    []rune
}
