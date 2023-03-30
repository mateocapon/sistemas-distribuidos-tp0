from common.utils import Bet
import socket


STRING_TO_READ_SIZE = 2

CONFIRMATION = b'O'
ERROR = b'E'

# Protocol Packet to receive Bet:
# Each string is sended with a 2 byte len of string
# data and the string after it to avoid short reads.

def receive_bet(client_sock) -> Bet:
    """
    Receives a single Bet from socket.
    A bet is a list of strings in the following order
    agency - first_name - last_name - document - birthdate - number
    """
    agency = receive_string(client_sock)
    first_name = receive_string(client_sock)
    last_name = receive_string(client_sock)
    document = receive_string(client_sock)
    birthdate = receive_string(client_sock)
    number = receive_string(client_sock)
    return Bet(agency, first_name, last_name, document, birthdate, number)

def receive_string(client_sock):
    """
    Receives a string len in 2 bytes, and then receives the whole string.
    """
    len_data = recvall(client_sock, STRING_TO_READ_SIZE)
    string_len = int.from_bytes(len_data, byteorder='big')
    name = recvall(client_sock, string_len)
    return name.decode('utf-8')

def send_confirmation(client_sock):
    """
    Send the client a 'bets stored' confirmation byte 
    """
    client_sock.sendall(CONFIRMATION)

def send_error(client_sock, error_msg):
    """
    Send error message to the client
    """
    msg = bytearray()
    msg += ERROR
    msg += len(error_msg).to_bytes(STRING_TO_READ_SIZE, "big")
    msg += error_msg.encode('utf-8')
    client_sock.sendall(msg)


def recvall(client_sock, n):
    """
    Recv all n bytes to avoid short read
    """
    data = b''
    while len(data) < n:
        received = client_sock.recv(n - len(data)) 
        if not received:
            raise OSError("No data received in recvall")
        data += received
    return data
