# Ejercicio  N°2

## Enunciado:
Modificar el cliente y el servidor para lograr que realizar cambios en el archivo de configuración no requiera un nuevo build de las imágenes de Docker para que los mismos sean efectivos. La configuración a través del archivo correspondiente (`config.ini` y `config.yaml`, dependiendo de la aplicación) debe ser inyectada en el container y persistida afuera de la imagen (hint: `docker volumes`).

## Solución:

Se elimina la linea que copia el archivo de configuración al construir la imagen en el Dockerfile del cliente. Además, se agregan los archivos `.dockerignore` tanto para el cliente como para el servidor, para que no se copien los archivos de configuración en el proceso de build de las imágenes.

Se modifica el DockerCompose, agregando a los archivos de configuración como volumes para el servidor y el cliente respectivamente. De esto modo, no es necesario hacer un rebuild de la imagen si uno quiere levantar un contenedor con un archivo de configuración modificado.  

