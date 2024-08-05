package UI

import (
	"gioui.org/layout"
	"gioui.org/widget"
)

var (
	Editor              = widget.Editor{SingleLine: true}
	ProgressIncrementer chan float32
	LineEditor          = &widget.Editor{SingleLine: true, Submit: true}
	TopLabel            = "Hello, Gio"
	List                = &widget.List{List: layout.List{Axis: layout.Vertical}}
	UsernameEditor      widget.Editor
	PasswordEditor      widget.Editor
	LoginButton         widget.Clickable
	LoginScreen         bool = true
	Button1             widget.Clickable
	Button2             widget.Clickable
	Button3             widget.Clickable
	ErrorDialogButton   widget.Clickable
	ShowErrorDialog     bool = false
)
