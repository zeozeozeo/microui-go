package microui

import "unsafe"

/*============================================================================
** commandlist
**============================================================================*/

// adds a new command with type cmd_type to command_list
func (ctx *Context) PushCommand(cmd_type int) *Command {
	cmd := Command{Type: cmd_type}
	size := unsafe.Sizeof(cmd)
	expect(uintptr(len(ctx.CommandList))*size+size < MU_COMMANDLIST_SIZE)
	cmd.Base.Type = cmd_type
	cmd.Base.Size = size
	cmd.Idx = len(ctx.CommandList)
	ctx.CommandList = append(ctx.CommandList, &cmd)
	return &cmd
}

// sets cmd to the next command in command_list, returns true if success
func (ctx *Context) NextCommand(cmd *Command) bool {
	if cmd.Idx+1 < len(ctx.CommandList) {
		*cmd = *ctx.CommandList[cmd.Idx+1]
		return true
	}
	return false
}

// pushes a new jump command to command_list
func (ctx *Context) PushJump(dst *Command) *Command {
	cmd := ctx.PushCommand(MU_COMMAND_JUMP)
	cmd.Jump.Dst = dst
	return cmd
}

// pushes a new clip command
func (ctx *Context) SetClip(rect MuRect) {
	cmd := ctx.PushCommand(MU_COMMAND_CLIP)
	cmd.Clip.Rect = rect
}

// pushes a new rect command
func (ctx *Context) DrawRect(rect MuRect, color MuColor) {
	rect2 := intersect_rects(rect, ctx.GetClipRect())
	if rect2.W > 0 && rect2.H > 0 {
		cmd := ctx.PushCommand(MU_COMMAND_RECT)
		cmd.Rect.Rect = rect2
		cmd.Rect.Color = color
	}
}

func (ctx *Context) DrawBox(rect MuRect, color MuColor) {
	ctx.DrawRect(Rect(rect.X+1, rect.Y, rect.W-2, 1), color)
	ctx.DrawRect(Rect(rect.X+1, rect.Y+rect.H-1, rect.W-2, 1), color)
	ctx.DrawRect(Rect(rect.X, rect.Y, 1, rect.H), color)
	ctx.DrawRect(Rect(rect.X+rect.W-1, rect.Y, 1, rect.H), color)
}

func (ctx *Context) DrawText(font Font, str string, pos MuVec2, color MuColor) {
	rect := Rect(pos.X, pos.Y, ctx.TextWidth(font, str), ctx.TextHeight(font))
	clipped := ctx.CheckClip(rect)
	if clipped == MU_CLIP_ALL {
		return
	}
	if clipped == MU_CLIP_PART {
		ctx.SetClip(ctx.GetClipRect())
	}
	// add command
	cmd := ctx.PushCommand(MU_COMMAND_TEXT)
	cmd.Text.Str = str
	cmd.Text.Pos = pos
	cmd.Text.Color = color
	cmd.Text.Font = font
	// reset clipping if it was set
	if clipped != 0 {
		ctx.SetClip(UnclippedRect)
	}
}

func (ctx *Context) DrawIcon(id int, rect MuRect, color MuColor) {
	// do clip command if the rect isn't fully contained within the cliprect
	clipped := ctx.CheckClip(rect)
	if clipped == MU_CLIP_ALL {
		return
	}
	if clipped == MU_CLIP_PART {
		ctx.SetClip(ctx.GetClipRect())
	}
	// do icon command
	cmd := ctx.PushCommand(MU_COMMAND_ICON)
	cmd.Icon.Id = id
	cmd.Icon.Rect = rect
	cmd.Icon.Color = color
	// reset clipping if it was set
	if clipped != 0 {
		ctx.SetClip(UnclippedRect)
	}
}
