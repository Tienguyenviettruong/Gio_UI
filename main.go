package main

import (
	"Gio_UI/UI"
	page "Gio_UI/UI/app"
	"Gio_UI/UI/app/Tabs"
	"Gio_UI/UI/app/Tree"
	"Gio_UI/UI/app/chart"
	"Gio_UI/UI/app/header"
	"Gio_UI/UI/app/importF"
	"Gio_UI/UI/app/menu"
	"Gio_UI/UI/app/table"
	"flag"
	"fmt"
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
	fmt.Println(UI.DBConnected)
	router := page.NewRouter()
	router.Register(0, header.New(&router))
	//router.Register(1, app2.New(&router))
	router.Register(2, importF.New(&router, th))
	//router.Register(3, dashboard.New(&router))
	//router.Register(2, textfield.New(&router))
	router.Register(3, menu.New(&router))
	router.Register(4, Tree.New(&router))
	router.Register(5, chart.New(&router))
	router.Register(6, Table.New(&router))
	router.Register(7, Tabs.New(&router))
	for {
		e := w.Event()
		switch e := e.(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			//if UI.LoginScreen {
			//	//UI.LayoutLogin(gtx, th)
			//	UI.DatabseLayout(gtx, th)
			//	//router.Layout(gtx, th)
			//} else {
			//	UI.DatabseLayout(gtx, th)
			//}
			if UI.DBConnected == true && UI.ShowLoginScreen == true {
				UI.LayoutLogin(gtx, th)
			} else {
				UI.DatabseLayout(gtx, th)
			}
			e.Frame(gtx.Ops)
		}
	}
}
