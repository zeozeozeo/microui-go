// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2024 The Ebitengine Authors

package main

import (
	"fmt"
	"image"
	"image/color"

	"github.com/ebitengine/microui"
)

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (g *Game) WriteLog(text string) {
	if len(g.logBuf) > 0 {
		g.logBuf += "\n"
	}
	g.logBuf += text
	g.logUpdated = true
}

func (g *Game) TestWindow() {
	if g.ctx.BeginWindow("Demo Window", image.Rect(40, 40, 340, 490)) {
		defer g.ctx.EndWindow()
		win := g.ctx.GetCurrentContainer()
		win.Rect.Max.X = win.Rect.Min.X + max(win.Rect.Dx(), 240)
		win.Rect.Max.Y = win.Rect.Min.Y + max(win.Rect.Dy(), 300)

		/* window info */
		if g.ctx.Header("Window Info") {
			win := g.ctx.GetCurrentContainer()
			g.ctx.LayoutRow(2, []int{54, -1}, 0)
			g.ctx.Label("Position:")
			g.ctx.Label(fmt.Sprintf("%d, %d", win.Rect.Min.X, win.Rect.Min.Y))
			g.ctx.Label("Size:")
			g.ctx.Label(fmt.Sprintf("%d, %d", win.Rect.Dx(), win.Rect.Dy()))
		}

		/* labels + buttons */
		if g.ctx.HeaderEx("Test Buttons", microui.OptExpanded) != 0 {
			g.ctx.LayoutRow(3, []int{100, -110, -1}, 0)
			g.ctx.Label("Test buttons 1:")
			if g.ctx.Button("Button 1") {
				g.WriteLog("Pressed button 1")
			}
			if g.ctx.Button("Button 2") {
				g.WriteLog("Pressed button 2")
			}
			g.ctx.Label("Test buttons 2:")
			if g.ctx.Button("Button 3") {
				g.WriteLog("Pressed button 3")
			}
			if g.ctx.Button("Popup") {
				g.ctx.OpenPopup("Test Popup")
			}
			if g.ctx.BeginPopup("Test Popup") != 0 {
				g.ctx.Button("Hello")
				g.ctx.Button("World")
				g.ctx.EndPopup()
			}
		}

		/* tree */
		if g.ctx.HeaderEx("Tree and Text", microui.OptExpanded) != 0 {
			g.ctx.LayoutRow(2, []int{140, -1}, 0)
			g.ctx.LayoutBeginColumn()
			if g.ctx.BeginTreeNode("Test 1") {
				if g.ctx.BeginTreeNode("Test 1a") {
					g.ctx.Label("Hello")
					g.ctx.Label("World")
					g.ctx.EndTreeNode()
				}
				if g.ctx.BeginTreeNode("Test 1b") {
					if g.ctx.Button("Button 1") {
						g.WriteLog("Pressed button 1")
					}
					if g.ctx.Button("Button 2") {
						g.WriteLog("Pressed button 2")
					}
					g.ctx.EndTreeNode()
				}
				g.ctx.EndTreeNode()
			}
			if g.ctx.BeginTreeNode("Test 2") {
				g.ctx.LayoutRow(2, []int{54, 54}, 0)
				if g.ctx.Button("Button 3") {
					g.WriteLog("Pressed button 3")
				}
				if g.ctx.Button("Button 4") {
					g.WriteLog("Pressed button 4")
				}
				if g.ctx.Button("Button 5") {
					g.WriteLog("Pressed button 5")
				}
				if g.ctx.Button("Button 6") {
					g.WriteLog("Pressed button 6")
				}
				g.ctx.EndTreeNode()
			}
			if g.ctx.BeginTreeNode("Test 3") {
				g.ctx.Checkbox("Checkbox 1", &g.checks[0])
				g.ctx.Checkbox("Checkbox 2", &g.checks[1])
				g.ctx.Checkbox("Checkbox 3", &g.checks[2])
				g.ctx.EndTreeNode()
			}
			g.ctx.LayoutEndColumn()

			g.ctx.LayoutBeginColumn()
			g.ctx.LayoutRow(1, []int{-1}, 0)
			g.ctx.Text("Lorem ipsum dolor sit amet, consectetur adipiscing " +
				"elit. Maecenas lacinia, sem eu lacinia molestie, mi risus faucibus " +
				"ipsum, eu varius magna felis a nulla.")
			g.ctx.LayoutEndColumn()
		}

		/* background color sliders */
		if g.ctx.HeaderEx("Background Color", microui.OptExpanded) != 0 {
			g.ctx.LayoutRow(2, []int{-78, -1}, 74)
			/* sliders */
			g.ctx.LayoutBeginColumn()
			g.ctx.LayoutRow(2, []int{46, -1}, 0)
			g.ctx.Label("Red:")
			g.ctx.Slider(&g.bg[0], 0, 255)
			g.ctx.Label("Green:")
			g.ctx.Slider(&g.bg[1], 0, 255)
			g.ctx.Label("Blue:")
			g.ctx.Slider(&g.bg[2], 0, 255)
			g.ctx.LayoutEndColumn()
			/* color preview */
			r := g.ctx.LayoutNext()
			g.ctx.DrawRect(r, color.RGBA{byte(g.bg[0]), byte(g.bg[1]), byte(g.bg[2]), 255})
			clr := fmt.Sprintf("#%02X%02X%02X", int(g.bg[0]), int(g.bg[1]), int(g.bg[2]))
			g.ctx.DrawControlText(clr, r, microui.ColorText, microui.OptAlignCenter)
		}
	}
}

func (g *Game) LogWindow() {
	if g.ctx.BeginWindow("Log Window", image.Rect(350, 40, 650, 240)) {
		defer g.ctx.EndWindow()
		/* output text panel */
		g.ctx.LayoutRow(1, []int{-1}, -25)
		g.ctx.BeginPanel("Log Output")
		panel := g.ctx.GetCurrentContainer()
		g.ctx.LayoutRow(1, []int{-1}, -1)
		g.ctx.Text(g.logBuf)
		g.ctx.EndPanel()
		if g.logUpdated {
			panel.Scroll.Y = panel.ContentSize.Y
			g.logUpdated = false
		}

		/* input textbox + submit button */
		var submitted bool
		g.ctx.LayoutRow(2, []int{-70, -1}, 0)
		if g.ctx.TextBox(&g.logSubmitBuf)&microui.ResSubmit != 0 {
			g.ctx.SetFocus(g.ctx.LastID)
			submitted = true
		}
		if g.ctx.Button("Submit") {
			submitted = true
		}
		if submitted {
			g.WriteLog(g.logSubmitBuf)
			g.logSubmitBuf = ""
		}
	}
}

func (g *Game) byteSlider(fvalue *float64, value *byte, low, high byte) int {
	*fvalue = float64(*value)
	res := g.ctx.SliderEx(fvalue, float64(low), float64(high), 0, "%.0f", microui.OptAlignCenter)
	*value = byte(*fvalue)

	return res
}

var (
	fcolors = [14]struct {
		R, G, B, A float64
	}{}
	colors = []struct {
		Label   string
		ColorID int
	}{
		{"text:", microui.ColorText},
		{"border:", microui.ColorBorder},
		{"windowbg:", microui.ColorWindowBG},
		{"titlebg:", microui.ColorTitleBG},
		{"titletext:", microui.ColorTitleText},
		{"panelbg:", microui.ColorPanelBG},
		{"button:", microui.ColorButton},
		{"buttonhover:", microui.ColorButtonHover},
		{"buttonfocus:", microui.ColorButtonFocus},
		{"base:", microui.ColorBase},
		{"basehover:", microui.ColorBaseHover},
		{"basefocus:", microui.ColorBaseFocus},
		{"scrollbase:", microui.ColorScrollBase},
		{"scrollthumb:", microui.ColorScrollThumb},
	}
)

func (g *Game) StyleWindow() {
	if g.ctx.BeginWindow("Style Editor", image.Rect(350, 250, 650, 490)) {
		sw := int(float64(g.ctx.GetCurrentContainer().Body.Dx()) * 0.14)
		g.ctx.LayoutRow(6, []int{80, sw, sw, sw, sw, -1}, 0)
		for _, c := range colors {
			g.ctx.Label(c.Label)
			g.byteSlider(&fcolors[c.ColorID].R, &g.ctx.Style.Colors[c.ColorID].R, 0, 255)
			g.byteSlider(&fcolors[c.ColorID].G, &g.ctx.Style.Colors[c.ColorID].G, 0, 255)
			g.byteSlider(&fcolors[c.ColorID].B, &g.ctx.Style.Colors[c.ColorID].B, 0, 255)
			g.byteSlider(&fcolors[c.ColorID].A, &g.ctx.Style.Colors[c.ColorID].A, 0, 255)
			g.ctx.DrawRect(g.ctx.LayoutNext(), g.ctx.Style.Colors[c.ColorID])
		}
		g.ctx.EndWindow()
	}
}

func (g *Game) ProcessFrame() {
	g.ctx.Begin()
	g.TestWindow()
	g.LogWindow()
	g.StyleWindow()
	g.ctx.End()
}
