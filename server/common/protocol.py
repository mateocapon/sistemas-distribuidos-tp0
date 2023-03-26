from common.utils import Bet
import socket


STRING_TO_READ_SIZE = 2


# Protocol Packet to receive Bet:
# Each string is sended with a 2 byte len of string
# data and the string after it to avoid short reads.

def receive_bet(client_sock) -> Bet:
    agency = receive_string(client_sock)
    first_name = receive_string(client_sock)
    last_name = receive_string(client_sock)
    document = receive_string(client_sock)
    birthdate = receive_string(client_sock)
    number = receive_string(client_sock)
    return Bet(agency, first_name, last_name, document, birthdate, number)

def receive_string(client_sock):
    len_data = recvall(client_sock, STRING_TO_READ_SIZE)
    string_len = int.from_bytes(len_data, byteorder='big')
    name = recvall(client_sock, string_len)
    return name.decode('utf-8')

def send_confirmation(client_sock):
    client_sock.sendall(b'O')


def recvall(client_sock, n):
    """ recv all n bytes to avoid short read"""
    data = b''
    while len(data) < n:
        data += client_sock.recv(n - len(data))
    return data
