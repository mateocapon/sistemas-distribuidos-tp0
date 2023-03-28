import socket
import logging
import signal
import multiprocessing as mp
from common.clienthandler import handle_client_connection
from common._bets_loaded_counter import count_loaded_bets

class Server:
    def __init__(self, port, listen_backlog, number_clients, n_workers = 4):
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

        self.clients_accepted_queue = mp.Queue()
        self.load_bets_queue = mp.Queue()
        self.waiting_winner_queue = mp.Queue()

        self._workers = [mp.Process(target=handle_client_connection, 
                                    args=(self.clients_accepted_queue, 
                                          self.load_bets_queue, 
                                          self.waiting_winner_queue)) 
                                    for i in range(n_workers)]

        self._bets_loaded_counter = mp.Process(target=count_loaded_bets, 
                                               args=(self.clients_accepted_queue, 
                                                     self.load_bets_queue, 
                                                     self.waiting_winner_queue,
                                                     number_clients))
        signal.signal(signal.SIGTERM, self.__stop_accepting)

    def run(self):
        """
        Dummy Server loop

        Server that accept a new connections and establishes a
        communication with a client. After client with communucation
        finishes, servers starts to accept new connections again
        """
        for worker in self._workers:
            worker.start()

        while self._server_active:
            client_sock = self.__accept_new_connection()
            if client_sock:
                self.__handle_client_connection(client_sock)
            elif self._server_active:
                self.__stop_accepting()
        for _, client_sock in self.waiting_winner_cli.items():
            client_sock.close()

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

