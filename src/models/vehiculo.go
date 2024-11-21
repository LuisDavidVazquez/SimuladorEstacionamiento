package models

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/storage"
)

type Vehiculo struct {
	id              int
	tiempoLim       time.Duration
	espacioAsignado int
	imagenEntrada   *canvas.Image
	imagenSalida    *canvas.Image
	imagenEspera    *canvas.Image
}

func NuevoVehiculo(id int) *Vehiculo {
	imagenEntrada := canvas.NewImageFromURI(storage.NewFileURI("./assets/carEntrada.png"))
	imagenSalida := canvas.NewImageFromURI(storage.NewFileURI("./assets/carSalida.png"))
	imagenEspera := canvas.NewImageFromURI(storage.NewFileURI("./assets/carEspera.png"))

	return &Vehiculo{
		id:              id,
		tiempoLim:       time.Duration(rand.Intn(50)+50) * time.Second,
		espacioAsignado: 0,
		imagenEntrada:   imagenEntrada,
		imagenSalida:    imagenSalida,
		imagenEspera:    imagenEspera,
	}
}

func (v *Vehiculo) Ingresar(p *Estacionamiento, contenedor *fyne.Container) {
	// Envía el ID del vehículo al canal de espacios en el estacionamiento.
	p.ObtenerEspacio() <- v.ObtenerID()
	// Adquiere el mutex
	p.ObtenerPuerta().Lock()

	// Obtiene el estado actual de ocupación de espacios en el estacionamiento.
	espacios := p.ObtenerArrayEspacios()
	const (
		columnasPorGrupo  = 10
		espacioHorizontal = 52
		espacioVertical   = 320
	)

	// Itera sobre los espacios para encontrar uno disponible.
	for i := 0; i < len(espacios); i++ {
		if !espacios[i] {
			// Marca el espacio como ocupado y asigna el número de espacio al vehículo.
			espacios[i] = true
			v.espacioAsignado = i

			// Calcula la posición en la interfaz gráfica según la fila y columna del espacio.
			fila := i / (columnasPorGrupo * 1)
			columna := i % (columnasPorGrupo * 1)
			if columna >= columnasPorGrupo {
				columna += 1
			}
			x := float32(50 + columna*espacioHorizontal)
			y := float32(50 + fila*espacioVertical)

			// Crea la imagen estacionado
			v.imagenEspera.Resize(fyne.NewSize(30, 50))
			v.imagenEspera.Move(fyne.NewPos(x, y))
			contenedor.Add(v.imagenEspera)
			contenedor.Refresh()
			break
		}
	}

	// Actualiza el estado de ocupación de espacios en el estacionamiento.
	p.EstablecerArrayEspacios(espacios)

	// Libera el mutex
	p.ObtenerPuerta().Unlock()
	contenedor.Refresh()
	fmt.Printf("Auto %d ocupó el lugar %d.\n", v.ObtenerID(), v.espacioAsignado)
	// Simula el tiempo que el vehículo permanece estacionado (5 segundos).
	time.Sleep(5 * time.Second)
}
func (v *Vehiculo) Salir(p *Estacionamiento, contenedor *fyne.Container) {
	// Recibe el espacio asignado al vehículo
	<-p.ObtenerEspacio()
	// Adquiere el mutex de la puerta
	p.ObtenerPuerta().Lock()
	espacios := p.ObtenerArrayEspacios()
	// Marca el espacio asignado como disponible.
	espacios[v.espacioAsignado] = false
	fmt.Printf("Auto %d salió. Espacio %d marcado como disponible.\n", v.ObtenerID(), v.espacioAsignado)
	p.EstablecerArrayEspacios(espacios)
	contenedor.Remove(v.imagenEspera)
	contenedor.Refresh()
	p.ObtenerPuerta().Unlock() //libera
	v.imagenSalida.Resize(fyne.NewSize(50, 30))
	v.imagenSalida.Move(fyne.NewPos(580, 160))
	contenedor.Add(v.imagenSalida)
	contenedor.Refresh()
	for i := 0; i < 100; i++ {
		v.imagenSalida.Move(fyne.NewPos(v.imagenSalida.Position().X-5, v.imagenSalida.Position().Y))
		time.Sleep(time.Millisecond * 15)
	}
	contenedor.Remove(v.imagenSalida)
	contenedor.Remove(v.imagenSalida)
	contenedor.Refresh()
}

func (v *Vehiculo) Iniciar(p *Estacionamiento, contenedor *fyne.Container, wg *sync.WaitGroup) {
	v.Avanzar(100, contenedor)
	v.Ingresar(p, contenedor)
	// Espera 10 segundos simulando la estancia del vehículo en el estacionamiento.
	time.Sleep(5 * time.Second)
	timer := time.NewTimer(5 * time.Second)

	// Selecciona el caso cuando el temporizador expire.
	select {
	case <-timer.C:
		contenedor.Remove(v.imagenEntrada)
		contenedor.Refresh()
		contenedor.Add(v.imagenSalida)
		contenedor.Refresh()
		p.ColaSalida(contenedor, v.imagenEntrada)
		v.Salir(p, contenedor)
		// Decrementa el contador de espera del WaitGroup(rutina termina).
		wg.Done()
	}
}

// Avanzar simula el avance del vehículo moviendo la imagen hacia abajo.
func (v *Vehiculo) Avanzar(pasos int, contenedor *fyne.Container) {
	for i := 0; i < pasos; i++ {
		v.imagenEntrada.Move(fyne.NewPos(v.imagenEntrada.Position().X+5, v.imagenEntrada.Position().Y))
		time.Sleep(time.Millisecond * 15)
	}
	contenedor.Remove(v.imagenEntrada)
	contenedor.Refresh()
}

func (v *Vehiculo) ObtenerID() int {
	return v.id
}

func (v *Vehiculo) ObtenerTiempoLim() time.Duration {
	return v.tiempoLim
}

func (v *Vehiculo) ObtenerImagenEntrada() *canvas.Image {
	return v.imagenEntrada
}
