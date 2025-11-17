package broadcast

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"time"

	structs "RaceAll/Broadcast/Structs"
)

const BroadcastingProtocolVersion = 4

type OutboundMessageTypes byte

const (
	REGISTER_COMMAND_APPLICATION   OutboundMessageTypes = 1
	UNREGISTER_COMMAND_APPLICATION OutboundMessageTypes = 9
	REQUEST_ENTRY_LIST             OutboundMessageTypes = 10
	REQUEST_TRACK_DATA             OutboundMessageTypes = 11
	CHANGE_HUD_PAGE                OutboundMessageTypes = 49
	CHANGE_FOCUS                   OutboundMessageTypes = 50
	INSTANT_REPLAY_REQUEST         OutboundMessageTypes = 51
	PLAY_MANUAL_REPLAY_HIGHLIGHT   OutboundMessageTypes = 52 // TODO (planned)
	SAVE_MANUAL_REPLAY_HIGHLIGHT   OutboundMessageTypes = 60 // TODO (planned)
)

type InboundMessageTypes byte

const (
	REGISTRATION_RESULT InboundMessageTypes = 1
	REALTIME_UPDATE     InboundMessageTypes = 2
	REALTIME_CAR_UPDATE InboundMessageTypes = 3
	ENTRY_LIST          InboundMessageTypes = 4
	TRACK_DATA          InboundMessageTypes = 5
	ENTRY_LIST_CAR      InboundMessageTypes = 6
	BROADCASTING_EVENT  InboundMessageTypes = 7
)

type (
	ConnectionStateChangedCallback func(state structs.ConnectionState)
	TrackDataUpdateCallback        func(sender string, track *structs.TrackData)
	EntryListUpdateCallback        func(sender string, car *structs.CarInfo)
	RealtimeUpdateCallback         func(sender string, update *structs.RealtimeUpdate)
	RealtimeCarUpdateCallback      func(sender string, carUpdate *structs.RealtimeCarUpdate)
	BroadcastingEventCallback      func(sender string, evt *structs.BroadcastingEvent)
)

type BroadcastingNetworkProtocol struct {
	ConnectionIdentifier string
	Send                 func([]byte)
	ConnectionId         int
	TrackMeters          float32

	entryListCars        []*structs.CarInfo
	lastEntryListRequest time.Time

	OnConnectionStateChanged ConnectionStateChangedCallback
	OnTrackDataUpdate        TrackDataUpdateCallback
	OnEntrylistUpdate        EntryListUpdateCallback
	OnRealtimeUpdate         RealtimeUpdateCallback
	OnRealtimeCarUpdate      RealtimeCarUpdateCallback
	OnBroadcastingEvent      BroadcastingEventCallback
}

// NewBroadcastingNetworkProtocol constructs a protocol instance
func NewBroadcastingNetworkProtocol(connectionIdentifier string, send func([]byte)) (*BroadcastingNetworkProtocol, error) {
	if connectionIdentifier == "" {
		return nil, errors.New("connectionIdentifier required")
	}
	if send == nil {
		return nil, errors.New("send callback required")
	}
	return &BroadcastingNetworkProtocol{
		ConnectionIdentifier: connectionIdentifier,
		Send:                 send,
		entryListCars:        make([]*structs.CarInfo, 0),
		lastEntryListRequest: time.Now(),
	}, nil
}

// ===================== Outbound (send) helpers =====================

// RequestConnection registers this client on ACC instance
func (p *BroadcastingNetworkProtocol) RequestConnection(displayName, connectionPassword string, msRealtimeUpdateInterval int, commandPassword string) {
	buf := &bytes.Buffer{}

	buf.WriteByte(byte(REGISTER_COMMAND_APPLICATION))
	buf.WriteByte(byte(BroadcastingProtocolVersion))

	writeString(buf, displayName)
	writeString(buf, connectionPassword)

	binary.Write(buf, binary.LittleEndian, int32(msRealtimeUpdateInterval))

	writeString(buf, commandPassword)

	p.Send(buf.Bytes())
}

// Disconnect unregisters the client
func (p *BroadcastingNetworkProtocol) Disconnect() {
	buf := &bytes.Buffer{}

	buf.WriteByte(byte(UNREGISTER_COMMAND_APPLICATION))

	p.Send(buf.Bytes())
}

func (p *BroadcastingNetworkProtocol) RequestEntryList() {
	buf := &bytes.Buffer{}

	buf.WriteByte(byte(REQUEST_ENTRY_LIST))

	binary.Write(buf, binary.LittleEndian, int32(p.ConnectionId))

	p.Send(buf.Bytes())
}

func (p *BroadcastingNetworkProtocol) RequestTrackData() {
	buf := &bytes.Buffer{}

	buf.WriteByte(byte(REQUEST_TRACK_DATA))

	binary.Write(buf, binary.LittleEndian, int32(p.ConnectionId))

	p.Send(buf.Bytes())
}

// SetFocus only sets car focus
func (p *BroadcastingNetworkProtocol) SetFocus(carIndex uint16) {
	p.setFocusInternal(&carIndex, nil, nil)
}

// SetCamera only sets camera
func (p *BroadcastingNetworkProtocol) SetCamera(cameraSet, camera string) {
	p.setFocusInternal(nil, &cameraSet, &camera)
}

// SetFocusCamera sets both car and camera
func (p *BroadcastingNetworkProtocol) SetFocusCamera(carIndex uint16, cameraSet, camera string) {
	p.setFocusInternal(&carIndex, &cameraSet, &camera)
}

func (p *BroadcastingNetworkProtocol) setFocusInternal(carIndex *uint16, cameraSet, camera *string) {
	buf := &bytes.Buffer{}

	buf.WriteByte(byte(CHANGE_FOCUS))

	binary.Write(buf, binary.LittleEndian, int32(p.ConnectionId))

	if carIndex == nil {
		buf.WriteByte(0)
	} else {
		buf.WriteByte(1)
		binary.Write(buf, binary.LittleEndian, *carIndex)
	}

	if cameraSet == nil || camera == nil || *cameraSet == "" || *camera == "" {
		buf.WriteByte(0)
	} else {
		buf.WriteByte(1)
		writeString(buf, *cameraSet)
		writeString(buf, *camera)
	}

	p.Send(buf.Bytes())
}

func (p *BroadcastingNetworkProtocol) RequestInstantReplay(startSessionTimeMS, durationMS float32, initialFocusedCarIndex int32, initialCameraSet, initialCamera string) {
	buf := &bytes.Buffer{}

	buf.WriteByte(byte(INSTANT_REPLAY_REQUEST))

	binary.Write(buf, binary.LittleEndian, int32(p.ConnectionId))
	binary.Write(buf, binary.LittleEndian, startSessionTimeMS)
	binary.Write(buf, binary.LittleEndian, durationMS)
	binary.Write(buf, binary.LittleEndian, initialFocusedCarIndex)

	writeString(buf, initialCameraSet)
	writeString(buf, initialCamera)

	p.Send(buf.Bytes())
}

func (p *BroadcastingNetworkProtocol) RequestHUDPage(hudPage string) {
	buf := &bytes.Buffer{}

	buf.WriteByte(byte(CHANGE_HUD_PAGE))

	binary.Write(buf, binary.LittleEndian, int32(p.ConnectionId))
	writeString(buf, hudPage)

	p.Send(buf.Bytes())
}

// RequestData convenience (temporal solution similar to C# TODO)
func (p *BroadcastingNetworkProtocol) RequestData() {
	p.RequestEntryList()
	p.RequestTrackData()
}

// ===================== Inbound parsing =====================

// ProcessMessage parses an inbound UDP payload
func (p *BroadcastingNetworkProtocol) ProcessMessage(data []byte) error {
	reader := bytes.NewReader(data)

	typ, err := reader.ReadByte()
	if err != nil {
		return err
	}

	switch InboundMessageTypes(typ) {
	case REGISTRATION_RESULT:
		return p.handleRegistrationResult(reader)
	case ENTRY_LIST:
		return p.handleEntryList(reader)
	case ENTRY_LIST_CAR:
		return p.handleEntryListCar(reader)
	case REALTIME_UPDATE:
		return p.handleRealtimeUpdate(reader)
	case REALTIME_CAR_UPDATE:
		return p.handleRealtimeCarUpdate(reader)
	case TRACK_DATA:
		return p.handleTrackData(reader)
	case BROADCASTING_EVENT:
		return p.handleBroadcastingEvent(reader)
	default:
		return fmt.Errorf("unknown inbound message type %d", typ)
	}
}

// ---- individual handlers ----
func (p *BroadcastingNetworkProtocol) handleRegistrationResult(r *bytes.Reader) error {
	var connectionId int32
	if err := binary.Read(r, binary.LittleEndian, &connectionId); err != nil {
		return err
	}

	p.ConnectionId = int(connectionId)

	successFlag, _ := r.ReadByte()
	readonlyFlag, _ := r.ReadByte()
	errMsg, err := readString(r)

	if err != nil {
		return err
	}

	state := structs.ConnectionState{
		ConnectionId:      int(connectionId),
		ConnectionSuccess: successFlag > 0,
		IsReadonly:        readonlyFlag == 0,
		Error:             errMsg,
	}

	if p.OnConnectionStateChanged != nil {
		p.OnConnectionStateChanged(state)
	}

	// Auto-request initial data if successful
	if state.ConnectionSuccess {
		p.RequestEntryList()
		p.RequestTrackData()
	}

	return nil
}

func (p *BroadcastingNetworkProtocol) handleEntryList(r *bytes.Reader) error {
	// Clear existing cache
	p.entryListCars = p.entryListCars[:0]

	// connectionId
	var connId int32
	if err := binary.Read(r, binary.LittleEndian, &connId); err != nil {
		return err
	}

	// carEntryCount
	var count uint16
	if err := binary.Read(r, binary.LittleEndian, &count); err != nil {
		return err
	}

	for i := 0; i < int(count); i++ {
		// We only receive indices first (see C# comment). Details follow in ENTRY_LIST_CAR messages.
		var carIndex uint16
		if err := binary.Read(r, binary.LittleEndian, &carIndex); err != nil {
			return err
		}
		p.entryListCars = append(p.entryListCars, &structs.CarInfo{CarIndex: carIndex, Drivers: []structs.DriverInfo{}})
	}

	return nil
}

func (p *BroadcastingNetworkProtocol) handleEntryListCar(r *bytes.Reader) error {
	var carID uint16
	if err := binary.Read(r, binary.LittleEndian, &carID); err != nil {
		return err
	}

	// Find car entry
	var car *structs.CarInfo

	for _, c := range p.entryListCars {
		if c.CarIndex == carID {
			car = c
			break
		}
	}

	if car == nil {
		return nil
	} // unknown car; ignore

	// Car model type
	modelType, _ := r.ReadByte() // Byte sized car model
	car.CarModelType = modelType
	teamName, err := readString(r)
	if err != nil {
		return err
	}

	car.TeamName = teamName
	var raceNumber int32
	if err := binary.Read(r, binary.LittleEndian, &raceNumber); err != nil {
		return err
	}

	car.RaceNumber = int(raceNumber)

	cupCategory, _ := r.ReadByte()
	car.CupCategory = cupCategory // Cup: Overall/Pro = 0, ProAm = 1, Am = 2, Silver = 3, National = 4

	currentDriverIdx, _ := r.ReadByte()
	car.CurrentDriverIndex = int(currentDriverIdx)

	var nationality uint16
	if err := binary.Read(r, binary.LittleEndian, &nationality); err != nil {
		return err
	}
	car.Nationality = structs.NationalityEnum(nationality)
	// Drivers

	driverCountByte, _ := r.ReadByte()
	driverCount := int(driverCountByte)
	car.Drivers = car.Drivers[:0]

	for i := 0; i < driverCount; i++ {
		firstName, err := readString(r)
		if err != nil {
			return err
		}
		lastName, err := readString(r)
		if err != nil {
			return err
		}
		shortName, err := readString(r)
		if err != nil {
			return err
		}
		catByte, _ := r.ReadByte()
		var nat uint16
		if err := binary.Read(r, binary.LittleEndian, &nat); err != nil {
			return err
		}
		car.Drivers = append(car.Drivers, structs.DriverInfo{
			FirstName:   firstName,
			LastName:    lastName,
			ShortName:   shortName,
			Category:    structs.DriverCategory(catByte), // Platinum = 3, Gold = 2, Silver = 1, Bronze = 0
			Nationality: structs.NationalityEnum(nat),
		})
	}

	if p.OnEntrylistUpdate != nil {
		p.OnEntrylistUpdate(p.ConnectionIdentifier, car)
	}

	return nil
}

func (p *BroadcastingNetworkProtocol) handleRealtimeUpdate(r *bytes.Reader) error {
	upd := &structs.RealtimeUpdate{}

	var evtIdx uint16
	binary.Read(r, binary.LittleEndian, &evtIdx)
	upd.EventIndex = int(evtIdx)

	var sessIdx uint16
	binary.Read(r, binary.LittleEndian, &sessIdx)
	upd.SessionIndex = int(sessIdx)
	sessType, _ := r.ReadByte()
	upd.SessionType = structs.RaceSessionType(sessType)
	phase, _ := r.ReadByte()
	upd.Phase = structs.SessionPhase(phase)

	var sessionTimeMS float32
	binary.Read(r, binary.LittleEndian, &sessionTimeMS)
	upd.SessionTime = time.Duration(sessionTimeMS) * time.Millisecond

	var sessionEndMS float32
	binary.Read(r, binary.LittleEndian, &sessionEndMS)
	upd.SessionEndTime = time.Duration(sessionEndMS) * time.Millisecond

	var focusedCar int32
	binary.Read(r, binary.LittleEndian, &focusedCar)
	upd.FocusedCarIndex = int(focusedCar)
	camSet, _ := readString(r)
	upd.ActiveCameraSet = camSet
	cam, _ := readString(r)
	upd.ActiveCamera = cam
	hudPage, _ := readString(r)
	upd.CurrentHudPage = hudPage
	replayFlag, _ := r.ReadByte()
	upd.IsReplayPlaying = replayFlag > 0

	if upd.IsReplayPlaying {
		binary.Read(r, binary.LittleEndian, &upd.ReplaySessionTime)
		binary.Read(r, binary.LittleEndian, &upd.ReplayRemainingTime)
	}

	var timeOfDayMS float32
	binary.Read(r, binary.LittleEndian, &timeOfDayMS)
	upd.TimeOfDay = time.Duration(timeOfDayMS) * time.Millisecond
	amb, _ := r.ReadByte()
	upd.AmbientTemp = amb
	trk, _ := r.ReadByte()
	upd.TrackTemp = trk
	cloudsByte, _ := r.ReadByte()
	upd.Clouds = float32(cloudsByte) / 10.0
	rainByte, _ := r.ReadByte()
	upd.RainLevel = float32(rainByte) / 10.0
	wetnessByte, _ := r.ReadByte()
	upd.Wetness = float32(wetnessByte) / 10.0

	upd.BestSessionLap = readLap(r)

	if p.OnRealtimeUpdate != nil {
		p.OnRealtimeUpdate(p.ConnectionIdentifier, upd)
	}

	return nil
}

func (p *BroadcastingNetworkProtocol) handleRealtimeCarUpdate(r *bytes.Reader) error {
	carUpd := &structs.RealtimeCarUpdate{}

	var carIdx uint16
	binary.Read(r, binary.LittleEndian, &carIdx)
	carUpd.CarIndex = int(carIdx)

	var drvIdx uint16
	binary.Read(r, binary.LittleEndian, &drvIdx)
	carUpd.DriverIndex = int(drvIdx) // Driver swap will make this change

	gear, _ := r.ReadByte()
	carUpd.Gear = int(gear) // -2 makes the R -1, N 0 and the rest as-is

	binary.Read(r, binary.LittleEndian, &carUpd.Heading)
	binary.Read(r, binary.LittleEndian, &carUpd.WorldPosX)
	binary.Read(r, binary.LittleEndian, &carUpd.WorldPosY)

	locByte, _ := r.ReadByte()
	carUpd.CarLocation = structs.CarLocationEnum(locByte) // - , Track, Pitlane, PitEntry, PitExit = 4

	var kmh int16
	binary.Read(r, binary.LittleEndian, &kmh)
	carUpd.Kmh = int(kmh)

	var position int16
	binary.Read(r, binary.LittleEndian, &position)
	carUpd.Position = int(position) // official P/Q/R position (1 based)

	var cupPos uint16
	binary.Read(r, binary.LittleEndian, &cupPos)
	carUpd.CupPosition = cupPos // official P/Q/R position (1 based)

	var trackPos int16
	binary.Read(r, binary.LittleEndian, &trackPos)
	carUpd.TrackPosition = int(trackPos) // position on track (1 based)

	binary.Read(r, binary.LittleEndian, &carUpd.SplinePosition) // track position between 0.0 and 1.0

	var delta int32
	binary.Read(r, binary.LittleEndian, &delta)
	carUpd.Delta = int(delta) // Realtime delta to best session lap

	var laps int32
	binary.Read(r, binary.LittleEndian, &laps)
	carUpd.Laps = int(laps)

	// Parse laps
	carUpd.BestSessionLap = readLap(r)
	carUpd.LastLap = readLap(r)
	carUpd.CurrentLap = readLap(r)

	driverCount, _ := r.ReadByte()
	carUpd.DriverCount = driverCount

	// Entry list desync detection
	var carEntry *structs.CarInfo
	for _, c := range p.entryListCars {
		if c.CarIndex == uint16(carUpd.CarIndex) {
			carEntry = c
			break
		}
	}

	if carEntry == nil || len(carEntry.Drivers) != int(carUpd.DriverCount) {
		if time.Since(p.lastEntryListRequest).Seconds() > 1 {
			p.lastEntryListRequest = time.Now()
			p.RequestEntryList()
		}
	}

	if p.OnRealtimeCarUpdate != nil {
		p.OnRealtimeCarUpdate(p.ConnectionIdentifier, carUpd)
	}

	return nil
}

func (p *BroadcastingNetworkProtocol) handleTrackData(r *bytes.Reader) error {

	trk := &structs.TrackData{}
	name, err := readString(r)
	if err != nil {
		return err
	}

	trk.TrackName = name
	var trackID int32

	if err := binary.Read(r, binary.LittleEndian, &trackID); err != nil {
		return err
	}

	trk.TrackId = int(trackID)

	binary.Read(r, binary.LittleEndian, &trk.TrackMeters)

	p.TrackMeters = trk.TrackMeters // maybe TrackMeters = trackData.TrackMeters > 0 ? trackData.TrackMeters : -1;

	// Camera sets: dictionary<string, list<string>>
	var camSetCount uint16
	if err := binary.Read(r, binary.LittleEndian, &camSetCount); err != nil {
		return err
	}

	trk.CameraSets = make(map[string][]string, camSetCount)

	for i := 0; i < int(camSetCount); i++ {
		camSetName, err := readString(r)
		if err != nil {
			return err
		}

		var camCount uint16
		if err := binary.Read(r, binary.LittleEndian, &camCount); err != nil {
			return err
		}

		cameras := make([]string, 0, camCount)

		for j := 0; j < int(camCount); j++ {
			cam, err := readString(r)
			if err != nil {
				return err
			}
			cameras = append(cameras, cam)
		}
		trk.CameraSets[camSetName] = cameras
	}
	// HUD pages
	var hudCount uint16
	if err := binary.Read(r, binary.LittleEndian, &hudCount); err != nil {
		return err
	}

	pages := make([]string, 0, hudCount)

	for i := 0; i < int(hudCount); i++ {
		page, err := readString(r)
		if err != nil {
			return err
		}
		pages = append(pages, page)
	}

	trk.HUDPages = pages

	if p.OnTrackDataUpdate != nil {
		p.OnTrackDataUpdate(p.ConnectionIdentifier, trk)
	}

	return nil
}

func (p *BroadcastingNetworkProtocol) handleBroadcastingEvent(r *bytes.Reader) error {
	evt := &structs.BroadcastingEvent{}

	typ, _ := r.ReadByte()

	evt.Type = structs.BroadcastingCarEventType(typ)

	msg, err := readString(r)
	if err != nil {
		return err
	}

	evt.Msg = msg

	var timeMS int32
	binary.Read(r, binary.LittleEndian, &timeMS)
	evt.TimeMs = int(timeMS)

	var carId int32
	binary.Read(r, binary.LittleEndian, &carId)
	evt.CarId = int(carId)

	// CarData may not be sent here in protocol; TODO if available
	if p.OnBroadcastingEvent != nil {
		p.OnBroadcastingEvent(p.ConnectionIdentifier, evt)
	}

	return nil
}


func readString(r *bytes.Reader) (string, error) {
	var length uint16
	if err := binary.Read(r, binary.LittleEndian, &length); err != nil {
		return "", err
	}

	bytesArr := make([]byte, length)
	if _, err := io.ReadFull(r, bytesArr); err != nil {
		return "", err
	}

	return string(bytesArr), nil
}

func writeString(buf *bytes.Buffer, s string) {
	binary.Write(buf, binary.LittleEndian, uint16(len(s)))
	buf.WriteString(s)
}

// readLap translates the C# ReadLap method
func readLap(r *bytes.Reader) *structs.LapInfo {

	lap := structs.NewLapInfo()

	var lapTime int32
	binary.Read(r, binary.LittleEndian, &lapTime)
	if lapTime != math.MaxInt32 {
		lapTimeInt := int(lapTime)
		lap.LaptimeMS = &lapTimeInt
	}

	var carIdx uint16
	binary.Read(r, binary.LittleEndian, &carIdx)
	lap.CarIndex = carIdx

	var drvIdx uint16
	binary.Read(r, binary.LittleEndian, &drvIdx)
	lap.DriverIndex = drvIdx

	splitCount, _ := r.ReadByte()
	for i := 0; i < int(splitCount); i++ {
		var split int32
		binary.Read(r, binary.LittleEndian, &split)
		if split != math.MaxInt32 {
			lap.Splits = append(lap.Splits, int(split))
		} else {
			lap.Splits = append(lap.Splits, 0)
		}
	}

	lap.IsInvalid = readBool(r)

	lap.IsValidForBest = readBool(r)

	isOut := readBool(r)

	isIn := readBool(r)

	if isOut {
		lap.Type = structs.Outlap
	} else if isIn {
		lap.Type = structs.Inlap
	} else {
		lap.Type = structs.Regular
	}

	for len(lap.Splits) < 3 {
		lap.Splits = append(lap.Splits, 0)
	}

	return lap
}

func readBool(r *bytes.Reader) bool {
	b, _ := r.ReadByte()
	return b > 0
}
