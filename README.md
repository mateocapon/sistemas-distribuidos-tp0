# Ejercicio N°5:

## Enunciado:
Modificar la lógica de negocio tanto de los clientes como del servidor para nuestro nuevo caso de uso.

#### Cliente
Emulará a una _agencia de quiniela_ que participa del proyecto. Existen 5 agencias. Deberán recibir como variables de entorno los campos que representan la apuesta de una persona: nombre, apellido, DNI, nacimiento, numero apostado (en adelante 'número'). Ej.: `NOMBRE=Santiago Lionel`, `APELLIDO=Lorca`, `DOCUMENTO=30904465`, `NACIMIENTO=1999-03-17` y `NUMERO=7574` respectivamente.

Los campos deben enviarse al servidor para dejar registro de la apuesta. Al recibir la confirmación del servidor se debe imprimir por log: `action: apuesta_enviada | result: success | dni: ${DNI} | numero: ${NUMERO}`.

#### Servidor
Emulará a la _central de Lotería Nacional_. Deberá recibir los campos de la cada apuesta desde los clientes y almacenar la información mediante la función `store_bet(...)` para control futuro de ganadores. La función `store_bet(...)` es provista por la cátedra y no podrá ser modificada por el alumno.
Al persistir se debe imprimir por log: `action: apuesta_almacenada | result: success | dni: ${DNI} | numero: ${NUMERO}`.

#### Comunicación:
Se deberá implementar un módulo de comunicación entre el cliente y el servidor donde se maneje el envío y la recepción de los paquetes, el cual se espera que contemple:
* Definición de un protocolo para el envío de los mensajes.
* Serialización de los datos.
* Correcta separación de responsabilidades entre modelo de dominio y capa de comunicación.
* Correcto empleo de sockets, incluyendo manejo de errores y evitando los fenómenos conocidos como [_short read y short write_](https://cs61.seas.harvard.edu/site/2018/FileDescriptors/).
* Límite máximo de paquete de 8kB.


## Solución:
Se implementa un modulo de protocolo tanto en el cliente como en el servidor. El protocolo de envio de un Bet consta de enviar cada atributo de la apuesta como un string, cuyo tamaño se anticipa con dos bytes de un entero sin signo en big endian. Ejemplifico el funcionamiento del protocolo de comunicación, usando la apuesta dada en el docker compose por el cliente 1.

- CLI_ID=1
- CLI_FIRSTNAME=Mateo
- CLI_LASTNAME=Capon Blanquer
- CLI_DOCUMENT=42496666
- CLI_BIRTHDATE=2000-03-03
- CLI_NUMBER=1234 

El packete se serializa en el orden de los atributos: id, nombre, apellido, documento, fecha de nacimiento y número. Por ejemplo, para serializar el atributo apellido, se toma el largo del apellido = len("Capon Blanquer") = 14. El largo se pasa a una tira de bytes ordenada en big endian de dos bytes, y luego se le adjunta el string correspondiente. Para este atributo queda la tira de bytes como:

- tira bytes del apellido = "0x00 0x0E 'Capon Blanquer'"
Este procedimiento se repite para todos los atributos y se almacena en una sola tira de bytes ordenada por atributo, la cual es enviada al servidor.


Por su lado, el servidor deserializa este mensaje, y en caso de que la apuesta pueda ejecutarse, envia el byte b'O' al cliente, para notificarle la ejecución exitosa. Si no se puede almacenar la apuesta, por ejemplo, por un error de formato en los datos que pasa el cliente, se envia un byte con un codigo de error b'E', seguido del largo de un string en dos bytes big endian y un string con el mensje correspondiente al error.

