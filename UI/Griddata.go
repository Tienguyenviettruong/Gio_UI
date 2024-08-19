package UI

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"image/color"
)

func LayoutFlexGrid(gtx layout.Context, th *material.Theme) error {
	// Define the number of columns.
	columns := 3

	// Create a flex layout.
	flex := layout.Flex{
		Axis: layout.Vertical,
	}

	// Create a list of flex items.
	var items []layout.FlexChild
	for i := 0; i < 100; i++ { // Example with 9 items
		index := i
		items = append(items, layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			// Define the style for each item.
			item := material.Button(th, &widget.Clickable{}, "Item "+string(index))
			item.Background = color.NRGBA{R: 200, G: 200, B: 200, A: 255}
			return item.Layout(gtx)
		}))
	}

	// Layout items in a grid-like fashion using flex.
	for i := 0; i < len(items); i += columns {
		end := i + columns
		if end > len(items) {
			end = len(items)
		}
		flex.Layout(gtx, items[i:end]...)
	}
	return nil
}
