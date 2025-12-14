package main

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

// GenerarCochesPorCategoria genera coches con distribución específica por categoría
func GenerarCochesPorCategoria(numA, numB, numC int) ([]CocheSimulacion, []string) {
	total := numA + numB + numC
	coches := make([]CocheSimulacion, total)
	tipos := make([]string, total)

	idx := 0

	// Generar coches Categoría A (mecánica)
	for i := 0; i < numA; i++ {
		tipoIncidencia := "mecánica"
		vehiculo := Vehiculo{
			Matricula:    fmt.Sprintf("%04dA%02d", idx+1, i+1),
			Marca:        "Test-A",
			Modelo:       "Categoría A",
			Incidencias:  tipoIncidencia,
			FechaEntrada: time.Now().Format("2006-01-02"),
		}
		coches[idx] = NuevoCocheSimulacion(idx+1, vehiculo, tipoIncidencia)
		tipos[idx] = tipoIncidencia
		idx++
	}

	// Generar coches Categoría B (eléctrica)
	for i := 0; i < numB; i++ {
		tipoIncidencia := "eléctrica"
		vehiculo := Vehiculo{
			Matricula:    fmt.Sprintf("%04dB%02d", idx+1, i+1),
			Marca:        "Test-B",
			Modelo:       "Categoría B",
			Incidencias:  tipoIncidencia,
			FechaEntrada: time.Now().Format("2006-01-02"),
		}
		coches[idx] = NuevoCocheSimulacion(idx+1, vehiculo, tipoIncidencia)
		tipos[idx] = tipoIncidencia
		idx++
	}

	// Generar coches Categoría C (carrocería)
	for i := 0; i < numC; i++ {
		tipoIncidencia := "carrocería"
		vehiculo := Vehiculo{
			Matricula:    fmt.Sprintf("%04dC%02d", idx+1, i+1),
			Marca:        "Test-C",
			Modelo:       "Categoría C",
			Incidencias:  tipoIncidencia,
			FechaEntrada: time.Now().Format("2006-01-02"),
		}
		coches[idx] = NuevoCocheSimulacion(idx+1, vehiculo, tipoIncidencia)
		tipos[idx] = tipoIncidencia
		idx++
	}

	return coches, tipos
}

// ejecutarSimulacionConDistribucion ejecuta una simulación con distribución específica
func ejecutarSimulacionConDistribucion(t *testing.T, nombre string, numA, numB, numC, numPlazas, numMecanicos int, usarRWMutex bool) time.Duration {
	total := numA + numB + numC

	implementacion := "WaitGroup"
	if usarRWMutex {
		implementacion = "RWMutex"
	}

	t.Logf("\n========================================")
	t.Logf("Test: %s", nombre)
	t.Logf("Implementación: %s", implementacion)
	t.Logf("Distribución: A=%d, B=%d, C=%d (Total: %d)", numA, numB, numC, total)
	t.Logf("Recursos: Plazas=%d, Mecánicos=%d", numPlazas, numMecanicos)
	t.Logf("========================================\n")

	inicio := time.Now()

	if usarRWMutex {
		LogInicioSimulacion("RWMutex", total, numPlazas, numMecanicos)
		simulador := NuevoSimuladorRWMutex(numPlazas, numMecanicos)
		coches, tiposIncidencias := GenerarCochesPorCategoria(numA, numB, numC)

		var wg sync.WaitGroup
		for i := range coches {
			wg.Add(1)
			go simulador.Fase1Llegada(&coches[i], tiposIncidencias[i], &wg)
			time.Sleep(100 * time.Millisecond)
		}
		wg.Wait()
		LogFinSimulacion(simulador.tiempoInicio, total)
	} else {
		LogInicioSimulacion("WaitGroup", total, numPlazas, numMecanicos)
		simulador := NuevoSimuladorWaitGroup(numPlazas, numMecanicos)
		coches, tiposIncidencias := GenerarCochesPorCategoria(numA, numB, numC)

		var wg sync.WaitGroup
		for i := range coches {
			wg.Add(1)
			go simulador.Fase1LlegadaWG(&coches[i], tiposIncidencias[i], &wg)
			time.Sleep(100 * time.Millisecond)
		}
		wg.Wait()
		LogFinSimulacion(simulador.tiempoInicio, total)
	}

	duracion := time.Since(inicio)
	t.Logf("✓ Simulación completada en: %v\n", duracion)

	return duracion
}

// TestComparativa1_Balanceado - Test 1: A=10, B=10, C=10
func TestComparativa1_Balanceado(t *testing.T) {
	t.Log("\n╔════════════════════════════════════════════════════════╗")
	t.Log("║  TEST 1: DISTRIBUCIÓN BALANCEADA (A=10, B=10, C=10)   ║")
	t.Log("╚════════════════════════════════════════════════════════╝")

	numPlazas := 10
	numMecanicos := 5

	// Ejecutar con RWMutex
	duracionRWMutex := ejecutarSimulacionConDistribucion(t, "Comparativa 1 - Balanceado", 10, 10, 10, numPlazas, numMecanicos, true)

	// Pequeña pausa entre simulaciones
	time.Sleep(2 * time.Second)

	// Ejecutar con WaitGroup
	duracionWaitGroup := ejecutarSimulacionConDistribucion(t, "Comparativa 1 - Balanceado", 10, 10, 10, numPlazas, numMecanicos, false)

	// Resumen
	t.Log("\n┌─────────────────────────────────────────────────────┐")
	t.Log("│           RESUMEN TEST 1 - BALANCEADO               │")
	t.Log("├─────────────────────────────────────────────────────┤")
	t.Logf("│ RWMutex:    %v", duracionRWMutex)
	t.Logf("│ WaitGroup:  %v", duracionWaitGroup)
	t.Logf("│ Diferencia: %v", abs(duracionRWMutex-duracionWaitGroup))
	t.Log("└─────────────────────────────────────────────────────┘\n")
}

// TestComparativa2_AHeavy - Test 2: A=20, B=5, C=5
func TestComparativa2_AHeavy(t *testing.T) {
	t.Log("\n╔════════════════════════════════════════════════════════╗")
	t.Log("║  TEST 2: DISTRIBUCIÓN A-HEAVY (A=20, B=5, C=5)        ║")
	t.Log("╚════════════════════════════════════════════════════════╝")

	numPlazas := 10
	numMecanicos := 5

	// Ejecutar con RWMutex
	duracionRWMutex := ejecutarSimulacionConDistribucion(t, "Comparativa 2 - A-Heavy", 20, 5, 5, numPlazas, numMecanicos, true)

	// Pequeña pausa entre simulaciones
	time.Sleep(2 * time.Second)

	// Ejecutar con WaitGroup
	duracionWaitGroup := ejecutarSimulacionConDistribucion(t, "Comparativa 2 - A-Heavy", 20, 5, 5, numPlazas, numMecanicos, false)

	// Resumen
	t.Log("\n┌─────────────────────────────────────────────────────┐")
	t.Log("│           RESUMEN TEST 2 - A-HEAVY                  │")
	t.Log("├─────────────────────────────────────────────────────┤")
	t.Logf("│ RWMutex:    %v", duracionRWMutex)
	t.Logf("│ WaitGroup:  %v", duracionWaitGroup)
	t.Logf("│ Diferencia: %v", abs(duracionRWMutex-duracionWaitGroup))
	t.Log("└─────────────────────────────────────────────────────┘\n")
}

// TestComparativa3_CHeavy - Test 3: A=5, B=5, C=20
func TestComparativa3_CHeavy(t *testing.T) {
	t.Log("\n╔════════════════════════════════════════════════════════╗")
	t.Log("║  TEST 3: DISTRIBUCIÓN C-HEAVY (A=5, B=5, C=20)        ║")
	t.Log("╚════════════════════════════════════════════════════════╝")

	numPlazas := 10
	numMecanicos := 5

	// Ejecutar con RWMutex
	duracionRWMutex := ejecutarSimulacionConDistribucion(t, "Comparativa 3 - C-Heavy", 5, 5, 20, numPlazas, numMecanicos, true)

	// Pequeña pausa entre simulaciones
	time.Sleep(2 * time.Second)

	// Ejecutar con WaitGroup
	duracionWaitGroup := ejecutarSimulacionConDistribucion(t, "Comparativa 3 - C-Heavy", 5, 5, 20, numPlazas, numMecanicos, false)

	// Resumen
	t.Log("\n┌─────────────────────────────────────────────────────┐")
	t.Log("│           RESUMEN TEST 3 - C-HEAVY                  │")
	t.Log("├─────────────────────────────────────────────────────┤")
	t.Logf("│ RWMutex:    %v", duracionRWMutex)
	t.Logf("│ WaitGroup:  %v", duracionWaitGroup)
	t.Logf("│ Diferencia: %v", abs(duracionRWMutex-duracionWaitGroup))
	t.Log("└─────────────────────────────────────────────────────┘\n")
}

// TestComparativaCompleta - Ejecuta todos los tests y muestra tabla comparativa final
func TestComparativaCompleta(t *testing.T) {
	t.Log("\n╔══════════════════════════════════════════════════════════════╗")
	t.Log("║        COMPARATIVA COMPLETA - TODOS LOS ESCENARIOS           ║")
	t.Log("╚══════════════════════════════════════════════════════════════╝")

	type Resultado struct {
		Nombre              string
		DistA, DistB, DistC int
		DuracionRWMutex     time.Duration
		DuracionWaitGroup   time.Duration
	}

	resultados := []Resultado{}
	numPlazas := 10
	numMecanicos := 5

	// Test 1: Balanceado
	t.Log("\n>>> Ejecutando Test 1: Balanceado (10/10/10)")
	durRW1 := ejecutarSimulacionConDistribucion(t, "Test 1", 10, 10, 10, numPlazas, numMecanicos, true)
	time.Sleep(2 * time.Second)
	durWG1 := ejecutarSimulacionConDistribucion(t, "Test 1", 10, 10, 10, numPlazas, numMecanicos, false)
	resultados = append(resultados, Resultado{"Test 1: Balanceado", 10, 10, 10, durRW1, durWG1})
	time.Sleep(2 * time.Second)

	// Test 2: A-Heavy
	t.Log("\n>>> Ejecutando Test 2: A-Heavy (20/5/5)")
	durRW2 := ejecutarSimulacionConDistribucion(t, "Test 2", 20, 5, 5, numPlazas, numMecanicos, true)
	time.Sleep(2 * time.Second)
	durWG2 := ejecutarSimulacionConDistribucion(t, "Test 2", 20, 5, 5, numPlazas, numMecanicos, false)
	resultados = append(resultados, Resultado{"Test 2: A-Heavy", 20, 5, 5, durRW2, durWG2})
	time.Sleep(2 * time.Second)

	// Test 3: C-Heavy
	t.Log("\n>>> Ejecutando Test 3: C-Heavy (5/5/20)")
	durRW3 := ejecutarSimulacionConDistribucion(t, "Test 3", 5, 5, 20, numPlazas, numMecanicos, true)
	time.Sleep(2 * time.Second)
	durWG3 := ejecutarSimulacionConDistribucion(t, "Test 3", 5, 5, 20, numPlazas, numMecanicos, false)
	resultados = append(resultados, Resultado{"Test 3: C-Heavy", 5, 5, 20, durRW3, durWG3})

	// Tabla comparativa final
	t.Log("\n╔══════════════════════════════════════════════════════════════════════════════╗")
	t.Log("║                    TABLA COMPARATIVA FINAL                                   ║")
	t.Log("╠══════════════════════════════════════════════════════════════════════════════╣")
	t.Log("║ Test              │ Dist.    │ RWMutex      │ WaitGroup    │ Diferencia     ║")
	t.Log("╠═══════════════════╪══════════╪══════════════╪══════════════╪════════════════╣")

	for _, r := range resultados {
		diff := abs(r.DuracionRWMutex - r.DuracionWaitGroup)
		t.Logf("║ %-17s │ %2d/%2d/%2d │ %-12v │ %-12v │ %-14v ║",
			r.Nombre,
			r.DistA, r.DistB, r.DistC,
			r.DuracionRWMutex.Round(time.Millisecond),
			r.DuracionWaitGroup.Round(time.Millisecond),
			diff.Round(time.Millisecond))
	}

	t.Log("╚══════════════════════════════════════════════════════════════════════════════╝")

	// Análisis
	t.Log("\n╔══════════════════════════════════════════════════════════════════════════════╗")
	t.Log("║                              ANÁLISIS                                        ║")
	t.Log("╚══════════════════════════════════════════════════════════════════════════════╝")
	t.Log("│ Categoría A (mecánica):   Prioridad 3, Tiempo 5s por fase                   │")
	t.Log("│ Categoría B (eléctrica):  Prioridad 2, Tiempo 3s por fase                   │")
	t.Log("│ Categoría C (carrocería): Prioridad 1, Tiempo 1s por fase                   │")
	t.Log("└──────────────────────────────────────────────────────────────────────────────┘")
}

// Función auxiliar para calcular valor absoluto de duración
func abs(d time.Duration) time.Duration {
	if d < 0 {
		return -d
	}
	return d
}

// Benchmarks para medición de rendimiento

func BenchmarkBalanceado_RWMutex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		simulador := NuevoSimuladorRWMutex(10, 5)
		coches, tiposIncidencias := GenerarCochesPorCategoria(10, 10, 10)
		var wg sync.WaitGroup
		for j := range coches {
			wg.Add(1)
			go simulador.Fase1Llegada(&coches[j], tiposIncidencias[j], &wg)
		}
		wg.Wait()
	}
}

func BenchmarkBalanceado_WaitGroup(b *testing.B) {
	for i := 0; i < b.N; i++ {
		simulador := NuevoSimuladorWaitGroup(10, 5)
		coches, tiposIncidencias := GenerarCochesPorCategoria(10, 10, 10)
		var wg sync.WaitGroup
		for j := range coches {
			wg.Add(1)
			go simulador.Fase1LlegadaWG(&coches[j], tiposIncidencias[j], &wg)
		}
		wg.Wait()
	}
}

func BenchmarkAHeavy_RWMutex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		simulador := NuevoSimuladorRWMutex(10, 5)
		coches, tiposIncidencias := GenerarCochesPorCategoria(20, 5, 5)
		var wg sync.WaitGroup
		for j := range coches {
			wg.Add(1)
			go simulador.Fase1Llegada(&coches[j], tiposIncidencias[j], &wg)
		}
		wg.Wait()
	}
}

func BenchmarkAHeavy_WaitGroup(b *testing.B) {
	for i := 0; i < b.N; i++ {
		simulador := NuevoSimuladorWaitGroup(10, 5)
		coches, tiposIncidencias := GenerarCochesPorCategoria(20, 5, 5)
		var wg sync.WaitGroup
		for j := range coches {
			wg.Add(1)
			go simulador.Fase1LlegadaWG(&coches[j], tiposIncidencias[j], &wg)
		}
		wg.Wait()
	}
}

func BenchmarkCHeavy_RWMutex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		simulador := NuevoSimuladorRWMutex(10, 5)
		coches, tiposIncidencias := GenerarCochesPorCategoria(5, 5, 20)
		var wg sync.WaitGroup
		for j := range coches {
			wg.Add(1)
			go simulador.Fase1Llegada(&coches[j], tiposIncidencias[j], &wg)
		}
		wg.Wait()
	}
}

func BenchmarkCHeavy_WaitGroup(b *testing.B) {
	for i := 0; i < b.N; i++ {
		simulador := NuevoSimuladorWaitGroup(10, 5)
		coches, tiposIncidencias := GenerarCochesPorCategoria(5, 5, 20)
		var wg sync.WaitGroup
		for j := range coches {
			wg.Add(1)
			go simulador.Fase1LlegadaWG(&coches[j], tiposIncidencias[j], &wg)
		}
		wg.Wait()
	}
}
