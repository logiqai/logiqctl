package discover

import (
	"fmt"
	"time"

	"github.com/dustin/go-humanize"

	"github.com/logiqai/logiqctl/services"
	"github.com/rivo/tview"
)

func getNamespacesList() *tview.List {
	nsView := tview.NewList()
	nameSpaces := services.GetNamespacesMock()

	for i, ns := range nameSpaces {
		lastSeen := humanize.Time(time.Unix(ns.LastSeen, 0))
		description := fmt.Sprintf("Updated %s", lastSeen)
		nsView.AddItem(ns.Namespace, description, rune(97+i), nil)
	}

	return nsView
}

func updateApplicationsList(namespace string, appView *tview.List) {
	appView.Clear()
	applications := services.GetApplicationsV2Mock(namespace)
	for i, app := range applications {
		lastSeen := humanize.Time(time.Unix(app.LastSeen, 0))
		description := fmt.Sprintf("Updated %s", lastSeen)
		appView.AddItem(app.Name, description, rune(97+i), nil)
	}
}

func updateProcessList(namespace string, appName string, processView *tview.List) {
	processView.Clear()
	processes := services.GetProcessesMock(namespace, appName)
	for i, process := range processes {
		lastSeen := humanize.Time(time.Unix(process.LastSeen, 0))
		description := fmt.Sprintf("Updated %s", lastSeen)
		processView.AddItem(process.ProcID, description, rune(97+i), nil)
	}
}
