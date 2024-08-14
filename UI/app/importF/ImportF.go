package importF

import (
	page "Gio_UI/UI/app"
	"Gio_UI/UI/icon"
	"fmt"
	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"image"
	"image/color"
	"log"
)

type (
	C = layout.Context
	D = layout.Dimensions
)
type Page struct {
	importButton widget.Clickable
	rowsInput    widget.Editor
	colsInput    widget.Editor
	*material.Theme
	*page.Router
}

func New(router *page.Router, th *material.Theme) *Page {
	return &Page{
		Router: router,
		Theme:  th,
	}
}

func (p *Page) NavItem() component.NavItem {
	return component.NavItem{
		Name: "Import",
		Icon: icon.ImportIcon,
	}
}

func (p *Page) Actions() []component.AppBarAction {
	return []component.AppBarAction{}
}

func (p *Page) Overflow() []component.OverflowAction {
	return []component.OverflowAction{}
}

func (p *Page) Layout(gtx C, th *material.Theme) D {
	if p.importButton.Clicked(gtx) {
		go p.openNewWindow()
	}

	return layout.Flex{}.Layout(gtx,
		layout.Flexed(1, func(gtx C) D {
			return layout.NE.Layout(gtx, func(gtx C) D {
				return layout.Inset{
					Top:    unit.Dp(8),
					Bottom: unit.Dp(8),
					Left:   unit.Dp(16),
					Right:  unit.Dp(16),
				}.Layout(gtx, func(gtx C) D {
					return layout.Flex{}.Layout(gtx,
						layout.Flexed(1, func(gtx C) D {
							return material.Editor(th, &p.rowsInput, "Rows").Layout(gtx)
						}),
						layout.Flexed(1, func(gtx C) D {
							return material.Editor(th, &p.colsInput, "Columns").Layout(gtx)
						}),
						layout.Flexed(1, func(gtx C) D {
							btn := material.Button(th, &p.importButton, "Edit")
							btn.TextSize = unit.Sp(12) // Giảm kích thước font của text
							return btn.Layout(gtx)
						}),
					)
				})
			})
		}),
	)
}

func (p *Page) openNewWindow() {
	go func() {
		rows, cols := 0, 0
		if _, err := fmt.Sscanf(p.rowsInput.Text(), "%d", &rows); err != nil {
			log.Println("Invalid number of rows")
			return
		}
		if _, err := fmt.Sscanf(p.colsInput.Text(), "%d", &cols); err != nil {
			log.Println("Invalid number of columns")
			return
		}

		newWindow := new(app.Window)
		newWindow.Option(app.Size(unit.Dp(800), unit.Dp(600)))
		if err := p.run(newWindow, rows, cols); err != nil {
			log.Fatal(err)
		}
	}()
}

func (p *Page) run(w *app.Window, rows, cols int) error {
	var ops op.Ops
	th := material.NewTheme()
	for {
		e := w.Event()
		switch e := e.(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			//layout.Flex{}.Layout(gtx,
			//	layout.Flexed(1, func(gtx C) D {
			//		return layout.Inset{Left: unit.Dp(108)}.Layout(gtx, func(gtx C) D {
			//			return material.Body1(p.Theme, "This is the new window").Layout(gtx)
			//		})
			//	}),
			//)

			p.drawTable(gtx, th, rows, cols)
			e.Frame(gtx.Ops)
		}
	}
	return nil
}
func (p *Page) drawTable(gtx layout.Context, th *material.Theme, rows, cols int) {
	lineColor := color.NRGBA{R: 0, G: 0, B: 255, A: 255} // Màu xanh dương

	rowHeight := float32(60)

	screenWidth := float32(gtx.Constraints.Max.X)
	cellWidth := screenWidth / float32(cols)

	screenHeight := rowHeight * float32(rows)

	strokeWidth := float32(1.0)

	for i := 0; i <= rows; i++ {
		y := float32(i) * rowHeight
		drawLine(gtx, 0, y, screenWidth, y, lineColor, strokeWidth)
	}

	for j := 0; j <= cols; j++ {
		x := float32(j) * cellWidth
		drawLine(gtx, x, 0, x, screenHeight, lineColor, strokeWidth)
	}
	drawText(gtx, fmt.Sprintf("Cell %d,%d", 1, 1), 0, 0, th)
	//drawText(gtx, fmt.Sprintf("Cell %d,%d", 1, 2), 400, 0, th)
	drawText(gtx, fmt.Sprintf("Cell %d,%d", 2, 1), 0, 60, th)
	drawText(gtx, fmt.Sprintf("Cell %d,%d", 2, 2), 400, 0, th)
	drawText(gtx, fmt.Sprintf("Cell %d,%d", 3, 1), 400, 60, th)
	//for i := 0; i < rows; i++ {
	//	for j := 0; j < cols; j++ {
	//		x := float32(j) * cellWidth
	//		y := float32(i) * rowHeight
	//		drawText(gtx, fmt.Sprintf("Cell %d,%d", i+1, j+1), x, y, th)
	//		fmt.Println("Cell %d,%d", i+1, j+1, x, y)
	//	}
	//}
}

func drawLine(gtx layout.Context, x1, y1, x2, y2 float32, clr color.NRGBA, width float32) {
	var path clip.Path
	path.Begin(gtx.Ops)
	path.Move(f32.Pt(x1, y1))
	path.Line(f32.Pt(x2-x1, y2-y1))
	paint.FillShape(gtx.Ops, clr, clip.Stroke{Path: path.End(), Width: width}.Op())
}
func drawText(gtx layout.Context, txt string, x, y float32, th *material.Theme) {
	// Tạo một đối tượng label
	label := material.Label(th, unit.Sp(14), txt)

	// Dịch chuyển vị trí để văn bản được vẽ tại góc trên bên trái của ô
	op.Offset(image.Pt(int(x), int(y))).Add(gtx.Ops)
	label.Layout(gtx)
}
