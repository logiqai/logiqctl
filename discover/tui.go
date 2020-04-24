package discover

import (
	"fmt"

	"github.com/logiqai/logiqctl/services"
	"github.com/rivo/tview"
)

type State struct {
	Namespace   string
	Application string
	Process     string
}

var currentState State

func RunDiscovery() {

	app := tview.NewApplication()
	//newPrimitive := func() tview.Primitive {
	//	return tview.NewTextView().SetText("logiqctl dfgd dsgfdf s sdf").SetBorder(true).SetTitle(" logiqctl ")
	//}

	nsView := getNamespacesList()
	nsView.SetBorder(true)
	nsView.SetTitle(" NameSpaces ")

	appsView := tview.NewList()
	appsView.SetBorder(true)
	appsView.SetTitle(" Applications ")

	nsView.SetSelectedFunc(func(i int, ns string, des string, r rune) {
		updateApplicationsList(ns, appsView)
	})

	procView := tview.NewList()
	procView.SetBorder(true)
	procView.SetTitle(" Processes ")

	appsView.SetSelectedFunc(func(i int, app string, des string, r rune) {

	})

	grid := tview.NewGrid().
		SetRows(0).
		SetColumns(0, 0, 0).
		AddItem(nsView, 0, 0, 1, 1, 0, 0, true).
		AddItem(appsView, 0, 1, 1, 1, 0, 0, false).
		AddItem(procView, 0, 2, 1, 1, 0, 0, false)

	if err := app.SetRoot(grid, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

func getNamespacesList() *tview.List {
	nsView := tview.NewList()
	nameSpaces := services.GetNamespaces()

	for i, ns := range nameSpaces {
		description := fmt.Sprintf("%d-%d", ns.LastSeen, ns.FirstSeen)
		nsView.AddItem(ns.Namespace, description, rune(65+i), nil)
	}

	return nsView
}

func updateApplicationsList(namespace string, appView *tview.List) {
	appView.Clear()
	applications := services.GetApplicationsV2(namespace)
	for i, app := range applications {
		description := fmt.Sprintf("%d-%d", app.LastSeen, app.FirstSeen)
		appView.AddItem(app.Namespace, description, rune(65+i), nil)
	}
}
