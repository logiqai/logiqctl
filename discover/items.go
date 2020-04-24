package discover

import (
	"fmt"

	"github.com/logiqai/logiqctl/services"
	"github.com/rivo/tview"
)

func getNamespacesList() *tview.List {
	nsView := tview.NewList()
	nameSpaces := services.GetNamespaces()

	for i, ns := range nameSpaces {
		description := fmt.Sprintf("%d-%d", ns.LastSeen, ns.FirstSeen)
		nsView.AddItem(ns.Namespace, description, rune(97+i), nil)
	}

	return nsView
}

func updateApplicationsList(namespace string, appView *tview.List) {
	appView.Clear()
	applications := services.GetApplicationsV2(namespace)
	for i, app := range applications {
		description := fmt.Sprintf("%d-%d", app.LastSeen, app.FirstSeen)
		appView.AddItem(app.Name, description, rune(97+i), nil)
	}
}

func updateProcessList(namespace string, appName string, processView *tview.List) {
	processView.Clear()
	applications := services.GetProcesses(namespace, appName)
	for i, app := range applications {
		description := fmt.Sprintf("%d-%d", app.LastSeen, app.FirstSeen)
		processView.AddItem(app.ProcID, description, rune(97+i), nil)
	}
}
