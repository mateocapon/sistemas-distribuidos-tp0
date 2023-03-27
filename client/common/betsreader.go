package common

import (
    "net"
    "encoding/csv"
    "fmt"
    "os"
    "io"
    log "github.com/sirupsen/logrus"
)

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
    protocol.sendBetsIntention(conn)
    reader := csv.NewReader(file)
    reader.Comma = ','
    reader.FieldsPerRecord = 5

    for {
        betData, err := reader.Read()
        if err != nil {
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
        if len(bets) == b.BatchSize {
            err := protocol.sendBetsChunk(conn, bets, b.ID)
            if err != nil {
                return err
            }
            confirmation, err := protocol.recvConfirmation(conn)
            if err != nil {
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

