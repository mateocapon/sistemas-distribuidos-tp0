package common

import (
    "net"
    "encoding/csv"
    "fmt"
    "os"
    "io"
    log "github.com/sirupsen/logrus"
)

type Bet struct {
    FirstName     string
    LastName      string
    Document      string
    Birthdate     string
    Number        string
}

type BetsReader struct {
    ID            string
    BatchSize int
}

func NewBetsReader(ID string, BatchSize int) *BetsReader {
    betsreader := &BetsReader{
        ID: ID,
        BatchSize: BatchSize,
    }
    return betsreader
}

func (b *BetsReader) processBets(conn net.Conn, protocol *Protocol) error {
    var bets []Bet
    filename := fmt.Sprintf("agency-%s.csv", b.ID)
    file, err := os.Open(filename)
    if err != nil {
        return err
    }
    defer file.Close()

    reader := csv.NewReader(file)
    reader.Comma = ','
    reader.FieldsPerRecord = 5

    for {
        betData, err := reader.Read()
        if err != nil {
            log.Infof("Se llega al EOF")
            if err == io.EOF {
                break
            }
            return err
        }

        bet := Bet{
            FirstName: betData[0],
            LastName:  betData[1],
            Document:  betData[2],
            Birthdate: betData[3],
            Number:    betData[4],
        }
        bets = append(bets, bet)
        log.Infof("bet: %s", bet.Document)
        if len(bets) == b.BatchSize {
            err := protocol.sendBetsChunk(conn, bets, b.ID)
            if err != nil {
                log.Infof("Errror")
                return err
            }
            confirmation, err := protocol.recvConfirmation(conn)
            if err != nil {
                log.Infof("Errroor")
                return err
            }
            if confirmation {
                log.Infof("action: chunk_enviado | result: success | agency: %s", 
                    b.ID,
                )
            }
            bets = nil
        }
    }

    err = protocol.sendBetsLastChunk(conn, bets, b.ID)
    if err != nil {
        log.Infof("Errasdsfaroor")
        return err
    }
    confirmation, err := protocol.recvConfirmation(conn)
    if confirmation {
        log.Infof("action: chunk_enviado | result: success | agency: %s", 
            b.ID,
        )
    }
    return err
}

