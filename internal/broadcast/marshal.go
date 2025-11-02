package broadcast

import (
	"encoding/binary"
	"io"
	"math"
)

const (
	InvalidLapTime    int32 = math.MaxInt32
	InvalidSectorTime int32 = math.MaxInt32
)

func writeUint8(w io.Writer, value uint8) error {
	return binary.Write(w, binary.LittleEndian, value)
}

func writeInt32(w io.Writer, value int32) error {
	return binary.Write(w, binary.LittleEndian, value)
}

func writeString(w io.Writer, s string) error {
	length := uint16(len(s))
	if err := binary.Write(w, binary.LittleEndian, length); err != nil {
		return err
	}
	_, err := w.Write([]byte(s))
	return err
}

func readUint8(r io.Reader) (uint8, error) {
	var value uint8
	err := binary.Read(r, binary.LittleEndian, &value)
	return value, err
}

func readInt8(r io.Reader) (int8, error) {
	var value int8
	err := binary.Read(r, binary.LittleEndian, &value)
	return value, err
}

func readUint16(r io.Reader) (uint16, error) {
	var value uint16
	err := binary.Read(r, binary.LittleEndian, &value)
	return value, err
}

func readInt32(r io.Reader) (int32, error) {
	var value int32
	err := binary.Read(r, binary.LittleEndian, &value)
	return value, err
}

func readFloat32(r io.Reader) (float32, error) {
	var value float32
	err := binary.Read(r, binary.LittleEndian, &value)
	return value, err
}

func readString(r io.Reader) (string, error) {
	var length uint16
	if err := binary.Read(r, binary.LittleEndian, &length); err != nil {
		return "", err
	}

	buffer := make([]byte, length)
	if err := binary.Read(r, binary.LittleEndian, &buffer); err != nil {
		return "", err
	}

	return string(buffer), nil
}

func readLap(r io.Reader) (LapInfo, error) {
	lap := LapInfo{}

	lapTimeMs, err := readInt32(r)
	if err != nil {
		return lap, err
	}

	if lapTimeMs == InvalidLapTime {
		lap.LaptimeMS = nil
	} else {
		lap.LaptimeMS = &lapTimeMs
	}

	carIndex, err := readUint16(r)
	if err != nil {
		return lap, err
	}
	lap.CarIndex = carIndex

	driverIndex, err := readUint16(r)
	if err != nil {
		return lap, err
	}
	lap.DriverIndex = driverIndex

	splitCount, err := readUint8(r)
	if err != nil {
		return lap, err
	}

	for i := uint8(0); i < splitCount && i < 3; i++ {
		splitTime, err := readInt32(r)
		if err != nil {
			return lap, err
		}

		if splitTime == InvalidSectorTime {
			lap.Splits[i] = nil
		} else {
			lap.Splits[i] = &splitTime
		}
	}

	isInvalid, err := readUint8(r)
	if err != nil {
		return lap, err
	}
	lap.IsInvalid = isInvalid > 0

	isValidForBest, err := readUint8(r)
	if err != nil {
		return lap, err
	}
	lap.IsValidForBest = isValidForBest > 0

	isOutlap, err := readUint8(r)
	if err != nil {
		return lap, err
	}

	isInlap, err := readUint8(r)
	if err != nil {
		return lap, err
	}

	if isOutlap > 0 {
		lap.Type = LapTypeOutlap
	} else if isInlap > 0 {
		lap.Type = LapTypeInlap
	} else {
		lap.Type = LapTypeRegular
	}

	return lap, nil
}

func MarshalRegistrationRequest(displayName, connectionPassword string, msRealtimeUpdateInterval int32, commandPassword string) ([]byte, error) {
	buffer := GetBuffer()
	defer PutBuffer(buffer)

	if err := writeUint8(buffer, byte(OutboundRegisterCommandApplication)); err != nil {
		return nil, NewError("MarshalRegistrationRequest", err)
	}
	if err := writeUint8(buffer, BroadcastingProtocolVersion); err != nil {
		return nil, NewError("MarshalRegistrationRequest", err)
	}
	if err := writeString(buffer, displayName); err != nil {
		return nil, NewError("MarshalRegistrationRequest", err)
	}
	if err := writeString(buffer, connectionPassword); err != nil {
		return nil, NewError("MarshalRegistrationRequest", err)
	}
	if err := writeInt32(buffer, msRealtimeUpdateInterval); err != nil {
		return nil, NewError("MarshalRegistrationRequest", err)
	}
	if err := writeString(buffer, commandPassword); err != nil {
		return nil, NewError("MarshalRegistrationRequest", err)
	}

	result := make([]byte, buffer.Len())
	copy(result, buffer.Bytes())
	return result, nil
}

func MarshalDisconnectRequest(connectionId int32) ([]byte, error) {
	buffer := GetBuffer()
	defer PutBuffer(buffer)

	if err := writeUint8(buffer, byte(OutboundUnregisterCommandApplication)); err != nil {
		return nil, NewError("MarshalDisconnectRequest", err)
	}
	if err := writeInt32(buffer, connectionId); err != nil {
		return nil, NewError("MarshalDisconnectRequest", err)
	}

	result := make([]byte, buffer.Len())
	copy(result, buffer.Bytes())
	return result, nil
}

func MarshalEntryListRequest(connectionId int32) ([]byte, error) {
	buffer := GetBuffer()
	defer PutBuffer(buffer)

	if err := writeUint8(buffer, byte(OutboundRequestEntryList)); err != nil {
		return nil, NewError("MarshalEntryListRequest", err)
	}
	if err := writeInt32(buffer, connectionId); err != nil {
		return nil, NewError("MarshalEntryListRequest", err)
	}

	result := make([]byte, buffer.Len())
	copy(result, buffer.Bytes())
	return result, nil
}

func MarshalTrackDataRequest(connectionId int32) ([]byte, error) {
	buffer := GetBuffer()
	defer PutBuffer(buffer)

	if err := writeUint8(buffer, byte(OutboundRequestTrackData)); err != nil {
		return nil, NewError("MarshalTrackDataRequest", err)
	}
	if err := writeInt32(buffer, connectionId); err != nil {
		return nil, NewError("MarshalTrackDataRequest", err)
	}

	result := make([]byte, buffer.Len())
	copy(result, buffer.Bytes())
	return result, nil
}

func MarshalSetFocusRequest(connectionId int32, carIndex *uint16, cameraSet, camera *string) ([]byte, error) {
	buffer := GetBuffer()
	defer PutBuffer(buffer)

	if err := writeUint8(buffer, byte(OutboundChangeFocus)); err != nil {
		return nil, NewError("MarshalSetFocusRequest", err)
	}
	if err := writeInt32(buffer, connectionId); err != nil {
		return nil, NewError("MarshalSetFocusRequest", err)
	}

	if carIndex != nil {
		if err := writeUint8(buffer, 1); err != nil {
			return nil, NewError("MarshalSetFocusRequest", err)
		}
		if err := binary.Write(buffer, binary.LittleEndian, *carIndex); err != nil {
			return nil, NewError("MarshalSetFocusRequest", err)
		}
	} else {
		if err := writeUint8(buffer, 0); err != nil {
			return nil, NewError("MarshalSetFocusRequest", err)
		}
	}

	if cameraSet != nil && camera != nil {
		if err := writeUint8(buffer, 1); err != nil {
			return nil, NewError("MarshalSetFocusRequest", err)
		}
		if err := writeString(buffer, *cameraSet); err != nil {
			return nil, NewError("MarshalSetFocusRequest", err)
		}
		if err := writeString(buffer, *camera); err != nil {
			return nil, NewError("MarshalSetFocusRequest", err)
		}
	} else {
		if err := writeUint8(buffer, 0); err != nil {
			return nil, NewError("MarshalSetFocusRequest", err)
		}
	}

	result := make([]byte, buffer.Len())
	copy(result, buffer.Bytes())
	return result, nil
}

func MarshalInstantReplayRequest(connectionId int32, startSessionTime, durationMS float32, initialFocusedCarIndex int32, initialCameraSet, initialCamera string) ([]byte, error) {
	buffer := GetBuffer()
	defer PutBuffer(buffer)

	if err := writeUint8(buffer, byte(OutboundInstantReplayRequest)); err != nil {
		return nil, NewError("MarshalInstantReplayRequest", err)
	}
	if err := writeInt32(buffer, connectionId); err != nil {
		return nil, NewError("MarshalInstantReplayRequest", err)
	}
	if err := binary.Write(buffer, binary.LittleEndian, startSessionTime); err != nil {
		return nil, NewError("MarshalInstantReplayRequest", err)
	}
	if err := binary.Write(buffer, binary.LittleEndian, durationMS); err != nil {
		return nil, NewError("MarshalInstantReplayRequest", err)
	}
	if err := writeInt32(buffer, initialFocusedCarIndex); err != nil {
		return nil, NewError("MarshalInstantReplayRequest", err)
	}
	if err := writeString(buffer, initialCameraSet); err != nil {
		return nil, NewError("MarshalInstantReplayRequest", err)
	}
	if err := writeString(buffer, initialCamera); err != nil {
		return nil, NewError("MarshalInstantReplayRequest", err)
	}

	result := make([]byte, buffer.Len())
	copy(result, buffer.Bytes())
	return result, nil
}

func MarshalHUDPageRequest(connectionId int32, hudPage string) ([]byte, error) {
	buffer := GetBuffer()
	defer PutBuffer(buffer)

	if err := writeUint8(buffer, byte(OutboundChangeHUDPage)); err != nil {
		return nil, NewError("MarshalHUDPageRequest", err)
	}
	if err := writeInt32(buffer, connectionId); err != nil {
		return nil, NewError("MarshalHUDPageRequest", err)
	}
	if err := writeString(buffer, hudPage); err != nil {
		return nil, NewError("MarshalHUDPageRequest", err)
	}

	result := make([]byte, buffer.Len())
	copy(result, buffer.Bytes())
	return result, nil
}
