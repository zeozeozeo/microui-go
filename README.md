# ![microui](https://user-images.githubusercontent.com/3920290/75171571-be83c500-5723-11ea-8a50-504cc2ae1109.png)

A tiny, portable, immediate-mode UI library ported to Go (as of commit [05d7b46c9cf650dd0c5fbc83a9bebf87c80d02a5](https://github.com/rxi/microui/tree/05d7b46c9cf650dd0c5fbc83a9bebf87c80d02a5))

# API changes

-   Functions and structs are renamed to be PascalCase and the prefix `mu_` is removed, like this:

    > `mu_push_command` -> `PushCommand` > `mu_begin_treenode_ex` -> `BeginTreeNodeEx` > `mu_get_clip_rect` -> `GetClipRect`

-   Every function that takes `mu_Context` (`Context`) instead has a `Context` reciever, so `Button(ctx, label)` becomes `ctx.Button(label)`
-   Stacks are now slices with variable length, `append` is used for `push` and `slice = slice[:len(slice)-1]` is used for `pop`
-   `mu_Font` (`Font`) is `interface{}`, since it doesn't store any font data. You can use `reflect` if you want to store values inside it
-   The library is split into separate files instead of one file 
-   The library is ~1300 lines of code in total

## Additional functions:

-   `NewContext`, which is a helper for creating a new `Context`
-   `ctx.Render`, which calls a function for every command inside the command list, then clears it

# Notes

The library expects the user to provide input and handle the resultant drawing commands, it does not do any drawing itself.