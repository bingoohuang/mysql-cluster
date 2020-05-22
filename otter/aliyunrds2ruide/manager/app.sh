#!/bin/bash

rm -fr ./bin/*.pid
cp ./otter.properties ./conf/otter.properties
cp ./startup.sh ./bin/startup.sh
sh ./bin/startup.sh
