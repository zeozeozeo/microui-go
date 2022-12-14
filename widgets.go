package microui

func (ctx *Context) Button(label string) bool {
	return ctx.ButtonEx(label, 0, MU_OPT_ALIGNCENTER) != 0
}

func (ctx *Context) TextBox(buf *string) int {
	return ctx.TextBoxEx(buf, 0)
}

func (ctx *Context) Slider(value *mu_Real, lo, hi mu_Real) int {
	return ctx.SliderEx(value, lo, hi, 0, MU_SLIDER_FMT, MU_OPT_ALIGNCENTER)
}

func (ctx *Context) Number(value *mu_Real, step mu_Real) int {
	return ctx.NumberEx(value, step, MU_SLIDER_FMT, MU_OPT_ALIGNCENTER)
}

func (ctx *Context) Header(label string) bool {
	return ctx.HeaderEx(label, 0) != 0
}

func (ctx *Context) BeginTreeNode(label string) bool {
	return ctx.BeginTreeNodeEx(label, 0) != 0
}

func (ctx *Context) BeginWindow(title string, rect Rect) bool {
	return ctx.BeginWindowEx(title, rect, 0) != 0
}

func (ctx *Context) BeginPanel(name string) {
	ctx.BeginPanelEx(name, 0)
}
