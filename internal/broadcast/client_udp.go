package broadcast

import (
	"fmt"
	"net"
	"sync"
)

type ACCUdpRemoteClient struct {
	client                   *net.UDPConn
	listenerTask             chan struct{}
	MessageHandler           *BroadcastingNetworkProtocol
	IpPort                   string
	DisplayName              string
	ConnectionPassword       string
	CommandPassword          string
	MsRealtimeUpdateInterval int

	disposedValue bool
	mu            sync.Mutex
}

func NewACCUdpRemoteClient(ip string, port int, displayName, connectionPassword, commandPassword string, msRealtimeUpdateInterval int) (*ACCUdpRemoteClient, error) {
	c := &ACCUdpRemoteClient{
		DisplayName:              displayName,
		ConnectionPassword:       connectionPassword,
		CommandPassword:          commandPassword,
		MsRealtimeUpdateInterval: msRealtimeUpdateInterval,
		listenerTask:             make(chan struct{}),
	}

	c.IpPort = fmt.Sprintf("%s:%d", ip, port)

	messageHandler, err := NewBroadcastingNetworkProtocol(c.IpPort, c.Send)
	if err != nil {
		return nil, err
	}
	c.MessageHandler = messageHandler

	localAddr := &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 0}

	serverAddr := &net.UDPAddr{IP: net.ParseIP(ip), Port: port}

	client, err := net.DialUDP("udp", localAddr, serverAddr)
	if err != nil {
		return nil, err
	}
	c.client = client

	go c.ConnectAndRun()

	return c, nil
}

func (c *ACCUdpRemoteClient) RequestData() {
	c.MessageHandler.RequestData()
}

func (c *ACCUdpRemoteClient) Send(payload []byte) {
	if c.client != nil {
		c.client.Write(payload)
	}
}

func (c *ACCUdpRemoteClient) Shutdown() {
	go func() {
		c.ShutdownAsync()
	}()
}

func (c *ACCUdpRemoteClient) ShutdownAsync() {
	if c.client == nil {
		return
	}

	if c.listenerTask != nil {
		c.MessageHandler.Disconnect()
		c.client.Close()
		c.client = nil
		<-c.listenerTask
	}
}

func (c *ACCUdpRemoteClient) ConnectAndRun() {
	defer close(c.listenerTask)

	c.MessageHandler.RequestConnection(c.DisplayName, c.ConnectionPassword, c.MsRealtimeUpdateInterval, c.CommandPassword)

	buffer := make([]byte, 65536)
	for c.client != nil {
		n, err := c.client.Read(buffer)

		if err != nil {
			if c.client == nil {
				break
			}
			continue
		}

		if n > 0 {
			data := buffer[:n]
			c.MessageHandler.ProcessMessage(data)
		}
	}
}

func (c *ACCUdpRemoteClient) Dispose() {
	c.dispose(true)
}

func (c *ACCUdpRemoteClient) dispose(disposing bool) {
	if !c.disposedValue {
		if disposing {
			if c.client != nil {
				c.client.Close()
				c.client = nil
			}
		}
		c.disposedValue = true
	}
}

func (c *ACCUdpRemoteClient) Close() error {
	c.Dispose()
	return nil
}
