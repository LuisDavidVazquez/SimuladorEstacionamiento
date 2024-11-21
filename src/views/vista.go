package views

import (
	"SIMULADOR/src/scenes"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

type VistaPrincipal struct{}

func NuevaVistaPrincipal() *VistaPrincipal {
	return &VistaPrincipal{}
}

func (v *VistaPrincipal) Ejecutar() {
	myApp := app.New()
	ventana := myApp.NewWindow("Estacionamiento")
	ventana.CenterOnScreen()
	ventana.SetFixedSize(true)
	ventana.Resize(fyne.NewSize(700, 400))
	escenaPrincipal := scenes.NuevaEscenaPrincipal(ventana)
	escenaPrincipal.Mostrar()
	ventana.ShowAndRun()
}
