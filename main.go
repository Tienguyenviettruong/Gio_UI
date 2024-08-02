package main

import (
	"bytes"
	"flag"
	"fmt"
	"gioui.org/font"
	"gioui.org/op/paint"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"time"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/gpu/headless"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

var (
	screenshot = flag.String("screenshot", "", "save a screenshot to a file and exit")
	disable    = flag.Bool("disable", false, "disable all widgets")
)

type iconAndTextButton struct {
	theme  *material.Theme
	button *widget.Clickable
	icon   *widget.Icon
	word   string
}

var (
	editor              = widget.Editor{SingleLine: true}
	icon                *widget.Icon
	progressIncrementer chan float32
	lineEditor          = &widget.Editor{
		SingleLine: true,
		Submit:     true,
	}
	topLabel = "Hello, Gio"
	list     = &widget.List{
		List: layout.List{
			Axis: layout.Vertical,
		},
	}
)

func main() {
	flag.Parse()
	editor.SetText("Long text here")
	ic, err := widget.NewIcon(icons.ContentAdd)
	if err != nil {
		log.Fatal(err)
	}
	icon = ic
	progressIncrementer = make(chan float32)
	if *screenshot != "" {
		if err := saveScreenshot(*screenshot); err != nil {
			fmt.Fprintf(os.Stderr, "failed to save screenshot: %v\n", err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	go func() {
		for {
			time.Sleep(time.Second)
			progressIncrementer <- 0.1
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

func saveScreenshot(f string) error {
	const scale = 1.5
	sz := image.Point{X: 800 * scale, Y: 600 * scale}
	w, err := headless.NewWindow(sz.X, sz.Y)
	if err != nil {
		return err
	}
	gtx := layout.Context{
		Ops: new(op.Ops),
		Metric: unit.Metric{
			PxPerDp: scale,
			PxPerSp: scale,
		},
		Constraints: layout.Exact(sz),
	}
	th := material.NewTheme()
	th.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
	kitchen(gtx, th)
	w.Frame(gtx.Ops)
	return nil
}

func loop(w *app.Window) error {
	th := material.NewTheme()
	var ops op.Ops
	for {
		e := w.Event()
		switch e := e.(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			if loginScreen {
				layoutLogin(gtx, th)
			} else {
				kitchen(gtx, th)
			}
			e.Frame(gtx.Ops)
		}
	}
}

type (
	D = layout.Dimensions
	C = layout.Context
)

var (
	usernameEditor    widget.Editor
	passwordEditor    widget.Editor
	loginButton       widget.Clickable
	loginScreen       bool = true
	button1           widget.Clickable
	button2           widget.Clickable
	button3           widget.Clickable
	logo              paint.ImageOp
	showErrorDialog   bool = false
	errorDialogButton widget.Clickable
)

func layoutLogin(gtx layout.Context, th *material.Theme) layout.Dimensions {
	drawImageBackground(gtx)
	file, err := os.Open("asset/decore.png")
	if err != nil {
		// handle error
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		// handle error
	}
	// Convert the image to an op.ImageOp
	opImage := paint.NewImageOp(img)

	buttonStyle := material.ButtonStyle{
		Background: color.NRGBA{R: 0x3b, G: 0x3f, B: 0x4c, A: 0xff}, // Dark blue background
		//TextColor:  color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}, // White text color
	}

	widgets := []layout.Widget{
		func(gtx C) D {

			return layout.Flex{}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					l := material.H3(th, topLabel)
					//l.State = topLabelState
					return l.Layout(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return widget.Image{
						Src: opImage,
					}.Layout(gtx)
				}),
			)
		}, func(gtx C) D {
			return layout.Spacer{Height: unit.Dp(16)}.Layout(gtx) // Add spacer to create space between the inputs
		},
		func(gtx C) D {
			return layout.Center.Layout(gtx, func(gtx C) D {
				return layout.Inset{Left: unit.Dp(16), Right: unit.Dp(16)}.Layout(gtx, func(gtx C) D {
					gtx.Constraints.Min.X = gtx.Dp(unit.Dp(240)) // Set minimum width for the username editor
					e := material.Editor(th, &usernameEditor, "Username")
					e.Font.Style = font.Italic
					border := widget.Border{Color: color.NRGBA{A: 0xff}, CornerRadius: unit.Dp(6), Width: unit.Dp(0.5)}
					return border.Layout(gtx, func(gtx C) D {
						return layout.UniformInset(unit.Dp(8)).Layout(gtx, e.Layout)
					})
				})
			})
		},
		//func(gtx C) D {
		//	return layout.Spacer{Height: unit.Dp(16)}.Layout(gtx) // Add spacer to create space between the inputs
		//},
		func(gtx C) D {
			return layout.Inset{Left: unit.Dp(16), Right: unit.Dp(16)}.Layout(gtx, func(gtx C) D {
				gtx.Constraints.Min.X = gtx.Dp(unit.Dp(240)) // Set maximum width for the password editor
				e := material.Editor(th, &passwordEditor, "Password")
				e.Font.Style = font.Italic
				border := widget.Border{Color: color.NRGBA{A: 0xff}, CornerRadius: unit.Dp(6), Width: unit.Dp(0.5)}
				return border.Layout(gtx, func(gtx C) D {
					return layout.UniformInset(unit.Dp(8)).Layout(gtx, e.Layout)
				})
			})
		},
		//func(gtx C) D {
		//	return layout.Spacer{Height: unit.Dp(16)}.Layout(gtx) // Add spacer to create space between the inputs
		//},
		func(gtx C) D {
			return layout.Inset{Left: unit.Dp(108)}.Layout(gtx, func(gtx C) D { // Center the login button
				btn := material.Button(th, &loginButton, "Sign In")
				btn.Background = buttonStyle.Background
				//btn.TextColor = buttonStyle.TextColor
				if loginButton.Clicked(gtx) {
					if usernameEditor.Text() == "1" && passwordEditor.Text() == "1" {
						loginScreen = false
					} else {
						showErrorDialog = true
					}
				}
				return btn.Layout(gtx)
			})
		},
	}
	if showErrorDialog {
		dims := layout.Center.Layout(gtx, func(gtx C) D {
			return widget.Border{
				Color:        color.NRGBA{R: 0x60, G: 0x60, B: 0x60, A: 0xFF},
				CornerRadius: unit.Dp(4),
				Width:        unit.Dp(2),
			}.Layout(gtx, func(gtx C) D {
				return layout.UniformInset(unit.Dp(20)).Layout(gtx, func(gtx C) D {
					return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							return material.Body1(th, "Sai tên đăng nhập hoặc mật khẩu.").Layout(gtx)
						}),
						layout.Rigid(func(gtx C) D {
							return layout.Inset{Top: unit.Dp(20)}.Layout(gtx, func(gtx C) D {
								btn := material.Button(th, &errorDialogButton, "OK")
								if errorDialogButton.Clicked(gtx) {
									showErrorDialog = false
								}
								return btn.Layout(gtx)
							})
						}),
					)
				})
			})
		})
		if showErrorDialog {
			return dims
		}
	}
	return layout.Center.Layout(gtx, func(gtx C) D {
		return material.List(th, list).Layout(gtx, len(widgets), func(gtx C, i int) D {
			return layout.UniformInset(unit.Dp(16)).Layout(gtx, widgets[i])
		})
	})
}

func kitchen(gtx layout.Context, th *material.Theme) {
	// Màn hình chính hiện tại, nội dung ở đây sẽ được hiển thị sau khi đăng nhập thành công
	layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return material.H6(th, "Main Screen").Layout(gtx)
	})
}

func drawImageBackground(gtx layout.Context) {
	data, err := ioutil.ReadFile("background.png")
	if err != nil {
		log.Fatal(err)
	}

	// Tạo một đối tượng image từ dữ liệu ảnh
	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		log.Fatal(err)
	}

	// Tạo một hình ảnh mới với kích thước của khu vực vẽ
	dst := image.NewRGBA(image.Rect(0, 0, gtx.Constraints.Max.X, gtx.Constraints.Max.Y))

	// Vẽ ảnh nền lên vùng đích
	draw.Draw(dst, dst.Bounds(), img, image.Point{}, draw.Over)

	// Vẽ hình ảnh lên layout
	paint.NewImageOp(dst).Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
}
