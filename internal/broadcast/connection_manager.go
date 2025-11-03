package broadcast

import (
	"RaceAll/internal/logger"
	"RaceAll/internal/sharedmemory"
	"context"
	"sync"
	"time"
)

// ConnectionManager maneja la reconexión automática del broadcast
type ConnectionManager struct {
	service          *Service
	sharedMemService *sharedmemory.Service
	ctx              context.Context
	cancel           context.CancelFunc
	wg               sync.WaitGroup
	mu               sync.RWMutex
	isMonitoring     bool
	lastPacketID     int32
	config           Config
}

// NewConnectionManager crea un nuevo gestor de conexiones
func NewConnectionManager(config Config, smService *sharedmemory.Service) *ConnectionManager {
	return &ConnectionManager{
		config:           config,
		sharedMemService: smService,
		lastPacketID:     0,
	}
}

// Start inicia el monitoreo de conexión
func (cm *ConnectionManager) Start() error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if cm.isMonitoring {
		return nil
	}

	cm.ctx, cm.cancel = context.WithCancel(context.Background())
	cm.isMonitoring = true

	// Iniciar servicio de broadcast
	cm.service = NewService(cm.config)

	// Iniciar goroutine de monitoreo
	cm.wg.Add(1)
	go cm.monitorConnection()

	logger.Info("Connection manager started")
	return nil
}

// Stop detiene el monitoreo
func (cm *ConnectionManager) Stop() {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if !cm.isMonitoring {
		return
	}

	cm.cancel()
	cm.wg.Wait()

	if cm.service != nil {
		cm.service.Stop()
	}

	cm.isMonitoring = false
	logger.Info("Connection manager stopped")
}

// monitorConnection monitorea el estado del juego y reconecta cuando es necesario
// Basado en Race Element: PageGraphicsTracker.cs
func (cm *ConnectionManager) monitorConnection() {
	defer cm.wg.Done()

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-cm.ctx.Done():
			return
		case <-ticker.C:
			cm.checkAndReconnect()
		}
	}
}

// checkAndReconnect verifica el estado del juego y reconecta si es necesario
func (cm *ConnectionManager) checkAndReconnect() {
	if cm.sharedMemService == nil {
		return
	}

	// Leer datos de shared memory
	data, err := cm.sharedMemService.GetLatestData()
	if err != nil {
		// Si no podemos leer shared memory, desconectar
		if cm.service != nil && cm.service.IsRunning() {
			logger.Info("Cannot read shared memory, disconnecting broadcast")
			cm.service.Stop()
		}
		return
	}

	graphics := data.Graphics
	physics := data.Physics

	isGameRunning := graphics.Status != sharedmemory.ACOff
	isBroadcastConnected := cm.service != nil && cm.service.IsRunning()

	// Caso 1: Juego está en OFF/MENU
	if graphics.Status == sharedmemory.ACOff {
		cm.handleGameOff(physics.PacketId, isBroadcastConnected)
		return
	}

	// Caso 2: Juego está activo (LIVE, PAUSE, REPLAY)
	if isGameRunning {
		if !isBroadcastConnected {
			logger.Info("Game is running, connecting broadcast")
			cm.connectBroadcast()
		}
	}
}

// handleGameOff maneja el caso cuando el juego está en OFF/MENU
// Race Element detecta si el juego está realmente cerrado o solo pausado
// comparando PacketIDs de physics
func (cm *ConnectionManager) handleGameOff(currentPacketID int32, isBroadcastConnected bool) {
	// Si el PacketID es 0 o menor que el anterior, el juego está realmente cerrado
	if currentPacketID <= cm.lastPacketID || currentPacketID == 0 {
		if isBroadcastConnected {
			logger.Info("Game closed detected, disconnecting broadcast")
			cm.service.Stop()
			cm.lastPacketID = 0
		}
	} else {
		// PacketID está aumentando, significa que el juego está pausado pero activo
		if !isBroadcastConnected {
			logger.Info("Game paused but active, connecting broadcast")
			// Esperar un poco antes de reconectar
			time.Sleep(1 * time.Second)
			cm.connectBroadcast()
		}
	}

	cm.lastPacketID = currentPacketID
}

// connectBroadcast conecta el servicio de broadcast
func (cm *ConnectionManager) connectBroadcast() {
	if cm.service == nil {
		cm.service = NewService(cm.config)
	}

	err := cm.service.Start()
	if err != nil {
		logger.Errorf("Failed to connect broadcast: %v", err)
	} else {
		logger.Info("Broadcast connected successfully")
	}
}

// GetService devuelve el servicio de broadcast
func (cm *ConnectionManager) GetService() *Service {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.service
}

// IsConnected verifica si el broadcast está conectado
func (cm *ConnectionManager) IsConnected() bool {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.service != nil && cm.service.IsRunning()
}
