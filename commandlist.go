// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2024 The Ebitengine Authors

package microui

import (
	"image"
	"image/color"
)

/*============================================================================
** commandlist
**============================================================================*/

// adds a new command with type cmd_type to command_list
func (ctx *Context) PushCommand(cmd_type int) *Command {
	cmd := Command{
		Type: cmd_type,
	}
	//expect(uintptr(len(ctx.CommandList))*size+size < MU_COMMANDLIST_SIZE)
	cmd.Base.Type = cmd_type
	cmd.Idx = len(ctx.CommandList)
	ctx.CommandList = append(ctx.CommandList, &cmd)
	return &cmd
}

// sets cmd to the next command in command_list, returns true if success
func (ctx *Context) NextCommand(cmd **Command) bool {
	if len(ctx.CommandList) == 0 {
		return false
	}
	if *cmd == nil {
		*cmd = ctx.CommandList[0]
	} else {
		*cmd = ctx.CommandList[(*cmd).Idx+1]
	}

	for (*cmd).Idx < len(ctx.CommandList) {
		if (*cmd).Type != CommandJump {
			return true
		}
		idx := (*cmd).Jump.DstIdx
		if idx > len(ctx.CommandList)-1 {
			break
		}
		*cmd = ctx.CommandList[idx]
	}
	return false
}

// pushes a new jump command to command_list
func (ctx *Context) PushJump(dstIdx int) int {
	cmd := ctx.PushCommand(CommandJump)
	cmd.Jump.DstIdx = dstIdx
	return len(ctx.CommandList) - 1
}

// pushes a new clip command
func (ctx *Context) SetClip(rect Rect) {
	cmd := ctx.PushCommand(CommandClip)
	cmd.Clip.Rect = rect
}

// pushes a new rect command
func (ctx *Context) DrawRect(rect Rect, color color.Color) {
	rect2 := intersect_rects(rect, ctx.GetClipRect())
	if rect2.W > 0 && rect2.H > 0 {
		cmd := ctx.PushCommand(CommandRect)
		cmd.Rect.Rect = rect2
		cmd.Rect.Color = color
	}
}

func (ctx *Context) DrawBox(rect Rect, color color.Color) {
	ctx.DrawRect(NewRect(rect.X+1, rect.Y, rect.W-2, 1), color)
	ctx.DrawRect(NewRect(rect.X+1, rect.Y+rect.H-1, rect.W-2, 1), color)
	ctx.DrawRect(NewRect(rect.X, rect.Y, 1, rect.H), color)
	ctx.DrawRect(NewRect(rect.X+rect.W-1, rect.Y, 1, rect.H), color)
}

func (ctx *Context) DrawText(font Font, str string, pos image.Point, color color.Color) {
	rect := NewRect(pos.X, pos.Y, ctx.TextWidth(font, str), ctx.TextHeight(font))
	clipped := ctx.CheckClip(rect)
	if clipped == ClipAll {
		return
	}
	if clipped == ClipPart {
		ctx.SetClip(ctx.GetClipRect())
	}
	// add command
	cmd := ctx.PushCommand(CommandText)
	cmd.Text.Str = str
	cmd.Text.Pos = pos
	cmd.Text.Color = color
	cmd.Text.Font = font
	// reset clipping if it was set
	if clipped != 0 {
		ctx.SetClip(unclippedRect)
	}
}

func (ctx *Context) DrawIcon(id int, rect Rect, color color.Color) {
	// do clip command if the rect isn't fully contained within the cliprect
	clipped := ctx.CheckClip(rect)
	if clipped == ClipAll {
		return
	}
	if clipped == ClipPart {
		ctx.SetClip(ctx.GetClipRect())
	}
	// do icon command
	cmd := ctx.PushCommand(CommandIcon)
	cmd.Icon.ID = id
	cmd.Icon.Rect = rect
	cmd.Icon.Color = color
	// reset clipping if it was set
	if clipped != 0 {
		ctx.SetClip(unclippedRect)
	}
}
