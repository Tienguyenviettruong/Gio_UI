package UI

import (
	connect "Gio_UI/UI/app/database"
	"fmt"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	_ "gorm.io/gorm"
	"image/color"
	"os"
	"strconv"
)

var (
	radioOptions = []string{"Oracle", "MYSQL", "postgresSQL"}
	radioGroup   widget.Enum
)

//var DBConnected bool
//var ShowLoginScreen bool = false

func DatabseLayout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	drawImageBackground(gtx)
	file, err := os.Open("asset/decore.png")
	if err != nil {
		// handle error
	}
	defer file.Close()

	return layout.Center.Layout(gtx, func(gtx C) D {
		gtx.Constraints.Max.X = gtx.Dp(unit.Dp(640))
		return widget.Border{
			Color:        color.NRGBA{R: 0x60, G: 0x60, B: 0x60, A: 0xFF},
			CornerRadius: unit.Dp(8),
			Width:        unit.Dp(1),
		}.Layout(gtx, func(gtx C) D {
			return layout.Flex{
				Axis: layout.Vertical,
			}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					inset := layout.Inset{Left: unit.Dp(40), Right: unit.Dp(40), Top: unit.Dp(20)}
					return inset.Layout(gtx, func(gtx C) D {
						// Bố trí label "Connect database" và RadioGroup cạnh nhau
						return layout.Flex{
							Axis: layout.Horizontal,
						}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								e := material.Label(th, 20, "Connect database")
								e.TextSize = 20
								e.Font.Typeface = "Go Mono"
								return e.Layout(gtx)
							}),
							layout.Rigid(layout.Spacer{Width: unit.Dp(80)}.Layout), // Khoảng cách giữa Label và RadioGroup
							layout.Rigid(func(gtx C) D {
								// Tạo RadioGroup với các tùy chọn A, B, C
								var options []layout.FlexChild
								for _, option := range radioOptions {
									opt := option
									options = append(options, layout.Rigid(func(gtx C) D {
										return material.RadioButton(th, &radioGroup, opt, opt).Layout(gtx)
									}))
								}
								return layout.Flex{Axis: layout.Horizontal}.Layout(gtx, options...)
							}),
						)
					})
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Left: unit.Dp(40), Right: unit.Dp(40), Top: unit.Dp(20)}.Layout(gtx, func(gtx C) D {
						Host.Alignment = inputAlignment
						return Host.Layout(gtx, th, "Host")
					})
				}),

				layout.Rigid(func(gtx C) D {
					switch radioGroup.Value {
					case "Oracle":
						Port.SetText("1521")
					case "MYSQL":
						Port.SetText("3306")
					case "postgresSQL":
						Port.SetText("5432")
					}

					return layout.Inset{Left: unit.Dp(40), Right: unit.Dp(40)}.Layout(gtx, func(gtx C) D {
						Port.Alignment = inputAlignment
						return Port.Layout(gtx, th, "Port")
					})
				}),

				layout.Rigid(func(gtx C) D {
					return layout.Inset{Left: unit.Dp(40), Right: unit.Dp(40)}.Layout(gtx, func(gtx C) D {
						Servicename.Alignment = inputAlignment
						return Servicename.Layout(gtx, th, "Servicename")
					})
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Left: unit.Dp(40), Right: unit.Dp(40)}.Layout(gtx, func(gtx C) D {
						User.Alignment = inputAlignment
						return User.Layout(gtx, th, "User")
					})
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Left: unit.Dp(40), Right: unit.Dp(40)}.Layout(gtx, func(gtx C) D {
						return layout.Stack{}.Layout(gtx,
							layout.Stacked(func(gtx C) D {
								return Password.Layout(gtx, th, "Password")
							}),
							layout.Stacked(func(gtx C) D {
								if ShowPasswordBtn.Clicked(gtx) {
									if Password.Editor.Mask == 0 {
										Password.Editor.Mask = '*'
									} else {
										Password.Editor.Mask = 0
									}
								}
								var icon *widget.Icon
								if Password.Editor.Mask == 0 {
									icon = iconVisibilityOff
								} else {
									icon = iconVisibility
								}
								return layout.Inset{Left: unit.Dp(510), Top: unit.Dp(8)}.Layout(gtx, func(gtx C) D {
									e := material.IconButton(th, &ShowPasswordBtn, icon, "Toggle visibility")
									e.Color = th.Fg
									e.Background = color.NRGBA{A: 0}
									return e.Layout(gtx)
								})
							}),
						)
					})
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Flex{
						Axis:      layout.Horizontal,
						Alignment: layout.Middle,
					}.Layout(
						gtx,
						layout.Flexed(1, layout.Spacer{}.Layout),
						layout.Rigid(func(gtx C) D {
							inset := layout.Inset{Right: unit.Dp(40), Top: unit.Dp(10)}
							return inset.Layout(gtx, func(gtx C) D {
								btn := material.Button(th, &ConnectBtn, "Connect")
								if ConnectBtn.Clicked(gtx) {
									host := Host.Text()
									port, _ := strconv.Atoi(Port.Text())
									user := User.Text()
									password := Password.Text()
									service := Servicename.Text()
									a, err := connect.ConnectDBOracle(host, port, user, password, service)
									//a, err := ConnectDB()
									fmt.Println(a)
									if err != nil {
										fmt.Println("Connect database failed:", err)
									} else {
										DBConnected = a
										ShowLoginScreen = true
									}
								}
								return btn.Layout(gtx)
							})
						}),
					)
				}),
			)
		})
	})
}

//var db *gorm.DB
//
//func ConnectDB() (bool, error) {
//	host := Host.Text()
//	port, _ := strconv.Atoi(Port.Text())
//	user := User.Text()
//	password := Password.Text()
//	service := Servicename.Text()
//	dataSourceName := fmt.Sprintf("%s/%s@%s:%d/%s", user, password, host, port, service)
//	fmt.Println(dataSourceName)
//	dbConnection, err := gorm.Open(oracle.Open(dataSourceName), &gorm.Config{})
//	if err != nil {
//		return false, err
//	}
//
//	if err := dbConnection.Raw("select 1").Error; err != nil {
//		return false, err
//	}
//
//	db = dbConnection
//	return true, nil
//}
