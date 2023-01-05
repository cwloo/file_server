#!/usr/bin/env bash

BASE_DIR=${PWD}/..
BIN_DIR=${BASE_DIR}/deploy/bin

execlist=(
    loader
    http_gate
    file_server
    file_client
)

for i in ${!execlist[@]}; do
  execname=${execlist[$i]}
  if [[ -f "${BIN_DIR}/${execname}" ]]; then
    ps -ef|grep ${execname}
  fi
done