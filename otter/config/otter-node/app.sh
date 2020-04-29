#!/bin/bash

cp /otter/file /bin/
cp /otter/aria2c /bin/
cp /otter/otter.properties /app/conf/otter.properties
cp /otter/startup.sh /app/bin/startup.sh
echo $NID > /app/conf/nid;
rm -fr /app/log/*

sh /app/bin/startup.sh
