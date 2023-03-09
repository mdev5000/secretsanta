#!/bin/bash

if ! docker-compose -f ./docker/dc-uitesting.yml up -d; then
  echo "Failed to start docker"
  docker-compose -f ./docker/dc-uitesting.yml down
  exit 1
fi

testPath="$(pwd)/tests/ui"
docker run --network=docker_ui-test-network -v "${testPath}:/tests" mcr.microsoft.com/playwright:v1.31.0-focal /tests/run-docker.sh
docker-compose -f ./docker/dc-uitesting.yml down