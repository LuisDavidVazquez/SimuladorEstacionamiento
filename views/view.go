package views

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

type View struct{}

func NewView() *View {
	return &View{}
}

func (v *View) RunView() {
	myApp := app.New()
	window := myApp.NewWindow("Simulador Estacionamiento")
	window.CenterOnScreen()
	window.SetFixedSize(true)
	window.Resize(fyne.NewSize(800, 600))
	//mainView := scenes.NewMainView(window)
	//mainView.Show()
	window.ShowAndRun()
}
