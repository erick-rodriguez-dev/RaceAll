package integration_test

import (
	"os"
	"testing"
	"time"

	"RaceAll/internal/broadcast"
	"RaceAll/internal/logger"
)

func TestBroadcastConnection(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	accAddress := os.Getenv("ACC_SERVER_ADDRESS")
	if accAddress == "" {
		accAddress = "127.0.0.1:9000"
		t.Logf("Using default ACC address: %s", accAddress)
	}

	logger.Init(logger.Config{
		Level:  logger.LevelDebug,
		Pretty: true,
		Output: os.Stdout,
	})

	client := broadcast.NewClient(
		accAddress,
		"Integration Test Client",
		"asd",
		"",
		250,
		*logger.Get(),
	)

	connectionReceived := false
	trackDataReceived := false
	entryListReceived := false

	client.OnConnectionStateChanged = func(state broadcast.ConnectionState) {
		t.Logf("Connection state: success=%v, connectionId=%d, readonly=%v, error=%s",
			state.Success, state.ConnectionId, state.IsReadonly, state.ErrorMsg)

		if state.Success {
			connectionReceived = true
		} else {
			t.Logf("Connection failed: %s", state.ErrorMsg)
		}
	}

	client.OnTrackDataUpdate = func(trackData broadcast.TrackData) {
		t.Logf("Track data received: %s (ID: %d, Length: %dm)",
			trackData.TrackName, trackData.TrackId, trackData.TrackMeters)
		trackDataReceived = true
	}

	client.OnEntrylistUpdate = func(carInfo broadcast.CarInfo) {
		t.Logf("Car #%d: %s (Model: %d)",
			carInfo.RaceNumber, carInfo.TeamName, carInfo.CarModelType)
		entryListReceived = true
	}

	client.SetTimeout(2 * time.Second)

	err := client.Connect()
	if err != nil {
		t.Skipf("Cannot connect to ACC server at %s: %v (ACC might not be running)", accAddress, err)
		return
	}
	defer client.Disconnect()

	done := make(chan struct{})
	go func() {
		time.Sleep(5 * time.Second)
		close(done)
	}()

	go func() {
		if err := client.Listen(); err != nil {
			t.Logf("Listen error: %v", err)
		}
	}()

	<-done
	client.Disconnect()

	if !connectionReceived {
		t.Error("Did not receive connection state callback")
	}

	t.Logf("Track data received: %v", trackDataReceived)
	t.Logf("Entry list received: %v", entryListReceived)
}

func TestBroadcastCommands(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	accAddress := os.Getenv("ACC_SERVER_ADDRESS")
	if accAddress == "" {
		accAddress = "127.0.0.1:9000"
	}

	logger.Init(logger.Config{
		Level:  logger.LevelInfo,
		Pretty: true,
		Output: os.Stdout,
	})

	client := broadcast.NewClient(
		accAddress,
		"Command Test Client",
		"asd",
		"",
		250,
		*logger.Get(),
	)

	connectionSuccess := false
	client.OnConnectionStateChanged = func(state broadcast.ConnectionState) {
		if state.Success {
			connectionSuccess = true
			t.Logf("Connected successfully with ID: %d", state.ConnectionId)
		}
	}

	client.SetTimeout(2 * time.Second)

	err := client.Connect()
	if err != nil {
		t.Skipf("Cannot connect to ACC: %v", err)
		return
	}
	defer client.Disconnect()

	time.Sleep(1 * time.Second)

	if !connectionSuccess {
		t.Skip("Connection not successful, skipping command tests")
		return
	}

	t.Run("RequestEntryList", func(t *testing.T) {
		err := client.RequestEntryList()
		if err != nil {
			t.Errorf("RequestEntryList() error = %v", err)
		}
	})

	t.Run("RequestTrackData", func(t *testing.T) {
		err := client.RequestTrackData()
		if err != nil {
			t.Errorf("RequestTrackData() error = %v", err)
		}
	})

	time.Sleep(2 * time.Second)
}

func TestBroadcastReconnection(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	accAddress := os.Getenv("ACC_SERVER_ADDRESS")
	if accAddress == "" {
		accAddress = "127.0.0.1:9000"
	}

	logger.Init(logger.Config{
		Level:  logger.LevelWarn,
		Pretty: true,
		Output: os.Stdout,
	})

	client := broadcast.NewClient(
		accAddress,
		"Reconnection Test",
		"asd",
		"",
		250,
		*logger.Get(),
	)

	client.SetTimeout(2 * time.Second)

	err := client.Connect()
	if err != nil {
		t.Skipf("Cannot connect to ACC: %v", err)
		return
	}

	time.Sleep(1 * time.Second)

	err = client.Disconnect()
	if err != nil {
		t.Errorf("Disconnect() error = %v", err)
	}

	time.Sleep(500 * time.Millisecond)

	err = client.Connect()
	if err != nil {
		t.Errorf("Reconnect failed: %v", err)
	}

	time.Sleep(1 * time.Second)

	client.Disconnect()
}
