package discover

import (
	"fmt"
	"strconv"
	"time"

	"github.com/logiqai/logiqctl/api/v1/query"

	"github.com/logiqai/logiqctl/services"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type State struct {
	Namespace     string
	Application   string
	Process       string
	Day           int
	Hour          int
	Minute        int
	tableRowCount int
}

var currentState *State

// TODO 1. shortcut key be first char of name (second for repeated ones)
// TODO 2. Tab should switch focus to next box

func RunDiscovery() {
	currentState := &State{}
	app := tview.NewApplication()

	statusView := tview.NewTextView()

	statusView.SetTitle("| logiq.ai |").SetTitleAlign(2)
	statusView.SetBorder(true)
	statusView.SetText("Select a Namespace")

	text := make(chan string)
	defer close(text)
	buttonLabelText := make(chan string)
	defer close(buttonLabelText)
	syslogChan := make(chan query.SysLogMessage)
	defer close(syslogChan)

	logsTable := tview.NewTable()
	nsView := getNamespacesList()
	nsView.SetBorder(true)
	nsView.SetTitle("| NameSpaces |")
	defaultBorderColor := nsView.GetBorderColor()
	nsView.SetBorderColor(tcell.ColorIndianRed)

	appsView := tview.NewList()
	appsView.SetBorder(true)
	appsView.SetTitle("| Applications |")

	procView := tview.NewList()
	procView.SetBorder(true)
	procView.SetTitle("| Processes |")

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

	dayField := tview.NewInputField().
		SetLabel("Enter Day(s): ").
		SetPlaceholder("0").
		SetFieldWidth(10).
		SetAcceptanceFunc(tview.InputFieldInteger)

	hourField := tview.NewInputField().
		SetLabel("Enter Hour(s): ").
		SetPlaceholder("1").
		SetFieldWidth(10).
		SetAcceptanceFunc(tview.InputFieldInteger)

	minuteField := tview.NewInputField().
		SetLabel("Enter Minute(s): ").
		SetPlaceholder("0").
		SetFieldWidth(10).
		SetAcceptanceFunc(tview.InputFieldInteger)
	doneFunc := func(key tcell.Key) {
		days, _ := strconv.Atoi(dayField.GetText())
		hours, _ := strconv.Atoi(hourField.GetText())
		minutes, _ := strconv.Atoi(minuteField.GetText())

		t := time.Now()
		t = t.Add(time.Duration(days) * -24 * time.Hour)
		t = t.Add(-1 * time.Duration(hours) * time.Hour)
		t = t.Add(-1 * time.Duration(minutes) * time.Minute)

		ts := t.Format(time.RFC822)
		buttonLabelText <- ts
		text <- ts
	}

	dayField.SetDoneFunc(doneFunc)
	hourField.SetDoneFunc(doneFunc)
	minuteField.SetDoneFunc(doneFunc)
	button := tview.NewButton("Show Logs from (1 Hour)")

	go func() {
		for {
			time.Sleep(5 * time.Millisecond)
			select {
			case message := <-text:
				statusView.SetText(message)
				break
			case label := <-buttonLabelText:
				button.SetLabel(fmt.Sprintf("Show Logs from (%s)", label))
				break
			case sysLogMessage := <-syslogChan:
				var color = defaultBorderColor
				switch sysLogMessage.SeverityString {
				case "warning":
					color = tcell.ColorYellow
					break
				case "critical":
					color = tcell.ColorOrange
					break
				case "emergency":
					color = tcell.ColorRed
					break
				}
				logsTable.SetCell(currentState.tableRowCount, 0, tview.NewTableCell(sysLogMessage.Timestamp).SetTextColor(color))
				logsTable.SetCell(currentState.tableRowCount, 1, tview.NewTableCell(sysLogMessage.SeverityString).SetTextColor(color))
				logsTable.SetCell(currentState.tableRowCount, 2, tview.NewTableCell(sysLogMessage.FacilityString).SetTextColor(color))
				logsTable.SetCell(currentState.tableRowCount, 3, tview.NewTableCell(sysLogMessage.Message).SetMaxWidth(0).SetTextColor(color))
				currentState.tableRowCount++
				break
			}
		}
	}()

	//flex := tview.NewFlex().SetDirection(tview.FlexColumn).
	//	//AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
	//	AddItem(dayField, 0, 1, true).
	//	AddItem(hourField, 0, 1, false).
	//	AddItem(minuteField, 0, 1, false).
	//	AddItem(button, 0, 1, false)

	buttonGrid := tview.NewGrid().
		SetRows(1, 1, 1).
		AddItem(dayField, 1, 0, 1, 1, 0, 0, true).
		AddItem(hourField, 1, 1, 1, 1, 0, 0, true).
		AddItem(tview.NewBox(), 0, 1, 1, 1, 0, 0, true). //place holder
		AddItem(minuteField, 1, 2, 1, 1, 0, 0, true).
		AddItem(tview.NewBox(), 2, 1, 1, 1, 0, 0, true). //place holder
		AddItem(button, 0, 3, 3, 2, 0, 0, true)
	buttonGrid.SetBorder(true).SetTitle("| Select Relative Start Time |")

	grid := tview.NewGrid().
		//SetRows(0).
		//SetColumns(0, 0, 0).
		SetRows(3, 0, 5).
		AddItem(statusView, 0, 0, 1, 3, 0, 0, false).
		AddItem(nsView, 1, 0, 1, 1, 0, 0, true).
		AddItem(appsView, 1, 1, 1, 1, 0, 0, false).
		AddItem(procView, 1, 2, 1, 1, 0, 0, false).
		AddItem(buttonGrid, 2, 0, 1, 3, 0, 0, false)

	button.SetSelectedFunc(func() {
		grid.SetRows(3, 15, 5, 0, 1).
			AddItem(logsGrid, 3, 0, 1, 3, 0, 0, true).
			AddItem(tview.NewButton("show More"), 4, 0, 1, 3, 0, 0, true)

		services.QueryAndGetDataMock(currentState.Namespace, currentState.Application, currentState.Process, syslogChan)

		logsTable.SetSelectable(true, false)
		logsTable.SetTitle(fmt.Sprintf("| Logs for %s (Namespace), %s (Application) and %s (Process) |", currentState.Namespace, currentState.Application, currentState.Process)).SetBorder(true)

		app.SetFocus(logsTable)
		logsGrid.AddItem(logsTable, 0, 0, 1, 1, 0, 0, false)
		logsTable.SetBorderColor(tcell.ColorIndianRed)
		procView.SetBorderColor(defaultBorderColor)
	})

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
		text <- "Selected Process - selectProc"

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
