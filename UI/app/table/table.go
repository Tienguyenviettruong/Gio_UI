package Table

import (
	"Gio_UI/UI/icon"
	"database/sql"
	"fmt"
	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"image"
	"image/color"
	"log"
	"time"

	page "Gio_UI/UI/app"
	_ "github.com/mattn/go-sqlite3"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

type Page struct {
	widget.List
	*page.Router
	ReadBtn      widget.Clickable
	Data         []DataRow
	SearchEditor widget.Editor
}

type DataRow struct {
	ID       int
	Filename string
	Modified time.Time
	Folder   string
	EditBtn  widget.Clickable
}

// New constructs a Page with the provided router.
func New(router *page.Router) *Page {
	return &Page{
		Router: router,
	}
}

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

	// Handle reading data from SQLite when the button is clicked
	if p.ReadBtn.Clicked(gtx) {
		p.Data = p.readDataFromDB("UI/access.sqlite")
	}
	for i := range p.Data {
		if p.Data[i].EditBtn.Clicked(gtx) {
			log.Printf("Edit button clicked for ID: %d", p.Data[i].ID)
		}
	}
	return material.List(th, &p.List).Layout(gtx, len(p.Data)+2, func(gtx C, index int) D {
		return layout.UniformInset(unit.Dp(4)).Layout(gtx, func(gtx C) D {
			switch index {
			case 0:
				return p.LayoutReadButton(gtx, th)
			case 1:

				return p.LayoutColumnHeaders(gtx, th) // Hiển thị tên cột đầu tiên
			default:
				return p.LayoutDataRow(gtx, th, p.Data[index-2])
			}
		})
	})
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
		layout.Flexed(2, func(gtx C) D { // Flexed to make Filename wider
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

	// Draw header and horizontal line
	return layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx C) D {
			return headerDims
		}),
		layout.Stacked(func(gtx C) D {
			lineColor := color.NRGBA{R: 0, G: 0, B: 0, A: 255} // Black color
			lineHeight := unit.Dp(1)

			// Draw the horizontal line exactly under the header
			paint.ColorOp{Color: lineColor}.Add(gtx.Ops)
			defer clip.Rect{
				Max: image.Point{
					X: headerDims.Size.X, // Match the width of the header
					Y: int(lineHeight),
				},
			}.Push(gtx.Ops).Pop()
			paint.PaintOp{}.Add(gtx.Ops)

			return layout.Dimensions{
				Size: image.Point{
					X: headerDims.Size.X, // Match the width of the header
					Y: int(lineHeight),
				},
			}
		}),
	)
}

func (p *Page) LayoutReadButton(gtx C, th *material.Theme) D {
	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Flexed(1, func(gtx C) D {
			return layout.Center.Layout(gtx, func(gtx C) D {
				return layout.Inset{Left: unit.Dp(16), Right: unit.Dp(16), Top: unit.Dp(12)}.Layout(gtx, func(gtx C) D {
					gtx.Constraints.Min.X = gtx.Dp(unit.Dp(240))                   // Đặt chiều rộng tối thiểu cho ô nhập liệu
					editor := material.Editor(th, new(widget.Editor), "Search...") // Tạo Editor cho ô nhập liệu
					editor.Font.Style = font.Italic                                // Thiết lập font in nghiêng cho placeholder
					border := widget.Border{Color: color.NRGBA{A: 0xff}, CornerRadius: unit.Dp(6), Width: unit.Dp(0.5)}
					return border.Layout(gtx, func(gtx C) D {
						return layout.UniformInset(unit.Dp(8)).Layout(gtx, editor.Layout)
					})
				})
			})
		}),

		layout.Rigid(func(gtx C) D {
			// Layout for the "Search" button
			iconButton := material.IconButton(th, new(widget.Clickable), icon.SearchIcon, "Edit")
			iconButton.Color = color.NRGBA{A: 255}    // Set icon color
			iconButton.Size = unit.Dp(20)             // Set the size of the icon
			iconButton.Background = color.NRGBA{A: 0} // Transparent background
			return layout.UniformInset(unit.Dp(4)).Layout(gtx, iconButton.Layout)

			//return layout.UniformInset(unit.Dp(4)).Layout(gtx, func(gtx C) D {
			//	return material.IconButton(th, new(widget.Clickable), icon.SearchIcon, "Search").Layout(gtx)
			//})
		}),
		layout.Rigid(func(gtx C) D {
			// Layout for the "Read" button
			return layout.UniformInset(unit.Dp(4)).Layout(gtx,
				material.Button(th, &p.ReadBtn, "Read").Layout,
			)
		}),
	)
}

func (p *Page) LayoutDataRow(gtx C, th *material.Theme, row DataRow) D {
	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return layout.UniformInset(unit.Dp(4)).Layout(gtx,
				material.Body1(th, fmt.Sprintf("%d", row.ID)).Layout,
			)
		}),
		layout.Flexed(2, func(gtx C) D { // Flexed để Filename chiếm toàn bộ chiều ngang còn lại
			return layout.UniformInset(unit.Dp(4)).Layout(gtx,
				material.Body1(th, row.Filename).Layout,
			)
		}),
		layout.Flexed(1, func(gtx C) D { // Flexed để Modified chiếm toàn bộ chiều ngang còn lại
			return layout.UniformInset(unit.Dp(4)).Layout(gtx,
				material.Body1(th, row.Modified.Format("2006-01-02 15:04:05")).Layout,
			)
		}),
		//layout.Flexed(1, func(gtx C) D { // Flexed để Modified chiếm toàn bộ chiều ngang còn lại
		//	return layout.UniformInset(unit.Dp(4)).Layout(gtx,
		//		material.Body1(th, row.Folder).Layout,
		//	)
		//}),
		layout.Rigid(func(gtx C) D {
			// Create the icon button
			iconButton := material.IconButton(th, &row.EditBtn, icon.EditIcon, "Edit")
			iconButton.Color = color.NRGBA{A: 255}     // Set icon color
			iconButton.Size = unit.Dp(18)              // Set the size of the icon
			iconButton.Background = color.NRGBA{A: 80} // Transparent background
			return layout.UniformInset(unit.Dp(4)).Layout(gtx, iconButton.Layout)
		}),
	)
}

func (p *Page) readDataFromDB(dbFile string) []DataRow {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, filename, modified FROM tbl_fileinfo")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var data []DataRow
	for rows.Next() {
		var row DataRow
		var modifiedStr string
		if err := rows.Scan(&row.ID, &row.Filename, &modifiedStr); err != nil {
			log.Fatal(err)
		}
		row.Modified, err = time.Parse("2006-01-02T15:04:05Z", modifiedStr)
		if err != nil {
			log.Fatal(err)
		}
		data = append(data, row)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return data
}
func LayoutSearchField(gtx layout.Context, th *material.Theme, searchEditor *widget.Editor) layout.Dimensions {
	return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{Left: unit.Dp(16), Right: unit.Dp(16)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Min.X = gtx.Dp(unit.Dp(240)) // Đặt chiều rộng tối thiểu cho trường tìm kiếm
			editor := material.Editor(th, searchEditor, "Search...")
			editor.Font.Style = font.Italic // Đặt kiểu font cho text là in nghiêng
			border := widget.Border{
				Color:        color.NRGBA{A: 0xff}, // Đặt màu cho đường viền
				CornerRadius: unit.Dp(6),           // Đặt độ cong cho các góc
				Width:        unit.Dp(0.5),         // Đặt độ dày cho đường viền
			}
			return border.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.UniformInset(unit.Dp(8)).Layout(gtx, editor.Layout)
			})
		})
	})
}
