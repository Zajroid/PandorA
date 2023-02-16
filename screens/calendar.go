package screens

import (
	"fmt"
	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"image/color"
	"time"
)

const (
	dateSize = 20
	daySize  = 15
	formSize = 40
)

func makeCellContent(date *canvas.Text, isToday bool) *fyne.Container {
	top := fyne.NewContainerWithLayout(layout.NewHBoxLayout(), layout.NewSpacer(), date)

	line := canvas.NewLine(color.White)
	line.StrokeWidth = 1

	rec := canvas.NewRectangle(color.Transparent)
	rec.SetMinSize(fyne.NewSize(formSize, formSize))
	form := fyne.NewContainerWithLayout(layout.NewMaxLayout(), rec)

	var back *canvas.Rectangle
	if isToday {
		back = canvas.NewRectangle(color.RGBA{R: 50, G: 50, A: 10})
	} else {
		back = canvas.NewRectangle(color.Transparent)
	}

	return fyne.NewContainerWithLayout(
		layout.NewMaxLayout(),
		back,
		fyne.NewContainerWithLayout(layout.NewVBoxLayout(), top, form, line),
	)
}

func makeCell(t time.Time, isToday bool) (cell *fyne.Container) {
	_, _, day := t.Date()
	date := canvas.NewText(fmt.Sprint(day), color.White)
	date.TextSize = dateSize

	switch t.Weekday() {
	case 0:
		date.Color = color.RGBA{R: 255, A: 255}
	case 6:
		date.Color = color.RGBA{B: 255, A: 50}
	}

	return makeCellContent(date, isToday)
}

func makedWeekDayLabel() (week *fyne.Container) {
	days := make([]*canvas.Text, 0, 7)
	week = fyne.NewContainerWithLayout(layout.NewGridLayoutWithColumns(7))

	days = append(days, canvas.NewText("Sun", color.RGBA{R: 255, A: 255}))
	for _, text := range []string{"Mon.", "Tue.", "Wed.", "Thu.", "Fri."} {
		days = append(days, canvas.NewText(text, color.White))
	}
	days = append(days, canvas.NewText("Sat.", color.RGBA{B: 255, A: 50}))

	for _, day := range days {
		day.Alignment = fyne.TextAlignCenter
		day.TextSize = 15
		week.AddObject(
			fyne.NewContainerWithLayout(layout.NewVBoxLayout(), day, canvas.NewLine(color.White)),
		)
	}
	return
}

func makeMonthGrid(year int, month time.Month) (monthGrid *fyne.Container) {
	loc, _ := time.LoadLocation("Local")
	day := time.Date(year, month, 1, 0, 0, 0, 0, loc)
	now := time.Now()

	count := int(day.Weekday())
	padding := make([]fyne.CanvasObject, 0, 6)
	for i := 0; i < count; i++ {
		padding = append(
			padding,
			makeCellContent(&canvas.Text{Color: color.Transparent, Text: "", TextSize: dateSize}, false),
		)
	}

	cells := make([]fyne.CanvasObject, 0, 42)
	cells = append(cells, padding...)
	for day.Month() == month {
		y, m, d := now.Date()
		isToday := (y == day.Year() && m == day.Month() && d == day.Day())
		cells = append(cells, makeCell(day, isToday))
		day = day.AddDate(0, 0, 1)
	}

	return fyne.NewContainerWithLayout(
		layout.NewVBoxLayout(),
		makedWeekDayLabel(),
		fyne.NewContainerWithLayout(layout.NewGridLayout(7), cells...),
	)
}

func SemesterScreen(base fyne.Window) fyne.CanvasObject {
	now := time.Now()
	year := now.Year()
	month := int(now.Month())

	var semester *widget.TabContainer
	switch month {
	case 4, 5, 6, 7, 8, 9:
		semester = widget.NewTabContainer(
			widget.NewTabItem("April", makeMonthGrid(year, time.April)),
			widget.NewTabItem("May", makeMonthGrid(year, time.May)),
			widget.NewTabItem("June", makeMonthGrid(year, time.June)),
			widget.NewTabItem("July", makeMonthGrid(year, time.July)),
			widget.NewTabItem("August", makeMonthGrid(year, time.August)),
			widget.NewTabItem("September", makeMonthGrid(year, time.September)),
		)
		semester.SelectTabIndex(month - 4)
	case 10, 11, 12, 1, 2, 3:
		semester = widget.NewTabContainer(
			widget.NewTabItem("October", makeMonthGrid(year, time.October)),
			widget.NewTabItem("November", makeMonthGrid(year, time.November)),
			widget.NewTabItem("December", makeMonthGrid(year, time.December)),
			widget.NewTabItem("January", makeMonthGrid(year, time.January)),
			widget.NewTabItem("February", makeMonthGrid(year, time.February)),
			widget.NewTabItem("March", makeMonthGrid(year, time.March)),
		)
		semester.SelectTabIndex((month + 2) % 12)
	}

	toolbar := widget.NewToolbar(
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.SettingsIcon(), func() {
			showUpdateForm(base)
		}),
	)

	return fyne.NewContainerWithLayout(
		layout.NewBorderLayout(toolbar, nil, nil, nil),
		toolbar,
		semester,
	)
}
