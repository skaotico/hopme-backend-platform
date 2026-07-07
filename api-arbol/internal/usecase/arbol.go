package usecase

import (
	"log/slog"

	"c4-arbol/internal/domain/port"
)

type arbolUseCase struct {
	repo port.ArbolRepository
	log  *slog.Logger
}

// NewArbolUseCase crea una nueva instancia del caso de uso de Árbol.
func NewArbolUseCase(repo port.ArbolRepository, log *slog.Logger) port.ArbolUseCase {
	return &arbolUseCase{
		repo: repo,
		log:  log.With(slog.String("component", "arbol_usecase")),
	}
}
