package UI

import (
	"Gio_UI/UI/icon"
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
)

var (
	navItems = []navItem{
		{item: component.NavItem{Name: "Home", Icon: icon.HomeIcon}},
		{item: component.NavItem{Name: "Profile", Icon: icon.EditIcon}},
		{item: component.NavItem{Name: "Settings"}},
	}
	selected int
)

type navItem struct {
	item      component.NavItem
	clickable widget.Clickable
}

func LayoutToolbar(gtx layout.Context, th *material.Theme) layout.Dimensions {
	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			var navList widget.List

			navList.Layout(gtx, len(navItems), func(gtx layout.Context, index int) layout.Dimensions {
				ni := &navItems[index]
				//button := material.Button(th, &ni.clickable, ni.item.Name)
				button := material.IconButton(th, &ni.clickable, icon.MenuIcon, "test")
				if ni.clickable.Clicked(gtx) {
					selected = index
				}
				return button.Layout(gtx)
			})

			// Optional: display selected nav item
			//layout.Flex{}.Layout(gtx,
			//	layout.Rigid(material.H6(th, "Selected: "+navItems[selected].item.Name).Layout),
			//)

			return layout.Dimensions{}
		}),
	)
}
