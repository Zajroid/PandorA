package screens

import (
	"fmt"
	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"image/color"
	"pandora/data"
)

func showUpdateForm(base fyne.Window) {
	top := canvas.NewText(
		"Enter your ECS_ID and Password. If blank, the values set previously will be used.",
		color.White,
	)
	top.Alignment = fyne.TextAlignCenter

	ecsID := widget.NewEntry()
	ecsID.PlaceHolder = "ECS-ID"
	password := widget.NewPasswordEntry()
	password.PlaceHolder = "Password"

	content := fyne.NewContainerWithLayout(
		layout.NewVBoxLayout(),
		top,
		ecsID,
		password,
	)

	dialog.ShowCustomConfirm(
		"Update Information",
		"Update",
		"Cancel",
		content,
		func(ok bool) {
			if ok {
				update(ecsID, password, base)
			}
		},
		base,
	)
}

func update(ecsID, password *widget.Entry, base fyne.Window) {
	if ecsID.Text != "" || password.Text != "" {
		if err := data.WriteAccountInfo(ecsID.Text, password.Text); err != nil {
			dialog.ShowError(err, base)
		}
	}

	id, pass, err := data.ReadAccountInfo()
	if err != nil {
		dialog.ShowError(err, base)
	}

	fmt.Println(id, pass)
}
