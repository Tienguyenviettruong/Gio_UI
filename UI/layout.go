package UI

import (
	"gioui.org/font"
	"image"
	"image/color"
	"os"

	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

func LayoutLogin(gtx layout.Context, th *material.Theme) layout.Dimensions {
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
		func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					l := material.H3(th, TopLabel)
					return l.Layout(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return widget.Image{
						Src: opImage,
					}.Layout(gtx)
				}),
			)
		},
		func(gtx layout.Context) layout.Dimensions {
			return layout.Spacer{Height: unit.Dp(16)}.Layout(gtx)
		},
		func(gtx layout.Context) layout.Dimensions {
			return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{Left: unit.Dp(16), Right: unit.Dp(16)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					gtx.Constraints.Min.X = gtx.Dp(unit.Dp(240))
					e := material.Editor(th, &UsernameEditor, "Username")
					e.Font.Style = font.Italic
					border := widget.Border{Color: color.NRGBA{A: 0xff}, CornerRadius: unit.Dp(6), Width: unit.Dp(0.5)}
					return border.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return layout.UniformInset(unit.Dp(8)).Layout(gtx, e.Layout)
					})
				})
			})
		},
		func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Left: unit.Dp(16), Right: unit.Dp(16)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Min.X = gtx.Dp(unit.Dp(240))
				e := material.Editor(th, &PasswordEditor, "Password")
				e.Font.Style = font.Italic
				border := widget.Border{Color: color.NRGBA{A: 0xff}, CornerRadius: unit.Dp(6), Width: unit.Dp(0.5)}
				return border.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.UniformInset(unit.Dp(8)).Layout(gtx, e.Layout)
				})
			})
		},
		func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Left: unit.Dp(108)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				btn := material.Button(th, &LoginButton, "Sign In")
				btn.Background = buttonStyle.Background
				if LoginButton.Clicked(gtx) {
					if UsernameEditor.Text() == "1" && PasswordEditor.Text() == "1" {
						LoginScreen = false
					} else {
						ShowErrorDialog = true
					}
				}
				return btn.Layout(gtx)
			})
		},
	}

	if ShowErrorDialog {
		dims := layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return widget.Border{
				Color:        color.NRGBA{R: 0x60, G: 0x60, B: 0x60, A: 0xFF},
				CornerRadius: unit.Dp(4),
				Width:        unit.Dp(2),
			}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.UniformInset(unit.Dp(20)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return material.Body1(th, "Sai tên đăng nhập hoặc mật khẩu.").Layout(gtx)
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return layout.Inset{Top: unit.Dp(20)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
								btn := material.Button(th, &ErrorDialogButton, "OK")
								if ErrorDialogButton.Clicked(gtx) {
									ShowErrorDialog = false
								}
								return btn.Layout(gtx)
							})
						}),
					)
				})
			})
		})
		return dims
	}

	return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return material.List(th, List).Layout(gtx, len(widgets), func(gtx layout.Context, i int) layout.Dimensions {
			return layout.UniformInset(unit.Dp(16)).Layout(gtx, widgets[i])
		})
	})
}

func Kitchen(gtx layout.Context, th *material.Theme) {
	layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return material.H6(th, "Main Screen").Layout(gtx)
	})
}
