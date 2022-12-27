#!/usr/bin/env bash

source ./proto_dir.cfg

for ((i = 0; i < ${#all_proto[*]}; i++)); do
  proto=${all_proto[$i]}
  protoc -I./ -I../ -I../../ -I../../../ -I../../../../ --go_out=plugins=grpc:. $proto
  echo "protoc --go_out=plugins=grpc:." $proto
done
echo "ok"

i=0
for file in $(find ./github.com -name   "*.go"); do
    filelist[i]=$file
    i=`expr $i + 1`
    echo 'src=' $file
done

for ((i = 0; i < ${#filelist[*]}; i++)); do
  proto=${filelist[$i]}
  parent=${proto%/*}
  echo 'cp' $proto '=>' ./${parent##*/}/${proto##*/}
  cp $proto  ./${parent##*/}/${proto##*/}
done

rm -rf ./github.com
