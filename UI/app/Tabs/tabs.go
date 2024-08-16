package Tabs

import (
	"fmt"
	"gioui.org/f32"
	"image"
	"image/color"
	"math"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"

	page "Gio_UI/UI/app"
	"Gio_UI/UI/icon"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

var (
	tabs Tabs
	//slider Slider
)

type Tabs struct {
	list     layout.List
	tabs     []Tab
	selected int
}

type Tab struct {
	btn   widget.Clickable
	Title string
}

func init() {
	for i := 1; i <= 10; i++ {
		tabs.tabs = append(tabs.tabs,
			Tab{Title: fmt.Sprintf("Tab %d", i)},
		)
	}
}

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
		Name: "Tabs",
		Icon: icon.MenuIcon,
	}
}

func (p *Page) Layout(gtx C, th *material.Theme) D {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return tabs.list.Layout(gtx, len(tabs.tabs), func(gtx C, tabIdx int) D {
				t := &tabs.tabs[tabIdx]
				if t.btn.Clicked(gtx) {
					//if tabs.selected < tabIdx {
					//	slider.PushLeft()
					//} else if tabs.selected > tabIdx {
					//	slider.PushRight()
					//}
					tabs.selected = tabIdx
				}
				var tabWidth int
				return layout.Stack{Alignment: layout.S}.Layout(gtx,
					layout.Stacked(func(gtx C) D {
						dims := material.Clickable(gtx, &t.btn, func(gtx C) D {
							return layout.UniformInset(unit.Dp(12)).Layout(gtx,
								material.H6(th, t.Title).Layout,
							)
						})
						tabWidth = dims.Size.X
						return dims
					}),
					layout.Stacked(func(gtx C) D {
						if tabs.selected != tabIdx {
							return layout.Dimensions{}
						}
						tabHeight := gtx.Dp(unit.Dp(4))
						tabRect := image.Rect(0, 0, tabWidth, tabHeight)
						paint.FillShape(gtx.Ops, th.Palette.ContrastBg, clip.Rect(tabRect).Op())
						return layout.Dimensions{
							Size: image.Point{X: tabWidth, Y: tabHeight},
						}
					}),
				)
			})
		}),
		layout.Flexed(1, func(gtx C) D {
			fill(gtx, dynamicColor(tabs.selected), dynamicColor(tabs.selected+1))
			return layout.Center.Layout(gtx,
				material.H1(th, fmt.Sprintf("Tab content #%d", tabs.selected+1)).Layout,
			)
		}),
	)
}
func fill(gtx layout.Context, col1, col2 color.NRGBA) {
	dr := image.Rectangle{Max: gtx.Constraints.Min}
	paint.FillShape(gtx.Ops,
		color.NRGBA{R: 0, G: 0, B: 0, A: 0xFF},
		clip.Rect(dr).Op(),
	)

	col2.R = byte(float32(col2.R))
	col2.G = byte(float32(col2.G))
	col2.B = byte(float32(col2.B))
	paint.LinearGradientOp{
		Stop1:  f32.Pt(float32(dr.Min.X), 0),
		Stop2:  f32.Pt(float32(dr.Max.X), 0),
		Color1: col1,
		Color2: col2,
	}.Add(gtx.Ops)
	defer clip.Rect(dr).Push(gtx.Ops).Pop()
	paint.PaintOp{}.Add(gtx.Ops)
}

func dynamicColor(i int) color.NRGBA {
	sn, cs := math.Sincos(float64(i) * math.Phi)
	return color.NRGBA{
		R: 0xA0 + byte(0x30*sn),
		G: 0xA0 + byte(0x30*cs),
		B: 0xD0,
		A: 0xFF,
	}
}
