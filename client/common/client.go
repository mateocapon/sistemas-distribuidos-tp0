package common

import (
	"net"
    "os"
    "os/signal"
    "syscall"

	log "github.com/sirupsen/logrus"
)

// ClientConfig Configuration used by the client
type ClientConfig struct {
	ID            string
	ServerAddress string
	MaxPackageSize int
	FirstName      string
	LastName      string
	Document      string
	Birthdate     string
	Number        string
}

// Client Entity that encapsulates how
type Client struct {
	config ClientConfig
	conn   net.Conn
}

// NewClient Initializes a new client receiving the configuration
// as a parameter
func NewClient(config ClientConfig) *Client {
	client := &Client{
		config: config,
	}
	return client
}

// CreateClientSocket Initializes client socket. In case of
// failure, error is printed in stdout/stderr and exit 1
// is returned
func (c *Client) createClientSocket() error {
	conn, err := net.Dial("tcp", c.config.ServerAddress)
	if err != nil {
		log.Fatalf(
	        "action: connect | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
	}
	c.conn = conn
	return nil
}

// StartClientLoop Send messages to the client until some time threshold is met
func (c *Client) StartClient() {
    signalChan := make(chan os.Signal, 1)
    signal.Notify(signalChan, syscall.SIGTERM)

	// Create the connection the server in every loop iteration. Send an
	c.createClientSocket()
	bet := Bet{
		ID:            c.config.ID,
		FirstName:     c.config.FirstName,
		LastName:	   c.config.LastName,
		Document:	   c.config.Document,
		Birthdate:	   c.config.Birthdate,
		Number:        c.config.Number,
	}
	protocol := NewProtocol(c.config.MaxPackageSize)
	err := protocol.sendBet(c.conn, bet)
	if err != nil {
		log.Errorf("action: send_message | result: fail | client_id: %v | error: %v",
            c.config.ID,
			err,
		)
		c.conn.Close()
		return
	}
	confirmation, err := protocol.recvConfirmation(c.conn)		
	c.conn.Close()

	if err != nil {
		log.Errorf("action: receive_message | result: fail | client_id: %v | error: %v",
            c.config.ID,
			err,
		)
		return
	}
	if confirmation {
		log.Infof("action: apuesta_enviada | result: success | dni: %s | numero: %s", 
		    bet.Document,
		    bet.Number,
		)
	}
	log.Infof("action: client_finished | result: success | client_id: %v", c.config.ID)
}
