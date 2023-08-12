#!/bin/bash

if ! make uitest.up; then
  echo "Failed to start docker"
  make uitest.down
  exit 1
fi

testPath="$(pwd)/tests/ui"
scriptPath="$(pwd)/scripts"
seccompFile="${scriptPath}/seccomp_profile.json"

if ! docker run -it -v "${testPath}:/tests" \
  --network="docker_ui-test-network" \
  --rm --ipc=host \
  --user "$UID" --security-opt seccomp="${seccompFile}" \
  mcr.microsoft.com/playwright:v1.37.0-focal \
  /tests/run-docker.sh; then
    echo "UI tests failed"
    exit 1
fi

make uitest.down
