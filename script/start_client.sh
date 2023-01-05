#!/usr/bin/env bash

BASE_DIR=${PWD}/..
BIN_DIR=${BASE_DIR}/deploy/bin

cd ${BIN_DIR}

if [[ -f "${BIN_DIR}/loader" ]]; then
  cd ${BIN_DIR}
  ./loader --dir-level=2 --conf-name=clientConfig/conf.ini &
  cd ${BASE_DIR}
else
  echo -e ${RED_PREFIX}"error"${COLOR_SUFFIX} ${YELLOW_PREFIX}${BIN_DIR}"/loader"${COLOR_SUFFIX}"\n"
fi