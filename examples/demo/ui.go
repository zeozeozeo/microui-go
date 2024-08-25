package main

import (
	"fmt"

	"github.com/Zyko0/microui-ebitengine"
)

var (
	logBuf       string
	logSubmitBuf string
	logUpdated   bool
	bg           = [3]microui.Mu_Real{90, 95, 100}
	checks       = [3]bool{true, false, true}
)

func WriteLog(text string) {
	if len(logBuf) > 0 {
		logBuf += "\n"
	}
	logBuf += text
	logUpdated = true
}

func TestWindow(ctx *microui.Context) {
	if ctx.BeginWindow("Demo Window", microui.NewRect(40, 40, 300, 450)) {
		defer ctx.EndWindow()
		win := ctx.GetCurrentContainer()
		win.Rect.W = max(win.Rect.W, 240)
		win.Rect.H = max(win.Rect.H, 300)

		/* window info */
		if ctx.Header("Window Info") {
			win := ctx.GetCurrentContainer()
			ctx.LayoutRow(2, []int{54, -1}, 0)
			ctx.Label("Position:")
			ctx.Label(fmt.Sprintf("%d, %d", win.Rect.X, win.Rect.Y))
			ctx.Label("Size:")
			ctx.Label(fmt.Sprintf("%d, %d", win.Rect.W, win.Rect.H))
		}

		/* labels + buttons */
		if ctx.HeaderEx("Test Buttons", microui.MU_OPT_EXPANDED) > 0 {
			ctx.LayoutRow(3, []int{86, -110, -1}, 0)
			ctx.Label("Test buttons 1:")
			if ctx.Button("Button 1") {
				WriteLog("Pressed button 1")
			}
			if ctx.Button("Button 2") {
				WriteLog("Pressed button 2")
			}
			ctx.Label("Test buttons 2:")
			if ctx.Button("Button 3") {
				WriteLog("Pressed button 3")
			}
			if ctx.Button("Popup") {
				ctx.OpenPopup("Test Popup")
			}
			if ctx.BeginPopup("Test Popup") > 0 {
				ctx.Button("Hello")
				ctx.Button("World")
				ctx.EndPopup()
			}
		}

		/* tree */
		if ctx.HeaderEx("Tree and Text", microui.MU_OPT_EXPANDED) > 0 {
			ctx.LayoutRow(2, []int{140, -1}, 0)
			ctx.LayoutBeginColumn()
			if ctx.BeginTreeNode("Test 1") {
				if ctx.BeginTreeNode("Test 1a") {
					ctx.Label("Hello")
					ctx.Label("World")
					ctx.EndTreeNode()
				}
				if ctx.BeginTreeNode("Test 1b") {
					if ctx.Button("Button 1") {
						WriteLog("Pressed button 1")
					}
					if ctx.Button("Button 2") {
						WriteLog("Pressed button 2")
					}
					ctx.EndTreeNode()
				}
				ctx.EndTreeNode()
			}
			if ctx.BeginTreeNode("Test 2") {
				ctx.LayoutRow(2, []int{54, 54}, 0)
				if ctx.Button("Button 3") {
					WriteLog("Pressed button 3")
				}
				if ctx.Button("Button 4") {
					WriteLog("Pressed button 4")
				}
				if ctx.Button("Button 5") {
					WriteLog("Pressed button 5")
				}
				if ctx.Button("Button 6") {
					WriteLog("Pressed button 6")
				}
				ctx.EndTreeNode()
			}
			if ctx.BeginTreeNode("Test 3") {
				ctx.Checkbox("Checkbox 1", &checks[0])
				ctx.Checkbox("Checkbox 2", &checks[1])
				ctx.Checkbox("Checkbox 3", &checks[2])
				ctx.EndTreeNode()
			}
			ctx.LayoutEndColumn()

			ctx.LayoutBeginColumn()
			ctx.LayoutRow(1, []int{-1}, 0)
			ctx.Text("Lorem ipsum dolor sit amet, consectetur adipiscing " +
				"elit. Maecenas lacinia, sem eu lacinia molestie, mi risus faucibus " +
				"ipsum, eu varius magna felis a nulla.")
			ctx.LayoutEndColumn()
		}

		/* background color sliders */
		if ctx.HeaderEx("Background Color", microui.MU_OPT_EXPANDED) > 0 {
			ctx.LayoutRow(2, []int{-78, -1}, 74)
			/* sliders */
			ctx.LayoutBeginColumn()
			ctx.LayoutRow(2, []int{46, -1}, 0)
			ctx.Label("Red:")
			ctx.Slider(&bg[0], 0, 255)
			ctx.Label("Green:")
			ctx.Slider(&bg[1], 0, 255)
			ctx.Label("Blue:")
			ctx.Slider(&bg[2], 0, 255)
			ctx.LayoutEndColumn()
			/* color preview */
			r := ctx.LayoutNext()
			ctx.DrawRect(r, microui.NewColor(uint8(bg[0]), uint8(bg[1]), uint8(bg[2]), 255))
			clr := fmt.Sprintf("#%02X%02X%02X", int(bg[0]), int(bg[1]), int(bg[2]))
			ctx.DrawControlText(clr, r, microui.MU_COLOR_TEXT, microui.MU_OPT_ALIGNCENTER)
		}
	}
}

func LogWindow(ctx *microui.Context) {
	if ctx.BeginWindow("Log Window", microui.NewRect(350, 40, 300, 200)) {
		defer ctx.EndWindow()
		/* output text panel */
		ctx.LayoutRow(1, []int{-1}, -25)
		ctx.BeginPanel("Log Output")
		panel := ctx.GetCurrentContainer()
		ctx.LayoutRow(1, []int{-1}, -1)
		ctx.Text(logBuf)
		ctx.EndPanel()
		if logUpdated {
			panel.Scroll.Y = panel.ContentSize.Y
			logUpdated = false
		}

		/* input textbox + submit button */
		var submitted bool
		ctx.LayoutRow(2, []int{-70, -1}, 0)
		if ctx.TextBox(&logSubmitBuf)&microui.MU_RES_SUBMIT != 0 {
			ctx.SetFocus(ctx.LastID)
			submitted = true
		}
		if ctx.Button("Submit") {
			submitted = true
		}
		if submitted {
			WriteLog(logSubmitBuf)
			logSubmitBuf = ""
		}
	}
}

func ProcessFrame(ctx *microui.Context) {
	ctx.Begin()
	LogWindow(ctx)
	TestWindow(ctx)
	ctx.End()
}
