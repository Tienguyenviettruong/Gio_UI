package toolbar

import (
	"Gio_UI/UI"
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"

	page "Gio_UI/UI/app"
	//alo "gioui.org/example/component/applayout"
	"Gio_UI/UI/icon"
)

//
//type (
//	C = layout.Context
//	D = layout.Dimensions
//)

// Page holds the state for a page demonstrating the features of
// the NavDrawer component.
type Page struct {
	nonModalDrawer widget.Bool
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
		Name: "Setting",
		Icon: icon.SettingsIcon,
	}
}

func (p *Page) Layout(gtx UI.C, th *material.Theme) UI.D {
	p.List.Axis = layout.Vertical
	return material.List(th, &p.List).Layout(gtx, 1, func(gtx UI.C, _ int) UI.D {
		return layout.Flex{
			Alignment: layout.Middle,
			Axis:      layout.Vertical,
		}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return UI.DetailRow{}.Layout(gtx,
					material.Body1(th, "Use non-modal drawer").Layout,
					func(gtx UI.C) UI.D {
						if p.nonModalDrawer.Update(gtx) {
							p.Router.NonModalDrawer = p.nonModalDrawer.Value
							if p.nonModalDrawer.Value {
								p.Router.NavAnim.Appear(gtx.Now)
							} else {
								p.Router.NavAnim.Disappear(gtx.Now)
							}
						}
						return material.Switch(th, &p.nonModalDrawer, "Use Non-Modal Navigation Drawer").Layout(gtx)
					})
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return UI.DetailRow{}.Layout(gtx,
					material.Body1(th, "Drag to Close").Layout,
					material.Body2(th, "You can close the modal nav drawer by dragging it to the left.").Layout)
			}),
		)
	})
}
