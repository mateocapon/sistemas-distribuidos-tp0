# Ejercicio  N°6

## Enunciado:
Modificar los clientes para que envíen varias apuestas a la vez (modalidad conocida como procesamiento por _chunks_ o _batchs_). La información de cada agencia será simulada por la ingesta de su archivo numerado correspondiente, provisto por la cátedra dentro de `.data/datasets.zip`.
Los _batchs_ permiten que el cliente registre varias apuestas en una misma consulta, acortando tiempos de transmisión y procesamiento. La cantidad de apuestas dentro de cada _batch_ debe ser configurable.
El servidor, por otro lado, deberá responder con éxito solamente si todas las apuestas del _batch_ fueron procesadas correctamente.

## Solución:
Se agrega el modulo `betsreader` en el cliente. Este modulo se encarga de leer su archivo y apoyarse en el protocolo para enviar los chunks de apuestas.
El protocolo por su lado fue modificado para que envie múltiples apuestas en un solo paquete. Si el cliente dejará de enviar apuestas, se envia un FLAG dentro del último paquete de apuestas, notificandole al servidor de esta situación.

Por otro lado, se monta el archivo correspondiente a cada cliente como un volumen de docker. Es necesario **descomprimir el archivo en el path.data/dataset.zip**. Si no se hace este paso, no se podrán levantar los clientes.
