#!/bin/bash

# Check if WALLET_PASSWORD is not set
if [ -z "$WALLET_PASSWORD" ]; then
  echo "WALLET_PASSWORD is not set, please set it and run the script again."
	echo "Create a wallet with coins and set WALLET_PASSWORD."
  exit 1
fi

export GITHUB_ACTIONS=false

task generate
task build
task run &

cd api/test
robot robot_tests
