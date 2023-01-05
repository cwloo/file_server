#!/usr/bin/env bash

source ${PWD}/Makelist

BASE_DIR=${PWD}
export BIN_DIR=${BASE_DIR}/../deploy/bin

# '@'/'*'
n=${#makelist[@]}
echo -e ${BLUE_PREFIX}"total"${COLOR_SUFFIX} ${RED_PREFIX}${n}${COLOR_SUFFIX}"\n"

# for ((i = 0; i < ${n}; i++)); do
# for ((i = 0; i < ${#makelist[@]}; i++)); do
for i in ${!makelist[@]}; do
  echo -e ${BLUE_PREFIX}${i}${COLOR_SUFFIX}${SKY_BLUE_PREFIX}"="${COLOR_SUFFIX}${PURPLE_PREFIX}${BASE_DIR}/${makelist[${i}]}${COLOR_SUFFIX}
done

if [[ ! -d ${BIN_DIR} ]]; then
  echo -e ${BLUE_PREFIX}"mkdir -p"${COLOR_SUFFIX} ${SKY_BLUE_PREFIX}${BIN_DIR}${COLOR_SUFFIX}
  mkdir -p ${BIN_DIR}
else
  echo -e ""
fi

c=0
# for ((i = 0; i < ${#makelist[@]}; i++)); do
for i in ${!makelist[@]}; do
  makepath=${makelist[$i]}
#for makepath in ${makelist[@]}; do
  export MAKE_DIR=${BASE_DIR}/${makepath}
  if [[ -f "${MAKE_DIR}/Makefile" ]]; then
    echo -e "\n"
    make subsystem
    c=`expr $c + 1`
  else
    echo -e ${RED_PREFIX}"error"${COLOR_SUFFIX} ${YELLOW_PREFIX}${MAKE_DIR}"/Makefile"${COLOR_SUFFIX}"\n"
  fi
done

echo -e ${BLUE_PREFIX}"succ"${COLOR_SUFFIX} ${RED_PREFIX}${c}${COLOR_SUFFIX} "${BLUE_PREFIX}ok${COLOR_SUFFIX}\n"