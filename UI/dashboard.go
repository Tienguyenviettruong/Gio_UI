package UI

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
)

func LayoutToolbar(gtx layout.Context, th *material.Theme) layout.Dimensions {
	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		//layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		//	return material.Button(th, &Button1, "Button 1").Layout(gtx)
		//}),
		//layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		//	return material.Button(th, &Button2, "Button 2").Layout(gtx)
		//}),
		//layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		//	return material.Button(th, &Button3, "Button 3").Layout(gtx)
		//}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			var navList widget.List
			//var selected int

			navItems := []component.NavItem{
				component.NavItem{Name: "Home", Icon: nil},
				//component.NavItem{Name: "Profile", Icon: nil},
				//component.NavItem{Name: "Settings", Icon: nil},
			}

			for _, item := range navItems {
				navList.Layout(gtx, 1, func(gtx layout.Context, index int) layout.Dimensions {
					button := material.Button(th, &widget.Clickable{}, item.Name)
					//if widget.Clickable{}.Clicked(gtx) {
					//	selected = i
					//}
					return button.Layout(gtx)
				})
			}

			return layout.Dimensions{}
		}),
	)
}
