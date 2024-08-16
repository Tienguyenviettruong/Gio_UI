package UI

import (
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget"
	"gioui.org/x/component"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

var (
	Editor                        = widget.Editor{SingleLine: true}
	ProgressIncrementer           chan float32
	LineEditor                    = &widget.Editor{SingleLine: true, Submit: true}
	TopLabel                      = "Hello, Gio"
	List                          = &widget.List{List: layout.List{Axis: layout.Vertical}}
	UsernameEditor                widget.Editor
	PasswordEditor                widget.Editor
	LoginButton                   widget.Clickable
	LoginScreen                   bool = true
	ConnectBtn                    widget.Clickable
	Button2                       widget.Clickable
	Button3                       widget.Clickable
	ErrorDialogButton             widget.Clickable
	ShowErrorDialog               bool = false
	User, Host, Port, Servicename component.TextField
	inputAlignment                text.Alignment
	Password                      = component.TextField{Editor: widget.Editor{Mask: '*'}}
	iconVisibility                *widget.Icon // Icon hiển thị mật khẩu
	iconVisibilityOff             *widget.Icon
	ShowPasswordBtn               widget.Clickable
	DBConnected                   bool = false
	ShowLoginScreen               bool = false
)

func init() {
	// Nạp biểu tượng
	iconVisibility, _ = widget.NewIcon(icons.ActionVisibility)
	iconVisibilityOff, _ = widget.NewIcon(icons.ActionVisibilityOff)
}
