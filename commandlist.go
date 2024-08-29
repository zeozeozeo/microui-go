// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2024 The Ebitengine Authors

package microui

import (
	"image"
	"image/color"
)

// pushCommand adds a new command with type cmd_type to command_list
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
		if (*cmd).typ != commandJump {
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

// pushJump pushes a new jump command to command_list
func (ctx *Context) pushJump(dstIdx int) int {
	cmd := ctx.pushCommand(commandJump)
	cmd.jump.dstIdx = dstIdx
	return len(ctx.commandList) - 1
}

// SetClip pushes a new clip command
func (ctx *Context) SetClip(rect image.Rectangle) {
	cmd := ctx.pushCommand(commandClip)
	cmd.clip.rect = rect
}

// DrawRect pushes a new rect command
func (ctx *Context) DrawRect(rect image.Rectangle, color color.Color) {
	rect2 := rect.Intersect(ctx.GetClipRect())
	if rect2.Dx() > 0 && rect2.Dy() > 0 {
		cmd := ctx.pushCommand(commandRect)
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

func (ctx *Context) DrawText(str string, pos image.Point, color color.Color) {
	rect := image.Rect(pos.X, pos.Y, pos.X+textWidth(str), pos.Y+textHeight())
	clipped := ctx.CheckClip(rect)
	if clipped == ClipAll {
		return
	}
	if clipped == ClipPart {
		ctx.SetClip(ctx.GetClipRect())
	}
	// add command
	cmd := ctx.pushCommand(commandText)
	cmd.text.str = str
	cmd.text.pos = pos
	cmd.text.color = color
	// reset clipping if it was set
	if clipped != 0 {
		ctx.SetClip(unclippedRect)
	}
}

func (ctx *Context) DrawIcon(icon Icon, rect image.Rectangle, color color.Color) {
	// do clip command if the rect isn't fully contained within the cliprect
	clipped := ctx.CheckClip(rect)
	if clipped == ClipAll {
		return
	}
	if clipped == ClipPart {
		ctx.SetClip(ctx.GetClipRect())
	}
	// do icon command
	cmd := ctx.pushCommand(commandIcon)
	cmd.icon.icon = icon
	cmd.icon.rect = rect
	cmd.icon.color = color
	// reset clipping if it was set
	if clipped != 0 {
		ctx.SetClip(unclippedRect)
	}
}
