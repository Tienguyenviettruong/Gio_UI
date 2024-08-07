package main

import (
	"Gio_UI/UI"
	page "Gio_UI/UI/app"
	header "Gio_UI/UI/app/header"
	app2 "Gio_UI/UI/app/toolbar"
	"flag"
	//"gioui.org/example/component/pages/appbar"
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
	router := page.NewRouter()
	router.Register(0, header.New(&router))
	router.Register(1, app2.New(&router))
	//router.Register(2, textfield.New(&router))
	//router.Register(3, menu.New(&router))
	//router.Register(4, discloser.New(&router))
	//router.Register(5, about.New(&router))
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
				router.Layout(gtx, th)
			}
			e.Frame(gtx.Ops)
		}
	}
}
