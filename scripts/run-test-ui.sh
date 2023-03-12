#!/bin/bash

if ! docker-compose -f ./docker/dc-uitesting.yml up -d; then
  echo "Failed to start docker"
  docker-compose -f ./docker/dc-uitesting.yml down
  exit 1
fi

testPath="$(pwd)/tests/ui"
scriptPath="$(pwd)/scripts"
seccompFile="${scriptPath}/seccomp_profile.json"

docker run -it -v "${testPath}:/tests" \
  --rm --ipc=host \
  --user "$UID" --security-opt seccomp="${seccompFile}" \
  mcr.microsoft.com/playwright:v1.29.1-focal /tests/run-docker.sh

docker-compose -f ./docker/dc-uitesting.yml down