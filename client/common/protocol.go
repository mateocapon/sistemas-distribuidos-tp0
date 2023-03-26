package common

import (
    "net"
    "errors"
    "bufio"
    "bytes"
    "fmt"
    "encoding/binary"
)



type Protocol struct {
    maxPackageSize int
}

func NewProtocol(maxPackageSize int) *Protocol {
    protocol := &Protocol{
        maxPackageSize: maxPackageSize,
    }
    return protocol
}


const LEN_STRING = 2
const LEN_BETS = 2
const SIMPLE_CHUNK = 'C'
const LAST_CHUNK = 'L'


func (p *Protocol) sendBets(conn net.Conn, bets []Bet, ID string, typeChunk byte) error {
    var data [][]byte 
    data = append(data, []byte{typeChunk}, toBigEndian(len(bets)), toBigEndian(len(ID)), []byte(ID))
    for _, bet := range bets {
        data = append(data, 
                     toBigEndian(len(bet.FirstName)), []byte(bet.FirstName),
                     toBigEndian(len(bet.LastName)), []byte(bet.LastName),
                     toBigEndian(len(bet.Document)), []byte(bet.Document),
                     toBigEndian(len(bet.Birthdate)), []byte(bet.Birthdate),
                     toBigEndian(len(bet.Number)), []byte(bet.Number))
    }
    joined := bytes.Join(data, []byte(""))
    if len(joined) > p.maxPackageSize {
        return errors.New(fmt.Sprintf("Package of size %d is too big", len(data))) 
    }
    return writeAll(conn, joined)
}


func (p *Protocol) sendBetsChunk(conn net.Conn, bets []Bet, ID string) error {
    return p.sendBets(conn, bets, ID, SIMPLE_CHUNK)
}

func (p *Protocol) sendBetsLastChunk(conn net.Conn, bets []Bet, ID string) error {
    return p.sendBets(conn, bets, ID, LAST_CHUNK)
}

func writeAll(conn net.Conn, data []byte) error {
    totalBytes := len(data)
    bytesWritten := 0
    for bytesWritten < totalBytes {
        n, err := conn.Write(data[bytesWritten:])
        if err != nil {
            return err
        }
        bytesWritten += n
    }
    return nil
}

func (p *Protocol) recvConfirmation(conn net.Conn) (bool, error) {
    read, err := bufio.NewReader(conn).ReadByte()
    if err != nil || read != 'O' {
        return false, err
    }
    return true, nil
}


func toBigEndian(number int) []byte {
    data := make([]byte, 2)
    binary.BigEndian.PutUint16(data, uint16(number))
    return data 
}
