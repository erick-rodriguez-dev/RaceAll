package broadcast

import (
	"bytes"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

// Protocol maneja el protocolo de comunicación con ACC
type Protocol struct {
	connectionId         int32
	connectionIdentifier string
	sendFunc             func([]byte) error
	logger               zerolog.Logger

	// Caché de entry list para evitar grandes paquetes UDP
	entryListCars  map[uint16]*CarInfo
	entryListMutex sync.RWMutex

	lastEntryListRequest time.Time

	// Callbacks
	OnConnectionStateChanged func(ConnectionState)
	OnTrackDataUpdate        func(TrackData)
	OnEntrylistUpdate        func(CarInfo)
	OnRealtimeUpdate         func(RealtimeUpdate)
	OnRealtimeCarUpdate      func(RealtimeCarUpdate)
	OnBroadcastingEvent      func(BroadcastingEvent)
}

func NewProtocol(connectionIdentifier string, sendFunc func([]byte) error, logger zerolog.Logger) *Protocol {
	return &Protocol{
		connectionIdentifier: connectionIdentifier,
		sendFunc:             sendFunc,
		logger:               logger,
		entryListCars:        make(map[uint16]*CarInfo),
		lastEntryListRequest: time.Now(),
	}
}

func (p *Protocol) ProcessMessage(data []byte) error {
	if len(data) == 0 {
		return nil
	}

	reader := bytes.NewReader(data)

	msgType, err := readUint8(reader)
	if err != nil {
		p.logger.Error().Err(err).Msg("Error leyendo tipo de mensaje")
		return err
	}

	messageType := InboundMessageType(msgType)

	switch messageType {
	case InboundRegistrationResult:
		return p.handleRegistrationResult(reader)
	case InboundEntryList:
		return p.handleEntryList(reader)
	case InboundEntryListCar:
		return p.handleEntryListCar(reader)
	case InboundTrackData:
		return p.handleTrackData(reader)
	case InboundRealtimeUpdate:
		return p.handleRealtimeUpdate(reader)
	case InboundRealtimeCarUpdate:
		return p.handleRealtimeCarUpdate(reader)
	case InboundBroadcastingEvent:
		return p.handleBroadcastingEvent(reader)
	default:
		p.logger.Warn().Msgf("Tipo de mensaje desconocido: %d", msgType)
	}

	return nil
}

func (p *Protocol) handleRegistrationResult(reader *bytes.Reader) error {
	state, err := UnmarshalRegistrationResult(reader)
	if err != nil {
		p.logger.Error().Err(err).Msg("Error al deserializar resultado de registro")
		return err
	}

	p.connectionId = state.ConnectionId

	p.logger.Info().
		Int32("connectionId", state.ConnectionId).
		Bool("success", state.Success).
		Bool("readonly", state.IsReadonly).
		Str("error", state.ErrorMsg).
		Msg("Resultado de registro")

	if p.OnConnectionStateChanged != nil {
		p.OnConnectionStateChanged(state)
	}

	if state.Success {
		p.RequestEntryList()
		p.RequestTrackData()
	}

	return nil
}

func (p *Protocol) handleEntryList(reader *bytes.Reader) error {
	connectionId, carIndexes, err := UnmarshalEntryList(reader)
	if err != nil {
		p.logger.Error().Err(err).Msg("Error al deserializar entry list")
		return err
	}

	p.logger.Debug().
		Int32("connectionId", connectionId).
		Int("carCount", len(carIndexes)).
		Msg("Entry list recibida")

	// Limpiar el caché y preparar para recibir detalles de cada auto
	p.entryListMutex.Lock()
	p.entryListCars = make(map[uint16]*CarInfo)
	for _, carIndex := range carIndexes {
		p.entryListCars[carIndex] = &CarInfo{CarIndex: carIndex}
	}
	p.entryListMutex.Unlock()

	return nil
}

func (p *Protocol) handleEntryListCar(reader *bytes.Reader) error {
	carInfo, err := UnmarshalEntryListCar(reader)
	if err != nil {
		p.logger.Error().Err(err).Msg("Error al deserializar entry list car")
		return err
	}

	p.logger.Debug().
		Uint16("carIndex", carInfo.CarIndex).
		Str("teamName", carInfo.TeamName).
		Int32("raceNumber", carInfo.RaceNumber).
		Msg("Entry list car recibida")

	p.entryListMutex.Lock()
	p.entryListCars[carInfo.CarIndex] = &carInfo
	p.entryListMutex.Unlock()

	if p.OnEntrylistUpdate != nil {
		p.OnEntrylistUpdate(carInfo)
	}

	return nil
}

func (p *Protocol) handleTrackData(reader *bytes.Reader) error {
	connectionId, trackData, err := UnmarshalTrackData(reader)
	if err != nil {
		p.logger.Error().Err(err).Msg("Error al deserializar track data")
		return err
	}

	p.logger.Info().
		Int32("connectionId", connectionId).
		Str("trackName", trackData.TrackName).
		Int32("trackMeters", trackData.TrackMeters).
		Msg("Track data recibida")

	if p.OnTrackDataUpdate != nil {
		p.OnTrackDataUpdate(trackData)
	}

	return nil
}

func (p *Protocol) handleRealtimeUpdate(reader *bytes.Reader) error {
	update, err := UnmarshalRealtimeUpdate(reader)
	if err != nil {
		p.logger.Error().Err(err).Msg("Error al deserializar realtime update")
		return err
	}

	if p.OnRealtimeUpdate != nil {
		p.OnRealtimeUpdate(update)
	}

	return nil
}

func (p *Protocol) handleRealtimeCarUpdate(reader *bytes.Reader) error {
	update, err := UnmarshalRealtimeCarUpdate(reader)
	if err != nil {
		p.logger.Error().Err(err).Msg("Error al deserializar realtime car update")
		return err
	}

	p.entryListMutex.RLock()
	carInfo, exists := p.entryListCars[update.CarIndex]
	p.entryListMutex.RUnlock()

	if !exists || carInfo == nil || len(carInfo.Drivers) != int(update.DriverCount) {
		if time.Since(p.lastEntryListRequest) > time.Second {
			p.logger.Warn().
				Uint16("carIndex", update.CarIndex).
				Msg("Auto desconocido, solicitando nueva entry list")
			p.lastEntryListRequest = time.Now()
			p.RequestEntryList()
		}
		return nil
	}

	if p.OnRealtimeCarUpdate != nil {
		p.OnRealtimeCarUpdate(update)
	}

	return nil
}

func (p *Protocol) handleBroadcastingEvent(reader *bytes.Reader) error {
	event, err := UnmarshalBroadcastingEvent(reader)
	if err != nil {
		p.logger.Error().Err(err).Msg("Error al deserializar broadcasting event")
		return err
	}

	p.entryListMutex.RLock()
	if carInfo, exists := p.entryListCars[uint16(event.CarId)]; exists {
		event.CarData = carInfo
	}
	p.entryListMutex.RUnlock()

	if p.OnBroadcastingEvent != nil {
		p.OnBroadcastingEvent(event)
	}

	return nil
}

func (p *Protocol) GetCarInfo(carIndex uint16) (*CarInfo, bool) {
	p.entryListMutex.RLock()
	defer p.entryListMutex.RUnlock()

	carInfo, exists := p.entryListCars[carIndex]
	return carInfo, exists
}

func (p *Protocol) RequestConnection(displayName, connectionPassword string, msRealtimeUpdateInterval int32, commandPassword string) error {
	data, err := MarshalRegistrationRequest(displayName, connectionPassword, msRealtimeUpdateInterval, commandPassword)
	if err != nil {
		p.logger.Error().Err(err).Msg("Error al crear solicitud de registro")
		return err
	}

	return p.sendFunc(data)
}

func (p *Protocol) Disconnect() error {
	data, err := MarshalDisconnectRequest(p.connectionId)
	if err != nil {
		p.logger.Error().Err(err).Msg("Error al crear solicitud de desconexión")
		return err
	}

	return p.sendFunc(data)
}

func (p *Protocol) RequestEntryList() error {
	data, err := MarshalEntryListRequest(p.connectionId)
	if err != nil {
		p.logger.Error().Err(err).Msg("Error al crear solicitud de entry list")
		return err
	}

	p.logger.Debug().Int32("connectionId", p.connectionId).Msg("Solicitando entry list")
	return p.sendFunc(data)
}

func (p *Protocol) RequestTrackData() error {
	data, err := MarshalTrackDataRequest(p.connectionId)
	if err != nil {
		p.logger.Error().Err(err).Msg("Error al crear solicitud de track data")
		return err
	}

	p.logger.Debug().Int32("connectionId", p.connectionId).Msg("Solicitando track data")
	return p.sendFunc(data)
}

func (p *Protocol) SetFocus(carIndex uint16) error {
	data, err := MarshalSetFocusRequest(p.connectionId, &carIndex, nil, nil)
	if err != nil {
		p.logger.Error().Err(err).Msg("Error al crear solicitud de cambio de foco")
		return err
	}

	return p.sendFunc(data)
}

func (p *Protocol) SetCamera(cameraSet, camera string) error {
	data, err := MarshalSetFocusRequest(p.connectionId, nil, &cameraSet, &camera)
	if err != nil {
		p.logger.Error().Err(err).Msg("Error al crear solicitud de cambio de cámara")
		return err
	}

	return p.sendFunc(data)
}

func (p *Protocol) SetFocusAndCamera(carIndex uint16, cameraSet, camera string) error {
	data, err := MarshalSetFocusRequest(p.connectionId, &carIndex, &cameraSet, &camera)
	if err != nil {
		p.logger.Error().Err(err).Msg("Error al crear solicitud de cambio de foco y cámara")
		return err
	}

	return p.sendFunc(data)
}

func (p *Protocol) RequestInstantReplay(startSessionTime, durationMS float32, initialFocusedCarIndex int32, initialCameraSet, initialCamera string) error {
	data, err := MarshalInstantReplayRequest(p.connectionId, startSessionTime, durationMS, initialFocusedCarIndex, initialCameraSet, initialCamera)
	if err != nil {
		p.logger.Error().Err(err).Msg("Error al crear solicitud de instant replay")
		return err
	}

	return p.sendFunc(data)
}

func (p *Protocol) RequestHUDPage(hudPage string) error {
	data, err := MarshalHUDPageRequest(p.connectionId, hudPage)
	if err != nil {
		p.logger.Error().Err(err).Msg("Error al crear solicitud de cambio de página HUD")
		return err
	}

	return p.sendFunc(data)
}
