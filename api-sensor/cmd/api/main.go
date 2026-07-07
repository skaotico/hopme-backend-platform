package main

// @title           Homelab Platform - API Sensor
// @version         1.0
// @description     Microservicio de gestión de sensores (iot.sensor, readings, alerts). Parte de la Homelab Collector Platform.
// @contact.name    Skaotico
// @host            localhost:9096
// @BasePath        /api/v1
// @schemes         http

import "c4-sensor/internal/bootstrap"

func main() {
	bootstrap.Run()
}
