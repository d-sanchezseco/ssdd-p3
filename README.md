# ssdd-p3
# Simulación de Taller de Reparación de Vehículos

Sistema de simulación concurrente para un taller de reparación de vehículos implementado en Go, que modela el flujo de trabajo completo desde la llegada hasta la entrega de vehículos.

## 📋 Descripción

Este proyecto simula un taller de reparación con múltiples fases secuenciales, gestión de prioridades y dos implementaciones diferentes de concurrencia (RWMutex y WaitGroup).

### Fases del Proceso

1. **Llegada**: El vehículo llega y ocupa una plaza de espera
2. **Reparación**: Un mecánico atiende el vehículo según su especialidad
3. **Limpieza**: El vehículo pasa por limpieza
4. **Entrega**: El vehículo es entregado al cliente

### Categorías de Vehículos

| Categoría | Tipo de Incidencia | Prioridad | Tiempo por Fase |
|-----------|-------------------|-----------|-----------------|
| **A** | Mecánica | 3 (Alta) | 5 segundos |
| **B** | Eléctrica | 2 (Media) | 3 segundos |
| **C** | Carrocería | 1 (Baja) | 1 segundo |

## 🚀 Ejecución

### Compilar y Ejecutar

```bash
go run .
```

### Ejecutar Tests Comparativos

```bash
# Todos los tests comparativos
go test -v -run TestComparativa

# Test específico
go test -v -run TestComparativa1_Balanceado
go test -v -run TestComparativa2_AHeavy
go test -v -run TestComparativa3_CHeavy
go test -v -run TestComparativaCompleta

# Ejecutar benchmarks
go test -bench=. -benchmem
```

## 📊 Resultados de Tests Comparativos

Se han ejecutado tres escenarios de prueba con diferentes distribuciones de categorías de vehículos:

### Tabla Comparativa

```
╔══════════════════════════════════════════════════════════════════════════════╗
║                    TABLA COMPARATIVA FINAL                                   ║
╠══════════════════════════════════════════════════════════════════════════════╣
║ Test              │ Dist.    │ RWMutex      │ WaitGroup    │ Diferencia     ║
╠═══════════════════╪══════════╪══════════════╪══════════════╪════════════════╣
║ Test 1: Balanceado │ 10/10/10 │ 44.197s      │ 40.582s      │ 3.615s         ║
║ Test 2: A-Heavy   │ 20/ 5/ 5 │ 1m0.417s     │ 52.466s      │ 7.951s         ║
║ Test 3: C-Heavy   │  5/ 5/20 │ 29.581s      │ 26.047s      │ 3.535s         ║
╚══════════════════════════════════════════════════════════════════════════════╝
```

### Análisis de Resultados

#### Test 1: Distribución Balanceada (A=10, B=10, C=10)
- **Tiempo**: ~40-46 segundos
- **Observación**: WaitGroup fue ~9% más rápido que RWMutex
- **Conclusión**: Distribución equilibrada permite buen aprovechamiento de recursos

#### Test 2: A-Heavy (A=20, B=5, C=5) ⏱️
- **Tiempo**: ~1 minuto (el más lento)
- **Observación**: WaitGroup fue ~13% más rápido
- **Conclusión**: Alta concentración de tareas de categoría A (5s por fase) incrementa significativamente el tiempo total

#### Test 3: C-Heavy (A=5, B=5, C=20) ⚡
- **Tiempo**: ~26 segundos (el más rápido)
- **Observación**: Diferencia mínima entre implementaciones (84ms)
- **Conclusión**: Mayoría de tareas de categoría C (1s por fase) acelera el procesamiento

### Conclusiones Generales

✅ **Todos los tests pasaron exitosamente** (4/4)  
✅ El sistema maneja correctamente diferentes distribuciones de carga  
✅ La priorización funciona adecuadamente (A > B > C)  
✅ Ambas implementaciones (RWMutex y WaitGroup) son viables  
✅ WaitGroup mostró mejor rendimiento general, especialmente en escenarios A-Heavy  

**Impacto de la distribución**:
- C-Heavy: ~26s (más rápido)
- Balanceado: ~40-46s
- A-Heavy: ~60s (más lento)

## 🏗️ Estructura del Proyecto

```
entrega-3/
├── main.go                      # Punto de entrada principal
├── tipos.go                     # Definición de tipos y estructuras
├── simulacion.go                # Lógica de simulación base
├── simulacion_rwmutex.go        # Implementación con RWMutex
├── simulacion_waitgroup.go      # Implementación con WaitGroup
├── prioridades.go               # Sistema de colas con prioridad
├── comparativas_test.go         # Tests comparativos
├── logger.go                    # Sistema de logging
├── menu_simulacion.go           # Menú interactivo
├── incidencias.go               # Gestión de incidencias
├── mecanico.go                  # Gestión de mecánicos
├── vehiculo.go                  # Gestión de vehículos
├── clientes.go                  # Gestión de clientes
├── auxiliar.go                  # Funciones auxiliares
└── readme.md                    # Este archivo
```

## 🔧 Implementaciones de Concurrencia

### RWMutex
- Utiliza `sync.RWMutex` para control de acceso a recursos compartidos
- Permite múltiples lecturas simultáneas
- Bloqueo exclusivo para escrituras

### WaitGroup
- Utiliza canales buffered para gestión de recursos
- `sync.WaitGroup` para sincronización de goroutines
- Enfoque más orientado a paso de mensajes

## 📝 Características Principales

- ✨ Sistema de prioridades basado en categorías
- 🔄 Procesamiento concurrente de múltiples vehículos
- 📊 Logging detallado de todas las fases
- ⚡ Dos implementaciones de concurrencia para comparación
- 🧪 Suite completa de tests comparativos
- 📈 Benchmarks de rendimiento

## 👥 Autores

Proyecto desarrollado para la asignatura de Sistemas Distribuidos.

## 📄 Licencia

Este proyecto es de uso académico.

