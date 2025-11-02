package broadcast

import (
	"context"
	"fmt"
	"sync"

	"RaceAll/internal/logger"
)

type Service struct {
	client      *Client
	ctx         context.Context
	cancel      context.CancelFunc
	wg          sync.WaitGroup
	isRunning   bool
	mu          sync.RWMutex
	subscribers []chan BroadcastMessage
	subMu       sync.RWMutex
	config      Config
}

type Config struct {
	Host            string
	Port            int
	DisplayName     string
	Password        string
	CommandPassword string
	UpdateMS        int32
}

type BroadcastMessage struct {
	Type    string
	Payload interface{}
}

func DefaultConfig() Config {
	return Config{
		Host:            "127.0.0.1",
		Port:            9000,
		DisplayName:     "RaceAll",
		Password:        "asd",
		CommandPassword: "",
		UpdateMS:        100,
	}
}

func NewService(config Config) *Service {
	if config.UpdateMS == 0 {
		config.UpdateMS = 100
	}

	return &Service{
		config:      config,
		subscribers: make([]chan BroadcastMessage, 0),
	}
}

func (s *Service) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.isRunning {
		return nil
	}

	// Create client
	address := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
	log := *logger.Get()

	s.client = NewClient(
		address,
		s.config.DisplayName,
		s.config.Password,
		s.config.CommandPassword,
		s.config.UpdateMS,
		log,
	)

	// Set callbacks
	s.setupCallbacks()

	s.ctx, s.cancel = context.WithCancel(context.Background())
	s.isRunning = true

	// Connect in background
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()

		err := s.client.ConnectAndListen()
		if err != nil {
			logger.Errorf("Broadcast client error: %v", err)
			s.mu.Lock()
			s.isRunning = false
			s.mu.Unlock()
		}
	}()

	logger.Info("Broadcast service started")
	return nil
}

func (s *Service) setupCallbacks() {
	s.client.OnConnectionStateChanged = func(state ConnectionState) {
		s.notifySubscribers(BroadcastMessage{
			Type:    "ConnectionState",
			Payload: state,
		})
	}

	s.client.OnTrackDataUpdate = func(track TrackData) {
		s.notifySubscribers(BroadcastMessage{
			Type:    "TrackData",
			Payload: track,
		})
	}

	s.client.OnEntrylistUpdate = func(car CarInfo) {
		s.notifySubscribers(BroadcastMessage{
			Type:    "EntryList",
			Payload: car,
		})
	}

	s.client.OnRealtimeUpdate = func(update RealtimeUpdate) {
		s.notifySubscribers(BroadcastMessage{
			Type:    "RealtimeUpdate",
			Payload: update,
		})
	}

	s.client.OnRealtimeCarUpdate = func(update RealtimeCarUpdate) {
		s.notifySubscribers(BroadcastMessage{
			Type:    "RealtimeCarUpdate",
			Payload: update,
		})
	}

	s.client.OnBroadcastingEvent = func(event BroadcastingEvent) {
		s.notifySubscribers(BroadcastMessage{
			Type:    "BroadcastingEvent",
			Payload: event,
		})
	}
}

func (s *Service) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.isRunning {
		return
	}

	if s.client != nil {
		s.client.Disconnect()
	}

	s.cancel()
	s.wg.Wait()

	s.isRunning = false
	logger.Info("Broadcast service stopped")
}

func (s *Service) notifySubscribers(msg BroadcastMessage) {
	s.subMu.RLock()
	defer s.subMu.RUnlock()

	for _, ch := range s.subscribers {
		select {
		case ch <- msg:
		default:
			// Skip if channel is full
		}
	}
}

func (s *Service) Subscribe() <-chan BroadcastMessage {
	s.subMu.Lock()
	defer s.subMu.Unlock()

	ch := make(chan BroadcastMessage, 10)
	s.subscribers = append(s.subscribers, ch)
	return ch
}

func (s *Service) Unsubscribe(ch <-chan BroadcastMessage) {
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

func (s *Service) SetFocus(carIndex uint16) error {
	if s.client == nil {
		return fmt.Errorf("not connected")
	}
	return s.client.SetFocus(carIndex)
}

func (s *Service) SetCamera(cameraSet, camera string) error {
	if s.client == nil {
		return fmt.Errorf("not connected")
	}
	return s.client.SetCamera(cameraSet, camera)
}

func (s *Service) RequestInstantReplay(startTime, duration float32, initialFocusedCarIndex int32, cameraSet, camera string) error {
	if s.client == nil {
		return fmt.Errorf("not connected")
	}
	return s.client.RequestInstantReplay(startTime, duration, initialFocusedCarIndex, cameraSet, camera)
}

func (s *Service) RequestHUDPage(hudPage string) error {
	if s.client == nil {
		return fmt.Errorf("not connected")
	}
	return s.client.RequestHUDPage(hudPage)
}

func (s *Service) GetCarInfo(carIndex uint16) (*CarInfo, bool) {
	if s.client == nil {
		return nil, false
	}
	return s.client.GetCarInfo(carIndex)
}
