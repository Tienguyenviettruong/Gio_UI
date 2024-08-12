package main

import (
	"Gio_UI/UI"
	page "Gio_UI/UI/app"
	"Gio_UI/UI/app/Tree"
	"Gio_UI/UI/app/chart"
	"Gio_UI/UI/app/header"
	"Gio_UI/UI/app/importF"
	"Gio_UI/UI/app/menu"
	"Gio_UI/UI/app/table"
	app2 "Gio_UI/UI/app/toolbar"
	"Gio_UI/UI/icon"
	"flag"
	"gioui.org/x/component"
	"image/color"

	"gioui.org/widget/material"
	"log"
	"os"
	"time"

	"gioui.org/app"
	"gioui.org/op"
	"gioui.org/unit"
)

func main() {
	flag.Parse()
	UI.ProgressIncrementer = make(chan float32)

	go func() {
		for {
			time.Sleep(time.Second)
			UI.ProgressIncrementer <- 0.1
		}
	}()

	go func() {
		w := new(app.Window)
		w.Option(app.Size(unit.Dp(1200), unit.Dp(800)))
		if err := loop(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func loop(w *app.Window) error {
	th := material.NewTheme()
	var ops op.Ops
	menuItems := []component.MenuItemStyle{
		{Icon: icon.SettingsIcon,
			Label:      material.Label(th, unit.Sp(16), "Menu 1"),
			HoverColor: color.NRGBA{R: 0xFF, G: 0x00, B: 0x00, A: 0xFF}, // Màu nền đỏ khi hover
			//LabelInset: layout.Inset{Top: unit.Dp(10), Bottom: unit.Dp(10), Left: unit.Dp(10), Right: unit.Dp(10)},
		},
		{
			Label:      material.Label(th, unit.Sp(16), "Menu 2"),
			HoverColor: color.NRGBA{R: 0x00, G: 0xFF, B: 0x00, A: 0xFF}, // Màu nền xanh lá khi hover
			//LabelInset: layout.Inset{Top: unit.Dp(10), Bottom: unit.Dp(10), Left: unit.Dp(10), Right: unit.Dp(10)},
		},
		{
			Label:      material.Label(th, unit.Sp(16), "Menu 3"),
			HoverColor: color.NRGBA{R: 0x00, G: 0x00, B: 0xFF, A: 0xFF}, // Màu nền xanh dương khi hover

			//LabelInset: layout.Inset{Top: unit.Dp(10), Bottom: unit.Dp(10), Left: unit.Dp(10), Right: unit.Dp(10)},
		},
	}
	router := page.NewRouter()
	router.Register(0, header.New(&router))
	router.Register(1, app2.New(&router))
	router.Register(2, importF.New(&router, th))
	//router.Register(3, dashboard.New(&router))
	//router.Register(2, textfield.New(&router))
	router.Register(3, menu.New(&router))
	router.Register(4, Tree.New(&router))
	router.Register(5, chart.New(&router))
	router.Register(6, Table.New(&router))
	for {
		e := w.Event()
		switch e := e.(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			if UI.LoginScreen {
				//UI.LayoutLogin(gtx, th)

				router.Layout(gtx, th)
			} else {
				//UI.DisplayMenu(gtx, th, menuItems)
			}
			e.Frame(gtx.Ops)
		}
	}
}
