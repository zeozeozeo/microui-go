package microui

const MU_VERSION = "2.01"

const (
	MU_COMMANDLIST_SIZE    = 256 * 1024
	MU_ROOTLIST_SIZE       = 32
	MU_CONTAINERSTACK_SIZE = 32
	MU_CLIPSTACK_SIZE      = 32
	MU_IDSTACK_SIZE        = 32
	MU_LAYOUTSTACK_SIZE    = 16
	MU_CONTAINERPOOL_SIZE  = 48
	MU_TREENODEPOOL_SIZE   = 48
	MU_MAX_WIDTHS          = 16
)

const (
	MU_REAL_FMT   = "%.3g"
	MU_SLIDER_FMT = "%.2f"
	MU_MAX_FMT    = 127
)

var default_style Style = Style{
	Font:          nil,
	Size:          MuVec2{68, 10},
	Padding:       5,
	Spacing:       4,
	Indent:        24,
	TitleHeight:   24,
	ScrollbarSize: 12,
	ThumbSize:     8,
	Colors: [14]MuColor{
		{230, 230, 230, 255}, // MU_COLOR_TEXT
		{25, 25, 25, 255},    // MU_COLOR_BORDER
		{50, 50, 50, 255},    // MU_COLOR_WINDOWBG
		{25, 25, 25, 255},    // MU_COLOR_TITLEBG
		{240, 240, 240, 255}, // MU_COLOR_TITLETEXT
		{0, 0, 0, 0},         // MU_COLOR_PANELBG
		{75, 75, 75, 255},    // MU_COLOR_BUTTON
		{95, 95, 95, 255},    // MU_COLOR_BUTTONHOVER
		{115, 115, 115, 255}, // MU_COLOR_BUTTONFOCUS
		{30, 30, 30, 255},    // MU_COLOR_BASE
		{35, 35, 35, 255},    // MU_COLOR_BASEHOVER
		{40, 40, 40, 255},    // MU_COLOR_BASEFOCUS
		{43, 43, 43, 255},    // MU_COLOR_SCROLLBASE
		{30, 30, 30, 255},    // MU_COLOR_SCROLLTHUMB
	},
}

var (
	UnclippedRect = MuRect{0, 0, 0x1000000, 0x1000000}
)

const (
	HASH_INITIAL = 2166136261 // 32bit fnv-1a hash
)
