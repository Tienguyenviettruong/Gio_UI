package Table

import (
	"Gio_UI/UI/icon"
	"database/sql"
	"fmt"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
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
	ReadBtn widget.Clickable
	Data    []DataRow
}

type DataRow struct {
	ID       int
	Filename string
	Modified time.Time
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
		Icon: icon.VisibilityIcon,
	}
}

func (p *Page) Layout(gtx C, th *material.Theme) D {
	p.List.Axis = layout.Vertical

	// Handle reading data from SQLite when the button is clicked
	if p.ReadBtn.Clicked(gtx) {
		p.Data = p.readDataFromDB("UI/access.sqlite")
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

// LayoutColumnHeaders hiển thị các tiêu đề cột với chữ in đậm
func (p *Page) LayoutColumnHeaders(gtx C, th *material.Theme) D {
	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			label := material.Body1(th, "ID")
			label.Font.Weight = 800 // In đậm
			return layout.UniformInset(unit.Dp(4)).Layout(gtx, label.Layout)
		}),
		layout.Flexed(1, func(gtx C) D {
			label := material.Body1(th, "Filename")
			label.Font.Weight = 800 // In đậm
			return layout.UniformInset(unit.Dp(4)).Layout(gtx, label.Layout)
		}),
		layout.Flexed(1, func(gtx C) D {
			label := material.Body1(th, "Modified")
			label.Font.Weight = 800 // In đậm
			return layout.UniformInset(unit.Dp(4)).Layout(gtx, label.Layout)
		}),
	)
}

func (p *Page) LayoutReadButton(gtx C, th *material.Theme) D {
	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
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
		layout.Flexed(1, func(gtx C) D { // Flexed để Filename chiếm toàn bộ chiều ngang còn lại
			return layout.UniformInset(unit.Dp(4)).Layout(gtx,
				material.Body1(th, row.Filename).Layout,
			)
		}),
		layout.Flexed(1, func(gtx C) D { // Flexed để Modified chiếm toàn bộ chiều ngang còn lại
			return layout.UniformInset(unit.Dp(4)).Layout(gtx,
				material.Body1(th, row.Modified.Format("2006-01-02 15:04:05")).Layout,
			)
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
