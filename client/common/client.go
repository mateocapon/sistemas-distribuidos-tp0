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
    BatchSize      int
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
    protocol := NewProtocol(c.config.MaxPackageSize)
    betsreader := NewBetsReader(c.config.ID, c.config.BatchSize)
    err := betsreader.processBets(c.conn, protocol)
    c.conn.Close()
    if err != nil {
        log.Errorf("action: send_bets | result: fail | client_id: %v | error: %v",
            c.config.ID,
            err,
        )
        return
    }
    c.askForWinners()
    log.Infof("action: client_finished | result: success | client_id: %v", c.config.ID)
}

// Creates a new connection to server. Notifies server to get all winners in agency.
func (c *Client) askForWinners() {
    c.createClientSocket()
    protocol := NewProtocol(c.config.MaxPackageSize)
    n_winners, err := protocol.receiveWinners(c.conn, c.config.ID)
    c.conn.Close()
    if err != nil {
        log.Errorf("action: receive_winners | result: fail | client_id: %v | error: %v",
            c.config.ID,
            err,
        )
        return
    }
    log.Infof("action: consulta_ganadores | result: success | cant_ganadores: %d", n_winners)
}
