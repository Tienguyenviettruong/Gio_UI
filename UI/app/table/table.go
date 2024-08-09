package table

import (
	"Gio_UI/UI"
	page "Gio_UI/UI/app"
	"Gio_UI/UI/icon"
	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"image/color"
	//"gioui.org/paint"
	"gioui.org/op"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

type Page struct {
	generateTableButton widget.Clickable
	*page.Router
}

func New(router *page.Router) *Page {
	return &Page{
		Router: router,
	}
}

var _ page.Page = &Page{}

func (p *Page) Actions() []component.AppBarAction {
	return []component.AppBarAction{}
}

func (p *Page) Overflow() []component.OverflowAction {
	return []component.OverflowAction{}
}

func (p *Page) NavItem() component.NavItem {
	return component.NavItem{
		Name: "Table",
		Icon: icon.TableIcon, // Bạn cần tự định nghĩa icon này
	}
}

func (p *Page) Layout(gtx C, th *material.Theme) D {
	if p.generateTableButton.Clicked(gtx) {
		// Chức năng tạo bảng không cần thiết phải sử dụng ở đây
	}

	return layout.Flex{}.Layout(gtx,
		layout.Flexed(1, func(gtx C) D {
			return layout.UniformInset(unit.Dp(16)).Layout(gtx, func(gtx C) D {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return UI.DetailRow{}.Layout(gtx,
							material.Body1(th, "Table").Layout,
							func(gtx UI.C) UI.D {
								return material.Button(th, &p.generateTableButton, "Generate Table").Layout(gtx)
							})
					}),
					layout.Rigid(func(gtx C) D {
						return drawTable(gtx, th)
					}),
				)
			})
		}),
	)
}

func drawTable(gtx C, th *material.Theme) D {
	const (
		numRows    = 4
		numCols    = 3
		cellWidth  = 100
		cellHeight = 30
	)

	tableWidth := numCols * cellWidth
	tableHeight := numRows * cellHeight

	var ops op.Ops
	// Draw table background
	defer clip.Rect{Max: f32.Pt(float32(tableWidth), float32(tableHeight))}.Push(&ops).Pop()
	paint.FillShape(&ops, color.NRGBA{R: 255, G: 255, B: 255, A: 255}, clip.Rect{Max: f32.Pt(float32(tableWidth), float32(tableHeight))}.Op())

	// Draw table grid
	for i := 0; i <= numRows; i++ {
		y := float32(i * cellHeight)
		clip.Rect{Max: f32.Pt(float32(tableWidth), y+1)}.Add(&ops)
		paint.FillShape(&ops, color.NRGBA{R: 0, G: 0, B: 0, A: 255}, clip.Rect{Max: f32.Pt(float32(tableWidth), y+1)}.Op())
	}

	for j := 0; j <= numCols; j++ {
		x := float32(j * cellWidth)
		clip.Rect{Max: f32.Pt(x+1, float32(tableHeight))}.Add(&ops)
		paint.FillShape(&ops, color.NRGBA{R: 0, G: 0, B: 0, A: 255}, clip.Rect{Max: f32.Pt(x+1, float32(tableHeight))}.Op())
	}

	// Draw table content
	for i := 0; i < numRows; i++ {
		for j := 0; j < numCols; j++ {
			x := float32(j * cellWidth)
			y := float32(i * cellHeight)
			drawText(&ops, x+cellWidth/2, y+cellHeight/2, "Cell", th)
		}
	}

	return layout.Dimensions{Size: f32.Pt(float32(tableWidth), float32(tableHeight))}
}

func drawText(ops *op.Ops, x, y float32, text string, th *material.Theme) {
	// Define text style
	style := material.Body1(th, text)
	style.Color = color.Black
	gtx := layout.Context{Ops: ops}

	// Layout text
	txtLayout := layout.Center.Layout(gtx, func(gtx C) D {
		return material.Body1(th, text).Layout(gtx)
	})

	// Draw text
	textOp := text.NewLayoutOps(txtLayout)
	textOp.Add(ops)
}
