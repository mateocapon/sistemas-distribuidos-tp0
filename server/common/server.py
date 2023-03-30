import socket
import logging
import signal
from common.utils import Bet, store_bets
from common.protocol import receive_bets_chunk, send_confirmation, send_error


class Server:
    def __init__(self, port, listen_backlog):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)
        self._server_active = True
        self._client_sock = None
        self._client_active = False
        signal.signal(signal.SIGTERM, self.__stop_accepting)


    def run(self):
        """
        Dummy Server loop

        Server that accept a new connections and establishes a
        communication with a client. After client with communucation
        finishes, servers starts to accept new connections again
        """
        while self._server_active:
            self._client_sock = self.__accept_new_connection()
            if self._client_sock:
                self._client_active = True
                self.__handle_client_connection()
            elif self._server_active:
                self.__stop_accepting()  

    def __handle_client_connection(self):
        """
        Read message from a specific client socket and closes the socket

        If a problem arises in the communication with the client, the
        client socket will also be closed
        """
        try:
            addr = self._client_sock.getpeername()
            more_chunks = True
            while more_chunks:
                more_chunks, bets, agency = receive_bets_chunk(self._client_sock)
                store_bets(bets)
                logging.info(f'action: apuesta_almacenada | result: success | agency: {agency} | n: {len(bets)} | active: {more_chunks}')
                send_confirmation(self._client_sock)
        except ValueError as e:
            send_error(self._client_sock, f'error: {e}')
            logging.error(f'action: receive_bets | result: fail | error: {e}')
        except OSError as e:
            logging.error(f'action: receive_message | result: fail | error: {e}')
        finally:
            # It might have been closed by Signal handler.
            if self._client_active:
                # This does not avoid calling twice close() to _client_sock, but avoids
                # the signal handler to call close if the main thread is inside close().
                self._client_active = False
                self._client_sock.close()
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
        except:
            logging.error(f'action: stop_server | result: fail | error: unknown')

        try:
            if self._client_active:
                self._client_active = False
                self._client_sock.shutdown(socket.SHUT_WR)
                self._client_sock.close()
        except OSError as e:
            logging.error(f'action: stop_client_sock | result: fail | error: {e}')
        except:
            logging.error(f'action: stop_client_sock | result: fail | error: unknown')
