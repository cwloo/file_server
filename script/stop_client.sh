#!/usr/bin/env bash

BASE_DIR=${PWD}/..
BIN_DIR=${BASE_DIR}/deploy/bin

stoplist=(
    file_client
)

kill -9 $(pgrep -f './loader --dir-level=2 --conf-name=clientConfig/conf.ini')

for i in ${!stoplist[@]}; do
  execname=${stoplist[$i]}
  kill -9 $(pgrep -f ${execname})
  #kill -9 $(pgrep -f '^${execname}$')
done