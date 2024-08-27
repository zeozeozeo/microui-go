package microui

// calls nextCmdFunc for every command in the command list, clears it when done.
// equivalent to calling `ctx.NextCommand` in a loop
func (ctx *Context) Render(nextCmdFunc func(cmd *Command)) {
	var cmd *Command
	for ctx.NextCommand(&cmd) {
		nextCmdFunc(cmd)
	}
	ctx.CommandList = nil
}
