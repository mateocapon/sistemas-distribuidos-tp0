import socket
import logging
import signal
from common.utils import Bet, store_bets, load_bets, has_won
from common.protocol import receive_bets_chunk, send_confirmation, send_error, get_client_intention
from common.protocol import SEND_BETS_INTENTION, GET_WINNER_INTENTION, send_winners, receive_agency_id

class Server:
    def __init__(self, port, listen_backlog, number_clients):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)
        self._server_active = True
        self.number_clients = number_clients
        # dict of "agency_ID: client_sock" 
        self.waiting_winner_cli = {}
        # set of agencies that have completed the storage of bets.
        self.agencies_stored_bets = set()
        signal.signal(signal.SIGTERM, self.__stop_accepting)

    def run(self):
        """
        Dummy Server loop

        Server that accept a new connections and establishes a
        communication with a client. After client with communucation
        finishes, servers starts to accept new connections again
        """
        while self._server_active:
            client_sock = self.__accept_new_connection()
            if client_sock:
                self.__handle_client_connection(client_sock)
            elif self._server_active:
                self.__stop_accepting()
        for _, client_sock in self.waiting_winner_cli.items():
            client_sock.close()

    def __handle_client_connection(self, client_sock):
        """
        Read message from a specific client socket and closes the socket

        If a problem arises in the communication with the client, the
        client socket will also be closed
        """
        persist_connection = False
        try:
            addr = client_sock.getpeername()
            client_intention = get_client_intention(client_sock)
            if client_intention == SEND_BETS_INTENTION:
                self.__receive_bets(client_sock)
            elif client_intention == GET_WINNER_INTENTION:
                persist_connection = self.__get_winner(client_sock)
            else:
                logging.error(f'action: get_client_intention | result: fail | error: intention_not_valid')
        except OSError as e:
            logging.error(f'action: receive_message | result: fail | error: {e}')
        finally:
            if not persist_connection:
                client_sock.close()
            logging.info(f'action: close_client | result: success | ip: {addr[0]}')


    def __accept_new_connection(self):
        """
        Accept new connections

        Function blocks until a connection to a client is made.
        Then connection created is printed and returned
        """

        # Connection arrived
        try:
            logging.info('action: accept_connections | result: in_progress')
            c, addr = self._server_socket.accept()
            logging.info(f'action: accept_connections | result: success | ip: {addr[0]}')
            return c
        except OSError as e:
            if self._server_active:
                logging.error(f'action: accept_connections | result: fail | error: {e}')
            return False

    def __stop_accepting(self, *args):
        """
        Closes server socket in order to stop the server gracefully. 
        """
        logging.info('action: stop_server | result: in_progress')
        self._server_active = False
        try:
            self._server_socket.shutdown(socket.SHUT_WR)
            self._server_socket.close()
            logging.info('action: stop_server | result: success')
            logging.info('action: release_server_socketfd | result: success')
        except OSError as e:
            logging.error(f'action: stop_server | result: fail | error: {e}')

    def __receive_bets(self, client_sock):
        """
        Receives all the bets from one agency and stores them
        in the bets file.
        """
        try:
            more_chunks = True
            while more_chunks:
                more_chunks, bets, agency = receive_bets_chunk(client_sock)
                store_bets(bets)
                logging.info(f'action: apuesta_almacenada | result: success | agency: {agency} | n: {len(bets)} | active: {more_chunks}')
                send_confirmation(client_sock)
            self.agencies_stored_bets.add(int(agency))
            if len(self.agencies_stored_bets) == self.number_clients:
                logging.info(f'action: sorteo | result: success')
                self.__send_winners()
        except ValueError as e:
            send_error(client_sock, f'error: {e}')
            logging.error(f'action: receive_bets | result: fail | error: {e}')

    def __get_winner(self, client_sock):
        """
        If there is a winner, sends the winners Document to all clients waiting.
        """
        agency_id = int(receive_agency_id(client_sock))
        if agency_id in self.waiting_winner_cli:
            self.waiting_winner_cli[agency_id].close()
        self.waiting_winner_cli[agency_id] = client_sock
        if len(self.agencies_stored_bets) == self.number_clients:
            self.__send_winners()
            return False
        return True

    def __send_winners(self):
        """
        Send the winner to all clients waiting.
        """
        documents_to_send = {}
        for agency in self.waiting_winner_cli:
            documents_to_send[agency] = []
        bets = load_bets()
        for bet in bets:
            if has_won(bet) and bet.agency in documents_to_send:
               documents_to_send[bet.agency].append(bet.document)
        for agency, documents in documents_to_send.items():
            send_winners(self.waiting_winner_cli[agency], documents)
        logging.info(f'action: send_winners | result: success')
        for _, client_sock in self.waiting_winner_cli.items():
            client_sock.close()
        self.waiting_winner_cli = {}
