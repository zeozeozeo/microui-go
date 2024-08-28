// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2024 The Ebitengine Authors

package microui

import (
	"image"
	"image/color"
)

type ID uintptr

// TODO: Replace muRect with image.Rectangle.

type muRect struct {
	X, Y, W, H int
}

func newMuRect(x, y, w, h int) muRect {
	return muRect{x, y, w, h}
}

func rectFromRectangle(r image.Rectangle) muRect {
	return muRect{r.Min.X, r.Min.Y, r.Dx(), r.Dy()}
}

func (r muRect) rectangle() image.Rectangle {
	return image.Rect(r.X, r.Y, r.X+r.W, r.Y+r.H)
}

type PoolItem struct {
	ID         ID
	LastUpdate int
}

type BaseCommand struct {
	Type int
}

type JumpCommand struct {
	Base   BaseCommand
	DstIdx int
}

type ClipCommand struct {
	Base BaseCommand
	Rect image.Rectangle
}

type RectCommand struct {
	Base  BaseCommand
	Rect  image.Rectangle
	Color color.Color
}

type Font interface{} // Font is interface{}, microui does not manage fonts

type TextCommand struct {
	Base  BaseCommand
	Font  Font
	Pos   image.Point
	Color color.Color
	Str   string
}

type IconCommand struct {
	Base  BaseCommand
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

type Command struct {
	Type int
	Idx  int
	Base BaseCommand // type 0 (TODO)
	Jump JumpCommand // type 1
	Clip ClipCommand // type 2
	Rect RectCommand // type 3
	Text TextCommand // type 4
	Icon IconCommand // type 5
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

	CommandList    []*Command
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
