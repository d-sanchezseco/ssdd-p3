package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type SimuladorRWMutex struct {
	numPlazas       int
	numMecanicos    int
	plazasLibres    int
	mecanicosLibres int

	mutexPlazas    sync.RWMutex
	mutexMecanicos sync.RWMutex

	colaPlazas    *GestorColaPrioridadSimulacion
	colaMecanicos *GestorColaPrioridadSimulacion
	colaLimpieza  *GestorColaPrioridadSimulacion
	colaEntrega   *GestorColaPrioridadSimulacion

	tiempoInicio time.Time
}

func NuevoSimuladorRWMutex(numPlazas, numMecanicos int) *SimuladorRWMutex {
	return &SimuladorRWMutex{
		numPlazas:       numPlazas,
		numMecanicos:    numMecanicos,
		plazasLibres:    numPlazas,
		mecanicosLibres: numMecanicos,
		colaPlazas:      NuevoGestorColaPrioridadSimulacion(),
		colaMecanicos:   NuevoGestorColaPrioridadSimulacion(),
		colaLimpieza:    NuevoGestorColaPrioridadSimulacion(),
		colaEntrega:     NuevoGestorColaPrioridadSimulacion(),
		tiempoInicio:    time.Now(),
	}
}

func (s *SimuladorRWMutex) Fase1Llegada(coche *CocheSimulacion, tipoIncidencia string, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		s.mutexPlazas.Lock()
		if s.plazasLibres > 0 {
			s.plazasLibres--
			s.mutexPlazas.Unlock()
			break
		}
		s.mutexPlazas.Unlock()
		time.Sleep(100 * time.Millisecond)
	}

	coche.FaseActual = FaseLlegada
	LogFaseSimulacion(s.tiempoInicio, *coche, tipoIncidencia, Entrando)

	time.Sleep(coche.TiempoFase)

	LogFaseSimulacion(s.tiempoInicio, *coche, tipoIncidencia, Saliendo)

	wg.Add(1)
	go s.Fase2Reparacion(coche, tipoIncidencia, wg)
}

func (s *SimuladorRWMutex) Fase2Reparacion(coche *CocheSimulacion, tipoIncidencia string, wg *sync.WaitGroup) {
	defer wg.Done()

	s.colaMecanicos.Agregar(*coche)

	for {
		s.mutexMecanicos.Lock()
		if s.mecanicosLibres > 0 && !s.colaMecanicos.EstaVacia() {
			cocheEnCola, ok := s.colaMecanicos.Extraer()
			if ok && cocheEnCola.ID == coche.ID {
				s.mecanicosLibres--
				s.mutexMecanicos.Unlock()
				break
			} else if ok {
				s.colaMecanicos.Agregar(cocheEnCola)
			}
		}
		s.mutexMecanicos.Unlock()
		time.Sleep(100 * time.Millisecond)
	}

	coche.FaseActual = FaseReparacion
	LogFaseSimulacion(s.tiempoInicio, *coche, tipoIncidencia, Entrando)

	time.Sleep(coche.TiempoFase)

	LogFaseSimulacion(s.tiempoInicio, *coche, tipoIncidencia, Saliendo)

	s.mutexMecanicos.Lock()
	s.mecanicosLibres++
	s.mutexMecanicos.Unlock()

	wg.Add(1)
	go s.Fase3Limpieza(coche, tipoIncidencia, wg)
}

func (s *SimuladorRWMutex) Fase3Limpieza(coche *CocheSimulacion, tipoIncidencia string, wg *sync.WaitGroup) {
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
		time.Sleep(100 * time.Millisecond)
	}

	coche.FaseActual = FaseLimpieza
	LogFaseSimulacion(s.tiempoInicio, *coche, tipoIncidencia, Entrando)

	time.Sleep(coche.TiempoFase)

	LogFaseSimulacion(s.tiempoInicio, *coche, tipoIncidencia, Saliendo)

	wg.Add(1)
	go s.Fase4Entrega(coche, tipoIncidencia, wg)
}

func (s *SimuladorRWMutex) Fase4Entrega(coche *CocheSimulacion, tipoIncidencia string, wg *sync.WaitGroup) {
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
		time.Sleep(100 * time.Millisecond)
	}

	coche.FaseActual = FaseEntrega
	LogFaseSimulacion(s.tiempoInicio, *coche, tipoIncidencia, Entrando)

	time.Sleep(coche.TiempoFase)

	LogFaseSimulacion(s.tiempoInicio, *coche, tipoIncidencia, Saliendo)

	s.mutexPlazas.Lock()
	s.plazasLibres++
	s.mutexPlazas.Unlock()
}

func EjecutarSimulacionRWMutex(numCoches, numPlazas, numMecanicos int) {
	LogInicioSimulacion("RWMutex", numCoches, numPlazas, numMecanicos)

	simulador := NuevoSimuladorRWMutex(numPlazas, numMecanicos)

	coches, tiposIncidencias := GenerarCochesSimulacion(numCoches)

	var wg sync.WaitGroup

	for i := range coches {
		wg.Add(1)
		go simulador.Fase1Llegada(&coches[i], tiposIncidencias[i], &wg)

		time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
	}

	wg.Wait()

	LogFinSimulacion(simulador.tiempoInicio, numCoches)
}

func GenerarCochesSimulacion(n int) ([]CocheSimulacion, []string) {
	tiposIncidencias := []string{"mecánica", "eléctrica", "carrocería"}
	coches := make([]CocheSimulacion, n)
	tipos := make([]string, n)

	for i := 0; i < n; i++ {
		tipoAleatorio := tiposIncidencias[rand.Intn(len(tiposIncidencias))]
		tipos[i] = tipoAleatorio

		vehiculo := Vehiculo{
			Matricula:    fmt.Sprintf("%04dSIM", i+1),
			Marca:        "Simulación",
			Modelo:       "Modelo " + tipoAleatorio,
			Incidencias:  tipoAleatorio,
			FechaEntrada: time.Now().Format("2006-01-02"),
		}

		coches[i] = NuevoCocheSimulacion(i+1, vehiculo, tipoAleatorio)
	}

	return coches, tipos
}
