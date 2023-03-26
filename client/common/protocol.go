package common

import (
    "net"
    "errors"
    "bytes"
    "bufio"
    "fmt"
    "io"
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
const CONFIRMATION = 'O'
const ERROR = 'E'


// Send a bet through socket.
// Data in buffer:
// |2 bytes Big Endian len string | String | ... |
// Above protocol is for all items in a bet, with order:
// ID - FirstName - LastName - Document - BirthDate - Number
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

// writes all the content of the data in socket.
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


// waits for client confirmation. Throws the corresponding error.
func (p *Protocol) recvConfirmation(conn net.Conn) (bool, error) {
    reader := bufio.NewReader(conn)
    read, err := reader.ReadByte()
    if err != nil {
        return false, err
    }
    if read == ERROR {
        s, err := readString(reader)
        if err != nil {
            return false, err
        }
        return false, errors.New(fmt.Sprintf("server-error-response: %s", s))
    }
    return true, nil
}

// reads a String from socket, reading first the len of it 
// in an uint16 type to avoid short reads.
func readString(reader *bufio.Reader) (string, error) {
    lenString := make([]byte, LEN_STRING)
    if _, err := io.ReadFull(reader, lenString); err != nil {
        return "", err
    }
    length := binary.BigEndian.Uint16(lenString)
    stringData := make([]byte, length)
    if _, err := io.ReadFull(reader, stringData); err != nil {
        return "", err
    }
    return string(stringData), nil
}

// Passes int as uint16 to Big Endian.
// int must fit in uint16. 
func toBigEndian(number int) []byte {
    data := make([]byte, 2)
    binary.BigEndian.PutUint16(data, uint16(number))
    return data 
}
