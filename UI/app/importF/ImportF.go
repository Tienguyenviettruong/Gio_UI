package importF

import (
	page "Gio_UI/UI/app"
	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"log"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

type Page struct {
	importButton widget.Clickable
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
		Icon: p.Icon.CheckBoxChecked,
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
					btn := material.Button(th, &p.importButton, "Edit")
					btn.TextSize = unit.Sp(12) // Giảm kích thước font của text
					return btn.Layout(gtx)
				})
			})
		}),
	)
}

func (p *Page) openNewWindow() {
	go func() {
		newWindow := new(app.Window)
		newWindow.Option(app.Size(unit.Dp(800), unit.Dp(600)))
		if err := p.run(newWindow); err != nil {
			log.Fatal(err)
		}
	}()
}

func (p *Page) run(w *app.Window) error {
	var ops op.Ops
	for {
		e := w.Event()
		switch e := e.(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			layout.Flex{}.Layout(gtx,
				layout.Flexed(1, func(gtx C) D {
					return layout.Inset{Left: unit.Dp(108)}.Layout(gtx, func(gtx C) D {
						return material.Body1(p.Theme, "This is the new window").Layout(gtx)
					})
				}),
			)
			e.Frame(gtx.Ops)
		}
	}
	return nil
}
