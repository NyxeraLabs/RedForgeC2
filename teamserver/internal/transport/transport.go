package transport

// Transport defines the minimal interface for a C2 transport.
// Implementations may include HTTPS, WebSocket, DNS, ICMP, etc.
type Transport interface {
	// Send sends a raw payload to the configured destination.
	Send([]byte) error

	// Receive blocks until a message is available or an error occurs.
	Receive() ([]byte, error)

	// Close shuts down the transport gracefully.
	Close() error
}
