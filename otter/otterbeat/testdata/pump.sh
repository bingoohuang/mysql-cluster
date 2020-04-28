#!/usr/bin/env bash

godoNum=${GODO_NUM}
rows=$(( RANDOM % (10 - 5 + 1 ) + 5 ))

if ((godoNum == 1)); then
    echo 写3311行数$rows
    ~/go/bin/pump -d "MYSQL_PWD=root mysql -h127.0.0.1 -P3311 -uroot" -t aa.tr_f_db -r $rows

elif ((godoNum == 2)); then
    echo 写3312行数$rows
    ~/go/bin/pump -d "MYSQL_PWD=root mysql -h127.0.0.1 -P3312 -uroot" -t aa.tr_f_db -r $rows

elif ((godoNum == 3)); then
    echo 检查总数
    ~/go/bin/pump -d "MYSQL_PWD=root mysql -h127.0.0.1 -P3311 -uroot" -s "select count(*) from aa.tr_f_db"    
    ~/go/bin/pump -d "MYSQL_PWD=root mysql -h127.0.0.1 -P3312 -uroot" -s "select count(*) from aa.tr_f_db"    
fi
