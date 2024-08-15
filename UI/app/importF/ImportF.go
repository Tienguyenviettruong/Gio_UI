package importF

import (
	page "Gio_UI/UI/app"
	"Gio_UI/UI/icon"
	"fmt"
	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"github.com/sqweek/dialog"
	"image"
	"image/color"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

type (
	C = layout.Context
	D = layout.Dimensions
)
type Page struct {
	DrawButton           widget.Clickable
	rowsInput, colsInput widget.Editor
	*material.Theme
	*page.Router
	progressBar                                widget.Float
	uploadComplete                             bool
	selectedFile                               string
	confirmButton, cancelButton, fileImportBtn widget.Clickable
	showUploadDialog                           bool
	Note                                       component.TextField
	inputAlignment                             text.Alignment
	chooseFileButton                           widget.Clickable
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
		Icon: icon.ImportIcon,
	}
}

func (p *Page) Actions() []component.AppBarAction {
	return []component.AppBarAction{}
}

func (p *Page) Overflow() []component.OverflowAction {
	return []component.OverflowAction{}
}

func (p *Page) Layout(gtx C, th *material.Theme) D {
	if p.DrawButton.Clicked(gtx) {
		go p.openNewWindow()
	}

	if p.fileImportBtn.Clicked(gtx) {
		p.showUploadDialog = true
		//go p.chooseFile(gtx)
	}
	if p.showUploadDialog {
		return p.LayoutUpload(gtx, th)
	}
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return layout.Flex{}.Layout(gtx,
				layout.Flexed(1, func(gtx C) D {
					return material.Editor(th, &p.rowsInput, "Rows").Layout(gtx)
				}),
				layout.Flexed(1, func(gtx C) D {
					return material.Editor(th, &p.colsInput, "Columns").Layout(gtx)
				}),
				layout.Rigid(func(gtx C) D {
					btn := material.Button(th, &p.DrawButton, "Edit")
					btn.TextSize = unit.Sp(12)
					return btn.Layout(gtx)
				}),
			)
		}),

		layout.Rigid(func(gtx C) D {
			return layout.Inset{
				Left:   unit.Dp(8),
				Right:  unit.Dp(8),
				Top:    unit.Dp(8),
				Bottom: unit.Dp(8),
			}.Layout(gtx, func(gtx C) D {
				// Tạo border cho button
				border := widget.Border{
					Color:        color.NRGBA{R: 0x60, G: 0x60, B: 0x60, A: 0xFF},
					CornerRadius: unit.Dp(8),
					Width:        unit.Dp(1),
				}

				return layout.Stack{}.Layout(gtx,
					layout.Stacked(func(gtx C) D {
						// Vẽ border xung quanh
						return border.Layout(gtx, func(gtx C) D {
							item := component.MenuItem(th, &p.fileImportBtn, "Import")
							item.Icon = icon.ImportIcon
							item.Label.TextSize = 12
							item.IconSize = 14
							return item.Layout(gtx)
						})
					}),
				)
			})
		}),
		layout.Rigid(func(gtx C) D {
			if p.uploadComplete {
				return material.Body1(th, "Upload Complete").Layout(gtx)
			}
			return layout.Inset{Top: unit.Dp(8)}.Layout(gtx, func(gtx C) D {
				return material.ProgressBar(th, p.progressBar.Value).Layout(gtx)
			})
		}),
	)
}

func (p *Page) openNewWindow() {
	go func() {
		rows, cols := 0, 0
		if _, err := fmt.Sscanf(p.rowsInput.Text(), "%d", &rows); err != nil {
			log.Println("Invalid number of rows")
			return
		}
		if _, err := fmt.Sscanf(p.colsInput.Text(), "%d", &cols); err != nil {
			log.Println("Invalid number of columns")
			return
		}

		newWindow := new(app.Window)
		newWindow.Option(app.Size(unit.Dp(800), unit.Dp(600)))
		if err := p.run(newWindow, rows, cols); err != nil {
			log.Fatal(err)
		}
	}()
}

func (p *Page) run(w *app.Window, rows, cols int) error {
	var ops op.Ops
	th := material.NewTheme()
	for {
		e := w.Event()
		switch e := e.(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			//layout.Flex{}.Layout(gtx,
			//	layout.Flexed(1, func(gtx C) D {
			//		return layout.Inset{Left: unit.Dp(108)}.Layout(gtx, func(gtx C) D {
			//			return material.Body1(p.Theme, "This is the new window").Layout(gtx)
			//		})
			//	}),
			//)

			p.drawTable(gtx, th, rows, cols)
			e.Frame(gtx.Ops)
		}
	}
	return nil
}
func (p *Page) drawTable(gtx layout.Context, th *material.Theme, rows, cols int) {
	lineColor := color.NRGBA{R: 0, G: 0, B: 255, A: 255} // Màu xanh dương

	rowHeight := float32(60)

	screenWidth := float32(gtx.Constraints.Max.X)
	cellWidth := screenWidth / float32(cols)

	screenHeight := rowHeight * float32(rows)

	strokeWidth := float32(1.0)

	for i := 0; i <= rows; i++ {
		y := float32(i) * rowHeight
		drawLine(gtx, 0, y, screenWidth, y, lineColor, strokeWidth)
	}

	for j := 0; j <= cols; j++ {
		x := float32(j) * cellWidth
		drawLine(gtx, x, 0, x, screenHeight, lineColor, strokeWidth)
	}
	drawText(gtx, fmt.Sprintf("Cell %d,%d", 1, 1), 0, 0, th)
	//drawText(gtx, fmt.Sprintf("Cell %d,%d", 1, 2), 400, 0, th)
	drawText(gtx, fmt.Sprintf("Cell %d,%d", 2, 1), 0, 60, th)
	drawText(gtx, fmt.Sprintf("Cell %d,%d", 2, 2), 400, 0, th)
	drawText(gtx, fmt.Sprintf("Cell %d,%d", 3, 1), 400, 60, th)
	//for i := 0; i < rows; i++ {
	//	for j := 0; j < cols; j++ {
	//		x := float32(j) * cellWidth
	//		y := float32(i) * rowHeight
	//		drawText(gtx, fmt.Sprintf("Cell %d,%d", i+1, j+1), x, y, th)
	//		fmt.Println("Cell %d,%d", i+1, j+1, x, y)
	//	}
	//}
}

func drawLine(gtx layout.Context, x1, y1, x2, y2 float32, clr color.NRGBA, width float32) {
	var path clip.Path
	path.Begin(gtx.Ops)
	path.Move(f32.Pt(x1, y1))
	path.Line(f32.Pt(x2-x1, y2-y1))
	paint.FillShape(gtx.Ops, clr, clip.Stroke{Path: path.End(), Width: width}.Op())
}
func drawText(gtx layout.Context, txt string, x, y float32, th *material.Theme) {
	label := material.Label(th, unit.Sp(14), txt)

	op.Offset(image.Pt(int(x), int(y))).Add(gtx.Ops)
	label.Layout(gtx)
}
func (p *Page) chooseFile(gtx layout.Context) {
	file, err := dialog.File().Title("Select a file").Load()
	if err != nil {
		log.Println("Error selecting file:", err)
		return
	}
	p.selectedFile = file
	go p.uploadFile(gtx, file)
}

func (p *Page) uploadFile(gtx layout.Context, filePath string) {
	p.progressBar.Value = 0
	p.uploadComplete = false
	for i := 0; i <= 100; i++ {
		p.progressBar.Value = float32(i) / 100
		time.Sleep(1 * time.Millisecond)
	}

	p.uploadComplete = true

	dst := filepath.Join("UI/Uploads", filepath.Base(filePath))
	fmt.Println(dst)
	input, err := os.Open(filePath)
	if err != nil {
		log.Println("Error opening file:", err)
		return
	}
	defer input.Close()

	output, err := os.Create(dst)
	if err != nil {
		log.Println("Error creating file:", err)
		return
	}
	defer output.Close()

	if _, err := io.Copy(output, input); err != nil {
		log.Println("Error copying file:", err)
	}
}
func (p *Page) LayoutUpload(gtx C, th *material.Theme) D {
	backgroundColor := color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	//var chooseFileButton widget.Clickable
	paint.ColorOp{Color: backgroundColor}.Add(gtx.Ops)
	defer clip.Rect{
		Max: image.Point{X: gtx.Constraints.Max.X, Y: gtx.Constraints.Max.Y},
	}.Push(gtx.Ops).Pop()
	paint.PaintOp{}.Add(gtx.Ops)

	return layout.Inset{Left: unit.Dp(80), Right: unit.Dp(80), Top: unit.Dp(120), Bottom: unit.Dp(180)}.Layout(gtx, func(gtx C) D {
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
					e := material.Label(th, 20, "Upload Example")
					e.TextSize = 20
					e.Font.Typeface = "Go Mono"
					return inset.Layout(gtx, func(gtx C) D {
						return e.Layout(gtx)
					})
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Left: unit.Dp(40), Right: unit.Dp(40)}.Layout(gtx, func(gtx C) D {
						img := paint.NewImageOp(loadImage("asset/Uploadfile.png"))

						return layout.Stack{}.Layout(gtx,
							// Lớp hiển thị ảnh nền
							layout.Expanded(func(gtx C) D {
								return widget.Image{
									Src:   img,
									Scale: 1,
									Fit:   widget.Contain,
								}.Layout(gtx)
							}),
							layout.Expanded(func(gtx C) D {
								return layout.Center.Layout(gtx, func(gtx C) D {
									if p.chooseFileButton.Clicked(gtx) {
										p.chooseFile(gtx)
									}
									btn := material.Button(p.Theme, &p.chooseFileButton, "         ")
									btn.Background = color.NRGBA{R: 0, G: 0, B: 0, A: 0}
									btn.Inset = layout.UniformInset(unit.Dp(0))
									return btn.Layout(gtx)
								})
							}),
						)
					})
				}),

				//layout.Rigid(func(gtx C) D {
				//	return layout.Inset{Left: unit.Dp(40), Right: unit.Dp(40)}.Layout(gtx, func(gtx C) D {
				//		img := paint.NewImageOp(loadImage("asset/Uploadfile.png"))
				//		return widget.Image{
				//			Src:   img,
				//			Scale: 1,
				//			Fit:   widget.Contain,
				//		}.Layout(gtx)
				//	})
				//}),

				layout.Rigid(func(gtx C) D {
					return layout.Inset{Left: unit.Dp(40), Right: unit.Dp(40), Top: unit.Dp(20)}.Layout(gtx, func(gtx C) D {
						p.Note.Alignment = p.inputAlignment
						return p.Note.Layout(gtx, th, "note")
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
							inset := layout.Inset{Right: unit.Dp(20)}
							return inset.Layout(gtx, func(gtx C) D {
								btn := material.Button(th, &p.confirmButton, "Update")
								if p.confirmButton.Clicked(gtx) {
									// Xử lý xác nhận chỉnh sửa
								}
								return btn.Layout(gtx)
							})
						}),
						layout.Rigid(func(gtx C) D {
							inset := layout.Inset{Right: unit.Dp(40)} // Cách lề phải 20 dp
							return inset.Layout(gtx, func(gtx C) D {
								btn := material.Button(th, &p.cancelButton, "Cancel")
								btn.Background = color.NRGBA{R: 255, G: 200, B: 100, A: 255}
								if p.cancelButton.Clicked(gtx) {
									p.showUploadDialog = false
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
func loadImage(path string) image.Image {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		log.Fatal(err)
	}

	return img
}
