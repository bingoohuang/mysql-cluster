#!/bin/bash

#set -x
#echo on

moreArgs="${*:2}"

map="up:docker-compose -f otter.yml up, \
     upd:docker-compose -f otter.yml up -d, \
     exec:docker-compose -f otter.yml exec ${moreArgs}, \
     ps: docker-compose -f otter.yml ps, \
     rm: docker-compose -f otter.yml rm -fsv, \
     reup: docker-compose -f otter.yml rm -fsv; docker-compose -f otter.yml up;, \
     reupd: docker-compose -f otter.yml rm -fsv; docker-compose -f otter.yml up -d;"

getRealCommand(){
   local dataMap=$1
   local key=$2
   NL="\n"
   if [[ $uname -eq "Darwin" ]]; then
       NL=$'\\\n'
   fi
   echo $dataMap | sed "s|,|${NL}|g"|sed 's/^[ \t]*//g'|grep -E "^$key:"|awk -F":" '{print $2}'|sed 's/^[ \t]*//g;/^[ \t]*$/d'
}     	

realCommand=`getRealCommand "$map" $1`
echo $realCommand

[ -z "$realCommand" ] && echo "command not exits" && exit 1

bash -c "$realCommand"


#if [[ "$1" == "up" ]]; then
 # docker-compose -f otter.yml up
#elif [[ "$1" == "upd" ]]; then
#  docker-compose -f otter.yml up -d
#elif [[ "$1" == "reup" ]]; then
 # docker-compose -f otter.yml rm -fsv
  #docker-compose -f otter.yml up
#elif [[ "$1" == "reupd" ]]; then
 # docker-compose -f otter.yml rm -fsv
  #docker-compose -f otter.yml up -d
#elif [[ "$1" == "rm" ]]; then
 # docker-compose -f otter.yml rm -fsv
#elif [[ "$1" == "ps" ]]; then
 # docker-compose -f otter.yml ps
#elif [[ "$1" == "exec" ]]; then
 # docker-compose -f otter.yml exec ${moreArgs}
#else
 # echo "$0 up|upd|reup|reupd|rm"
#fi
