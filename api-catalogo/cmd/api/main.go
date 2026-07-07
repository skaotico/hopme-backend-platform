package main

import "c4-catalogo/internal/bootstrap"

// @title           Homelab Collector Platform - Catalogo API
// @version         1.0
// @description     Servicio de Catálogo de Homelab Collector Platform.
// @host            localhost:9094
// @BasePath        /
func main() {
	bootstrap.Run()
}

