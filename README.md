# ![microui](https://user-images.githubusercontent.com/3920290/75171571-be83c500-5723-11ea-8a50-504cc2ae1109.png)

A tiny, portable, immediate-mode UI library ported to Go (as of commit [0850aba860959c3e75fb3e97120ca92957f9d057](https://github.com/rxi/microui/tree/0850aba860959c3e75fb3e97120ca92957f9d057), v2.02)

# API changes

-   Functions and structs are renamed to be PascalCase and the prefix `mu_` is removed, like this:

    > `mu_push_command` -> `PushCommand`

    > `mu_begin_treenode_ex` -> `BeginTreeNodeEx`

    > `mu_get_clip_rect` -> `GetClipRect`

-   Every function that takes `mu_Context` (`Context`) instead has a `Context` reciever, so `Button(ctx, label)` becomes `ctx.Button(label)`
-   Stacks are now slices with variable length, `append` is used for `push` and `slice = slice[:len(slice)-1]` is used for `pop`
-   `mu_Font` (`Font`) is `interface{}`, since it doesn't store any font data. You can use `reflect` if you want to store values inside it
-   All pointer-based commands (`MU_COMMAND_JUMP`) and the `Command` struct have been reworked to use indices
-   The `mu_Real` type has been replaced with `float32` because Go does not allow implicit casting of identical type aliases 
-   The library is split into separate files instead of one file
-   The library is ~1300 lines of code in total

## Additional functions:

-   `NewContext`, which is a helper for creating a new `Context`
-   `ctx.Render`, which calls a function for every command inside the command list, then clears it

# Integrations, demos, renderers

* [Ebitengine](https://ebitengine.org/) rendering backend + demo port: [zeozeozeo/ebitengine-microui-go](https://github.com/zeozeozeo/ebitengine-microui-go)
    ![microui demo running in Ebitengine](https://github.com/zeozeozeo/ebitengine-microui-go/blob/main/screenshots/demo.png?raw=true)
* Official Ebitengine fork and integration efforts: [ebitengine/microui](https://github.com/ebitengine/microui)

# Notes

The library expects the user to provide input and handle the resultant drawing commands, it does not do any drawing/tessellation itself.

# Credits

Thank you [@rxi](https://github.com/rxi) for creating this awesome library and thank you [@Zyko0](https://github.com/Zyko0) for contributing numerous fixes to this Go port.
