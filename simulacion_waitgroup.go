package main

import (
	"math/rand"
	"sync"
	"time"
)

type SimuladorWaitGroup struct {
	numPlazas    int
	numMecanicos int

	canalPlazas    chan struct{}
	canalMecanicos chan struct{}

	colaPlazas    *GestorColaPrioridadSimulacion
	colaMecanicos *GestorColaPrioridadSimulacion
	colaLimpieza  *GestorColaPrioridadSimulacion
	colaEntrega   *GestorColaPrioridadSimulacion

	tiempoInicio time.Time
}

func NuevoSimuladorWaitGroup(numPlazas, numMecanicos int) *SimuladorWaitGroup {
	canalPlazas := make(chan struct{}, numPlazas)
	canalMecanicos := make(chan struct{}, numMecanicos)

	for i := 0; i < numPlazas; i++ {
		canalPlazas <- struct{}{}
	}
	for i := 0; i < numMecanicos; i++ {
		canalMecanicos <- struct{}{}
	}

	return &SimuladorWaitGroup{
		numPlazas:      numPlazas,
		numMecanicos:   numMecanicos,
		canalPlazas:    canalPlazas,
		canalMecanicos: canalMecanicos,
		colaPlazas:     NuevoGestorColaPrioridadSimulacion(),
		colaMecanicos:  NuevoGestorColaPrioridadSimulacion(),
		colaLimpieza:   NuevoGestorColaPrioridadSimulacion(),
		colaEntrega:    NuevoGestorColaPrioridadSimulacion(),
		tiempoInicio:   time.Now(),
	}
}

func (s *SimuladorWaitGroup) Fase1LlegadaWG(coche *CocheSimulacion, tipoIncidencia string, wg *sync.WaitGroup) {
	defer wg.Done()

	s.colaPlazas.Agregar(*coche)

	for {
		if !s.colaPlazas.EstaVacia() {
			cocheEnCola, ok := s.colaPlazas.Extraer()
			if ok && cocheEnCola.ID == coche.ID {
				select {
				case <-s.canalPlazas:
					goto PlazaObtenida
				default:
					s.colaPlazas.Agregar(cocheEnCola)
				}
			} else if ok {
				s.colaPlazas.Agregar(cocheEnCola)
			}
		}
		time.Sleep(50 * time.Millisecond)
	}

PlazaObtenida:
	coche.FaseActual = FaseLlegada
	LogFaseSimulacion(s.tiempoInicio, *coche, tipoIncidencia, Entrando)

	time.Sleep(coche.TiempoFase)

	LogFaseSimulacion(s.tiempoInicio, *coche, tipoIncidencia, Saliendo)

	wg.Add(1)
	go s.Fase2ReparacionWG(coche, tipoIncidencia, wg)
}

func (s *SimuladorWaitGroup) Fase2ReparacionWG(coche *CocheSimulacion, tipoIncidencia string, wg *sync.WaitGroup) {
	defer wg.Done()

	s.colaMecanicos.Agregar(*coche)

	for {
		if !s.colaMecanicos.EstaVacia() {
			cocheEnCola, ok := s.colaMecanicos.Extraer()
			if ok && cocheEnCola.ID == coche.ID {
				select {
				case <-s.canalMecanicos:
					goto MecanicoObtenido
				default:
					s.colaMecanicos.Agregar(cocheEnCola)
				}
			} else if ok {
				s.colaMecanicos.Agregar(cocheEnCola)
			}
		}
		time.Sleep(50 * time.Millisecond)
	}

MecanicoObtenido:
	coche.FaseActual = FaseReparacion
	LogFaseSimulacion(s.tiempoInicio, *coche, tipoIncidencia, Entrando)

	time.Sleep(coche.TiempoFase)

	LogFaseSimulacion(s.tiempoInicio, *coche, tipoIncidencia, Saliendo)

	s.canalMecanicos <- struct{}{}

	wg.Add(1)
	go s.Fase3LimpiezaWG(coche, tipoIncidencia, wg)
}

func (s *SimuladorWaitGroup) Fase3LimpiezaWG(coche *CocheSimulacion, tipoIncidencia string, wg *sync.WaitGroup) {
	defer wg.Done()

	s.colaLimpieza.Agregar(*coche)

	for {
		if !s.colaLimpieza.EstaVacia() {
			cocheEnCola, ok := s.colaLimpieza.Extraer()
			if ok && cocheEnCola.ID == coche.ID {
				break
			} else if ok {
				s.colaLimpieza.Agregar(cocheEnCola)
			}
		}
		time.Sleep(50 * time.Millisecond)
	}

	coche.FaseActual = FaseLimpieza
	LogFaseSimulacion(s.tiempoInicio, *coche, tipoIncidencia, Entrando)

	time.Sleep(coche.TiempoFase)

	LogFaseSimulacion(s.tiempoInicio, *coche, tipoIncidencia, Saliendo)

	wg.Add(1)
	go s.Fase4EntregaWG(coche, tipoIncidencia, wg)
}

func (s *SimuladorWaitGroup) Fase4EntregaWG(coche *CocheSimulacion, tipoIncidencia string, wg *sync.WaitGroup) {
	defer wg.Done()

	s.colaEntrega.Agregar(*coche)

	for {
		if !s.colaEntrega.EstaVacia() {
			cocheEnCola, ok := s.colaEntrega.Extraer()
			if ok && cocheEnCola.ID == coche.ID {
				break
			} else if ok {
				s.colaEntrega.Agregar(cocheEnCola)
			}
		}
		time.Sleep(50 * time.Millisecond)
	}

	coche.FaseActual = FaseEntrega
	LogFaseSimulacion(s.tiempoInicio, *coche, tipoIncidencia, Entrando)

	time.Sleep(coche.TiempoFase)

	LogFaseSimulacion(s.tiempoInicio, *coche, tipoIncidencia, Saliendo)

	s.canalPlazas <- struct{}{}
}

func EjecutarSimulacionWaitGroup(numCoches, numPlazas, numMecanicos int) {
	LogInicioSimulacion("WaitGroup", numCoches, numPlazas, numMecanicos)

	simulador := NuevoSimuladorWaitGroup(numPlazas, numMecanicos)

	coches, tiposIncidencias := GenerarCochesSimulacion(numCoches)

	var wg sync.WaitGroup

	for i := range coches {
		wg.Add(1)
		go simulador.Fase1LlegadaWG(&coches[i], tiposIncidencias[i], &wg)

		time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
	}

	wg.Wait()

	LogFinSimulacion(simulador.tiempoInicio, numCoches)
}
