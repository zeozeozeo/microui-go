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
func (ctx *Context) pushCommand(cmd_type int) *command {
	cmd := command{
		typ: cmd_type,
	}
	//expect(uintptr(len(ctx.CommandList))*size+size < MU_COMMANDLIST_SIZE)
	cmd.base.typ = cmd_type
	cmd.idx = len(ctx.commandList)
	ctx.commandList = append(ctx.commandList, &cmd)
	return &cmd
}

// sets cmd to the next command in command_list, returns true if success
func (ctx *Context) nextCommand(cmd **command) bool {
	if len(ctx.commandList) == 0 {
		return false
	}
	if *cmd == nil {
		*cmd = ctx.commandList[0]
	} else {
		*cmd = ctx.commandList[(*cmd).idx+1]
	}

	for (*cmd).idx < len(ctx.commandList) {
		if (*cmd).typ != CommandJump {
			return true
		}
		idx := (*cmd).jump.dstIdx
		if idx > len(ctx.commandList)-1 {
			break
		}
		*cmd = ctx.commandList[idx]
	}
	return false
}

// pushes a new jump command to command_list
func (ctx *Context) pushJump(dstIdx int) int {
	cmd := ctx.pushCommand(CommandJump)
	cmd.jump.dstIdx = dstIdx
	return len(ctx.commandList) - 1
}

// pushes a new clip command
func (ctx *Context) SetClip(rect image.Rectangle) {
	cmd := ctx.pushCommand(CommandClip)
	cmd.clip.rect = rect
}

// pushes a new rect command
func (ctx *Context) DrawRect(rect image.Rectangle, color color.Color) {
	rect2 := rect.Intersect(ctx.GetClipRect())
	if rect2.Dx() > 0 && rect2.Dy() > 0 {
		cmd := ctx.pushCommand(CommandRect)
		cmd.rect.rect = rect2
		cmd.rect.color = color
	}
}

func (ctx *Context) DrawBox(rect image.Rectangle, color color.Color) {
	ctx.DrawRect(image.Rect(rect.Min.X+1, rect.Min.Y, rect.Max.X-1, rect.Min.Y+1), color)
	ctx.DrawRect(image.Rect(rect.Min.X+1, rect.Max.Y-1, rect.Max.X-1, rect.Max.Y), color)
	ctx.DrawRect(image.Rect(rect.Min.X, rect.Min.Y, rect.Min.X+1, rect.Max.Y), color)
	ctx.DrawRect(image.Rect(rect.Max.X-1, rect.Min.Y, rect.Max.X, rect.Max.Y), color)
}

func (ctx *Context) DrawText(font Font, str string, pos image.Point, color color.Color) {
	rect := image.Rect(pos.X, pos.Y, pos.X+ctx.TextWidth(font, str), pos.Y+ctx.TextHeight(font))
	clipped := ctx.CheckClip(rect)
	if clipped == ClipAll {
		return
	}
	if clipped == ClipPart {
		ctx.SetClip(ctx.GetClipRect())
	}
	// add command
	cmd := ctx.pushCommand(CommandText)
	cmd.text.str = str
	cmd.text.pos = pos
	cmd.text.color = color
	cmd.text.font = font
	// reset clipping if it was set
	if clipped != 0 {
		ctx.SetClip(unclippedRect)
	}
}

func (ctx *Context) DrawIcon(id int, rect image.Rectangle, color color.Color) {
	// do clip command if the rect isn't fully contained within the cliprect
	clipped := ctx.CheckClip(rect)
	if clipped == ClipAll {
		return
	}
	if clipped == ClipPart {
		ctx.SetClip(ctx.GetClipRect())
	}
	// do icon command
	cmd := ctx.pushCommand(CommandIcon)
	cmd.icon.id = id
	cmd.icon.rect = rect
	cmd.icon.color = color
	// reset clipping if it was set
	if clipped != 0 {
		ctx.SetClip(unclippedRect)
	}
}
