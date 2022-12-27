#!/usr/bin/env bash

source ./proto_dir.cfg

# '@'/'*'
n=${#protos[@]}
echo -e '\ntotal' ${n} '\n'

# for ((i = 0; i < ${n}; i++)); do
# for ((i = 0; i < ${#protos[@]}; i++)); do
for i in ${!protos[@]}; do
  echo ${i} '=' ${protos[${i}]}
done
echo -e ""

# for ((i = 0; i < ${#protos[@]}; i++)); do
# for i in ${!protos[@]}; do
  # proto=${protos[$i]}
for proto in ${protos[@]}; do
  protoc -I./ -I../ -I../../ -I../../../ -I../../../../ --go_out=plugins=grpc:. ${proto}
  echo "protoc --go_out=plugins=grpc:." ${proto}
done
echo -e "ok\n"

i=0
for target in $(find ./github.com -name   "*.go"); do
    list[${i}]=${target}
    i=`expr ${i} + 1`
    echo 'src=' ${target}
done
echo -e "ok\n"

# for ((i = 0; i < ${#list[@]}; i++)); do
# for i in ${!list[@]}; do
  # target=${list[$i]}
for target in ${list[@]}; do
  parent=${target%/*}
  echo 'cp' ${target} '=>' ./${parent##*/}/${target##*/}
  cp ${target}  ./${parent##*/}/${target##*/}
done

rm -rf ./github.com