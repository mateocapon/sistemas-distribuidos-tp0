# Ejercicio  N°7

## Enunciado:
Modificar los clientes para que notifiquen al servidor al finalizar con el envío de todas las apuestas y así proceder con el sorteo.
Inmediatamente después de la notificacion, los clientes consultarán la lista de ganadores del sorteo correspondientes a su agencia.
Una vez el cliente obtenga los resultados, deberá imprimir por log: `action: consulta_ganadores | result: success | cant_ganadores: ${CANT}`.

El servidor deberá esperar la notificación de las 5 agencias para considerar que se realizó el sorteo e imprimir por log: `action: sorteo | result: success`.
Luego de este evento, podrá verificar cada apuesta con las funciones `load_bets(...)` y `has_won(...)` y retornar los DNI de los ganadores de la agencia en cuestión. Antes del sorteo, no podrá responder consultas por la lista de ganadores.
Las funciones `load_bets(...)` y `has_won(...)` son provistas por la cátedra y no podrán ser modificadas por el alumno.

## Solución:
Se modifica el cliente para que cree dos conexiones con el servidor. Una primera para enviar las apuestas, y una segunda para pedir el resultado del sorteo. Esto facilita la extensión del código en el servidor. Permite que una agencia pregunte múltiples veces por el resultado de un sorteo, sin la necesidad de enviar más de una vez las apuestas.
Al realizar está modificación, se agrega un primer mensaje en el protocolo de un byte, indicando el tipo de operación que quiere hacer el cliente. Este puede ser enviar las apuestas, o pedir los ganadores.

Por su lado, el servidor mantiene un set que contiene las agencias que ya notificaron el envio de las apuestas. Cuando este set tiene 5 clientes (o la cantidad configurada), imprime por log que se realiza la apuesta y envia los resultados a todos los clientes que hicieron una conexión pidiendo los ganadores.

Siendo que son pocos los clientes conectados al sistema, y el tiempo de demora en cargar las apuestas es bajo, los clientes no hacen polling, sino que se mantiene la conexión hasta que el servidor envíe los resultados del sorteo. 

Siguiendo el protocolo de los dos anteriores puntos, el pedido por la apuesta consiste en un único mensaje que contiene un byte indicando que se piden los ganadores, dos bytes de un entero sin signo con el largo del string del identificador de la agencia. Y por último el respectivo string.

El servidor contesta con un único mensaje, el cual contiene dos bytes de un entero sin signo con la cantidad de ganadores. Y luego la lista de DNIs de los ganadores. Cada DNI se envía de igual modo que todas las anteriores cadenas de caracteres: dos bytes con el largo del string, y luego el documento.
