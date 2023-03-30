package common

import (
    "net"
    "os"
    "os/signal"
    "syscall"

    "sync"
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

func (c *Client) StartClient() {
    signalChan := make(chan os.Signal, 1)
    signal.Notify(signalChan, syscall.SIGTERM)
    var wg sync.WaitGroup
    connectionFinishedChan := make(chan bool)
    c.createClientSocket()
    
    wg.Add(1)
    go func() {
        protocol := NewProtocol(c.config.MaxPackageSize)
        betsreader := NewBetsReader(c.config.ID, c.config.BatchSize)
        err := betsreader.processBets(c.conn, protocol)
        if err != nil {
            log.Errorf("action: send_bets | result: fail | client_id: %v | error: %v",
                c.config.ID,
                err,
            )
        }
        connectionFinishedChan <- true
        wg.Done()
    }()
    finished := false
    select {
    case <-signalChan: 
        c.conn.Close()
        finished = true
    case <-connectionFinishedChan:
        c.conn.Close()
    }
    log.Infof("action: release_socketfd | result: success | client_id: %v",
            c.config.ID,
    )
    wg.Wait()
    if finished {
        return
    }
    c.createClientSocket()
    
    wg.Add(1)
    go func() {
        defer wg.Done()
        c.askForWinners()
        connectionFinishedChan <- true
    }()
    select {
    case <-signalChan: 
        c.conn.Close()
    case <-connectionFinishedChan:
        c.conn.Close()
    }
    log.Infof("action: release_socketfd | result: success | client_id: %v",
            c.config.ID,
    )
    wg.Wait()
}

// Creates a new connection to server. Notifies server to get all winners in agency.
func (c *Client) askForWinners() {
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
