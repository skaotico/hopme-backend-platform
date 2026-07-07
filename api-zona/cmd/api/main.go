package main

// @title           Homelab Platform - API Zona
// @version         1.0
// @description     Microservicio de gestión de zonas dentro de ecoparques. Parte de la Homelab Collector Platform.
// @contact.name    Skaotico
// @host            localhost:9092
// @BasePath        /api/v1
// @schemes         http

import "c4-zona/internal/bootstrap"

func main() {
	bootstrap.Run()
}
