#!/usr/bin/env bash
source ${PWD}/stop_server.sh

./stop_server.sh

BASE_DIR=${PWD}/..
BIN_DIR=${BASE_DIR}/deploy/bin

if [[ -f "${BIN_DIR}/loader" ]]; then
  cd ${BIN_DIR}
  ./loader --dir-level=2 --conf-name=config/conf.ini &
  cd ${BASE_DIR}
else
  echo -e ${RED_PREFIX}"error"${COLOR_SUFFIX} ${YELLOW_PREFIX}${BIN_DIR}"/loader"${COLOR_SUFFIX}"\n"
fi