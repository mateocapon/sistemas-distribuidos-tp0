### Ejercicio N°1.1

### Enunciado:
Definir un script (en el lenguaje deseado) que permita crear una definición de DockerCompose con una cantidad configurable de clientes.

### Solución:
Se crea un archivo en el directorio root del proyecto, llamado `create-docker-compose-N.sh`. El mismo se corre con:

```bash
./create-docker-compose-N.sh <N_CLIENTS>
```

Sobreescribe el archivo llamado `docker-compose-dev.yaml`, creando una definición de DockerCompose con N_CLIENTS clientes.

Recordar darle los permisos de ejecución al archivo con 

```bash
chmod a+x create-docker-compose-N.sh
````
