package main

import (
	"container/heap"
	"sync"
)

type ItemColaSimulacion struct {
	coche     CocheSimulacion
	prioridad int
	index     int
}

type ColaPrioridadSimulacion []*ItemColaSimulacion

func (cp ColaPrioridadSimulacion) Len() int { return len(cp) }

func (cp ColaPrioridadSimulacion) Less(i, j int) bool {

	if cp[i].prioridad == cp[j].prioridad {
		return cp[i].index < cp[j].index
	}

	return cp[i].prioridad > cp[j].prioridad
}

func (cp ColaPrioridadSimulacion) Swap(i, j int) {
	cp[i], cp[j] = cp[j], cp[i]
	cp[i].index = i
	cp[j].index = j
}

func (cp *ColaPrioridadSimulacion) Push(x interface{}) {
	n := len(*cp)
	item := x.(*ItemColaSimulacion)
	item.index = n
	*cp = append(*cp, item)
}

func (cp *ColaPrioridadSimulacion) Pop() interface{} {
	old := *cp
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.index = -1
	*cp = old[0 : n-1]
	return item
}

type GestorColaPrioridadSimulacion struct {
	cola   *ColaPrioridadSimulacion
	mu     sync.Mutex
	indice int
}

func NuevoGestorColaPrioridadSimulacion() *GestorColaPrioridadSimulacion {
	cp := &ColaPrioridadSimulacion{}
	heap.Init(cp)
	return &GestorColaPrioridadSimulacion{
		cola:   cp,
		indice: 0,
	}
}

func (g *GestorColaPrioridadSimulacion) Agregar(coche CocheSimulacion) {
	g.mu.Lock()
	defer g.mu.Unlock()

	item := &ItemColaSimulacion{
		coche:     coche,
		prioridad: coche.Prioridad,
		index:     g.indice,
	}
	g.indice++

	heap.Push(g.cola, item)
}

func (g *GestorColaPrioridadSimulacion) Extraer() (CocheSimulacion, bool) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.cola.Len() == 0 {
		return CocheSimulacion{}, false
	}

	item := heap.Pop(g.cola).(*ItemColaSimulacion)
	return item.coche, true
}

func (g *GestorColaPrioridadSimulacion) Tamaño() int {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.cola.Len()
}

func (g *GestorColaPrioridadSimulacion) EstaVacia() bool {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.cola.Len() == 0
}
