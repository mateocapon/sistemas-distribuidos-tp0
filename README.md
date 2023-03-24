# Ejercicio  N°4

## Enunciado:
Modificar servidor y cliente para que ambos sistemas terminen de forma _graceful_ al recibir la signal SIGTERM. Terminar la aplicación de forma _graceful_ implica que todos los _file descriptors_ (entre los que se encuentran archivos, sockets, threads y procesos) deben cerrarse correctamente antes que el thread de la aplicación principal muera. Loguear mensajes en el cierre de cada recurso (hint: Verificar que hace el flag `-t` utilizado en el comando `docker compose down`).

## Solución:

Se implementa un cierre polite del cliente y del servidor. Esto es, cuando se recibe la señal SIGTERM se permite que termine la comunicación entre el cliente y el servidor, no se cierran los peer sockets inmediatamente.

### Servidor
Se agrega el handler `__stop_accepting` en el server, de tal modo que se haga el shutdown y close del socket aceptador, cuando se recibe la signal SIGTERM. De este modo, se dejan de aceptar conexiones nuevas. Una vez que la conexión con el cliente actual termine, el programa finaliza. 

### Cliente
Se agrega un channel por el cual se notifica la llegada de la señal SIGTERM. Ante la misma, se espera que termine el loop y se finaliza el programa.


### Probar funcionamiento

En una terminal correr los siguientes comandos:

```bash
make docker-compose-up
docker compose -f docker-compose-dev.yaml down --timeout 10
```

En otra terminal, ver los logs con

```bash
make docker-compose-logs
```

El flag `-t` (`--timeout`) en el `docker compose down` permite que se esperen unos segundos dados por el número pasado como argumento, hasta el envio de la señal SIGKILL al contenedor. Esta señal finaliza el proceso sin permitirle hacer una cerrada gracefull de los recursos. Por lo tanto, no puede ser manejada por el programa. Al momento de ejecutar `docker compose down` se envia la señal SIGTERM la cual sí puede ser manejada.