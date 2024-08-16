package Tree

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

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
	Icon     widget.Image
	Children []TreeNode
	component.DiscloserState
	font font.Style
}

type (
	C = layout.Context
	D = layout.Dimensions
)

type Page struct {
	TreeNode
	widget.List
	*page.Router
	CustomDiscloserState component.DiscloserState
	font.Style

	// New fields for input and button
	Input   widget.Editor
	AddBtn  widget.Clickable
	ViewBtn widget.Clickable
}

func New(router *page.Router) *Page {
	page := &Page{
		Router: router,
		Input:  widget.Editor{SingleLine: true},
	}

	if err := page.LoadTreeData("C:/Users/truongtnv/Desktop/JobFIMS_New/Gio_UI/UI/app/root.json"); err != nil {
		fmt.Println("Error loading tree data:", err)
	}

	return page
}

// LoadTreeData reads JSON data from a file and converts it to TreeNode structure.
func (p *Page) LoadTreeData(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	var rootNode TreeNode
	if err := json.Unmarshal(data, &rootNode); err != nil {
		return err
	}

	p.TreeNode = rootNode
	return nil
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
		Icon: icon.FolderIcon,
	}
}

func (p *Page) LayoutTreeNode(gtx C, th *material.Theme, tn *TreeNode) D {
	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			//
			//e := material.Body1(th, tn.Text)
			//e.Font.Style = font.Italic
			item := component.MenuItem(th, &p.AddBtn, tn.Text)
			item.Icon = icon.FolderIcon
			item.IconColor = th.ContrastBg
			if len(tn.Children) == 0 {
				return layout.UniformInset(unit.Dp(2)).Layout(gtx, item.Layout)
			}
			children := make([]layout.FlexChild, 0, len(tn.Children))
			for i := range tn.Children {
				child := &tn.Children[i]
				children = append(children, layout.Rigid(
					func(gtx C) D {
						return p.LayoutTreeNode(gtx, th, child)
					}))
			}
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return component.SimpleDiscloser(th, &tn.DiscloserState).Layout(gtx, item.Layout,
						func(gtx C) D {
							return layout.Flex{Axis: layout.Vertical}.Layout(gtx, children...)
						})
				}),
			)
		}),
	)
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
