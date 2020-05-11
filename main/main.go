// Demo code for the InputField primitive.
package main

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()

	inputField := tview.NewInputField().
		SetLabel("Enter Days: ").
		SetPlaceholder("0").
		SetFieldWidth(10).
		SetAcceptanceFunc(tview.InputFieldInteger).
		SetDoneFunc(func(key tcell.Key) {
			app.Stop()
		})
	inputField1 := tview.NewInputField().
		SetLabel("Enter Hours: ").
		SetPlaceholder("0").
		SetFieldWidth(10).
		SetAcceptanceFunc(tview.InputFieldInteger).
		SetDoneFunc(func(key tcell.Key) {
			app.Stop()
		})
	inputField2 := tview.NewInputField().
		SetLabel("Enter minutes: ").
		SetPlaceholder("60").
		SetFieldWidth(10).
		SetAcceptanceFunc(tview.InputFieldInteger).
		SetDoneFunc(func(key tcell.Key) {
			app.Stop()
		})

	button := tview.NewButton("Show Logs")

	flex := tview.NewFlex().
		AddItem(tview.NewBox().SetBorder(true), 0, 1, true).
		AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(inputField, 0, 1, true).
			AddItem(inputField1, 0, 1, false).
			AddItem(inputField2, 0, 1, false).
			AddItem(button, 0, 1, false),
			0, 1, true)

	if err := app.SetRoot(flex, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
