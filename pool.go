// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2024 The Ebitengine Authors

package microui

func (ctx *Context) poolInit(items []poolItem, id ID) int {
	f := ctx.tick
	n := -1
	for i := 0; i < len(items); i++ {
		if items[i].lastUpdate < f {
			f = items[i].lastUpdate
			n = i
		}
	}
	expect(n > -1)
	items[n].id = id
	ctx.poolUpdate(items, n)
	return n
}

// returns the index of an ID in the pool. returns -1 if it is not found
func (ctx *Context) poolGet(items []poolItem, id ID) int {
	for i := 0; i < len(items); i++ {
		if items[i].id == id {
			return i
		}
	}
	return -1
}

func (ctx *Context) poolUpdate(items []poolItem, idx int) {
	items[idx].lastUpdate = ctx.tick
}
