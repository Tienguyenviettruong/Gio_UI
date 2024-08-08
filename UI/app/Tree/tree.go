package Tree

import (
	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"

	page "Gio_UI/UI/app"
	"Gio_UI/UI/icon"
)

type TreeNode struct {
	Text     string
	Children []TreeNode
	component.DiscloserState
	font font.Style
}

type (
	C = layout.Context
	D = layout.Dimensions
)

// Page holds the state for a page demonstrating the features of
// the AppBar component.
type Page struct {
	TreeNode
	widget.List
	*page.Router
	CustomDiscloserState component.DiscloserState
	font.Style

	// New fields for input and button
	Input  widget.Editor
	AddBtn widget.Clickable
}

// New constructs a Page with the provided router.
func New(router *page.Router) *Page {
	return &Page{
		Router: router,
		TreeNode: TreeNode{
			font: font.Italic,
			Text: "Expand Me",
			Children: []TreeNode{
				{
					Text: "Disclosers can be (expand me)...",
					Children: []TreeNode{
						{
							Text: "...nested to arbitrary depths.",
						},
						{
							Text: "There are also types available to customize the look and feel of the discloser:",
							Children: []TreeNode{
								{
									Text: "• DiscloserStyle lets you provide your own control instead of the default triangle used here.",
								},
								{
									Text: "• DiscloserArrowStyle lets you alter the presentation of the triangle used here, like changing its color, size, left/right anchoring, or margin.",
								},
							},
						},
					},
				},
			},
		},
		Input: widget.Editor{SingleLine: true},
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
		Name: "Tree",
		Icon: icon.VisibilityIcon,
	}
}

func (p *Page) LayoutTreeNode(gtx C, th *material.Theme, tn *TreeNode) D {
	e := material.Body1(th, tn.Text)
	e.Font.Style = font.Italic
	if len(tn.Children) == 0 {
		return layout.UniformInset(unit.Dp(2)).Layout(gtx, e.Layout)
	}
	children := make([]layout.FlexChild, 0, len(tn.Children))
	for i := range tn.Children {
		child := &tn.Children[i]
		children = append(children, layout.Rigid(
			func(gtx C) D {
				return p.LayoutTreeNode(gtx, th, child)
			}))
	}
	return component.SimpleDiscloser(th, &tn.DiscloserState).Layout(gtx, e.Layout,
		func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx, children...)
		})
}

func (p *Page) LayoutCustomDiscloser(gtx C, th *material.Theme) D {
	return component.Discloser(th, &p.CustomDiscloserState).Layout(gtx,
		func(gtx C) D {
			var l material.LabelStyle
			l = material.Body1(th, "+")
			l.Font.Style = font.Italic
			l.Font.Typeface = "Go Mono"
			if p.CustomDiscloserState.Visible() {
				l.Text = "-"
			}

			return layout.UniformInset(unit.Dp(2)).Layout(gtx, l.Layout)
		},
		material.Subtitle2(th, "Custom Control").Layout,
		material.Body2(th, "This control only took 9 lines of code.").Layout,
	)
}

func (p *Page) Layout(gtx C, th *material.Theme) D {
	p.List.Axis = layout.Vertical

	// Handle adding a new node
	if p.AddBtn.Clicked(gtx) {
		newNode := TreeNode{
			Text: p.Input.Text(),
		}
		p.TreeNode.Children = append(p.TreeNode.Children, newNode)
		p.Input.SetText("")
	}

	return material.List(th, &p.List).Layout(gtx, 3, func(gtx C, index int) D {
		return layout.UniformInset(unit.Dp(4)).Layout(gtx, func(gtx C) D {
			switch index {
			case 0:
				return p.LayoutTreeNode(gtx, th, &p.TreeNode)
			case 1:
				return p.LayoutCustomDiscloser(gtx, th)
			case 2:
				return p.LayoutAddNode(gtx, th)
			}
			return D{}
		})
	})
}

// LayoutAddNode displays the input field and the "Add" button
func (p *Page) LayoutAddNode(gtx C, th *material.Theme) D {
	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return layout.UniformInset(unit.Dp(4)).Layout(gtx,
				material.Editor(th, &p.Input, "Enter node text").Layout,
			)
		}),
		layout.Rigid(func(gtx C) D {
			return layout.UniformInset(unit.Dp(4)).Layout(gtx,
				material.Button(th, &p.AddBtn, "Add Node").Layout,
			)
		}),
	)
}
