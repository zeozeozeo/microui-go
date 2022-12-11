package microui

import "image/color"

type mu_Id uintptr
type mu_Real float32

type MuVec2 struct {
	X, Y int
}

type MuRect struct {
	X, Y, W, H int
}

type MuColor struct {
	R, G, B, A uint8
}

func (c *MuColor) ToRGBA() color.RGBA {
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
	Rect MuRect
}

type RectCommand struct {
	Base  BaseCommand
	Rect  MuRect
	Color MuColor
}

type Font interface{} // Font is interface{}, microui does not manage fonts

type TextCommand struct {
	Base  BaseCommand
	Font  Font
	Pos   MuVec2
	Color MuColor
	Str   string
}

type IconCommand struct {
	Base  BaseCommand
	Rect  MuRect
	Id    int
	Color MuColor
}

type Layout struct {
	Body      MuRect
	Next      MuRect
	Position  MuVec2
	Size      MuVec2
	Max       MuVec2
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
	Rect        MuRect
	Body        MuRect
	ContentSize MuVec2
	Scroll      MuVec2
	Zindex      int
	Open        bool
}

type Style struct {
	Font          Font
	Size          MuVec2
	Padding       int
	Spacing       int
	Indent        int
	TitleHeight   int
	ScrollbarSize int
	ThumbSize     int
	Colors        [MU_COLOR_MAX]MuColor
}

type Context struct {
	// callbacks

	TextWidth  func(font Font, str string) int
	TextHeight func(font Font) int
	DrawFrame  func(ctx *Context, rect MuRect, colorid int)

	// core state

	_style        Style
	Style         *Style
	Hover         mu_Id
	Focus         mu_Id
	LastID        mu_Id
	LastRect      MuRect
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
	ClipStack      []MuRect
	IdStack        []mu_Id
	LayoutStack    []Layout

	// retained state pools

	ContainerPool [MU_CONTAINERPOOL_SIZE]MuPoolItem
	Containers    [MU_CONTAINERPOOL_SIZE]Container
	TreeNodePool  [MU_TREENODEPOOL_SIZE]MuPoolItem

	// input state

	MousePos     MuVec2
	lastMousePos MuVec2
	MouseDelta   MuVec2
	ScrollDelta  MuVec2
	MouseDown    int
	MousePressed int
	KeyDown      int
	KeyPressed   int
	TextInput    []rune
}
