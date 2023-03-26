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


// Send all bets through socket in one packet.
// First bytes of the packet are:
// 1. typeChunk = SIMPLE_CHUNK or LAST_CHUNK to notify that is the last chunk to send.
// 2. number of bets in uint16_t big endian.
// 3. All bets.
// Data for all bets is send like the following.
// |2 bytes Big Endian len string | String | ... |
// Above protocol is for all items in a bet, ordered by:
// FirstName - LastName - Document - BirthDate - Number
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

// sends a SIMPLE_CHUNK: it is not the last chunk
func (p *Protocol) sendBetsChunk(conn net.Conn, bets []Bet, ID string) error {
    return p.sendBets(conn, bets, ID, SIMPLE_CHUNK)
}

// sends the LAST_CHUNK of bets.
func (p *Protocol) sendBetsLastChunk(conn net.Conn, bets []Bet, ID string) error {
    return p.sendBets(conn, bets, ID, LAST_CHUNK)
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
