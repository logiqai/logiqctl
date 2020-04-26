package discover

import (
	"fmt"
	"time"

	"github.com/logiqai/logiqctl/services"

	"github.com/gdamore/tcell"
	"github.com/logiqai/logiqctl/api/v1/query"

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

var aLine = []query.SysLogMessage{{
	AppName:        "AppName",
	FacilityString: "info",
	Message:        "shortcut key be first char of name (second for repeated ones), shortcut key be first char of name (second for repeated ones)",
	ProcID:         "283838",
	Timestamp:      time.RFC3339,
	Namespace:      "tito-home",
},
	{
		AppName:        "AppName",
		FacilityString: "info",
		Message:        "shortcut key be first char of name (second for repeated ones), shortcut key be first char of name (second for repeated ones)",
		ProcID:         "283838",
		Timestamp:      time.RFC3339,
		Namespace:      "tito-home",
	},
	{
		AppName:        "AppName",
		FacilityString: "info",
		Message:        "shortcut key be first char of name (second for repeated ones), shortcut key be first char of name (second for repeated ones)",
		ProcID:         "283838",
		Timestamp:      time.RFC3339,
		Namespace:      "tito-home",
	},
	{
		AppName:        "AppName",
		FacilityString: "info",
		Message:        "shortcut key be first char of name (second for repeated ones), shortcut key be first char of name (second for repeated ones)",
		ProcID:         "283838",
		Timestamp:      time.RFC3339,
		Namespace:      "tito-home",
	},
	{
		AppName:        "AppName",
		FacilityString: "info",
		Message:        "shortcut key be first char of name (second for repeated ones), shortcut key be first char of name (second for repeated ones)",
		ProcID:         "283838",
		Timestamp:      time.RFC3339,
		Namespace:      "tito-home",
	},
	{
		AppName:        "AppName",
		FacilityString: "info",
		Message:        "shortcut key be first char of name (second for repeated ones), shortcut key be first char of name (second for repeated ones)",
		ProcID:         "283838",
		Timestamp:      time.RFC3339,
		Namespace:      "tito-home",
	},
	{
		AppName:        "AppName",
		FacilityString: "info",
		Message:        "shortcut key be first char of name (second for repeated ones), shortcut key be first char of name (second for repeated ones), \n shortcut key be first char of name , shortcut key be first char of name , shortcut key be first char of name ",
		ProcID:         "283838",
		Timestamp:      time.RFC3339,
		Namespace:      "tito-home",
	},
}

func RunDiscovery() {
	currentState := &State{}
	app := tview.NewApplication()

	statusView := tview.NewTextView()
	statusView.SetBorder(true)
	statusView.SetText("Select a Namespace")

	text := make(chan string)

	go func() {
		for {
			select {
			case message := <-text:
				statusView.SetText(message)
			}
		}
	}()

	logsTable := tview.NewTable()
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
		text <- "Fetching Applications"
		updateApplicationsList(ns, appsView)
		app.SetFocus(appsView)
		appsView.SetBorderColor(tcell.ColorIndianRed)
		nsView.SetBorderColor(defaultBorderColor)
		text <- "Select a Application"
		procView.Clear()
	})

	logsGrid := tview.NewGrid()
	//logsGrid.SetBorder(true)
	//logsGrid.SetTitle(" Logs ")
	//logsFlex1 := tview.NewFlex().SetDirection(tview.FlexRow)

	grid := tview.NewGrid().
		//SetRows(0).
		//SetColumns(0, 0, 0).
		SetRows(0, 3).
		AddItem(nsView, 0, 0, 1, 1, 0, 0, true).
		AddItem(appsView, 0, 1, 1, 1, 0, 0, false).
		AddItem(procView, 0, 2, 1, 1, 0, 0, false).
		AddItem(statusView, 1, 0, 1, 3, 0, 0, false)

	appsView.SetSelectedFunc(func(i int, selectApp string, des string, r rune) {
		currentState.Application = selectApp
		currentState.Process = ""
		text <- "Fetching Processes"
		updateProcessList(currentState.Namespace, selectApp, procView)
		text <- "Select a process to display logs"
		app.SetFocus(procView)
		procView.SetBorderColor(tcell.ColorIndianRed)
		appsView.SetBorderColor(defaultBorderColor)
	})

	procView.SetSelectedFunc(func(i int, selectProc string, des string, r rune) {
		currentState.Process = selectProc
		//focus to shift to logs view

		grid.SetRows(18, 3, 0).AddItem(logsGrid, 2, 0, 1, 3, 0, 0, true)

		data := services.QueryAndGetDataMock(currentState.Namespace, currentState.Application, currentState.Process)

		logsTable.SetSelectable(true, false)
		logsTable.SetTitle(fmt.Sprintf(" Log for %s (Namespace), %s (Application) and %s (Process) ", currentState.Namespace, currentState.Application, currentState.Process)).SetBorder(true)
		for i, log := range data {
			var color = defaultBorderColor
			switch log.SeverityString {
			case "warning":
				color = tcell.ColorYellow
				break
			case "critical":
				color = tcell.ColorOrange
			case "emergency":
				color = tcell.ColorIndianRed
			}
			logsTable.SetCell(i, 0, tview.NewTableCell(log.Timestamp).SetTextColor(color))
			logsTable.SetCell(i, 1, tview.NewTableCell(log.SeverityString).SetTextColor(color))
			logsTable.SetCell(i, 2, tview.NewTableCell(log.FacilityString).SetTextColor(color))
			logsTable.SetCell(i, 3, tview.NewTableCell(log.Message).SetMaxWidth(0).SetTextColor(color))

		}
		app.SetFocus(logsTable)
		logsGrid.AddItem(logsTable, 0, 0, 1, 1, 0, 0, false)
		logsTable.SetBorderColor(tcell.ColorIndianRed)
		procView.SetBorderColor(defaultBorderColor)
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

	logsGrid.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab:

			//Move focus only if app is populated
			app.SetFocus(nsView)
			nsView.SetBorderColor(tcell.ColorIndianRed)
			logsTable.SetBorderColor(defaultBorderColor)
			//TODO rotate focus even if no select app
			//Ideally populate app for the highlighted ns even if not selected
			//TODO handle mouse focus

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
