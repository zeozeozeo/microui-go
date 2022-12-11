package microui

/*============================================================================
** input handlers
**============================================================================*/

func (ctx *Context) InputMouseMove(x, y int) {
	ctx.MousePos = Vec2(x, y)
}

func (ctx *Context) InputMouseDown(x, y int, btn int) {
	ctx.InputMouseMove(x, y)
	ctx.MouseDown |= btn
	ctx.MousePressed |= btn
}

func (ctx *Context) InputMouseUp(x, y int, btn int) {
	ctx.InputMouseMove(x, y)
	ctx.MouseDown &= ^btn
}

func (ctx *Context) InputScroll(x, y int) {
	ctx.ScrollDelta.X += x
	ctx.ScrollDelta.Y += y
}

func (ctx *Context) InputKeyDown(key int) {
	ctx.KeyPressed |= key
	ctx.KeyDown |= key
}

func (ctx *Context) InputKeyUp(key int) {
	ctx.KeyDown &= ^key
}

func (ctx *Context) InputText(text []rune) {
	ctx.TextInput = text
}
