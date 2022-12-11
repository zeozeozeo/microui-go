package microui

// calls nextCmdFunc for every command in the command list, clears it when done
func (ctx *Context) Render(nextCmdFunc func(cmd *Command)) {
	for _, cmd := range ctx.CommandList {
		nextCmdFunc(cmd)
	}
	ctx.CommandList = nil
}
