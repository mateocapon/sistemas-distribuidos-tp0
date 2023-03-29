# Ejercicio  N°2

## Enunciado:
Modificar el cliente y el servidor para lograr que realizar cambios en el archivo de configuración no requiera un nuevo build de las imágenes de Docker para que los mismos sean efectivos. La configuración a través del archivo correspondiente (`config.ini` y `config.yaml`, dependiendo de la aplicación) debe ser inyectada en el container y persistida afuera de la imagen (hint: `docker volumes`).

## Solución:

Se elimina la linea que copia el archivo de configuración al construir la imagen en el Dockerfile del cliente. Además, se agregan los archivos `.dockerignore` tanto para el cliente como para el servidor, para que no se copien los archivos de configuración en el proceso de build de las imágenes.

Se modifica el DockerCompose, agregando a los archivos de configuración como volumes para el servidor y el cliente respectivamente. De esto modo, no es necesario hacer un rebuild de la imagen si uno quiere levantar un contenedor con un archivo de configuración modificado.  

Si se quiere probar el funcionamiento, se puede levantar el DockerCompose dos veces, teniendo en la segunda vez los archivos de configuración modificados. Se verá que se usa el cache para cada layer de las imágenes, tal como se ve en la siguiente captura.

![image](https://user-images.githubusercontent.com/65830097/228651452-5fc1a1c8-b84b-4cff-b3e4-9b55c014271a.png)

Para verificar que las variables de configuración fueron modificadas, luego de ejecutar `make docker-compose-up`, observar los logs con 

```bash
make docker-compose-logs | grep config
```

Las siguientes capturas muestran la salida de este comando para una primera corrida de `make docker-compose-up`, y para una segunda, luego de modificar las variables SERVER_LISTEN_BACKLOG (en el servidor) y loop lapse (en el cliente).


![image](https://user-images.githubusercontent.com/65830097/228653198-6fef9f75-2f28-4a0b-a38e-70c324f0332b.png)


![image](https://user-images.githubusercontent.com/65830097/228652907-161f3ef4-9c5b-41b9-b596-d647d68f9cd4.png)

