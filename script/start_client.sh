#!/usr/bin/env bash
source ${PWD}/stop_client.sh

./stop_client.sh

BASE_DIR=${PWD}/..
BIN_DIR=${BASE_DIR}/deploy/bin

if [[ -f "${BIN_DIR}/loader" ]]; then
  cd ${BIN_DIR}
  ./loader --dir-level=2 --conf-name=clientConfig/conf.ini &
  cd ${BASE_DIR}
else
  echo -e ${RED_PREFIX}"error"${COLOR_SUFFIX} ${YELLOW_PREFIX}${BIN_DIR}"/loader"${COLOR_SUFFIX}"\n"
fi