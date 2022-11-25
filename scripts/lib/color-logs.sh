#!/bin/bash

green () { echo -e "\033[01;32m$1:\033[0m $2"; }

warn () { echo -e "\033[01;33m[$OS]WARNING:\033[0m $1"; }

fatal () { echo -e "\033[01;31m[$OS]ERROR:\033[0m $1" >&2; exit 1; }