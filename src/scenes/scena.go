package scenes

import (
	"SIMULADOR/src/models"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"gonum.org/v1/gonum/stat/distuv"
)

type EscenaPrincipal struct {
	ventana fyne.Window
}

func NuevaEscenaPrincipal(ventana fyne.Window) *EscenaPrincipal {
	return &EscenaPrincipal{
		ventana: ventana,
	}
}

func (s *EscenaPrincipal) Mostrar() {
	fondoEstacionamiento := canvas.NewImageFromFile("assets/fondo.jpg")
	fondoEstacionamiento.Resize(fyne.NewSize(690, 400))
	fondoEstacionamiento.Move(fyne.NewPos(0, 0))

	contenedor := container.NewWithoutLayout()
	contenedor.Add(fondoEstacionamiento)

	go s.Ejecutar()
	s.ventana.SetContent(contenedor)
}

func (s *EscenaPrincipal) Ejecutar() {
	p := models.NuevoEstacionamiento(make(chan int, 20), &sync.Mutex{})
	contenedor := s.ventana.Content().(*fyne.Container)
	var wg sync.WaitGroup //espera a que las rutinas terminen su ejecusion

	// Inicia 100 goroutines para simular vehículos ingresando al estacionamiento.
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(id int) {
			// Crea un nuevo vehículo con un ID único.
			vehiculo := models.NuevoVehiculo(id)
			imagen := vehiculo.ObtenerImagenEntrada()
			imagen.Resize(fyne.NewSize(50, 30))
			imagen.Move(fyne.NewPos(80, 210))
			contenedor.Add(imagen)
			contenedor.Refresh()
			vehiculo.Iniciar(p, contenedor, &wg)
			time.Sleep(time.Millisecond * 100)
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
