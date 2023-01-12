#!/bin/bash

set -e

echo Generating TLS certificate
docker-compose run generate-cert > /dev/null
docker-compose up --build
