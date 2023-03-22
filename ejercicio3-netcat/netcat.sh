#!/bin/bash

response=$(echo "Hi" | nc server "$SERVERPORT")

if [ "$response" == "Hi" ]; then
  echo "Server is working"
else
  echo "Server is not responding"
fi