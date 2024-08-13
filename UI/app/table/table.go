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
)

type (
	C = layout.Context
	D = layout.Dimensions
)

type Page struct {
	widget.List
	*page.Router
	ReadBtn                widget.Clickable
	SearchBtn              widget.Clickable
	Data                   []DataRow
	SearchEditor           widget.Editor
	CurrentPage            int
	RowsPerPage            int
	PreviousBtn            widget.Clickable
	NextBtn                widget.Clickable
	SearchID               string
	ShowDeleteConfirmation bool
	SelectedRowToDelete    *DataRow
	ConfirmDeleteBtn       widget.Clickable
	CancelDeleteBtn        widget.Clickable
	EditBtns               []widget.Clickable
	AddBtns                []widget.Clickable
	DelBtns                []widget.Clickable
	DeleteID               int // ID của bản ghi cần xóa
	ShowConfirm            bool
}

// New constructs a Page with the provided router.

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

	// Xử lý tìm kiếm và đọc dữ liệu
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
		}
		if p.AddBtns[i-start].Clicked(gtx) {
			log.Printf("Add button clicked for ID: %d", p.Data[i].ID)
		}
		if p.DelBtns[i-start].Clicked(gtx) {
			log.Printf("Delete button clicked for ID: %d", p.Data[i].ID)
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
