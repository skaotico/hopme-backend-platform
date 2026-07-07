package main

// @title           Homelab Platform - API Arbol
// @version         1.0
// @description     Microservicio de gestión de árboles (flora.arbol, historial y mediciones). Parte de la Homelab Collector Platform.
// @contact.name    Skaotico
// @host            localhost:9095
// @BasePath        /api/v1
// @schemes         http

import "c4-arbol/internal/bootstrap"

func main() {
	bootstrap.Run()
}
