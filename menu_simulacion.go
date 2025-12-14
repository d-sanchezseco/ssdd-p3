package main

import "fmt"

func menuSimulacionFases() {
	clearScreen()
	var option int
	var numCoches, numPlazas, numMecanicos, numSims int

	fmt.Println("===== SIMULACIÓN DE TALLER - 4 FASES =====")
	fmt.Println()
	fmt.Println("Esta simulación implementa el sistema de reparación con:")
	fmt.Println("- 4 fases secuenciales: Llegada, Reparación, Limpieza, Entrega")
	fmt.Println("- 3 categorías de prioridad:")
	fmt.Println("  * Categoría A (Mecánica): Prioridad Alta - 5s por fase")
	fmt.Println("  * Categoría B (Eléctrica): Prioridad Media - 3s por fase")
	fmt.Println("  * Categoría C (Carrocería): Prioridad Baja - 1s por fase")
	fmt.Println()

	ingresarDatos("número de coches a reparar", &numCoches)
	ingresarDatos("número de plazas disponibles", &numPlazas)
	ingresarDatos("número de mecánicos disponibles", &numMecanicos)
	ingresarDatos("número de simulaciones a ejecutar", &numSims)

	fmt.Println()
	fmt.Println("Seleccione el método de sincronización:")
	fmt.Println("1. RWMutex (Read-Write Mutex)")
	fmt.Println("2. WaitGroup (Canales como Semáforos)")
	fmt.Println("3. Comparar ambos")
	fmt.Print("Opción: ")
	fmt.Scanln(&option)

	fmt.Println()
	fmt.Println("Iniciando simulación...")
	fmt.Println()

	switch option {
	case 1:
		for i := 1; i <= numSims; i++ {
			if numSims > 1 {
				fmt.Printf("\n========== SIMULACIÓN %d/%d ==========\n\n", i, numSims)
			}
			EjecutarSimulacionRWMutex(numCoches, numPlazas, numMecanicos)
		}
	case 2:
		for i := 1; i <= numSims; i++ {
			if numSims > 1 {
				fmt.Printf("\n========== SIMULACIÓN %d/%d ==========\n\n", i, numSims)
			}
			EjecutarSimulacionWaitGroup(numCoches, numPlazas, numMecanicos)
		}
	case 3:
		for i := 1; i <= numSims; i++ {
			if numSims > 1 {
				fmt.Printf("\n========== SIMULACIÓN %d/%d ==========\n\n", i, numSims)
			}
			fmt.Println("--- Método 1: RWMutex ---")
			EjecutarSimulacionRWMutex(numCoches, numPlazas, numMecanicos)
			fmt.Println("\n--- Método 2: WaitGroup ---")
			EjecutarSimulacionWaitGroup(numCoches, numPlazas, numMecanicos)
		}
	default:
		fmt.Println("Opción no válida")
	}

	fmt.Println("\nPresione ENTER para volver al menú principal...")
	reader.ReadString('\n')
	reader.ReadString('\n')
}
