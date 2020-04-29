#!/bin/bash

rm -fr /app/bin/*.pid
cp /otter/otter.properties /app/conf/otter.properties
cp /otter/startup.sh /app/bin/startup.sh
sh /app/bin/startup.sh
