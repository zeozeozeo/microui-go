package microui

func drawFrame(ctx *Context, rect MuRect, colorid int) {
	if colorid == MU_COLOR_SCROLLBASE ||
		colorid == MU_COLOR_SCROLLTHUMB ||
		colorid == MU_COLOR_TITLEBG {
		return
	}

	// draw border
	if ctx.Style.Colors[MU_COLOR_BORDER].A != 0 {
		ctx.DrawBox(expand_rect(rect, 1), ctx.Style.Colors[MU_COLOR_BORDER])
	}
}

func initContext(ctx *Context) {
	ctx.DrawFrame = drawFrame
	ctx._style = default_style
	ctx.Style = &ctx._style
}

func NewContext() *Context {
	ctx := &Context{}
	initContext(ctx)
	return ctx
}
