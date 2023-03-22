# Ejercicio  N°3

## Enunciado:
Crear un script que permita verificar el correcto funcionamiento del servidor utilizando el comando `netcat` para interactuar con el mismo. Dado que el servidor es un EchoServer, se debe enviar un mensaje al servidor y esperar recibir el mismo mensaje enviado. Netcat no debe ser instalado en la máquina _host_ y no se puede exponer puertos del servidor para realizar la comunicación (hint: `docker network`).

## Solución:
Se crea una carpeta `ejercicio3-netcat` la cual cuenta con un Dockerfile, un script bash y un archivo de configuración. El Dockerfile descarga netcat desde una imagen base de alpine y copia el script. El script simplemente envia un mensaje tcp a la dirección y puerto del server, e imprime por consola el resultado del echo request.

Se agrega un comando al Makefile para que se pueda testear este ejercicio de un modo más simple utilizando:

```bash
make netcat-server
```

Esto hace el build de la imagen en `ejercicio3-netcat/Dockerfile`, y luego la corre en la red `tp0_testing_net`, agregando las variables de entorno del archivo de configuración, en este caso, sólamente el puerto del servidor.