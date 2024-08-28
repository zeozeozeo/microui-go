// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2024 The Ebitengine Authors

package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/ebitengine/microui"
)

type Game struct {
	ctx *microui.Context
}

func New() *Game {
	return &Game{
		ctx: microui.NewContext(),
	}
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}
	g.ProcessFrame()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.ctx.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 1280, 960
}

func main() {
	ebiten.SetWindowTitle("Ebitengine Microui Demo")
	ebiten.SetWindowSize(1280, 960)
	if err := ebiten.RunGame(New()); err != nil {
		log.Fatal("err: ", err)
	}
}
