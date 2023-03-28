# Ejercicio  N°6

## Enunciado:
Modificar los clientes para que envíen varias apuestas a la vez (modalidad conocida como procesamiento por _chunks_ o _batchs_). La información de cada agencia será simulada por la ingesta de su archivo numerado correspondiente, provisto por la cátedra dentro de `.data/datasets.zip`.
Los _batchs_ permiten que el cliente registre varias apuestas en una misma consulta, acortando tiempos de transmisión y procesamiento. La cantidad de apuestas dentro de cada _batch_ debe ser configurable.
El servidor, por otro lado, deberá responder con éxito solamente si todas las apuestas del _batch_ fueron procesadas correctamente.

## Solución:
Se agrega el modulo `betsreader` en el cliente. Este modulo se encarga de leer su archivo y apoyarse en el protocolo para enviar los chunks de apuestas.
El protocolo por su lado fue modificado para que envie múltiples apuestas en un solo paquete. Cuando el cliente está por finalizar el envio de apuestas, envia un FLAG dentro del último paquete de apuestas, notificandole al servidor de esta situación.

El protocolo de envio de cada chunk consta de un "header" el cual contiene un byte indicando el tipo de chunk (ultimo chunk o no), dos bytes indicando la cantidad de apuestas a leer, dos bytes con el largo del string del identificador y luego la cadena de caracteres que identifica a la agencia. Luego, en el "payload", se envian las apuestas de igual modo que en el ejercicio 5. 

Por otro lado, se monta el archivo correspondiente a cada cliente como un volumen de docker. Es necesario **descomprimir la carpeta en el path .data/dataset.zip**, quedando los archivos en el path `.data/dataset/agency-<N>.csv`. Si no se hace este paso, no se podrán levantar los clientes.
