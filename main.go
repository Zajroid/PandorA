package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/theme"
	"pandora/screens"
)

func main() {
	pandora := app.New()
	pandora.Settings().SetTheme(theme.DarkTheme())

	window := pandora.NewWindow("Pandora")
	object := screens.SemesterScreen(window)

	window.Resize(fyne.NewSize(700, 500))
	window.SetContent(object)
	window.ShowAndRun()
}
