package app

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type App struct {
	*tview.Application
}

func NewApp() *App {
	app := &App{
		Application: tview.NewApplication(),
	}

	app.setupBindings()

	root := tview.NewFlex()
	app.setupLayout(root)

	app.Application.SetRoot(root, true)

	return app
}

func (a *App) Run() error {
	return a.Application.Run()
}

func (a *App) setupBindings() {
	a.Application.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlC:
			a.Stop()
			return nil
		}
		return event
	})
}

func (a *App) setupLayout(root *tview.Flex) {
}
