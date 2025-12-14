package main

import "time"

type Fase string

const (
	FaseLlegada    Fase = "Llegada"
	FaseReparacion Fase = "Reparación"
	FaseLimpieza   Fase = "Limpieza"
	FaseEntrega    Fase = "Entrega"
)

type EstadoFase string

const (
	Entrando EstadoFase = "Entrando"
	Saliendo EstadoFase = "Saliendo"
)

type Categoria string

const (
	CategoriaA Categoria = "A"
	CategoriaB Categoria = "B"
	CategoriaC Categoria = "C"
)

type CocheSimulacion struct {
	ID         int
	Vehiculo   Vehiculo
	Categoria  Categoria
	Prioridad  int
	TiempoFase time.Duration
	FaseActual Fase
}

func ObtenerCategoriaPorTipo(tipo string) Categoria {
	switch tipo {
	case "mecánica":
		return CategoriaA
	case "eléctrica":
		return CategoriaB
	case "carrocería":
		return CategoriaC
	default:
		return CategoriaC
	}
}

func ObtenerPrioridadPorCategoria(cat Categoria) int {
	switch cat {
	case CategoriaA:
		return 3
	case CategoriaB:
		return 2
	case CategoriaC:
		return 1
	default:
		return 1
	}
}

func ObtenerTiempoPorCategoria(cat Categoria) time.Duration {
	switch cat {
	case CategoriaA:
		return 5 * time.Second
	case CategoriaB:
		return 3 * time.Second
	case CategoriaC:
		return 1 * time.Second
	default:
		return 1 * time.Second
	}
}

func NuevoCocheSimulacion(id int, vehiculo Vehiculo, tipoIncidencia string) CocheSimulacion {
	cat := ObtenerCategoriaPorTipo(tipoIncidencia)
	return CocheSimulacion{
		ID:         id,
		Vehiculo:   vehiculo,
		Categoria:  cat,
		Prioridad:  ObtenerPrioridadPorCategoria(cat),
		TiempoFase: ObtenerTiempoPorCategoria(cat),
		FaseActual: FaseLlegada,
	}
}
