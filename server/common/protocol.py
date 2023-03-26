from common.utils import Bet
import socket
import logging


UINT16_SIZE = 2

NORMAL_CHUNK = 'C'


def receive_bets_chunk(client_sock):
    type_chunk = chr(client_sock.recv(1)[0]) 
    more_chunks = type_chunk == NORMAL_CHUNK
    number_bets = receive_uint16(client_sock)
    agency = receive_string(client_sock)
    bets = []
    for i in range(number_bets):
        bets.append(receive_bet(agency, client_sock))
    return (more_chunks, bets, agency)


# Protocol Packet to receive Bet:
# Each string is sended with a 2 byte len of string
# data and the string after it to avoid short reads.

def receive_bet(agency, client_sock) -> Bet:
    first_name = receive_string(client_sock)
    last_name = receive_string(client_sock)
    document = receive_string(client_sock)
    birthdate = receive_string(client_sock)
    number = receive_string(client_sock)
    return Bet(agency, first_name, last_name, document, birthdate, number)

def receive_uint16(client_sock):
    len_data = recvall(client_sock, UINT16_SIZE)
    return int.from_bytes(len_data, byteorder='big')

def receive_string(client_sock):
    string_len = receive_uint16(client_sock)
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
