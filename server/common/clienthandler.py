import multiprocessing as mp
import socket
import logging
from common.utils import Bet, store_bets, load_bets, has_won
from common.protocol import receive_bets_chunk, send_confirmation, send_error, get_client_intention
from common.protocol import SEND_BETS_INTENTION, GET_WINNER_INTENTION, send_winners, receive_agency_id


JUST_ARRIVED = 'A'
GET_WINNER_VALIDATED = 'G'

def handle_client_connection(clients_queue, load_bets_queue, waiting_winner_queue, bets_file_lock):
    """
    Read message from a specific client socket and closes the socket

    If a problem arises in the communication with the client, the
    client socket will also be closed
    """
    persist_connection = False
    all_bets_loaded = False
    server_working = True
    while server_working:
        try:
            client_sock, status = clients_queue.get()
            persist_connection = False
            addr = client_sock.getpeername()
            if status == JUST_ARRIVED:
                status = get_client_intention(client_sock)
            if status == SEND_BETS_INTENTION:
                persist_connection = __receive_bets(client_sock, clients_queue, load_bets_queue, bets_file_lock)
            elif (status == GET_WINNER_VALIDATED) or (status == GET_WINNER_INTENTION and all_bets_loaded):
                # once one client is validated to get the winners, all agencies can get the winners.
                all_bets_loaded = True
                __send_winners(client_sock)
            elif status == GET_WINNER_INTENTION:
                # must be validated by bets loaded counter process.
                waiting_winner_queue.put(client_sock)
                persist_connection = True
            else:
                logging.error(f'action: get_client_intention | result: fail | error: intention_not_valid')
        except OSError as e:
            logging.error(f'action: receive_message | result: fail | error: {e}')
        finally:
            if not persist_connection:
                client_sock.close()
                logging.info(f'action: close_client | result: success | ip: {addr[0]}')

def __receive_bets(client_sock, clients_queue, load_bets_queue, bets_file_lock):
    """
    Receives one chunk of bets from one agency and stores them
    in the bets file. If it is the last chunk of bets from this
    agency, notifies the bets_loaded_counter queue.

    Return True if there are more bets to be processed.
    Else returns False
    """
    try:
        more_chunks, bets, agency = receive_bets_chunk(client_sock)
        with bets_file_lock:
            store_bets(bets)
        logging.info(f'action: apuesta_almacenada | result: success | agency: {agency} | n: {len(bets)}')
        send_confirmation(client_sock)
        if more_chunks:
            clients_queue.put((client_sock, SEND_BETS_INTENTION))
        else:
            load_bets_queue.put(agency)
    except ValueError as e:
        send_error(client_sock, f'error: {e}')
        logging.error(f'action: receive_bets | result: fail | error: {e}')
    return more_chunks



def __send_winners(client_sock):
    """
    Send the winner to the client.
    """
    agency_id = int(receive_agency_id(client_sock))
    bets = load_bets()
    winning_bets = filter(lambda bet: has_won(bet) and bet.agency == agency_id, bets)
    documents = list(map(lambda bet: bet.document, winning_bets))
    send_winners(client_sock, documents)
    logging.info(f'action: send_winners | result: success')
