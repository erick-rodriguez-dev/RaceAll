package integration_test

import (
	"runtime"
	"testing"
	"time"

	"RaceAll/internal/sharedmemory"
)

// TestSharedMemoryConnection prueba la conexión con la memoria compartida de ACC
// Nota: Este test solo funciona en Windows y requiere que ACC esté ejecutándose
func TestSharedMemoryConnection(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	if runtime.GOOS != "windows" {
		t.Skip("Shared memory tests only work on Windows")
	}

	reader := sharedmemory.NewSharedMemoryReader()
	if reader == nil {
		t.Fatal("NewSharedMemoryReader() returned nil")
	}

	// Intentar conectar
	err := reader.Connect()
	if err != nil {
		t.Skipf("Cannot connect to ACC shared memory (ACC might not be running): %v", err)
		return
	}
	defer reader.Disconnect()

	// Verificar que está conectado
	if !reader.IsConnected() {
		t.Error("IsConnected() = false after successful Connect()")
	}

	t.Log("Successfully connected to ACC shared memory")
}

// TestSharedMemoryReadPhysics prueba la lectura de datos físicos
func TestSharedMemoryReadPhysics(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	if runtime.GOOS != "windows" {
		t.Skip("Shared memory tests only work on Windows")
	}

	reader := sharedmemory.NewSharedMemoryReader()
	err := reader.Connect()
	if err != nil {
		t.Skipf("Cannot connect to ACC: %v", err)
		return
	}
	defer reader.Disconnect()

	// Leer datos de física
	physics, err := reader.ReadPhysics()
	if err != nil {
		t.Fatalf("ReadPhysics() error = %v", err)
	}

	if physics == nil {
		t.Fatal("ReadPhysics() returned nil")
	}

	// Verificar que los datos tienen sentido
	t.Logf("Physics data - PacketId: %d", physics.PacketId)
	t.Logf("Gas: %.2f, Brake: %.2f", physics.Gas, physics.Brake)
	t.Logf("Speed KMH: %.2f", physics.SpeedKmh)
	t.Logf("RPM: %d", physics.Rpms)
	t.Logf("Gear: %d", physics.Gear)

	// Valores básicos de sanidad
	if physics.Gear < -1 || physics.Gear > 8 {
		t.Errorf("Gear value %d seems out of range", physics.Gear)
	}

	if physics.Gas < 0 || physics.Gas > 1 {
		t.Errorf("Gas value %.2f should be between 0 and 1", physics.Gas)
	}

	if physics.Brake < 0 || physics.Brake > 1 {
		t.Errorf("Brake value %.2f should be between 0 and 1", physics.Brake)
	}
}

// TestSharedMemoryReadGraphics prueba la lectura de datos gráficos
func TestSharedMemoryReadGraphics(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	if runtime.GOOS != "windows" {
		t.Skip("Shared memory tests only work on Windows")
	}

	reader := sharedmemory.NewSharedMemoryReader()
	err := reader.Connect()
	if err != nil {
		t.Skipf("Cannot connect to ACC: %v", err)
		return
	}
	defer reader.Disconnect()

	// Leer datos gráficos
	graphics, err := reader.ReadGraphics()
	if err != nil {
		t.Fatalf("ReadGraphics() error = %v", err)
	}

	if graphics == nil {
		t.Fatal("ReadGraphics() returned nil")
	}

	t.Logf("Graphics data - PacketId: %d", graphics.PacketId)
	t.Logf("Status: %d", graphics.Status)
	t.Logf("Session: %d", graphics.Session)
	t.Logf("Current Time: %s", graphics.GetCurrentTime())
	t.Logf("Last Time: %s", graphics.GetLastTime())
	t.Logf("Best Time: %s", graphics.GetBestTime())
	t.Logf("Fuel: %.2f L", graphics.FuelXLap)
	t.Logf("Position: %d", graphics.Position)
	clockInt := int(graphics.Clock)
	t.Logf("Clock: %d:%d:%d", clockInt/3600, (clockInt%3600)/60, clockInt%60)
}

// TestSharedMemoryReadStatic prueba la lectura de datos estáticos
func TestSharedMemoryReadStatic(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	if runtime.GOOS != "windows" {
		t.Skip("Shared memory tests only work on Windows")
	}

	reader := sharedmemory.NewSharedMemoryReader()
	err := reader.Connect()
	if err != nil {
		t.Skipf("Cannot connect to ACC: %v", err)
		return
	}
	defer reader.Disconnect()

	// Leer datos estáticos
	static, err := reader.ReadStatic()
	if err != nil {
		t.Fatalf("ReadStatic() error = %v", err)
	}

	if static == nil {
		t.Fatal("ReadStatic() returned nil")
	}

	t.Logf("Static data:")
	t.Logf("SM Version: %s", static.GetSMVersion())
	t.Logf("AC Version: %s", static.GetACVersion())
	t.Logf("Player: %s %s (%s)", static.GetPlayerName(), static.GetPlayerSurname(), static.GetPlayerNick())
	t.Logf("Car: %s", static.GetCarModel())
	t.Logf("Track: %s - %s", static.GetTrack(), static.GetTrackConfiguration())
	t.Logf("Max RPM: %d", static.MaxRpm)
	t.Logf("Max Fuel: %.2f L", static.MaxFuel)
}

// TestSharedMemoryContinuousRead prueba la lectura continua
func TestSharedMemoryContinuousRead(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	if runtime.GOOS != "windows" {
		t.Skip("Shared memory tests only work on Windows")
	}

	reader := sharedmemory.NewSharedMemoryReader()
	err := reader.Connect()
	if err != nil {
		t.Skipf("Cannot connect to ACC: %v", err)
		return
	}
	defer reader.Disconnect()

	// Leer continuamente por 3 segundos
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	timeout := time.After(3 * time.Second)
	readCount := 0

	for {
		select {
		case <-timeout:
			t.Logf("Completed %d reads in 3 seconds", readCount)
			if readCount == 0 {
				t.Error("No successful reads")
			}
			return
		case <-ticker.C:
			physics, err := reader.ReadPhysics()
			if err != nil {
				t.Errorf("Read error: %v", err)
				return
			}

			graphics, err := reader.ReadGraphics()
			if err != nil {
				t.Errorf("Read graphics error: %v", err)
				return
			}

			if physics != nil && graphics != nil {
				readCount++
				if readCount%10 == 0 {
					t.Logf("Read #%d - Speed: %.1f km/h, Gear: %d, Position: %d",
						readCount, physics.SpeedKmh, physics.Gear, graphics.Position)
				}
			}
		}
	}
}

// TestSharedMemoryDisconnectReconnect prueba desconexión y reconexión
func TestSharedMemoryDisconnectReconnect(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	if runtime.GOOS != "windows" {
		t.Skip("Shared memory tests only work on Windows")
	}

	reader := sharedmemory.NewSharedMemoryReader()

	// Primera conexión
	err := reader.Connect()
	if err != nil {
		t.Skipf("Cannot connect to ACC: %v", err)
		return
	}

	if !reader.IsConnected() {
		t.Error("Should be connected after Connect()")
	}

	// Desconectar
	reader.Disconnect()

	if reader.IsConnected() {
		t.Error("Should not be connected after Disconnect()")
	}

	// Reconectar
	err = reader.Connect()
	if err != nil {
		t.Fatalf("Failed to reconnect: %v", err)
	}

	if !reader.IsConnected() {
		t.Error("Should be connected after reconnect")
	}

	// Verificar que podemos leer
	_, err = reader.ReadPhysics()
	if err != nil {
		t.Errorf("Failed to read after reconnect: %v", err)
	}

	reader.Disconnect()
}
