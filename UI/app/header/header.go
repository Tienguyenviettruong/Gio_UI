package header

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"

	alo "Gio_UI/UI"
	page "Gio_UI/UI/app"
	"Gio_UI/UI/icon"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

// Page holds the state for a page demonstrating the features of
// the AppBar component.
type Page struct {
	heartBtn, plusBtn, contextBtn          widget.Clickable
	exampleOverflowState, red, green, blue widget.Clickable
	bottomBar, customNavIcon               widget.Bool
	favorited                              bool
	widget.List
	*page.Router
	nonModalDrawer widget.Bool
}

// New constructs a Page with the provided router.
func New(router *page.Router) *Page {
	return &Page{
		Router: router,
	}
}

var _ page.Page = &Page{}

func (p *Page) Actions() []component.AppBarAction {
	return []component.AppBarAction{
		{
			OverflowAction: component.OverflowAction{
				Name: "Favorite",
				Tag:  &p.heartBtn,
			},
			Layout: func(gtx layout.Context, bg, fg color.NRGBA) layout.Dimensions {
				if p.heartBtn.Clicked(gtx) {
					p.favorited = !p.favorited
				}
				btn := component.SimpleIconButton(bg, fg, &p.heartBtn, icon.HeartIcon)
				btn.Background = bg
				if p.favorited {
					btn.Color = color.NRGBA{R: 255, G: 200, B: 100, A: 255}
				} else {
					btn.Color = fg
				}
				return btn.Layout(gtx)
			},
		},
		component.SimpleIconAction(&p.plusBtn, icon.PlusIcon,
			component.OverflowAction{
				Name: "Create",
				Tag:  &p.plusBtn,
			},
		),
	}
}

func (p *Page) Overflow() []component.OverflowAction {
	return []component.OverflowAction{
		{
			Name: "Example 1",
			Tag:  &p.exampleOverflowState,
		},
		{
			Name: "Example 2",
			Tag:  &p.exampleOverflowState,
		},
	}
}

func (p *Page) NavItem() component.NavItem {
	return component.NavItem{
		Name: "Setting",
		Icon: icon.SettingsIcon,
	}
}

const (
	settingNameColumnWidth    = 2
	settingDetailsColumnWidth = 4 - settingNameColumnWidth
)

func (p *Page) Layout(gtx C, th *material.Theme) D {
	p.List.Axis = layout.Vertical
	return material.List(th, &p.List).Layout(gtx, 1, func(gtx C, _ int) D {
		return layout.Flex{
			Alignment: layout.Middle,
			Axis:      layout.Vertical,
		}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return alo.DetailRow{}.Layout(gtx,
					material.Body1(th, "Use non-modal drawer").Layout,
					func(gtx alo.C) alo.D {
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
			layout.Rigid(func(gtx C) D {
				return alo.DetailRow{}.Layout(gtx, material.Body1(th, "Contextual App Bar").Layout, func(gtx C) D {
					if p.contextBtn.Clicked(gtx) {
						p.Router.AppBar.SetContextualActions(
							[]component.AppBarAction{
								component.SimpleIconAction(&p.red, icon.HeartIcon,
									component.OverflowAction{
										Name: "House",
										Tag:  &p.red,
									},
								),
							},
							[]component.OverflowAction{
								{
									Name: "foo",
									Tag:  &p.blue,
								},
								{
									Name: "bar",
									Tag:  &p.green,
								},
							},
						)
						p.Router.AppBar.StartContextual(gtx.Now, "Contextual Title")
					}
					return material.Button(th, &p.contextBtn, "Trigger").Layout(gtx)
				})
			}),
			layout.Rigid(func(gtx C) D {
				return alo.DetailRow{}.Layout(gtx,
					material.Body1(th, "Bottom App Bar").Layout,
					func(gtx C) D {
						if p.bottomBar.Update(gtx) {
							if p.bottomBar.Value {
								p.Router.ModalNavDrawer.Anchor = component.Bottom
								p.Router.AppBar.Anchor = component.Bottom
							} else {
								p.Router.ModalNavDrawer.Anchor = component.Top
								p.Router.AppBar.Anchor = component.Top
							}
							p.Router.BottomBar = p.bottomBar.Value
						}

						return material.Switch(th, &p.bottomBar, "Use Bottom App Bar").Layout(gtx)
					})
			}),
			layout.Rigid(func(gtx C) D {
				return alo.DetailRow{}.Layout(gtx,
					material.Body1(th, "Custom Navigation Icon").Layout,
					func(gtx C) D {
						if p.customNavIcon.Update(gtx) {
							if p.customNavIcon.Value {
								p.Router.AppBar.NavigationIcon = icon.HomeIcon
							} else {
								p.Router.AppBar.NavigationIcon = icon.MenuIcon
							}
						}
						return material.Switch(th, &p.customNavIcon, "Use Custom Navigation Icon").Layout(gtx)
					})
			}),
		)
	})
}
