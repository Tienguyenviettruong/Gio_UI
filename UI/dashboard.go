package UI

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"gioui.org/x/component"
)

func DisplayMenu(gtx layout.Context, th *material.Theme, menuItems []component.MenuItemStyle) layout.Dimensions {
	// Sử dụng Flex layout để xếp các mục theo chiều dọc
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		// Lặp qua các mục menu và hiển thị chúng
		func() []layout.FlexChild {
			children := make([]layout.FlexChild, len(menuItems))
			for i, menuItem := range menuItems {
				i := i // Capturing loop variable

				children[i] = layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					// Hiển thị từng mục menu
					return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return menuItem.Label.Layout(gtx)
					})
				})
			}
			return children
		}()...,
	)
}
