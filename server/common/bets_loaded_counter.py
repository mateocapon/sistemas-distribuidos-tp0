from common.clienthandler import GET_WINNER_VALIDATED
import logging

def count_loaded_bets(clients_queue, load_bets_queue, waiting_winner_queue, number_clients):
    # set of agencies that have completed the storage of bets.
    try:
        agencies_stored_bets = set()
        while len(agencies_stored_bets) < number_clients:
            agency = load_bets_queue.get()
            if not agency:
                for i in range(number_clients):
                    clients_queue.put((None, None))
                logging.info(f'action: stop_process_count | result: success')
                return
            agencies_stored_bets.add(int(agency))
        logging.info(f'action: sorteo | result: success')
        while True:
            client_sock = waiting_winner_queue.get()
            if not client_sock:
                for i in range(number_clients):
                    clients_queue.put((None, None))
                logging.info(f'action: stop_process_count | result: success')
                return
            # the lottery is done, client is ready to get results.
            clients_queue.put((client_sock, GET_WINNER_VALIDATED))
    except ValueError:
        logging.debug(f'action: get_or_put_to_queue | result: fail')
    except:
        logging.error(f'action: count_loaded_bets | result: fail')
