package chart

import (
	"Gio_UI/UI"
	page "Gio_UI/UI/app"
	"Gio_UI/UI/icon"
	"gioui.org/font"
	"gioui.org/layout"
	"github.com/fogleman/gg"
	"github.com/wcharczuk/go-chart/v2"
	"image/color"
	"path/filepath"

	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"os"
)

type (
	C = layout.Context
	D = layout.Dimensions
)
type Page struct {
	generateButton         widget.Clickable
	generatePieChartButton widget.Clickable
	generateChartButton    widget.Clickable
	generateTableButton    widget.Clickable
	chartData              chart.BarChart
	pieChartData           chart.PieChart
	LinGraphData           chart.Chart
	tableImage             string
	*page.Router
}

func New(router *page.Router) *Page {
	return &Page{
		Router: router,
		chartData: chart.BarChart{
			Background: chart.Style{
				Padding: chart.Box{
					Top: 40,
				},
			},
			Width:  380,
			Height: 200,
			Bars: []chart.Value{
				{Value: 5, Label: "A"},
				{Value: 10, Label: "B"},
				{Value: 15, Label: "C"},
			},
		},
		pieChartData: chart.PieChart{
			Values: []chart.Value{
				{Value: 30, Label: "X"},
				{Value: 45, Label: "Y"},
				{Value: 25, Label: "Z"},
			},
		},
		tableImage: "table_output.png",
		LinGraphData: chart.Chart{
			Background: chart.Style{
				Padding: chart.Box{
					Top:  20,
					Left: 260,
				},
			},
			Series: []chart.Series{
				chart.ContinuousSeries{
					Name:    "1",
					XValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
					YValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
				},

				chart.ContinuousSeries{
					Name:    "2",
					XValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
					YValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
				},

				chart.ContinuousSeries{
					Name:    "3",
					XValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
					YValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
				},

				chart.ContinuousSeries{
					Name:    "4",
					XValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
					YValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
				},

				chart.ContinuousSeries{
					Name:    "5",
					XValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
					YValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
				},

				chart.ContinuousSeries{
					Name:    "6",
					XValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
					YValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
				},

				chart.ContinuousSeries{
					Name:    "7",
					XValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
					YValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
				},

				chart.ContinuousSeries{
					Name:    "8",
					XValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
					YValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
				},

				chart.ContinuousSeries{
					Name:    "9",
					XValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
					YValues: []float64{5.0, 4.0, 3.0, 2.0, 1.0},
				},

				chart.ContinuousSeries{
					Name:    "10",
					XValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
					YValues: []float64{5.0, 4.0, 3.0, 2.0, 1.0},
				},

				chart.ContinuousSeries{
					Name:    "11",
					XValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
					YValues: []float64{5.0, 4.0, 3.0, 2.0, 1.0},
				},
			},
		},
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
		Name: "Chart",
		Icon: icon.ChartIcon, // Bạn cần tự định nghĩa icon này
	}
}

func (p *Page) Layout(gtx C, th *material.Theme) D {
	// Đường dẫn thư mục
	dirPath := "UI/Genimage"
	// Tạo thư mục nếu chưa tồn tại
	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		// Xử lý lỗi nếu không tạo được thư mục
		return D{}
	}

	if p.generateButton.Clicked(gtx) {
		filePath := filepath.Join(dirPath, "barchart.png")
		f, _ := os.Create(filePath)
		defer f.Close()
		p.chartData.Render(chart.PNG, f)
	}

	if p.generatePieChartButton.Clicked(gtx) {
		filePath := filepath.Join(dirPath, "piechart.png")
		f, _ := os.Create(filePath)
		defer f.Close()
		p.pieChartData.Render(chart.PNG, f)
	}

	if p.generateChartButton.Clicked(gtx) {
		p.LinGraphData.Elements = []chart.Renderable{
			chart.LegendLeft(&p.LinGraphData),
		}
		filePath := filepath.Join(dirPath, "LineGraph.png")
		f, _ := os.Create(filePath)
		defer f.Close()
		p.LinGraphData.Render(chart.PNG, f)
	}
	if p.generateTableButton.Clicked(gtx) {
		// Tạo một bảng đơn giản với thư viện gg
		const W = 400
		const H = 200

		dc := gg.NewContext(W, H)
		dc.SetColor(color.White)
		dc.Clear()

		dc.SetColor(color.Black)
		dc.DrawStringAnchored("Simple Table", W/2, 20, 0.5, 0.5)

		data := [][]string{
			{"Name", "Age", "City"},
			{"Alice", "30", "New York"},
			{"Bob", "25", "San Francisco"},
			{"Charlie", "35", "Los Angeles"},
		}

		cellWidth := float64(W / len(data[0]))
		cellHeight := float64(H-40) / float64(len(data))

		for i, row := range data {
			for j, str := range row {
				x := float64(j) * cellWidth
				y := float64(i)*cellHeight + 40
				dc.DrawRectangle(x, y, cellWidth, cellHeight)
				dc.Stroke()
				dc.DrawStringAnchored(str, x+cellWidth/2, y+cellHeight/2, 0.5, 0.5)
			}
		}

		filePath := filepath.Join(dirPath, p.tableImage)
		dc.SavePNG(filePath)
	}
	return layout.Flex{}.Layout(gtx,
		layout.Flexed(1, func(gtx C) D {
			return layout.UniformInset(unit.Dp(16)).Layout(gtx, func(gtx C) D {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return UI.DetailRow{}.Layout(gtx,
							material.Body1(th, "Bar chart").Layout,
							func(gtx UI.C) UI.D {
								return material.Button(th, &p.generateButton, "BarChart").Layout(gtx)
							})
					}),
					layout.Rigid(func(gtx C) D {
						gtx.Constraints.Min.X = gtx.Dp(unit.Dp(36))
						e := material.Body1(th, "Pie chart")
						e.Font.Style = font.Italic
						e.Color = color.NRGBA(chart.ColorAlternateBlue)
						return UI.DetailRow{}.Layout(gtx,
							e.Layout,
							func(gtx UI.C) UI.D {
								return material.Button(th, &p.generatePieChartButton, "PieChart").Layout(gtx)
							})
					}),
					layout.Rigid(func(gtx C) D {
						gtx.Constraints.Min.X = gtx.Dp(unit.Dp(36))
						e := material.Body1(th, "LineGraph")
						e.Font.Style = font.Italic
						e.Color = color.NRGBA(chart.ColorOrange)
						return UI.DetailRow{}.Layout(gtx,
							e.Layout,
							func(gtx UI.C) UI.D {
								return material.Button(th, &p.generateChartButton, "LineGraph").Layout(gtx)
							})
					}),
					layout.Rigid(func(gtx C) D {
						gtx.Constraints.Min.X = gtx.Dp(unit.Dp(36))
						e := material.Body1(th, "Generate Table")
						e.Font.Style = font.Italic
						e.Color = color.NRGBA(chart.ColorOrange)
						return UI.DetailRow{}.Layout(gtx,
							e.Layout,
							func(gtx UI.C) UI.D {
								return material.Button(th, &p.generateChartButton, "Generate Table").Layout(gtx)
							})
					}),
				)
			})
		}),
	)
}
