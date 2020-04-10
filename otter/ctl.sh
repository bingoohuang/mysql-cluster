#!/bin/bash

#set -x
#echo on

moreArgs="${*:2}"

if [[ "$1" == "up" ]]; then
  docker-compose -f otter.yml up
elif [[ "$1" == "upd" ]]; then
  docker-compose -f otter.yml up -d
elif [[ "$1" == "reup" ]]; then
  docker-compose -f otter.yml rm -fsv
  docker-compose -f otter.yml up
elif [[ "$1" == "reupd" ]]; then
  docker-compose -f otter.yml rm -fsv
  docker-compose -f otter.yml up -d
elif [[ "$1" == "rm" ]]; then
  docker-compose -f otter.yml rm -fsv
elif [[ "$1" == "ps" ]]; then
  docker-compose -f otter.yml ps
elif [[ "$1" == "exec" ]]; then
  docker-compose -f otter.yml exec ${moreArgs}
else
  echo "$0 up|upd|reup|reupd|rm"
fi
