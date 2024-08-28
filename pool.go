// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2024 The Ebitengine Authors

package microui

/*============================================================================
** pool
**============================================================================*/

func (ctx *Context) PoolInit(items []PoolItem, id ID) int {
	f := ctx.Frame
	n := -1
	for i := 0; i < len(items); i++ {
		if items[i].LastUpdate < f {
			f = items[i].LastUpdate
			n = i
		}
	}
	expect(n > -1)
	items[n].ID = id
	ctx.PoolUpdate(items, n)
	return n
}

// returns the index of an ID in the pool. returns -1 if it is not found
func (ctx *Context) PoolGet(items []PoolItem, id ID) int {
	for i := 0; i < len(items); i++ {
		if items[i].ID == id {
			return i
		}
	}
	return -1
}

func (ctx *Context) PoolUpdate(items []PoolItem, idx int) {
	items[idx].LastUpdate = ctx.Frame
}
