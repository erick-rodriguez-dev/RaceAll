package broadcast

import (
	"fmt"
	"io"
	"time"
)

func UnmarshalRegistrationResult(r io.Reader) (ConnectionState, error) {
	var state ConnectionState

	connectionId, err := readInt32(r)
	if err != nil {
		return state, NewError("UnmarshalRegistrationResult", fmt.Errorf("failed to read connectionId: %w", err))
	}
	state.ConnectionId = connectionId

	success, err := readUint8(r)
	if err != nil {
		return state, NewError("UnmarshalRegistrationResult", fmt.Errorf("failed to read success: %w", err))
	}
	state.Success = success > 0

	isReadonly, err := readUint8(r)
	if err != nil {
		return state, NewError("UnmarshalRegistrationResult", fmt.Errorf("failed to read isReadonly: %w", err))
	}
	state.IsReadonly = isReadonly == 0

	errMsg, err := readString(r)
	if err != nil {
		return state, NewError("UnmarshalRegistrationResult", fmt.Errorf("failed to read errorMsg: %w", err))
	}
	state.ErrorMsg = errMsg

	return state, nil
}

func UnmarshalEntryList(r io.Reader) (int32, []uint16, error) {
	connectionId, err := readInt32(r)
	if err != nil {
		return 0, nil, NewError("UnmarshalEntryList", fmt.Errorf("failed to read connectionId: %w", err))
	}

	carCount, err := readUint16(r)
	if err != nil {
		return connectionId, nil, NewError("UnmarshalEntryList", fmt.Errorf("failed to read carCount: %w", err))
	}

	if carCount > 200 {
		return connectionId, nil, NewError("UnmarshalEntryList",
			NewValidationError("carCount", carCount, "exceeds maximum of 200"))
	}

	carIndexes := make([]uint16, carCount)
	for i := uint16(0); i < carCount; i++ {
		carIndex, err := readUint16(r)
		if err != nil {
			return connectionId, nil, NewError("UnmarshalEntryList",
				fmt.Errorf("failed to read carIndex at position %d: %w", i, err))
		}

		if err := ValidateCarIndex(carIndex); err != nil {
			return connectionId, nil, NewError("UnmarshalEntryList", err)
		}

		carIndexes[i] = carIndex
	}

	return connectionId, carIndexes, nil
}

func UnmarshalEntryListCar(r io.Reader) (CarInfo, error) {
	var car CarInfo

	carIndex, err := readUint16(r)
	if err != nil {
		return car, NewError("UnmarshalEntryListCar", fmt.Errorf("failed to read carIndex: %w", err))
	}

	if err := ValidateCarIndex(carIndex); err != nil {
		return car, NewError("UnmarshalEntryListCar", err)
	}
	car.CarIndex = carIndex

	carModelType, err := readUint8(r)
	if err != nil {
		return car, NewError("UnmarshalEntryListCar", fmt.Errorf("failed to read carModelType: %w", err))
	}
	car.CarModelType = carModelType

	teamName, err := readString(r)
	if err != nil {
		return car, NewError("UnmarshalEntryListCar", fmt.Errorf("failed to read teamName: %w", err))
	}
	car.TeamName = teamName

	raceNumber, err := readInt32(r)
	if err != nil {
		return car, NewError("UnmarshalEntryListCar", fmt.Errorf("failed to read raceNumber: %w", err))
	}
	car.RaceNumber = raceNumber

	cupCategory, err := readUint8(r)
	if err != nil {
		return car, NewError("UnmarshalEntryListCar", fmt.Errorf("failed to read cupCategory: %w", err))
	}

	if cupCategory > 4 {
		return car, NewError("UnmarshalEntryListCar",
			NewValidationError("cupCategory", cupCategory, "must be between 0 and 4"))
	}
	car.CupCategory = cupCategory

	currentDriverIndex, err := readUint8(r)
	if err != nil {
		return car, NewError("UnmarshalEntryListCar", fmt.Errorf("failed to read currentDriverIndex: %w", err))
	}
	car.CurrentDriverIndex = currentDriverIndex

	nationality, err := readUint16(r)
	if err != nil {
		return car, NewError("UnmarshalEntryListCar", fmt.Errorf("failed to read nationality: %w", err))
	}
	car.Nationality = NationalityEnum(nationality)

	driverCount, err := readUint8(r)
	if err != nil {
		return car, NewError("UnmarshalEntryListCar", fmt.Errorf("failed to read driverCount: %w", err))
	}

	if driverCount > 10 {
		return car, NewError("UnmarshalEntryListCar",
			NewValidationError("driverCount", driverCount, "exceeds maximum of 10"))
	}

	if currentDriverIndex >= driverCount {
		return car, NewError("UnmarshalEntryListCar",
			NewValidationError("currentDriverIndex", currentDriverIndex,
				fmt.Sprintf("must be less than driverCount (%d)", driverCount)))
	}

	car.Drivers = make([]DriverInfo, driverCount)
	for i := uint8(0); i < driverCount; i++ {
		firstName, err := readString(r)
		if err != nil {
			return car, NewError("UnmarshalEntryListCar",
				fmt.Errorf("failed to read driver firstName at index %d: %w", i, err))
		}

		lastName, err := readString(r)
		if err != nil {
			return car, NewError("UnmarshalEntryListCar",
				fmt.Errorf("failed to read driver lastName at index %d: %w", i, err))
		}

		shortName, err := readString(r)
		if err != nil {
			return car, NewError("UnmarshalEntryListCar",
				fmt.Errorf("failed to read driver shortName at index %d: %w", i, err))
		}

		category, err := readUint8(r)
		if err != nil {
			return car, NewError("UnmarshalEntryListCar",
				fmt.Errorf("failed to read driver category at index %d: %w", i, err))
		}

		if err := ValidateDriverCategory(DriverCategory(category)); err != nil {
			return car, NewError("UnmarshalEntryListCar",
				fmt.Errorf("invalid driver category at index %d: %w", i, err))
		}

		driverNationality, err := readUint16(r)
		if err != nil {
			return car, NewError("UnmarshalEntryListCar",
				fmt.Errorf("failed to read driver nationality at index %d: %w", i, err))
		}

		car.Drivers[i] = DriverInfo{
			FirstName:   firstName,
			LastName:    lastName,
			ShortName:   shortName,
			Category:    DriverCategory(category),
			Nationality: NationalityEnum(driverNationality),
		}
	}

	return car, nil
}

func UnmarshalTrackData(r io.Reader) (int32, TrackData, error) {
	var trackData TrackData

	connectionId, err := readInt32(r)
	if err != nil {
		return 0, trackData, NewError("UnmarshalTrackData", fmt.Errorf("failed to read connectionId: %w", err))
	}

	trackName, err := readString(r)
	if err != nil {
		return connectionId, trackData, NewError("UnmarshalTrackData", fmt.Errorf("failed to read trackName: %w", err))
	}
	trackData.TrackName = trackName

	trackId, err := readInt32(r)
	if err != nil {
		return connectionId, trackData, NewError("UnmarshalTrackData", fmt.Errorf("failed to read trackId: %w", err))
	}
	trackData.TrackId = trackId

	trackMeters, err := readInt32(r)
	if err != nil {
		return connectionId, trackData, NewError("UnmarshalTrackData", fmt.Errorf("failed to read trackMeters: %w", err))
	}

	if trackMeters < 1000 || trackMeters > 25000 {
		return connectionId, trackData, NewError("UnmarshalTrackData",
			NewValidationError("trackMeters", trackMeters, "must be between 1000 and 25000"))
	}
	trackData.TrackMeters = trackMeters

	trackData.CameraSets = make(map[string][]string)

	cameraSetCount, err := readUint8(r)
	if err != nil {
		return connectionId, trackData, NewError("UnmarshalTrackData", fmt.Errorf("failed to read cameraSetCount: %w", err))
	}

	if cameraSetCount > 50 {
		return connectionId, trackData, NewError("UnmarshalTrackData",
			NewValidationError("cameraSetCount", cameraSetCount, "exceeds maximum of 50"))
	}

	for i := uint8(0); i < cameraSetCount; i++ {
		cameraSetName, err := readString(r)
		if err != nil {
			return connectionId, trackData, NewError("UnmarshalTrackData",
				fmt.Errorf("failed to read cameraSetName at index %d: %w", i, err))
		}

		cameraCount, err := readUint8(r)
		if err != nil {
			return connectionId, trackData, NewError("UnmarshalTrackData",
				fmt.Errorf("failed to read cameraCount at index %d: %w", i, err))
		}

		if cameraCount > 100 {
			return connectionId, trackData, NewError("UnmarshalTrackData",
				NewValidationError("cameraCount", cameraCount, "exceeds maximum of 100"))
		}

		cameras := make([]string, cameraCount)
		for j := uint8(0); j < cameraCount; j++ {
			cameraName, err := readString(r)
			if err != nil {
				return connectionId, trackData, NewError("UnmarshalTrackData",
					fmt.Errorf("failed to read cameraName at set %d, camera %d: %w", i, j, err))
			}
			cameras[j] = cameraName
		}

		trackData.CameraSets[cameraSetName] = cameras
	}

	hudPageCount, err := readUint8(r)
	if err != nil {
		return connectionId, trackData, NewError("UnmarshalTrackData", fmt.Errorf("failed to read hudPageCount: %w", err))
	}

	if hudPageCount > 50 {
		return connectionId, trackData, NewError("UnmarshalTrackData",
			NewValidationError("hudPageCount", hudPageCount, "exceeds maximum of 50"))
	}

	trackData.HUDPages = make([]string, hudPageCount)
	for i := uint8(0); i < hudPageCount; i++ {
		hudPage, err := readString(r)
		if err != nil {
			return connectionId, trackData, NewError("UnmarshalTrackData",
				fmt.Errorf("failed to read hudPage at index %d: %w", i, err))
		}
		trackData.HUDPages[i] = hudPage
	}

	return connectionId, trackData, nil
}

func UnmarshalRealtimeUpdate(r io.Reader) (RealtimeUpdate, error) {
	var update RealtimeUpdate

	eventIndex, err := readUint16(r)
	if err != nil {
		return update, NewError("UnmarshalRealtimeUpdate", fmt.Errorf("failed to read eventIndex: %w", err))
	}
	update.EventIndex = eventIndex

	sessionIndex, err := readUint16(r)
	if err != nil {
		return update, NewError("UnmarshalRealtimeUpdate", fmt.Errorf("failed to read sessionIndex: %w", err))
	}
	update.SessionIndex = sessionIndex

	sessionType, err := readUint8(r)
	if err != nil {
		return update, NewError("UnmarshalRealtimeUpdate", fmt.Errorf("failed to read sessionType: %w", err))
	}
	update.SessionType = RaceSessionType(sessionType)

	phase, err := readUint8(r)
	if err != nil {
		return update, NewError("UnmarshalRealtimeUpdate", fmt.Errorf("failed to read phase: %w", err))
	}
	update.Phase = SessionPhase(phase)

	sessionTime, err := readFloat32(r)
	if err != nil {
		return update, NewError("UnmarshalRealtimeUpdate", fmt.Errorf("failed to read sessionTime: %w", err))
	}
	update.SessionTime = time.Duration(sessionTime) * time.Millisecond

	sessionEndTime, err := readFloat32(r)
	if err != nil {
		return update, NewError("UnmarshalRealtimeUpdate", fmt.Errorf("failed to read sessionEndTime: %w", err))
	}
	update.SessionEndTime = time.Duration(sessionEndTime) * time.Millisecond

	focusedCarIndex, err := readInt32(r)
	if err != nil {
		return update, NewError("UnmarshalRealtimeUpdate", fmt.Errorf("failed to read focusedCarIndex: %w", err))
	}
	update.FocusedCarIndex = focusedCarIndex

	activeCameraSet, err := readString(r)
	if err != nil {
		return update, NewError("UnmarshalRealtimeUpdate", fmt.Errorf("failed to read activeCameraSet: %w", err))
	}
	update.ActiveCameraSet = activeCameraSet

	activeCamera, err := readString(r)
	if err != nil {
		return update, NewError("UnmarshalRealtimeUpdate", fmt.Errorf("failed to read activeCamera: %w", err))
	}
	update.ActiveCamera = activeCamera

	currentHudPage, err := readString(r)
	if err != nil {
		return update, NewError("UnmarshalRealtimeUpdate", fmt.Errorf("failed to read currentHudPage: %w", err))
	}
	update.CurrentHudPage = currentHudPage

	isReplayPlaying, err := readUint8(r)
	if err != nil {
		return update, NewError("UnmarshalRealtimeUpdate", fmt.Errorf("failed to read isReplayPlaying: %w", err))
	}
	update.IsReplayPlaying = isReplayPlaying > 0

	if update.IsReplayPlaying {
		replaySessionTime, err := readFloat32(r)
		if err != nil {
			return update, NewError("UnmarshalRealtimeUpdate", fmt.Errorf("failed to read replaySessionTime: %w", err))
		}
		update.ReplaySessionTime = replaySessionTime

		replayRemainingTime, err := readFloat32(r)
		if err != nil {
			return update, NewError("UnmarshalRealtimeUpdate", fmt.Errorf("failed to read replayRemainingTime: %w", err))
		}
		update.ReplayRemainingTime = replayRemainingTime
	}

	timeOfDay, err := readFloat32(r)
	if err != nil {
		return update, NewError("UnmarshalRealtimeUpdate", fmt.Errorf("failed to read timeOfDay: %w", err))
	}
	update.TimeOfDay = time.Duration(timeOfDay) * time.Millisecond

	ambientTemp, err := readUint8(r)
	if err != nil {
		return update, NewError("UnmarshalRealtimeUpdate", fmt.Errorf("failed to read ambientTemp: %w", err))
	}
	update.AmbientTemp = ambientTemp

	trackTemp, err := readUint8(r)
	if err != nil {
		return update, NewError("UnmarshalRealtimeUpdate", fmt.Errorf("failed to read trackTemp: %w", err))
	}
	update.TrackTemp = trackTemp

	clouds, err := readUint8(r)
	if err != nil {
		return update, NewError("UnmarshalRealtimeUpdate", fmt.Errorf("failed to read clouds: %w", err))
	}
	update.Clouds = float32(clouds) / 10.0

	rainLevel, err := readUint8(r)
	if err != nil {
		return update, NewError("UnmarshalRealtimeUpdate", fmt.Errorf("failed to read rainLevel: %w", err))
	}
	update.RainLevel = float32(rainLevel) / 10.0

	wetness, err := readUint8(r)
	if err != nil {
		return update, NewError("UnmarshalRealtimeUpdate", fmt.Errorf("failed to read wetness: %w", err))
	}
	update.Wetness = float32(wetness) / 10.0

	bestSessionLap, err := readLap(r)
	if err != nil {
		return update, NewError("UnmarshalRealtimeUpdate", fmt.Errorf("failed to read bestSessionLap: %w", err))
	}
	update.BestSessionLap = bestSessionLap

	return update, nil
}

func UnmarshalRealtimeCarUpdate(r io.Reader) (RealtimeCarUpdate, error) {
	var update RealtimeCarUpdate

	carIndex, err := readUint16(r)
	if err != nil {
		return update, NewError("UnmarshalRealtimeCarUpdate", fmt.Errorf("failed to read carIndex: %w", err))
	}
	if err := ValidateCarIndex(carIndex); err != nil {
		return update, NewError("UnmarshalRealtimeCarUpdate", err)
	}
	update.CarIndex = carIndex

	driverIndex, err := readUint16(r)
	if err != nil {
		return update, NewError("UnmarshalRealtimeCarUpdate", fmt.Errorf("failed to read driverIndex: %w", err))
	}
	update.DriverIndex = driverIndex

	driverCount, err := readUint8(r)
	if err != nil {
		return update, NewError("UnmarshalRealtimeCarUpdate", fmt.Errorf("failed to read driverCount: %w", err))
	}
	update.DriverCount = driverCount

	gear, err := readUint8(r)
	if err != nil {
		return update, NewError("UnmarshalRealtimeCarUpdate", fmt.Errorf("failed to read gear: %w", err))
	}
	update.Gear = int8(gear) - 2

	worldPosX, err := readFloat32(r)
	if err != nil {
		return update, NewError("UnmarshalRealtimeCarUpdate", fmt.Errorf("failed to read worldPosX: %w", err))
	}
	update.WorldPosX = worldPosX

	worldPosY, err := readFloat32(r)
	if err != nil {
		return update, NewError("UnmarshalRealtimeCarUpdate", fmt.Errorf("failed to read worldPosY: %w", err))
	}
	update.WorldPosY = worldPosY

	heading, err := readFloat32(r)
	if err != nil {
		return update, NewError("UnmarshalRealtimeCarUpdate", fmt.Errorf("failed to read heading: %w", err))
	}
	update.Heading = heading

	carLocation, err := readUint8(r)
	if err != nil {
		return update, NewError("UnmarshalRealtimeCarUpdate", fmt.Errorf("failed to read carLocation: %w", err))
	}
	update.CarLocation = CarLocationEnum(carLocation)

	kmh, err := readUint16(r)
	if err != nil {
		return update, NewError("UnmarshalRealtimeCarUpdate", fmt.Errorf("failed to read kmh: %w", err))
	}
	update.Kmh = kmh

	position, err := readUint16(r)
	if err != nil {
		return update, NewError("UnmarshalRealtimeCarUpdate", fmt.Errorf("failed to read position: %w", err))
	}
	update.Position = position

	cupPosition, err := readUint16(r)
	if err != nil {
		return update, NewError("UnmarshalRealtimeCarUpdate", fmt.Errorf("failed to read cupPosition: %w", err))
	}
	update.CupPosition = cupPosition

	trackPosition, err := readUint16(r)
	if err != nil {
		return update, NewError("UnmarshalRealtimeCarUpdate", fmt.Errorf("failed to read trackPosition: %w", err))
	}
	update.TrackPosition = trackPosition

	splinePosition, err := readFloat32(r)
	if err != nil {
		return update, NewError("UnmarshalRealtimeCarUpdate", fmt.Errorf("failed to read splinePosition: %w", err))
	}
	update.SplinePosition = splinePosition

	laps, err := readUint16(r)
	if err != nil {
		return update, NewError("UnmarshalRealtimeCarUpdate", fmt.Errorf("failed to read laps: %w", err))
	}
	update.Laps = laps

	delta, err := readInt32(r)
	if err != nil {
		return update, NewError("UnmarshalRealtimeCarUpdate", fmt.Errorf("failed to read delta: %w", err))
	}
	update.Delta = delta

	bestSessionLap, err := readLap(r)
	if err != nil {
		return update, NewError("UnmarshalRealtimeCarUpdate", fmt.Errorf("failed to read bestSessionLap: %w", err))
	}
	update.BestSessionLap = bestSessionLap

	lastLap, err := readLap(r)
	if err != nil {
		return update, NewError("UnmarshalRealtimeCarUpdate", fmt.Errorf("failed to read lastLap: %w", err))
	}
	update.LastLap = lastLap

	currentLap, err := readLap(r)
	if err != nil {
		return update, NewError("UnmarshalRealtimeCarUpdate", fmt.Errorf("failed to read currentLap: %w", err))
	}
	update.CurrentLap = currentLap

	return update, nil
}

func UnmarshalBroadcastingEvent(r io.Reader) (BroadcastingEvent, error) {
	var event BroadcastingEvent

	eventType, err := readUint8(r)
	if err != nil {
		return event, NewError("UnmarshalBroadcastingEvent", fmt.Errorf("failed to read eventType: %w", err))
	}
	event.Type = BroadcastingEventType(eventType)

	msg, err := readString(r)
	if err != nil {
		return event, NewError("UnmarshalBroadcastingEvent", fmt.Errorf("failed to read msg: %w", err))
	}
	event.Msg = msg

	timeMs, err := readInt32(r)
	if err != nil {
		return event, NewError("UnmarshalBroadcastingEvent", fmt.Errorf("failed to read timeMs: %w", err))
	}
	event.TimeMs = timeMs

	carId, err := readInt32(r)
	if err != nil {
		return event, NewError("UnmarshalBroadcastingEvent", fmt.Errorf("failed to read carId: %w", err))
	}
	event.CarId = carId

	return event, nil
}
