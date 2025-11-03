package main

import (
	"context"
	"fmt"

	"RaceAll/internal/broadcast"
	"RaceAll/internal/logger"
	"RaceAll/internal/sharedmemory"
)

// App struct
type App struct {
	ctx               context.Context
	sharedMemService  *sharedmemory.Service
	connectionManager *broadcast.ConnectionManager
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		sharedMemService: sharedmemory.NewService(),
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Inicializar logger con nivel Debug
	logger.Init(logger.DevelopmentConfig())

	// Iniciar servicio de shared memory
	if err := a.sharedMemService.Start(); err != nil {
		fmt.Printf("Error starting shared memory service: %v\n", err)
	}

	// Configurar e iniciar connection manager
	config := broadcast.DefaultConfig()
	a.connectionManager = broadcast.NewConnectionManager(config, a.sharedMemService)

	if err := a.connectionManager.Start(); err != nil {
		fmt.Printf("Error starting connection manager: %v\n", err)
	}
}

// shutdown is called at application termination
func (a *App) shutdown(ctx context.Context) {
	if a.connectionManager != nil {
		a.connectionManager.Stop()
	}

	if a.sharedMemService != nil {
		a.sharedMemService.Stop()
	}
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// GetBroadcastService devuelve el servicio de broadcast
func (a *App) GetBroadcastService() *broadcast.Service {
	if a.connectionManager == nil {
		return nil
	}
	return a.connectionManager.GetService()
}

// IsConnected verifica si el broadcast está conectado
func (a *App) IsConnected() bool {
	if a.connectionManager == nil {
		return false
	}
	return a.connectionManager.IsConnected()
}
