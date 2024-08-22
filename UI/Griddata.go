package UI

import (
	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"image"
	"image/color"
	"strconv"
	"time"
)

//type (
//	C = layout.Context
//	D = layout.Dimensions
//)

type FrameTiming struct {
	Start, End      time.Time
	FrameCount      int
	FramesPerSecond float64
	xxxx            string
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

//
//func loop(w *app.Window) error {
//	th := material.NewTheme()
//	th.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
//	var (
//		ops  op.Ops
//		grid component.GridState
//	)
//	timingWindow := time.Second
//	timings := []FrameTiming{}
//	frameCounter := 0
//	timingStart := time.Time{}
//	for {
//		switch e := w.Event().(type) {
//		case app.DestroyEvent:
//			return e.Err
//		case app.FrameEvent:
//			gtx := app.NewContext(&ops, e)
//			gtx.Execute(op.InvalidateCmd{})
//			if timingStart == (time.Time{}) {
//				timingStart = gtx.Now
//			}
//			if interval := gtx.Now.Sub(timingStart); interval >= timingWindow {
//				timings = append(timings, FrameTiming{
//					Start:           timingStart,
//					End:             gtx.Now,
//					FrameCount:      frameCounter,
//					FramesPerSecond: float64(frameCounter) / interval.Seconds(),
//				})
//				frameCounter = 0
//				timingStart = gtx.Now
//			}
//			LayoutTable(th, gtx, timings, &grid)
//			e.Frame(gtx.Ops)
//			frameCounter++
//		}
//	}
//}

var headingText = []string{"Start", "End", "Frames", "FPS", "Col4", "Col5"}

func LayoutTable(th *material.Theme, gtx C, timings []FrameTiming, grid *component.GridState) D {
	// Configure width based on available space and a minimum size.
	minSize := gtx.Dp(unit.Dp(200))
	border := widget.Border{
		Color: color.NRGBA{A: 255},
		Width: unit.Dp(0),
	}

	inset := layout.UniformInset(unit.Dp(0.5))

	// Configure a label styled to be a heading.
	headingLabel := material.Body1(th, "")
	headingLabel.Font.Weight = font.Bold
	headingLabel.Alignment = text.Middle
	headingLabel.MaxLines = 2

	// Configure a label styled to be a data element.
	dataLabel := material.Body1(th, "")
	dataLabel.Font.Typeface = "Go Mono"
	dataLabel.MaxLines = 0
	dataLabel.Alignment = text.Middle

	// Measure the height of a heading row.
	orig := gtx.Constraints
	gtx.Constraints.Min = image.Point{}
	macro := op.Record(gtx.Ops)
	dims := inset.Layout(gtx, headingLabel.Layout)
	_ = macro.Stop()
	gtx.Constraints = orig

	return component.Table(th, grid).Layout(gtx, len(timings), 6, // Change to 6 columns
		func(axis layout.Axis, index, constraint int) int {
			widthUnit := max(int(float32(constraint)/6), minSize) // Update division to 6
			switch axis {
			case layout.Horizontal:
				return int(widthUnit) // All columns have equal width
			case layout.Vertical:
				return dims.Size.Y // Row height
			default:
				return 0
			}
		},
		func(gtx C, col int) D {
			return border.Layout(gtx, func(gtx C) D {
				return inset.Layout(gtx, func(gtx C) D {
					if col < len(headingText) {
						headingLabel.Text = headingText[col]
					}
					return headingLabel.Layout(gtx)
				})
			})
		},
		func(gtx C, row, col int) D {
			return inset.Layout(gtx, func(gtx C) D {
				if row < len(timings) {
					timing := timings[row]
					switch col {
					case 0:
						dataLabel.Text = timing.Start.Format("15:04:05.000000")
					case 1:
						dataLabel.Text = timing.End.Format("15:04:05.000000")
					case 2:
						dataLabel.Text = strconv.Itoa(timing.FrameCount)
					case 3:
						dataLabel.Text = strconv.FormatFloat(timing.FramesPerSecond, 'f', 2, 64)
					case 4:
						dataLabel.Text = "Data for Col4" // Placeholder data
					case 5:
						dataLabel.Text = "Data for Col5" // Placeholder data
					default:
						dataLabel.Text = ""
					}
				}
				return dataLabel.Layout(gtx)
			})
		},
	)
}
