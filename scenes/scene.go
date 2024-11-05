package scenes

import (
	"SIMULADOR/models"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"gonum.org/v1/gonum/stat/distuv"
)

type MainView struct {
	window fyne.Window
}

func NewMainView(window fyne.Window) *MainView {
	return &MainView{
		window: window,
	}
}

func (s *MainView) Show() {
	background := canvas.NewImageFromFile("assets/fon2.jpg")
	background.Resize(fyne.NewSize(690, 400))
	background.Move(fyne.NewPos(0, 0))

	container := container.NewWithoutLayout()
	container.Add(background)

	go s.Run()

	vbox := fyne.NewContainerWithLayout(layout.NewVBoxLayout(),
		layout.NewSpacer(),
		layout.NewSpacer(),
		layout.NewSpacer(),
		layout.NewSpacer(),
	)

	// Agrega el contenedor vertical al contenedor principal.
	container.Add(vbox)
	s.window.SetContent(container)
}

func (s *MainView) Run() {
	p := models.NuevoEstacionamiento(make(chan int, 20), &sync.Mutex{})
	container := s.window.Content().(*fyne.Container)
	var wg sync.WaitGroup //espera a que las rutinas terminen su ejecusion

	// Inicia 100 goroutines para simular vehículos ingresando al estacionamiento.
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(id int) {
			// Crea un nuevo vehículo con un ID único.
			vehicle := models.NuevoVehiculo(id)
			image := vehicle.ObtenerImagenEntrada()
			image.Resize(fyne.NewSize(30, 50))
			image.Move(fyne.NewPos(40, -10))
			container.Add(image)
			container.Refresh()
			vehicle.Iniciar(p, container, &wg)
			time.Sleep(time.Millisecond * 200)
		}(i)
		poisson := generarPoisson(2)
		time.Sleep(time.Second * time.Duration(poisson))
	}
	// Espera a que todas las goroutines terminen y no pase mas alla
	wg.Wait()
}

// generarPoisson genera un número aleatorio según la distribución de Poisson con un parámetro lambda dado.
func generarPoisson(lambda float64) float64 {
	poisson := distuv.Poisson{Lambda: lambda, Src: nil}
	return poisson.Rand()
}
