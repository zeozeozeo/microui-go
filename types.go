package microui

import "image/color"

type mu_Id uintptr
type mu_Real float32

type Vec2 struct {
	X, Y int
}

type Rect struct {
	X, Y, W, H int
}

type Color struct {
	R, G, B, A uint8
}

func (c *Color) ToRGBA() color.RGBA {
	return color.RGBA{c.R, c.G, c.B, c.A}
}

type MuPoolItem struct {
	ID         mu_Id
	LastUpdate int
}

type BaseCommand struct {
	Type int
	Size uintptr
}

type JumpCommand struct {
	Base BaseCommand
	Dst  *Command
}

type ClipCommand struct {
	Base BaseCommand
	Rect Rect
}

type RectCommand struct {
	Base  BaseCommand
	Rect  Rect
	Color Color
}

type Font interface{} // Font is interface{}, microui does not manage fonts

type TextCommand struct {
	Base  BaseCommand
	Font  Font
	Pos   Vec2
	Color Color
	Str   string
}

type IconCommand struct {
	Base  BaseCommand
	Rect  Rect
	Id    int
	Color Color
}

type Layout struct {
	Body      Rect
	Next      Rect
	Position  Vec2
	Size      Vec2
	Max       Vec2
	Widths    [MU_MAX_WIDTHS]int
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
	Head, Tail  *Command
	Rect        Rect
	Body        Rect
	ContentSize Vec2
	Scroll      Vec2
	Zindex      int
	Open        bool
}

type Style struct {
	Font          Font
	Size          Vec2
	Padding       int
	Spacing       int
	Indent        int
	TitleHeight   int
	ScrollbarSize int
	ThumbSize     int
	Colors        [MU_COLOR_MAX]Color
}

type Context struct {
	// callbacks

	TextWidth  func(font Font, str string) int
	TextHeight func(font Font) int
	DrawFrame  func(ctx *Context, rect Rect, colorid int)

	// core state

	_style        Style
	Style         *Style
	Hover         mu_Id
	Focus         mu_Id
	LastID        mu_Id
	LastRect      Rect
	LastZindex    int
	UpdatedFocus  bool
	Frame         int
	HoverRoot     *Container
	NextHoverRoot *Container
	ScrollTarget  *Container
	NumberEditBuf string
	NumberEdit    mu_Id

	// stacks

	CommandList    []*Command
	RootList       []*Container
	ContainerStack []*Container
	ClipStack      []Rect
	IdStack        []mu_Id
	LayoutStack    []Layout

	// retained state pools

	ContainerPool [MU_CONTAINERPOOL_SIZE]MuPoolItem
	Containers    [MU_CONTAINERPOOL_SIZE]Container
	TreeNodePool  [MU_TREENODEPOOL_SIZE]MuPoolItem

	// input state

	MousePos     Vec2
	lastMousePos Vec2
	MouseDelta   Vec2
	ScrollDelta  Vec2
	MouseDown    int
	MousePressed int
	KeyDown      int
	KeyPressed   int
	TextInput    []rune
}
