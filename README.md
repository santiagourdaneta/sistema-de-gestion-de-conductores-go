# 🚗 Sistema de Gestión de Conductores Concurrente en Go

Este proyecto es una simulación de un sistema de gestión de conductores y pedidos, similar a plataformas como Uber, construido con Go. El objetivo principal es demostrar el uso práctico de la **concurrencia** a través de `goroutines` y `canales` para manejar múltiples tareas simultáneamente.

## 🚀 Características

* **Asignación de Pedidos:** El sistema asigna un pedido al conductor disponible más cercano.
* **Concurrencia:** Utiliza `goroutines` para simular la operación de varios conductores y pasajeros al mismo tiempo.
* **Sincronización Segura:** Emplea `canales` para la comunicación entre goroutines y `mutex` para proteger la lista de conductores, evitando condiciones de carrera.
* **Simulación de Viajes:** Cada viaje es una operación concurrente que simula un tiempo de duración aleatorio.

## ⚙️ Cómo Ejecutar el Proyecto

1.  Asegúrate de tener Go instalado en tu sistema.
2.  Clona el repositorio:
    ```bash
    git clone https://github.com/santiagourdaneta/sistema-de-gestion-de-conductores-go/
    cd sistema-de-gestion-de-conductores-go
    ```
3.  Inicializa el módulo de Go (si es la primera vez):
    ```bash
    go mod init sistema-de-gestion-de-conductores-go
    ```
4.  Ejecuta el programa:
    ```bash
    go run main.go
    ```

## ✅ Pruebas

El proyecto incluye pruebas unitarias y de integración para validar su funcionalidad.

* Para ejecutar todas las pruebas:
    ```bash
    go test
    ```

## 📈 Análisis de Rendimiento (Pprof)

Para analizar el rendimiento y detectar posibles cuellos de botella:

1.  Abre un terminal y ejecuta el programa con el servidor `pprof` habilitado (asegúrate de que `main.go` esté modificado para incluir `net/http/pprof`).
    ```bash
    go run main.go
    ```
2.  Abre un segundo terminal y utiliza `go tool pprof` para acceder al servidor y obtener un gráfico en el navegador:
    ```bash
    go tool pprof -http=:8080 localhost:6060/debug/pprof/heap
    ```
3.  Esto abrirá una interfaz web donde podrás analizar el uso de CPU, memoria, goroutines, etc.

---

## 🛠️ Estructura del Código

* `main.go`: Contiene la lógica principal de la simulación, incluyendo las estructuras de datos, las goroutines para pasajeros y el asignador de pedidos.
* `main_test.go`: Incluye las pruebas unitarias para las funciones principales y las pruebas de integración para el flujo de asignación de pedidos.

## ✍️ Contribuir

¡Las contribuciones son bienvenidas! Siéntete libre de abrir un *issue* o enviar un *pull request* con mejoras, correcciones o nuevas funcionalidades.

Labels & Tags
go golang concurrencia goroutines canales simulación sistema-distribuido mutex ejemplo proyecto-de-practica

Hashtags
#go #golang #concurrencia #goroutines #canales #sistemagestion #simulacion #concurrency #softwareengineering
