package analysis

import (
	page "Gio_UI/UI/app"
	"Gio_UI/UI/icon"
	"bytes"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"github.com/wcharczuk/go-chart/v2"
	"image"
	"image/color"
	"image/png"
	"log"
)

type Page struct {
	widget.List
	*page.Router
}

// New constructs a Page with the provided router.
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
		Name: "Analysis",
		Icon: icon.AnalysisIcon,
	}
}
func (p *Page) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	return layout.Flex{
		Axis:    layout.Vertical,
		Spacing: layout.SpaceEnd,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{
				Axis:    layout.Horizontal,
				Spacing: layout.SpaceBetween,
			}.Layout(gtx,
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return createInfoCard(th, gtx, "User Count", "2,000", "Total User Count 120,000", color.NRGBA{0x36, 0x9E, 0xDC, 0xFF})
				}),
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return createInfoCard(th, gtx, "Page Views", "20,000", "Total Page Views 500,000", color.NRGBA{0x4A, 0xBB, 0xD9, 0xFF})
				}),
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return createInfoCard(th, gtx, "Downloads", "8,000", "Total Downloads 120,000", color.NRGBA{0x53, 0x6D, 0xFE, 0xFF})
				}),
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return createInfoCard(th, gtx, "Usage", "5,000", "Total Usage 50,000", color.NRGBA{0x60, 0x9B, 0xCE, 0xFF})
				}),
			)
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			// Khu vực biểu đồ
			return layout.Inset{
				Top:    unit.Dp(16),
				Bottom: unit.Dp(16),
				Left:   unit.Dp(8),
				Right:  unit.Dp(8),
			}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				barchartImg := drawBarChart()
				imgOp := paint.NewImageOp(barchartImg)

				return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					img := widget.Image{
						Src:   imgOp,
						Scale: 1,
					}
					return img.Layout(gtx)
				})
			})
		}),
	)
}
func createInfoCard(th *material.Theme, gtx layout.Context, title, value, subtext string, bgColor color.NRGBA) layout.Dimensions {
	return layout.Inset{
		Top:    unit.Dp(8),
		Bottom: unit.Dp(8),
		Left:   unit.Dp(8),
		Right:  unit.Dp(8),
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		gtx.Constraints.Max.X = 240
		gtx.Constraints.Max.Y = 120

		return widget.Border{
			Color:        color.NRGBA{R: 0, G: 0, B: 0, A: 255},
			CornerRadius: unit.Dp(16),
			Width:        unit.Dp(2),
		}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Stack{}.Layout(gtx,
				layout.Expanded(func(gtx layout.Context) layout.Dimensions {
					return fillBackground(gtx, bgColor)
				}),
				layout.Stacked(func(gtx layout.Context) layout.Dimensions {
					return layout.UniformInset(unit.Dp(12)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{
							Axis:    layout.Vertical,
							Spacing: layout.SpaceBetween,
						}.Layout(gtx,
							layout.Rigid(material.Label(th, unit.Sp(14), title).Layout),
							layout.Rigid(material.Label(th, unit.Sp(24), value).Layout),
							layout.Rigid(material.Label(th, unit.Sp(12), subtext).Layout),
						)
					})
				}),
			)
		})
	})
}

func fillBackground(gtx layout.Context, color color.NRGBA) layout.Dimensions {
	rect := clip.Rect{Max: gtx.Constraints.Max}.Op()
	paint.FillShape(gtx.Ops, color, rect)
	return layout.Dimensions{Size: gtx.Constraints.Max}
}

func drawBarChart() image.Image {
	barChart := chart.BarChart{
		Title: "Sample Bar Chart",
		Background: chart.Style{
			Padding: chart.Box{
				Top:    40,
				Left:   20,
				Right:  20,
				Bottom: 20,
			},
		},
		Height:   200,
		Width:    400,
		BarWidth: 60,
		Bars: []chart.Value{
			{Value: 5, Label: "A"},
			{Value: 10, Label: "B"},
			{Value: 15, Label: "C"},
		},
	}

	var buffer bytes.Buffer
	err := barChart.Render(chart.PNG, &buffer)
	if err != nil {
		log.Fatalf("failed to render chart: %v", err)
	}

	img, err := png.Decode(&buffer)
	if err != nil {
		log.Fatalf("failed to decode png: %v", err)
	}

	return img
}
