#!/bin/bash

cp ./aria2c ./bin
cp ./otter.properties ./conf/otter.properties
cp ./startup.sh ./bin/startup.sh
echo $NID > ./conf/nid;
rm -fr ./log/*

sh ./bin/startup.sh
