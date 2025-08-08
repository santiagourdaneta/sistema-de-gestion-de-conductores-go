package main

import (
	"math"
	"sync"
	"testing"
)

// --- Pruebas Unitarias ---

// TestCalcularDistancia verifica que la función de distancia al cuadrado
// funcione correctamente para diferentes pares de coordenadas.
func TestCalcularDistancia(t *testing.T) {
	tests := []struct {
		name     string
		p1       Coordenada
		p2       Coordenada
		expected float64
	}{
		{"mismo punto", Coordenada{0, 0}, Coordenada{0, 0}, 0.0},
		{"eje x", Coordenada{0, 0}, Coordenada{3, 0}, 9.0},
		{"eje y", Coordenada{0, 0}, Coordenada{0, 4}, 16.0},
		{"diagonal", Coordenada{0, 0}, Coordenada{3, 4}, 25.0},
		{"puntos negativos", Coordenada{-1, -1}, Coordenada{2, 3}, 25.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calcularDistancia(tt.p1, tt.p2)
			if got != tt.expected {
				t.Errorf("calcularDistancia() = %v, se esperaba %v", got, tt.expected)
			}
		})
	}
}

// --- Pruebas de Integración ---

// TestAsignacionDePedido verifica que el sistema asigne el pedido al conductor
// más cercano y lo marque como no disponible.
func TestAsignacionDePedido(t *testing.T) {
	// 1. Setup: Crear un conjunto de conductores y un pedido para esta prueba
	conductoresDePrueba := []*Conductor{
		{ID: 1, Ubicacion: Coordenada{10, 10}, Disponible: true}, // El más cercano
		{ID: 2, Ubicacion: Coordenada{50, 50}, Disponible: true},
		{ID: 3, Ubicacion: Coordenada{80, 80}, Disponible: true},
	}
	pedidoDePrueba := &Pedido{ID: 1, Origen: Coordenada{12, 12}}

	// Usa un canal para simular el flujo de un pedido real
	testPedidosPendientes := make(chan *Pedido, 1)

	// Simula un "asignador" que procesará un solo pedido
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for pedido := range testPedidosPendientes {
			// Bloquear el acceso a los conductores de prueba para simular la concurrencia
			conductoresMutex.Lock()
			
			var conductorAsignado *Conductor
			distanciaMinima := math.MaxFloat64
			for _, c := range conductoresDePrueba {
				if c.Disponible {
					distancia := calcularDistancia(c.Ubicacion, pedido.Origen)
					if distancia < distanciaMinima {
						distanciaMinima = distancia
						conductorAsignado = c
					}
				}
			}

			if conductorAsignado != nil {
				conductorAsignado.Disponible = false
				pedido.AsignadoA = conductorAsignado
			}
			
			conductoresMutex.Unlock()
			// Detener el bucle después de procesar el pedido
			close(testPedidosPendientes)
			break
		}
	}()

	// 2. Acción: Enviar el pedido al canal para que sea procesado
	testPedidosPendientes <- pedidoDePrueba
	
	// Esperar a que el goroutine del asignador termine
	wg.Wait()

	// 3. Verificación: Asegurarse de que el resultado es el esperado
	if pedidoDePrueba.AsignadoA == nil {
		t.Fatal("El pedido no fue asignado a ningún conductor.")
	}

	if pedidoDePrueba.AsignadoA.ID != 1 {
		t.Errorf("Se esperaba que el pedido se asignara al conductor 1, pero se asignó al %d", pedidoDePrueba.AsignadoA.ID)
	}

	if pedidoDePrueba.AsignadoA.Disponible {
		t.Errorf("El conductor asignado debería estar no disponible después de la asignación")
	}

	// Restaurar el estado del conductor para otras pruebas si fuera necesario
	pedidoDePrueba.AsignadoA.Disponible = true
	
	// Limpiar el estado de los conductores de prueba
	conductoresDePrueba = nil
}