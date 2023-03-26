package common

import (
    "net"
    "errors"
    "bytes"
    "bufio"
    "fmt"
    "encoding/binary"
)

type Bet struct {
    ID            string
    FirstName     string
    LastName      string
    Document      string
    Birthdate     string
    Number        string
}

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

func (p *Protocol) sendBet(conn net.Conn, bet Bet) error {
    lenMessage := LEN_STRING + len(bet.ID) + LEN_STRING + len(bet.FirstName) +
                  LEN_STRING + len(bet.LastName) + LEN_STRING + len(bet.Document) +
                  LEN_STRING + len(bet.Birthdate) + LEN_STRING + len(bet.Number)
    if lenMessage > p.maxPackageSize {
        return errors.New(fmt.Sprintf("Package of size %d is too big", lenMessage)) 
    }

    data := [][]byte{toBigEndian(len(bet.ID)), []byte(bet.ID), 
                     toBigEndian(len(bet.FirstName)), []byte(bet.FirstName),
                     toBigEndian(len(bet.LastName)), []byte(bet.LastName),
                     toBigEndian(len(bet.Document)), []byte(bet.Document),
                     toBigEndian(len(bet.Birthdate)), []byte(bet.Birthdate),
                     toBigEndian(len(bet.Number)), []byte(bet.Number),
                    }
    joined := bytes.Join(data, []byte(""))
    return writeAll(conn, joined)
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
