package UI

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

var (
	menuItems = []string{
		"Auth API",
		"Certs API",
		"ID API",
		"Todo API",
		"WhoAmI",
	}
	menuList = widget.List{
		List: layout.List{
			Axis: layout.Vertical,
		},
	}
)

func Workspace(gtx layout.Context, th *material.Theme) layout.Dimensions {
	return layout.Flex{}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Left: unit.Dp(16), Top: unit.Dp(16)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return material.List(th, &menuList).Layout(gtx, len(menuItems), func(gtx layout.Context, index int) layout.Dimensions {
					return layout.UniformInset(unit.Dp(8)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						btn := material.Button(th, new(widget.Clickable), menuItems[index])
						return btn.Layout(gtx)
					})
				})
			})
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return material.H6(th, "Main Content Area").Layout(gtx)
			})
		}),
	)
}
