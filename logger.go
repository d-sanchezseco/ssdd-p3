package main

import (
	"fmt"
	"time"
)

func LogFaseSimulacion(tiempoInicio time.Time, coche CocheSimulacion, tipoIncidencia string, estado EstadoFase) {
	tiempoTranscurrido := time.Since(tiempoInicio).Seconds()

	fmt.Printf("Tiempo %.2f Coche %d Incidencia %s Fase %s Estado %s\n",
		tiempoTranscurrido,
		coche.ID,
		tipoIncidencia,
		coche.FaseActual,
		estado,
	)
}

func LogInicioSimulacion(metodo string, numCoches, numPlazas, numMecanicos int) {
	fmt.Println("=====================================")
	fmt.Printf("  SIMULACIÓN DE TALLER - %s\n", metodo)
	fmt.Println("=====================================")
	fmt.Printf("Coches a reparar: %d\n", numCoches)
	fmt.Printf("Plazas disponibles: %d\n", numPlazas)
	fmt.Printf("Mecánicos disponibles: %d\n", numMecanicos)
	fmt.Println("=====================================")
	fmt.Println()
}

func LogFinSimulacion(tiempoInicio time.Time, numCoches int) {
	tiempoTotal := time.Since(tiempoInicio).Seconds()
	fmt.Println()
	fmt.Println("=====================================")
	fmt.Printf("  SIMULACIÓN COMPLETADA\n")
	fmt.Println("=====================================")
	fmt.Printf("Tiempo total: %.2f segundos\n", tiempoTotal)
	fmt.Printf("Coches reparados: %d\n", numCoches)
	fmt.Println("=====================================")
}

func LogEstadisticasSimulacion(cochesPorCategoria map[Categoria]int, tiempoPromedio map[Categoria]float64) {
	fmt.Println()
	fmt.Println("===== ESTADÍSTICAS =====")
	fmt.Println("Coches por categoría:")
	for cat, count := range cochesPorCategoria {
		fmt.Printf("  Categoría %s: %d coches\n", cat, count)
	}
	fmt.Println()
	fmt.Println("Tiempo promedio por categoría:")
	for cat, tiempo := range tiempoPromedio {
		fmt.Printf("  Categoría %s: %.2f segundos\n", cat, tiempo)
	}
	fmt.Println("========================")
}
