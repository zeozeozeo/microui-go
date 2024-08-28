// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2024 The Ebitengine Authors

package microui

const (
	ClipPart = 1 + iota
	ClipAll
)

const (
	CommandJump = 1 + iota
	CommandClip
	CommandRect
	CommandText
	CommandIcon
)

const (
	ColorText = iota
	ColorBorder
	ColorWindowBG
	ColorTitleBG
	ColorTitleText
	ColorPanelBG
	ColorButton
	ColorButtonHover
	ColorButtonFocus
	ColorBase
	ColorBaseHover
	ColorBaseFocus
	ColorScrollBase
	ColorScrollThumb
	ColorMax = ColorScrollThumb
)

const (
	IconClose = 1 + iota
	IconCheck
	IconCollapsed
	IconExpanded
)

const (
	ResActive = (1 << 0)
	ResSubmit = (1 << 1)
	ResChange = (1 << 2)
)

const (
	OptAlignCenter = (1 << 0)
	OptAlignRight  = (1 << 1)
	OptNoInteract  = (1 << 2)
	OptNoFrame     = (1 << 3)
	OptNoResize    = (1 << 4)
	OptNoScroll    = (1 << 5)
	OptNoClose     = (1 << 6)
	OptNoTitle     = (1 << 7)
	OptHoldFocus   = (1 << 8)
	OptAutoSize    = (1 << 9)
	OptPopup       = (1 << 10)
	OptClosed      = (1 << 11)
	OptExpanded    = (1 << 12)
)

const (
	mouseLeft   = (1 << 0)
	mouseRight  = (1 << 1)
	mouseMiddle = (1 << 2)
)

const (
	MU_KEY_SHIFT     = (1 << 0)
	MU_KEY_CTRL      = (1 << 1)
	MU_KEY_ALT       = (1 << 2)
	MU_KEY_BACKSPACE = (1 << 3)
	MU_KEY_RETURN    = (1 << 4)
)

const (
	Relative = 1 + iota
	Absolute
)
