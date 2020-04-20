#!/bin/bash

#set -x
#echo on

moreArgs="${*:2}"

map='up:up, \
     upd:up -d, \
     exec:exec ${moreArgs}, \
     ps: ps, \
     rm: rm -fsv, \
     reup: rm -fsv;up;, \
     reupd: m -fsv;up -d;'

getRealCommand(){
   subCommand=$1
   moreArgs="${*:2}"
   NL="\n"
   if [[ $uname -eq "Darwin" ]]; then
       NL=$'\\\n'
   fi
   echo $map | sed "s|,|${NL}|g"|sed 's/^[ \t]*//g'|grep -E "^$subCommand:"|awk -F":" '{print $2}'|sed "s|;|${NL}|g"|sed 's/^[ \t]*//g;/^[ \t]*$/d'|sed 's/^/docker-compose -f otter.yml /'
}     	

realCommand=`getRealCommand $*`
echo $realCommand
#realCommand=${map["$subCommand"]}
#echo ""----"${!map[@]}"
#eho "command="$subCommand $realCommand$￥[ -z "$realCommand" ] && echo "command not exits" && exit 1

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
