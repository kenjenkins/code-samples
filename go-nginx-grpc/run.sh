#!/bin/bash

set -e

openssl req -x509 -nodes -days 30 -newkey rsa:2048 \
  -subj '/CN=localhost' -keyout ssl.key -out ssl.cert
docker-compose up --build
