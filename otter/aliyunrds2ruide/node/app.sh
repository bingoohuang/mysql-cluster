#!/bin/bash

cp ./aria2c ./bin
cp ./otter.properties ./conf/otter.properties
cp ./startup.sh ./bin/startup.sh

if [ ! -f ./conf/nid ]; then
  echo NID?
  read NID
  echo $NID > ./conf/nid;
fi

rm -fr ./log/*

sh ./bin/startup.sh
