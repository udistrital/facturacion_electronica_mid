# Facturacion Electronica MID

API MID para la gestion/requerimientos asociados a temas de facturacion electronica.

## Especificaciones Técnicas

### Tecnologías Implementadas y Versiones
* [Golang](https://github.com/udistrital/introduccion_oas/blob/master/instalacion_de_herramientas/golang.md)
* [BeeGo](https://github.com/udistrital/introduccion_oas/blob/master/instalacion_de_herramientas/beego.md)
* [Docker](https://docs.docker.com/engine/install/ubuntu/)

### Variables de Entorno
```shell
export FACTURACION_ELECTRONICA_MID_HTTPPORT=8080
export FACTURACION_ELECTRONICA_MID_RUN_MODE=dev
export FACTURACION_ELECTRONICA_MID_PROTOCOL_ADMIN=http
export FACTURACION_ELECTRONICA_MID_TERCEROS_SERVICE=
export FACTURACION_ELECTRONICA_MID_BUSSERVICIOS_WSO2_SERVICE=
export FACTURACION_ELECTRONICA_MID_WSO2_TERCERO_PAGO_PATH=
export FACTURACION_ELECTRONICA_MID_CONSULTAR_RECIBO_PATH=
export FACTURACION_ELECTRONICA_MID_SOFIA_SERVICE=

```
**NOTA:** Las variables se pueden ver en el fichero conf/app.conf.

### Ejecución del Proyecto
```shell
#1. Obtener el repositorio con Go
go get github.com/udistrital/facturacion_electronica_mid

#2. Moverse a la carpeta del repositorio
cd $GOPATH/src/github.com/udistrital/facturacion_electronica_mid

# 3. Moverse a la rama **develop**
git pull origin develop && git checkout develop

# 4. alimentar todas las variables de entorno que utiliza el proyecto.
FACTURACION_ELECTRONICA_MID_HTTPPORT=8080 ...

# 5. Ejecutar comandos para descargar dependencias
go mod init && go get -t

# 6. Ejecutar proyecto
bee run -downdoc=true -gendoc=true
```

### Ejecución Dockerfile
```shell
# docker build --tag=facturacion_electronica_mid . --no-cache
# docker run -p 80:80 facturacion_electronica_mid
```

### Ejecución Pruebas

Pruebas unitarias
```shell
# En Proceso
```

## Estado CI

| Develop | Relese 0.0.1 | Master |
| -- | -- | -- |
| [![Build Status](https://hubci.portaloas.udistrital.edu.co/api/badges/udistrital/facturacion_electronica_mid/status.svg?ref=refs/heads/develop)](https://hubci.portaloas.udistrital.edu.co/udistrital/facturacion_electronica_mid) | [![Build Status](https://hubci.portaloas.udistrital.edu.co/api/badges/udistrital/facturacion_electronica_mid/status.svg?ref=refs/heads/release/0.0.1)](https://hubci.portaloas.udistrital.edu.co/udistrital/facturacion_electronica_mid) | [![Build Status](https://hubci.portaloas.udistrital.edu.co/api/badges/udistrital/facturacion_electronica_mid/status.svg)](https://hubci.portaloas.udistrital.edu.co/udistrital/facturacion_electronica_mid) |

## Licencia

This file is part of facturacion_electronica_mid.

facturacion_electronica_mid is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

facturacion_electronica_mid is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with sga_mid. If not, see https://www.gnu.org/licenses/.
