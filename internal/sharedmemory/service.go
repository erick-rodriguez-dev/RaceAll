package sharedmemory

import (
	"context"
	"sync"
	"time"

	"RaceAll/internal/errors"
	"RaceAll/internal/logger"
)

type Service struct {
	reader               *SharedMemoryReader
	ctx                  context.Context
	cancel               context.CancelFunc
	wg                   sync.WaitGroup
	isRunning            bool
	mu                   sync.RWMutex
	subscribers          []chan TelemetryData
	subMu                sync.RWMutex
	reconnectEnabled     bool
	reconnectDelay       time.Duration
	consecutiveErrors    int
	maxConsecutiveErrors int
}

type TelemetryData struct {
	Physics  *Physics
	Graphics *Graphics
	Static   *Static
}

func NewService() *Service {
	return &Service{
		reader:               NewSharedMemoryReader(),
		subscribers:          make([]chan TelemetryData, 0),
		reconnectEnabled:     true,
		reconnectDelay:       5 * time.Second,
		maxConsecutiveErrors: 10,
	}
}

func (s *Service) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.isRunning {
		return nil
	}

	s.ctx, s.cancel = context.WithCancel(context.Background())
	s.isRunning = true

	s.wg.Add(1)
	go s.connectionLoop()

	logger.Info("Shared Memory service started")
	return nil
}

func (s *Service) connectionLoop() {
	defer s.wg.Done()

	for {
		select {
		case <-s.ctx.Done():
			return
		default:
			// Intentar conectar
			err := s.reader.Connect()
			if err != nil {
				logger.Warnf("Failed to connect to shared memory: %v. Retrying in %v...", err, s.reconnectDelay)
				time.Sleep(s.reconnectDelay)
				continue
			}

			logger.Info("Successfully connected to ACC shared memory")
			s.consecutiveErrors = 0

			// Iniciar el loop de lectura
			s.readLoop()

			// Si llegamos aquí, el readLoop terminó (desconexión)
			s.reader.Disconnect()

			if !s.reconnectEnabled {
				return
			}

			logger.Warn("Shared memory connection lost. Attempting to reconnect...")
			time.Sleep(s.reconnectDelay)
		}
	}
}

func (s *Service) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.isRunning {
		return
	}

	s.cancel()
	s.wg.Wait()
	s.reader.Disconnect()
	s.isRunning = false

	logger.Info("Shared Memory service stopped")
}

func (s *Service) readLoop() {
	ticker := time.NewTicker(16 * time.Millisecond) // ~60Hz
	defer ticker.Stop()

	for {
		select {
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			physics, err := s.reader.ReadPhysics()
			if err != nil {
				s.consecutiveErrors++
				if s.consecutiveErrors >= s.maxConsecutiveErrors {
					logger.Errorf("Too many consecutive read errors, reconnecting...")
					return
				}
				continue
			}

			graphics, err := s.reader.ReadGraphics()
			if err != nil {
				s.consecutiveErrors++
				if s.consecutiveErrors >= s.maxConsecutiveErrors {
					logger.Errorf("Too many consecutive read errors, reconnecting...")
					return
				}
				continue
			}

			static, err := s.reader.ReadStatic()
			if err != nil {
				s.consecutiveErrors++
				if s.consecutiveErrors >= s.maxConsecutiveErrors {
					logger.Errorf("Too many consecutive read errors, reconnecting...")
					return
				}
				continue
			}

			// Reset error counter on successful read
			s.consecutiveErrors = 0

			data := TelemetryData{
				Physics:  physics,
				Graphics: graphics,
				Static:   static,
			}

			s.notifySubscribers(data)
		}
	}
}

func (s *Service) notifySubscribers(data TelemetryData) {
	s.subMu.RLock()
	defer s.subMu.RUnlock()

	for _, ch := range s.subscribers {
		select {
		case ch <- data:
		default:
			// Skip if channel is full
		}
	}
}

func (s *Service) Subscribe() <-chan TelemetryData {
	s.subMu.Lock()
	defer s.subMu.Unlock()

	ch := make(chan TelemetryData, 10)
	s.subscribers = append(s.subscribers, ch)
	return ch
}

func (s *Service) Unsubscribe(ch <-chan TelemetryData) {
	s.subMu.Lock()
	defer s.subMu.Unlock()

	for i, sub := range s.subscribers {
		if sub == ch {
			s.subscribers = append(s.subscribers[:i], s.subscribers[i+1:]...)
			close(sub)
			break
		}
	}
}

func (s *Service) IsRunning() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.isRunning
}

func (s *Service) GetLatestData() (*TelemetryData, error) {
	if !s.reader.IsConnected() {
		return nil, NewError("GetLatestData", errors.ErrNotConnected)
	}

	physics, err := s.reader.ReadPhysics()
	if err != nil {
		return nil, err
	}

	graphics, err := s.reader.ReadGraphics()
	if err != nil {
		return nil, err
	}

	static, err := s.reader.ReadStatic()
	if err != nil {
		return nil, err
	}

	return &TelemetryData{
		Physics:  physics,
		Graphics: graphics,
		Static:   static,
	}, nil
}

// IsConnected returns true if the reader is connected to shared memory
func (s *Service) IsConnected() bool {
	return s.reader != nil && s.reader.IsConnected()
}

// SetReconnectEnabled enables or disables automatic reconnection
func (s *Service) SetReconnectEnabled(enabled bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.reconnectEnabled = enabled
}

// SetReconnectDelay sets the delay between reconnection attempts
func (s *Service) SetReconnectDelay(delay time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.reconnectDelay = delay
}
