#!/bin/bash

if [ -z "$1" ]
then
  echo "Usage: $0 <number of clients>"
  exit 1
fi

NUM_CLIENTS="$1"

echo "version: '3.9'
name: tp0
services:
  server:
    container_name: server
    image: server:latest
    entrypoint: python3 /main.py
    environment:
      - PYTHONUNBUFFERED=1
      - LOGGING_LEVEL=DEBUG
    networks:
      - testing_net
" > docker-compose-dev.yaml

for ((i=1;i<=$NUM_CLIENTS;i++)); do
echo "
  client$i:
    container_name: client$i
    image: client:latest
    entrypoint: /client
    environment:
      - CLI_ID=$i
      - CLI_LOG_LEVEL=DEBUG
    networks:
      - testing_net
    depends_on:
      - server
" >> docker-compose-dev.yaml
done

echo "
networks:
  testing_net:
    ipam:
      driver: default
      config:
        - subnet: 172.25.125.0/24
" >> docker-compose-dev.yaml

echo "Created docker-compose-dev.yaml with $NUM_CLIENTS clients."
