package discover

import (
	"github.com/gdamore/tcell"

	"github.com/rivo/tview"
)

type State struct {
	Namespace   string
	Application string
	Process     string
}

var currentState *State

// TODO 1. shortcut key be first char of name (second for repeated ones)
// TODO 2. Tab should switch focus to next box

func RunDiscovery() {
	currentState := &State{}
	app := tview.NewApplication()
	//newPrimitive := func() tview.Primitive {
	//	return tview.NewTextView().SetText("logiqctl dfgd dsgfdf s sdf").SetBorder(true).SetTitle(" logiqctl ")
	//}

	nsView := getNamespacesList()
	nsView.SetBorder(true)
	nsView.SetTitle(" NameSpaces ")
	defaultBorderColor := nsView.GetBorderColor()
	nsView.SetBorderColor(tcell.ColorIndianRed)

	appsView := tview.NewList()
	appsView.SetBorder(true)
	appsView.SetTitle(" Applications ")

	procView := tview.NewList()
	procView.SetBorder(true)
	procView.SetTitle(" Processes ")

	nsView.SetSelectedFunc(func(i int, ns string, des string, r rune) {
		currentState.Namespace = ns
		currentState.Application = ""
		currentState.Process = ""
		updateApplicationsList(ns, appsView)
		app.SetFocus(appsView)
		appsView.SetBorderColor(tcell.ColorIndianRed)
		nsView.SetBorderColor(defaultBorderColor)
		procView.Clear()
	})

	logsView := tview.NewTextView()
	logsView.SetBorder(true)
	logsView.SetTitle(" Logs ")

	grid := tview.NewGrid().
		//SetRows(0).
		//SetColumns(0, 0, 0).
		AddItem(nsView, 0, 0, 1, 1, 0, 0, true).
		AddItem(appsView, 0, 1, 1, 1, 0, 0, false).
		AddItem(procView, 0, 2, 1, 1, 0, 0, false)

	appsView.SetSelectedFunc(func(i int, selectApp string, des string, r rune) {
		currentState.Application = selectApp
		currentState.Process = ""
		updateProcessList(currentState.Namespace, selectApp, procView)
		app.SetFocus(procView)
		procView.SetBorderColor(tcell.ColorIndianRed)
		appsView.SetBorderColor(defaultBorderColor)
	})

	procView.SetSelectedFunc(func(i int, selectProc string, des string, r rune) {
		currentState.Process = selectProc
		//focus to shift to logs view
		grid.SetRows(20, 0).AddItem(logsView, 1, 0, 1, 3, 0, 0, true)
	})

	//Used to rotate Focus
	nsView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab:
			if currentState.Application != "" {
				//Move focus only if app is populated
				app.SetFocus(appsView)
				appsView.SetBorderColor(tcell.ColorIndianRed)
				nsView.SetBorderColor(defaultBorderColor)
				//TODO rotate focus even if no select app
				//Ideally populate app for the highlighted ns even if not selected
				//TODO handle mouse focus
			}
			break
		default:
			return event
		}
		return nil
	})

	appsView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab:
			if currentState.Process != "" {
				//Move focus only if process is populated
				app.SetFocus(procView)
				procView.SetBorderColor(tcell.ColorIndianRed)
				appsView.SetBorderColor(defaultBorderColor)
			} else {
				app.SetFocus(nsView)
				nsView.SetBorderColor(tcell.ColorIndianRed)
				appsView.SetBorderColor(defaultBorderColor)
			}
			break
		default:
			return event
		}
		return nil
	})

	procView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab:
			app.SetFocus(nsView)
			nsView.SetBorderColor(tcell.ColorIndianRed)
			procView.SetBorderColor(defaultBorderColor)
			break
		default:
			return event
		}
		return nil
	})

	if err := app.SetRoot(grid, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
