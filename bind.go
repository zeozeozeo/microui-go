// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2024 The Ebitengine Authors

package microui

import (
	"bytes"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font/gofont/goregular"
)

var (
	face text.Face
)

func init() {
	var err error

	src, err := text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
	if err != nil {
		panic(err)
	}

	face = &text.GoTextFace{
		Source: src,
		Size:   14,
	}
}

func (c *Context) updateInput() {
	cx, cy := ebiten.CursorPosition()
	c.inputMouseMove(cx, cy)
	if wx, wy := ebiten.Wheel(); wx != 0 || wy != 0 {
		c.inputScroll(int(wx*-30), int(wy*-30))
	}
	// TODO: Use exp/textinput.Field.
	chars := ebiten.AppendInputChars(nil)
	if len(chars) > 0 {
		c.inputText(chars)
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		c.inputMouseDown(cx, cy, ebiten.MouseButtonLeft)
	} else if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		c.inputMouseUp(cx, cy, ebiten.MouseButtonLeft)
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		c.inputMouseDown(cx, cy, ebiten.MouseButtonRight)
	} else if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonRight) {
		c.inputMouseUp(cx, cy, ebiten.MouseButtonRight)
	}
	for _, k := range []ebiten.Key{ebiten.KeyAlt, ebiten.KeyBackspace, ebiten.KeyControl, ebiten.KeyEnter, ebiten.KeyShift} {
		if inpututil.IsKeyJustPressed(k) {
			c.inputKeyDown(k)
		} else if inpututil.IsKeyJustReleased(k) {
			c.inputKeyUp(k)
		}
	}
}

func (c *Context) Draw(screen *ebiten.Image) {
	target := screen
	for _, cmd := range c.CommandList {
		switch cmd.Type {
		case CommandRect:
			vector.DrawFilledRect(
				target,
				float32(cmd.Rect.Rect.Min.X),
				float32(cmd.Rect.Rect.Min.Y),
				float32(cmd.Rect.Rect.Dx()),
				float32(cmd.Rect.Rect.Dy()),
				cmd.Rect.Color,
				false,
			)
		case CommandText:
			geom := ebiten.GeoM{}
			geom.Translate(
				float64(cmd.Text.Pos.X),
				float64(cmd.Text.Pos.Y),
			)
			cs := ebiten.ColorScale{}
			cs.ScaleWithColor(cmd.Text.Color)
			text.Draw(target, cmd.Text.Str, face, &text.DrawOptions{
				DrawImageOptions: ebiten.DrawImageOptions{
					GeoM:       geom,
					ColorScale: cs,
				},
			})
		case CommandIcon:
			// TODO: Draw icon images.
			vector.DrawFilledRect(
				target,
				float32(cmd.Icon.Rect.Min.X),
				float32(cmd.Icon.Rect.Min.Y),
				float32(cmd.Icon.Rect.Dx()),
				float32(cmd.Icon.Rect.Dy()),
				cmd.Icon.Color,
				false,
			)
		case CommandClip:
			target = screen.SubImage(image.Rect(
				cmd.Clip.Rect.Min.X,
				cmd.Clip.Rect.Min.Y,
				min(cmd.Clip.Rect.Max.X, screen.Bounds().Dx()),
				min(cmd.Clip.Rect.Max.Y, screen.Bounds().Dy()),
			)).(*ebiten.Image)
		}
	}
}
