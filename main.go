package main

import (
	"fmt"
	"math"
	"math/rand"
	"sync"
	"time"
	"net/http"
	_ "net/http/pprof" // The underscore means we import it for its side effects
)

// Estructuras para representar el sistema
type Coordenada struct {
	X int
	Y int
}

type Conductor struct {
	ID         int
	Ubicacion  Coordenada
	Disponible bool
}

type Pedido struct {
	ID        int
	Origen    Coordenada
	Destino   Coordenada
	AsignadoA *Conductor
}

// Canal para recibir los pedidos
var pedidosPendientes = make(chan *Pedido, 100)

// Mutex para proteger el acceso a la lista de conductores
var conductoresMutex sync.Mutex
var conductores []*Conductor

func main() {

	// Start an HTTP server for pprof to listen on
		go func() {
			fmt.Println("Pprof server listening on :6060")
			if err := http.ListenAndServe("localhost:6060", nil); err != nil {
				fmt.Println("Pprof server error:", err)
			}
		}()


	rand.Seed(time.Now().UnixNano())

	// 1. Inicializar conductores
	for i := 1; i <= 5; i++ {
		conductores = append(conductores, &Conductor{
			ID:         i,
			Ubicacion:  Coordenada{rand.Intn(100), rand.Intn(100)},
			Disponible: true,
		})
	}
	fmt.Println("¡Sistema de gestión de conductores inicializado! 🚗")

	// 2. Iniciar el goroutine del asignador de pedidos
	go asignarPedidos()

	// 3. Simular la llegada de pasajeros (pedidos)
	var wg sync.WaitGroup
	numPedidos := 100 // Puedes aumentar este número para pruebas de estrés

	for i := 1; i <= numPedidos; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			pedido := &Pedido{
				ID:      id,
				Origen:  Coordenada{rand.Intn(100), rand.Intn(100)},
				Destino: Coordenada{rand.Intn(100), rand.Intn(100)},
			}
			pedidosPendientes <- pedido // El pasajero crea un pedido y lo envía al canal
			fmt.Printf("📦 Pasajero %d ha solicitado un viaje desde (%d, %d) a (%d, %d)\n",
				pedido.ID, pedido.Origen.X, pedido.Origen.Y, pedido.Destino.X, pedido.Destino.Y)
		}(i)
		time.Sleep(time.Millisecond * 200) // Simular un tiempo entre pedidos
	}

	wg.Wait()
	close(pedidosPendientes) // Cerrar el canal cuando todos los pedidos han sido enviados
	time.Sleep(time.Second * 5) // Dar tiempo para que los viajes finalicen
	fmt.Println("\n¡Todos los pedidos han sido procesados y los viajes han finalizado! 🎉")
}

// `asignarPedidos` es el corazón del sistema, se ejecuta en su propia goroutine
func asignarPedidos() {
	for pedido := range pedidosPendientes {
		fmt.Printf("\n🔍 El sistema está buscando un conductor para el pedido %d...\n", pedido.ID)
		encontrarYAsignarConductor(pedido)
	}
}

// Encuentra el conductor más cercano y lo asigna al pedido
func encontrarYAsignarConductor(pedido *Pedido) {
	conductoresMutex.Lock()
	defer conductoresMutex.Unlock()

	var conductorAsignado *Conductor
	distanciaMinima := math.MaxFloat64 // Usar el valor máximo para la comparación inicial

	for _, c := range conductores {
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
		fmt.Printf("✅ Pedido %d asignado al conductor %d. Distancia al conductor: %.2f\n",
			pedido.ID, conductorAsignado.ID, math.Sqrt(distanciaMinima))
		go iniciarViaje(pedido) // Iniciar el viaje en una nueva goroutine
	} else {
		fmt.Printf("❌ No se encontró un conductor disponible para el pedido %d. Reintentando...\n", pedido.ID)
		// Aquí se podría implementar una lógica de reintento, por ejemplo,
		// enviando el pedido de vuelta al canal después de un breve delay.
	}
}

// `iniciarViaje` simula el viaje del conductor
func iniciarViaje(pedido *Pedido) {
	tiempoViaje := time.Duration(rand.Intn(5)+3) * time.Second
	fmt.Printf("🚗 Conductor %d ha aceptado el pedido %d y está en camino. El viaje durará aprox. %v\n",
		pedido.AsignadoA.ID, pedido.ID, tiempoViaje)
	time.Sleep(tiempoViaje)
	fmt.Printf("🏁 El viaje del pedido %d con el conductor %d ha finalizado.\n",
		pedido.ID, pedido.AsignadoA.ID)

	// Una vez que el viaje termina, el conductor vuelve a estar disponible
	conductoresMutex.Lock()
	pedido.AsignadoA.Ubicacion = pedido.Destino
	pedido.AsignadoA.Disponible = true
	conductoresMutex.Unlock()
}

// `calcularDistancia` calcula la distancia euclidiana al cuadrado entre dos puntos
func calcularDistancia(p1, p2 Coordenada) float64 {
	dx := float64(p1.X - p2.X)
	dy := float64(p1.Y - p2.Y)
	return dx*dx + dy*dy // Se devuelve la distancia al cuadrado para evitar el cálculo de la raíz, que es costoso.
}