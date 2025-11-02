package broadcast

import (
	"context"
	"fmt"
	"net"
	"time"

	"RaceAll/internal/errors"

	"github.com/rs/zerolog"
)

const (
	// ReadBufferSize is the size of the UDP read buffer
	ReadBufferSize = 32 * 1024

	// DefaultTimeout is the default timeout for connections
	DefaultTimeout = 5 * time.Second
)

// Client is the UDP client for connecting to the ACC broadcasting
type Client struct {
	logger   zerolog.Logger
	conn     *net.UDPConn
	protocol *Protocol

	// Configuration
	address                  string
	displayName              string
	connectionPassword       string
	commandPassword          string
	msRealtimeUpdateInterval int32
	timeout                  time.Duration

	// Control
	ctx    context.Context
	cancel context.CancelFunc

	// Public callbacks - delegated to protocol
	OnConnectionStateChanged func(ConnectionState)
	OnTrackDataUpdate        func(TrackData)
	OnEntrylistUpdate        func(CarInfo)
	OnRealtimeUpdate         func(RealtimeUpdate)
	OnRealtimeCarUpdate      func(RealtimeCarUpdate)
	OnBroadcastingEvent      func(BroadcastingEvent)
}

// NewClient create a new instance of the broadcasting client
func NewClient(address, displayName, connectionPassword, commandPassword string, msRealtimeUpdateInterval int32, logger zerolog.Logger) *Client {
	return &Client{
		logger:                   logger,
		address:                  address,
		displayName:              displayName,
		connectionPassword:       connectionPassword,
		commandPassword:          commandPassword,
		msRealtimeUpdateInterval: msRealtimeUpdateInterval,
		timeout:                  DefaultTimeout,
	}
}

// SetTimeout sets the timeout for the connection
func (c *Client) SetTimeout(timeout time.Duration) {
	c.timeout = timeout
}

// Connect establishes the connection with ACC
func (c *Client) Connect() error {
	// Resolve UDP address
	raddr, err := net.ResolveUDPAddr("udp", c.address)
	if err != nil {
		c.logger.Error().Err(err).Msg("Error resolving UDP address")
		return NewError("Connect", fmt.Errorf("failed to resolve address: %w", err))
	}

	// Create UDP connection
	c.conn, err = net.DialUDP("udp", nil, raddr)
	if err != nil {
		c.logger.Error().Err(err).Msg("Error creating UDP connection")
		return NewError("Connect", fmt.Errorf("failed to dial UDP: %w", err))
	}

	// Create protocol handler
	c.protocol = NewProtocol(c.address, c.send, c.logger)

	// Configure protocol callbacks
	c.protocol.OnConnectionStateChanged = func(state ConnectionState) {
		if c.OnConnectionStateChanged != nil {
			c.OnConnectionStateChanged(state)
		}
	}

	c.protocol.OnTrackDataUpdate = func(trackData TrackData) {
		if c.OnTrackDataUpdate != nil {
			c.OnTrackDataUpdate(trackData)
		}
	}

	c.protocol.OnEntrylistUpdate = func(carInfo CarInfo) {
		if c.OnEntrylistUpdate != nil {
			c.OnEntrylistUpdate(carInfo)
		}
	}

	c.protocol.OnRealtimeUpdate = func(update RealtimeUpdate) {
		if c.OnRealtimeUpdate != nil {
			c.OnRealtimeUpdate(update)
		}
	}

	c.protocol.OnRealtimeCarUpdate = func(update RealtimeCarUpdate) {
		if c.OnRealtimeCarUpdate != nil {
			c.OnRealtimeCarUpdate(update)
		}
	}

	c.protocol.OnBroadcastingEvent = func(event BroadcastingEvent) {
		if c.OnBroadcastingEvent != nil {
			c.OnBroadcastingEvent(event)
		}
	}

	// Create context for control
	c.ctx, c.cancel = context.WithCancel(context.Background())

	c.logger.Info().Str("address", c.address).Msg("Connected to ACC")

	// Request connection
	if err := c.protocol.RequestConnection(c.displayName, c.connectionPassword, c.msRealtimeUpdateInterval, c.commandPassword); err != nil {
		c.logger.Error().Err(err).Msg("Error requesting connection")
		return NewError("Connect", fmt.Errorf("failed to request connection: %w", err))
	}

	return nil
}

// Listen starts listening for messages from the server
func (c *Client) Listen() error {
	if c.conn == nil {
		return NewError("Listen", errors.ErrConnectionClosed)
	}

	buffer := make([]byte, ReadBufferSize)

	for {
		select {
		case <-c.ctx.Done():
			c.logger.Info().Msg("Deteniendo listener")
			return nil
		default:
			// Set read timeout
			if err := c.conn.SetReadDeadline(time.Now().Add(c.timeout)); err != nil {
				c.logger.Error().Err(err).Msg("Error setting deadline")
				return NewError("Listen", err)
			}

			// Read data
			n, err := c.conn.Read(buffer)
			if err != nil {
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					c.logger.Error().Msg("Timeout: ACC no respondió")
					return NewError("Listen", errors.ErrTimeout)
				}

				// If the context was canceled, it's not an error
				select {
				case <-c.ctx.Done():
					return nil
				default:
					c.logger.Error().Err(err).Msg("Error reading data")
					return NewError("Listen", err)
				}
			}

			if n > 0 {
				// Process message
				data := make([]byte, n)
				copy(data, buffer[:n])

				if err := c.protocol.ProcessMessage(data); err != nil {
					c.logger.Error().Err(err).Msg("Error processing message")
					// don't return the error, continue listening
				}
			}
		}
	}
}

// ConnectAndListen connects and starts listening in a single call
func (c *Client) ConnectAndListen() error {
	if err := c.Connect(); err != nil {
		return err
	}

	return c.Listen()
}

func (c *Client) send(data []byte) error {
	if c.conn == nil {
		return NewError("send", errors.ErrConnectionClosed)
	}

	n, err := c.conn.Write(data)
	if err != nil {
		c.logger.Error().Err(err).Msg("Error al enviar datos")
		return NewError("send", err)
	}

	if n != len(data) {
		return NewError("send", errors.ErrPartialWrite)
	}

	return nil
}

func (c *Client) Disconnect() error {
	if c.cancel != nil {
		c.cancel()
	}

	if c.protocol != nil {
		if err := c.protocol.Disconnect(); err != nil {
			c.logger.Warn().Err(err).Msg("Error al enviar mensaje de desconexión")
		}
	}

	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			c.logger.Warn().Err(err).Msg("Error al cerrar conexión UDP")
			return err
		}
		c.conn = nil
	}

	c.logger.Info().Msg("Desconectado de ACC")
	return nil
}

func (c *Client) GetCarInfo(carIndex uint16) (*CarInfo, bool) {
	if c.protocol == nil {
		return nil, false
	}
	return c.protocol.GetCarInfo(carIndex)
}

func (c *Client) RequestEntryList() error {
	if c.protocol == nil {
		return NewError("RequestEntryList", errors.ErrProtocolNotInitialized)
	}
	return c.protocol.RequestEntryList()
}

func (c *Client) RequestTrackData() error {
	if c.protocol == nil {
		return NewError("RequestTrackData", errors.ErrProtocolNotInitialized)
	}
	return c.protocol.RequestTrackData()
}

func (c *Client) SetFocus(carIndex uint16) error {
	if c.protocol == nil {
		return NewError("SetFocus", errors.ErrProtocolNotInitialized)
	}
	return c.protocol.SetFocus(carIndex)
}

func (c *Client) SetCamera(cameraSet, camera string) error {
	if c.protocol == nil {
		return NewError("SetCamera", errors.ErrProtocolNotInitialized)
	}
	return c.protocol.SetCamera(cameraSet, camera)
}

func (c *Client) SetFocusAndCamera(carIndex uint16, cameraSet, camera string) error {
	if c.protocol == nil {
		return NewError("SetFocusAndCamera", errors.ErrProtocolNotInitialized)
	}
	return c.protocol.SetFocusAndCamera(carIndex, cameraSet, camera)
}

func (c *Client) RequestInstantReplay(startSessionTime, durationMS float32, initialFocusedCarIndex int32, initialCameraSet, initialCamera string) error {
	if c.protocol == nil {
		return NewError("RequestInstantReplay", errors.ErrProtocolNotInitialized)
	}
	return c.protocol.RequestInstantReplay(startSessionTime, durationMS, initialFocusedCarIndex, initialCameraSet, initialCamera)
}

func (c *Client) RequestHUDPage(hudPage string) error {
	if c.protocol == nil {
		return NewError("RequestHUDPage", errors.ErrProtocolNotInitialized)
	}
	return c.protocol.RequestHUDPage(hudPage)
}
