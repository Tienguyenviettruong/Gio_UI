package main

import (
	"Gio_UI/UI"
	page "Gio_UI/UI/app"
	"Gio_UI/UI/app/Tabs"
	"Gio_UI/UI/app/Tree"
	"Gio_UI/UI/app/analysis"
	"Gio_UI/UI/app/chart"
	"Gio_UI/UI/app/header"
	"Gio_UI/UI/app/importF"
	"Gio_UI/UI/app/menu"
	"Gio_UI/UI/app/table"
	"flag"
	"fmt"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"gioui.org/x/outlay"
	"log"
	"os"
	"time"

	"gioui.org/app"
	"gioui.org/op"
	"gioui.org/unit"
)

type GridState struct {
	VScrollbar widget.Scrollbar
	HScrollbar widget.Scrollbar
	outlay.Grid
}

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
	var (
		ops  op.Ops
		grid component.GridState
	)
	timingWindow := time.Second
	timings := []UI.FrameTiming{}
	frameCounter := 0
	timingStart := time.Time{}
	fmt.Println(UI.DBConnected)
	router := page.NewRouter()
	router.Register(0, header.New(&router))
	//router.Register(1, app2.New(&router))
	router.Register(1, importF.New(&router, th))
	//router.Register(3, dashboard.New(&router))
	router.Register(2, menu.New(&router))
	router.Register(3, Tree.New(&router))
	router.Register(4, analysis.New(&router))
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
			gtx.Execute(op.InvalidateCmd{})
			if timingStart == (time.Time{}) {
				timingStart = gtx.Now
			}
			if interval := gtx.Now.Sub(timingStart); interval >= timingWindow {
				timings = append(timings, UI.FrameTiming{
					Start:           timingStart,
					End:             gtx.Now,
					FrameCount:      frameCounter,
					FramesPerSecond: float64(frameCounter) / interval.Seconds(),
				})
				frameCounter = 0
				timingStart = gtx.Now
			}
			if UI.LoginScreen {
				UI.LayoutTable(th, gtx, timings, &grid)
				//UI.LayoutDashboard(gtx, th)
				//UI.LayoutFlexGrid(gtx, th)
				//router.Layout(gtx, th)
			} else {
				UI.DatabseLayout(gtx, th)
			}
			//if UI.DBConnected == true && UI.ShowLoginScreen == true {
			//	UI.LayoutLogin(gtx, th)
			//} else {
			//	UI.DatabseLayout(gtx, th)
			//}
			e.Frame(gtx.Ops)
		}
	}
}
