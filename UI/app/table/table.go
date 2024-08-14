package Table

import (
	page "Gio_UI/UI/app"
	"Gio_UI/UI/icon"
	"fmt"
	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	_ "github.com/mattn/go-sqlite3"
	"image"
	"image/color"
	"log"
	"unicode"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

var _ page.Page = &Page{}

func (p *Page) Actions() []component.AppBarAction {
	return []component.AppBarAction{}
}

func (p *Page) Overflow() []component.OverflowAction {
	return []component.OverflowAction{}
}

func (p *Page) NavItem() component.NavItem {
	return component.NavItem{
		Name: "Table",
		Icon: icon.EditIcon,
	}
}

func (p *Page) Layout(gtx C, th *material.Theme) D {
	p.List.Axis = layout.Vertical
	if p.SearchEditor.Text() != "" && p.SearchBtn.Clicked(gtx) {
		p.SearchID = p.SearchEditor.Text()
		p.Data = p.searchDataByID("UI/access.sqlite", p.SearchID)
	} else if p.ReadBtn.Clicked(gtx) {
		p.Data = p.readDataFromDB("UI/access.sqlite")
	}

	start := p.CurrentPage * p.RowsPerPage
	end := start + p.RowsPerPage
	if end > len(p.Data) {
		end = len(p.Data)
	}

	for i := start; i < end; i++ {
		if p.EditBtns[i-start].Clicked(gtx) {
			log.Printf("Edit button clicked for ID: %d", p.Data[i].ID)
			p.SelectedRow = &p.Data[i] // Ghi nhận dòng cần chỉnh sửa
			p.ShowEditConfirmation = true
		}
		if p.AddBtns[i-start].Clicked(gtx) {
			log.Printf("Add button clicked for ID: %d", p.Data[i].ID)
		}
		if p.DelBtns[i-start].Clicked(gtx) {
			log.Printf("Delete button clicked for ID: %d", p.Data[i].ID)
			p.SelectedRow = &p.Data[i] // Ghi nhận dòng cần xóa
			p.ShowDeleteConfirmation = true
		}
	}

	listLayout := material.List(th, &p.List)

	content := listLayout.Layout(gtx, end-start+2, func(gtx C, index int) D {
		return layout.UniformInset(unit.Dp(4)).Layout(gtx, func(gtx C) D {
			switch index {
			case 0:
				return p.LayoutReadButton(gtx, th)
			case 1:
				return p.LayoutColumnHeaders(gtx, th)
			default:
				return p.LayoutDataRow(gtx, th, p.Data[start+index-2], index-2)
			}
		})
	})
	if p.ShowDeleteConfirmation {
		return p.LayoutDeleteConfirmation(gtx, th)
	}
	if p.ShowEditConfirmation {
		return p.LayoutEditConfirmation(gtx, th)
	}
	var pagination layout.Dimensions
	pagination = layout.SE.Layout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				totalLabel := fmt.Sprintf("Total: %d", len(p.Data))
				label := material.Body1(th, totalLabel)
				label.Font.Style = font.Italic
				label.TextSize = 14
				return layout.UniformInset(unit.Dp(8)).Layout(gtx, label.Layout)
			}),
			layout.Rigid(func(gtx C) D {
				if p.PreviousBtn.Clicked(gtx) {
					if p.CurrentPage > 0 {
						p.CurrentPage--
					}
				}
				iconButton := material.IconButton(th, &p.PreviousBtn, icon.PreviousIcon, "Previous")
				iconButton.Color = color.NRGBA{A: 255}
				iconButton.Size = unit.Dp(20)
				iconButton.Background = color.NRGBA{A: 0}
				return layout.UniformInset(unit.Dp(4)).Layout(gtx, iconButton.Layout)
			}),
			layout.Rigid(func(gtx C) D {
				currentPageLabel := fmt.Sprintf("Page %d", p.CurrentPage+1)
				label := material.Body1(th, currentPageLabel)
				label.Font.Style = font.Italic
				label.TextSize = 14
				return layout.UniformInset(unit.Dp(8)).Layout(gtx, label.Layout)
			}),
			layout.Rigid(func(gtx C) D {
				if p.NextBtn.Clicked(gtx) {
					if (p.CurrentPage+1)*p.RowsPerPage < len(p.Data) {
						p.CurrentPage++
					}
				}
				iconButton := material.IconButton(th, &p.NextBtn, icon.NextIcon, "Next")
				iconButton.Color = color.NRGBA{A: 255}
				iconButton.Size = unit.Dp(20)
				iconButton.Background = color.NRGBA{A: 0}
				return layout.UniformInset(unit.Dp(4)).Layout(gtx, iconButton.Layout)
			}),
		)
	})
	return layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx C) D {
			return content
		}),
		layout.Stacked(func(gtx C) D {
			return pagination
		}),
	)
}
func (p *Page) LayoutDeleteConfirmation(gtx C, th *material.Theme) D {
	prompt := "Are you sure you want to delete this record?"
	backgroundColor := color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}
	dims := layout.Center.Layout(gtx, func(gtx C) D {
		paint.ColorOp{Color: backgroundColor}.Add(gtx.Ops)
		defer clip.Rect{
			Max: image.Point{
				X: 400,
				Y: 120,
			},
		}.Push(gtx.Ops).Pop()
		paint.PaintOp{}.Add(gtx.Ops)
		return widget.Border{
			Color:        color.NRGBA{R: 0x60, G: 0x60, B: 0x60, A: 0xFF}, // Border color
			CornerRadius: unit.Dp(8),
			Width:        unit.Dp(2),
		}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{
				Top:    20,
				Bottom: 20,
				Left:   20,
				Right:  20,
			}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return material.Body1(th, prompt).Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{Top: unit.Dp(20)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									btn := material.Button(th, &p.ConfirmBtn, "Yes")
									if p.ConfirmBtn.Clicked(gtx) {
										if p.SelectedRow != nil {
											err := p.deleteDataByID("UI/access.sqlite", p.SelectedRow.ID)
											if err != nil {
												log.Printf("Error deleting data: %v", err)
											} else {
												p.Data = p.readDataFromDB("UI/access.sqlite")
											}
										}
										p.ShowDeleteConfirmation = false
										p.SelectedRow = nil
									}
									return layout.Inset{Right: unit.Dp(10)}.Layout(gtx, btn.Layout)
								}),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									btn := material.Button(th, &p.CancelBtn, "No")
									btn.Background = color.NRGBA{R: 255, G: 200, B: 100, A: 255}
									if p.CancelBtn.Clicked(gtx) {
										p.ShowDeleteConfirmation = false
										p.SelectedRow = nil
									}
									return btn.Layout(gtx)
								}),
							)
						})
					}),
				)
			})
		})
	})
	return dims
}

func (p *Page) LayoutColumnHeaders(gtx C, th *material.Theme) D {
	// Define header dimensions
	var headerDims layout.Dimensions
	headerDims = layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			label := material.Body1(th, "ID")
			label.Font.Weight = font.Bold // Bold text
			return layout.UniformInset(unit.Dp(4)).Layout(gtx, label.Layout)
		}),
		layout.Flexed(1, func(gtx C) D {
			label := material.Body1(th, "Filename")
			label.Font.Weight = font.Bold // Bold text
			return layout.UniformInset(unit.Dp(4)).Layout(gtx, label.Layout)
		}),
		layout.Flexed(1, func(gtx C) D {
			label := material.Body1(th, "Modified")
			label.Font.Weight = font.Bold // Bold text
			return layout.UniformInset(unit.Dp(4)).Layout(gtx, label.Layout)
		}),
		layout.Flexed(1, func(gtx C) D {
			label := material.Body1(th, "Folder")
			label.Font.Weight = font.Bold // Bold text
			return layout.UniformInset(unit.Dp(4)).Layout(gtx, label.Layout)
		}),
	)
	return layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx C) D {
			return headerDims
		}),
		layout.Stacked(func(gtx C) D {
			lineColor := color.NRGBA{R: 0, G: 0, B: 0, A: 255} // Black color
			lineHeight := unit.Dp(1)

			paint.ColorOp{Color: lineColor}.Add(gtx.Ops)
			defer clip.Rect{
				Max: image.Point{
					X: headerDims.Size.X,
					Y: int(lineHeight),
				},
			}.Push(gtx.Ops).Pop()
			paint.PaintOp{}.Add(gtx.Ops)

			return layout.Dimensions{
				Size: image.Point{
					X: headerDims.Size.X,
					Y: int(lineHeight),
				},
			}
		}),
	)
}

func (p *Page) LayoutReadButton(gtx C, th *material.Theme) D {
	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		// Sử dụng layout.Inset cho ô nhập
		layout.Flexed(1, func(gtx C) D {
			return layout.E.Layout(gtx, func(gtx C) D {
				return layout.Inset{Left: unit.Dp(16), Right: unit.Dp(4), Top: unit.Dp(10)}.Layout(gtx, func(gtx C) D {
					gtx.Constraints.Min.X = gtx.Dp(unit.Dp(240))
					editor := material.Editor(th, &p.SearchEditor, "Search...")
					editor.Font.Style = font.Italic
					editor.TextSize = 13
					border := widget.Border{Color: color.NRGBA{A: 0xff}, CornerRadius: unit.Dp(6), Width: unit.Dp(0.5)}
					return border.Layout(gtx, func(gtx C) D {
						return layout.UniformInset(unit.Dp(8)).Layout(gtx, editor.Layout)
					})
				})
			})
		}),
		layout.Rigid(func(gtx C) D {
			iconButton := material.IconButton(th, &p.SearchBtn, icon.SearchIcon, "Search")
			iconButton.Color = color.NRGBA{A: 255}
			iconButton.Size = unit.Dp(20)
			iconButton.Background = color.NRGBA{A: 0}
			return layout.UniformInset(unit.Dp(4)).Layout(gtx, iconButton.Layout)
		}),
		layout.Rigid(func(gtx C) D {
			return layout.UniformInset(unit.Dp(4)).Layout(gtx,
				material.Button(th, &p.ReadBtn, "Read").Layout,
			)
		}),
	)
}

func (p *Page) LayoutDataRow(gtx C, th *material.Theme, row DataRow, index int) D {
	dataRowDims := layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return layout.Inset{Right: unit.Dp(18), Top: unit.Dp(18)}.Layout(gtx,
				material.Body1(th, fmt.Sprintf("%d", row.ID)).Layout,
			)
		}),
		layout.Flexed(1, func(gtx C) D {
			e := material.Body1(th, row.Filename)
			e.TextSize = 12
			return layout.Inset{Top: unit.Dp(18)}.Layout(gtx, e.Layout)
		}),
		layout.Flexed(1, func(gtx C) D {
			e := material.Body1(th, row.Modified.Format("2006-01-02 15:04:05"))
			e.TextSize = 12
			return layout.Inset{Top: unit.Dp(18)}.Layout(gtx, e.Layout)
		}),
		layout.Flexed(1, func(gtx C) D {
			e := material.Body1(th, row.Folder)
			e.TextSize = 12
			return layout.Inset{Top: unit.Dp(18)}.Layout(gtx, e.Layout)
		}),
		layout.Rigid(func(gtx C) D {
			iconButton := material.IconButton(th, &p.EditBtns[index], icon.EditIcon, "Edit")
			iconButton.Color = color.NRGBA{A: 56}
			iconButton.Size = unit.Dp(18)
			iconButton.Background = color.NRGBA{A: 0}
			return layout.UniformInset(unit.Dp(4)).Layout(gtx, iconButton.Layout)
		}),
		layout.Rigid(func(gtx C) D {
			iconButton := material.IconButton(th, &p.AddBtns[index], icon.PlusIcon, "Add")
			iconButton.Color = color.NRGBA{A: 255}
			iconButton.Size = unit.Dp(18)
			iconButton.Background = color.NRGBA{A: 0}
			return layout.UniformInset(unit.Dp(4)).Layout(gtx, iconButton.Layout)
		}),
		layout.Rigid(func(gtx C) D {
			iconButton := material.IconButton(th, &p.DelBtns[index], icon.DeleteIcon, "Delete")
			iconButton.Color = color.NRGBA{A: 255}
			iconButton.Size = unit.Dp(18)
			iconButton.Background = color.NRGBA{A: 0}
			return layout.UniformInset(unit.Dp(4)).Layout(gtx, iconButton.Layout)
		}),
	)

	p.DrawHorizontalLine(gtx)

	return dataRowDims
}

func (p *Page) DrawHorizontalLine(gtx C) D {
	lineColor := color.NRGBA{R: 200, G: 200, B: 200, A: 255} // Màu đen
	lineHeight := unit.Dp(1)

	paint.ColorOp{Color: lineColor}.Add(gtx.Ops)
	defer clip.Rect{
		Max: image.Point{
			X: gtx.Constraints.Max.X,
			Y: int(lineHeight),
		},
	}.Push(gtx.Ops).Pop()
	paint.PaintOp{}.Add(gtx.Ops)

	return layout.Dimensions{
		Size: image.Point{
			X: gtx.Constraints.Max.X,
			Y: int(lineHeight),
		},
	}
}

func (p *Page) LayoutEditConfirmation(gtx C, th *material.Theme) D {
	backgroundColor := color.NRGBA{R: 255, G: 255, B: 255, A: 255}

	paint.ColorOp{Color: backgroundColor}.Add(gtx.Ops)
	defer clip.Rect{
		Max: image.Point{X: gtx.Constraints.Max.X, Y: gtx.Constraints.Max.Y},
	}.Push(gtx.Ops).Pop()
	paint.PaintOp{}.Add(gtx.Ops)

	return layout.Inset{Right: unit.Dp(80), Top: unit.Dp(120), Bottom: unit.Dp(180)}.Layout(gtx, func(gtx C) D {
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
					e := material.Label(th, 20, "Edit Example")
					e.TextSize = 20
					e.Font.Typeface = "Go Mono"
					return inset.Layout(gtx, func(gtx C) D {
						return e.Layout(gtx)
					})
				}),
				layout.Rigid(func(gtx C) D {
					if p.FilenameInput.Text() == "" {
						p.FilenameInput.SetText(p.SelectedRow.Filename)
					}
					return layout.Inset{Left: unit.Dp(40), Right: unit.Dp(40), Top: unit.Dp(20)}.Layout(gtx, func(gtx C) D {
						p.FilenameInput.Alignment = p.inputAlignment
						return p.FilenameInput.Layout(gtx, th, "FileName")
					})
				}),

				layout.Rigid(func(gtx C) D {
					if p.FolderInput.Text() == "" {
						p.FolderInput.SetText(p.SelectedRow.Folder)
					}
					return layout.Inset{Left: unit.Dp(40), Right: unit.Dp(40)}.Layout(gtx, func(gtx C) D {
						p.FolderInput.Alignment = p.inputAlignment
						return p.FolderInput.Layout(gtx, th, "Folder")
					})
				}),
				layout.Rigid(func(gtx C) D {
					inset := layout.Inset{Left: unit.Dp(40), Right: unit.Dp(40)}
					p.priceInput.Prefix = func(gtx C) D {
						th := *th
						th.Palette.Fg = color.NRGBA{R: 100, G: 100, B: 100, A: 255}
						return material.Label(&th, th.TextSize, "$").Layout(gtx)
					}
					p.priceInput.Suffix = func(gtx C) D {
						th := *th
						th.Palette.Fg = color.NRGBA{R: 100, G: 100, B: 100, A: 255}
						return material.Label(&th, th.TextSize, ".00").Layout(gtx)
					}
					p.priceInput.SingleLine = true
					p.priceInput.Alignment = p.inputAlignment
					return inset.Layout(gtx, func(gtx C) D {
						return p.priceInput.Layout(gtx, th, "Price")
					})
				}),
				layout.Rigid(func(gtx C) D {
					inset := layout.Inset{Left: unit.Dp(40), Right: unit.Dp(40)}
					if err := func() string {
						for _, r := range p.IDInput.Text() {
							if !unicode.IsDigit(r) {
								return "Must contain only digits"
							}
						}
						return ""
					}(); err != "" {
						p.IDInput.SetError(err)
					} else {
						p.IDInput.ClearError()
					}
					p.IDInput.SingleLine = true
					p.IDInput.Alignment = p.inputAlignment
					return inset.Layout(gtx, func(gtx C) D {
						return p.IDInput.Layout(gtx, th, "ID")
					})
				}),
				layout.Rigid(func(gtx C) D {
					inset := layout.Inset{Left: unit.Dp(40), Right: unit.Dp(40)}
					if p.tweetInput.TextTooLong() {
						p.tweetInput.SetError("Too many characters")
					} else {
						p.tweetInput.ClearError()
					}
					p.tweetInput.CharLimit = 128
					p.tweetInput.Helper = "Character count"
					p.tweetInput.Alignment = p.inputAlignment
					return inset.Layout(gtx, func(gtx C) D {
						return p.tweetInput.Layout(gtx, th, "Note")
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
								btn := material.Button(th, &p.ConfirmBtn, "Update")
								if p.ConfirmBtn.Clicked(gtx) {
									// Xử lý xác nhận chỉnh sửa
								}
								return btn.Layout(gtx)
							})
						}),
						layout.Rigid(func(gtx C) D {
							inset := layout.Inset{Right: unit.Dp(40)} // Cách lề phải 20 dp
							return inset.Layout(gtx, func(gtx C) D {
								btn := material.Button(th, &p.CancelBtn, "Cancel")
								btn.Background = color.NRGBA{R: 255, G: 200, B: 100, A: 255}
								if p.CancelBtn.Clicked(gtx) {
									p.ShowEditConfirmation = false
									p.SelectedRow = nil
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
