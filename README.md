# Ejercicio  N°8

## Enunciado:
Modificar el servidor para que permita aceptar conexiones y procesar mensajes en paralelo.
En este ejercicio es importante considerar los mecanismos de sincronización a utilizar para el correcto funcionamiento de la persistencia.

## Solución:
El protocolo de comunicación se mantiene del ejercicio 7. Sólo se hicieron modificaciones en el servidor para que se puedan procesar las consultas en paralelo.
Decidí lanzar tres tipos de procesos. El proceso main, es el encargado de aceptar clientes y agregarlos a una cola de clientes.
Por su lado, lanzo un conjunto de procesos, que ejecutan la función `handle_client_connection()`. Estos consumen la cola de clientes y los procesan dependiendo del estado de la conexión.

El estado de la conexión puede ser uno de los siguientes
- JUST_ARRIVED
- SEND_BETS
- GET_WINNER
- GET_WINNER_VALIDATED

Cuando un cliente quiere obtener el ganador de las apuestas, es posible que todavía no hayan notificado todas las otras agencias. Por lo tanto se lo agrega a otra cola la cual es consumida por un tercer proceso.

Este tercer proceso ejecuta la función `count_loaded_bets()` que consta de dos etapas. Una primera etapa que se bloquea en una cola esperando por ids de agencias que ya finalizaron de enviar las apuestas. Y luego, una vez que las agencias finalizan el envío, una segunda etapa donde agrega a los clientes que estaban esperando por un ganador devuelta a la cola del client_handler. En este caso, modificandole el estado a GET_WINNER_VALIDATED. Esto es, ya podrá obtener el resultado de la apuesta.

Esta situación se puede observar en el siguiente diagrama de secuencias, donde cada Actor es un proceso diferente.

![image](https://user-images.githubusercontent.com/65830097/228957706-3dce25af-14df-46d7-9a09-6a5418277076.png)


En función de sincronizar el acceso al archivo de apuestas, la pool de procesos comparte un Lock para escribir una apuesta en el archivo. Siendo que la lectura del archivo se hace solamente cuando la escritura finalizó, no se toma el lock para poder leer del archivo.

En cuánto a la sincronización entre procesos, tal como se comenta anteriormente, se realizó con colas bloqueantes.

Un punto que me parece interesante agregar, es que los procesos `handle_client_connection()` cuando se encuentran con una agencia que desea enviar las apuestas, procesa un chunk de las apuestas y vuelve a encolar al cliente. No procesa todas las apuestas en un loop. De este modo, existe un mayor fairness entre las agencias que desean enviar apuestas. Cada loop tiene solo una parte del procesamiento total que hay que hacer.