#!/usr/bin/env bash

# Just a simple Mysql Replication Health Check script I wrote. You can put this in a cron.
# from https://gist.github.com/OliverBailey/66ea8a143b027207932e5febc476a251

# login-path是MySQL5.6开始支持的新特性。通过mysql_config_editor工具将登录MySQL服务器的认证信息加密保存在~/.mylogin.cnf文件v中。
# 之后MySQL客户端工具可通过读取该加密文件连接MySQL，避免重复输入登录信息，避免敏感信息暴露。
# URL: https://opensourcedbms.com/dbms/passwordless-authentication-using-mysql_config_editor-with-mysql-5-6/
# mysql_config_editor print --all
# mysql_config_editor remove --login-path=mysql_login
# remove parts of your profile like hostname from existing profile, like host:
# mysql_config_editor remove --login-path=testbed1 --host
# mysql_config_editor set --login-path=mysql_login --host=127.0.0.1 --port=33061 --user=root --password
# mysql --login-path=mysql_login

### VARIABLES ###
EMAIL=""
SERVER=$(hostname)
MYSQL_CHECK=$(mysql --login-path=mysql_login -e "SHOW VARIABLES LIKE '%version%';" || echo 1)
STATUS_LINE=$(mysql --login-path=mysql_login -e "SHOW SLAVE STATUS\G")"1"
LAST_ERRNO=$(grep "Last_Errno" <<< "$STATUS_LINE" | awk '{ print $2 }')
SECONDS_BEHIND_MASTER=$( grep "Seconds_Behind_Master" <<< "$STATUS_LINE" | awk '{ print $2 }')
IO_IS_RUNNING=$(grep "Slave_IO_Running:" <<< "$STATUS_LINE" | awk '{ print $2 }')
SQL_IS_RUNNING=$(grep "Slave_SQL_Running:" <<< "$STATUS_LINE" | awk '{ print $2 }')
MASTER_LOG_FILE=$(grep " Master_Log_File" <<< "$STATUS_LINE" | awk '{ print $2 }')
RELAY_MASTER_LOG_FILE=$(grep "Relay_Master_Log_File" <<< "$STATUS_LINE" | awk '{ print $2 }')
ERRORS=()
SUBJECT="Errors on SQL Replication"
FILENAME=$(date +"%Y/%m/%d")
DATE=$(date +"%Y-%m-%d")
BACKUP="Backup complete"
### Run Some Checks ###

## Check if I can connect to Mysql ##
if [[ "$MYSQL_CHECK" == 1 ]]; then
    ERRORS=("${ERRORS[@]}" "Can't connect to MySQL (Check Pass)")
fi

## Check For Last Error ##
if [[ "$LAST_ERRNO" != 0 ]]; then
    LAST_ERROR=$(mysql --login-path=mysql_login -e "SHOW SLAVE STATUS\G" | grep "Last_Error" | awk '{ print $2 }')
    ERRORS=("${ERRORS[@]}" "Error when processing relay log (Last_Errno = $LAST_ERRNO)")
    ERRORS=("${ERRORS[@]}" "(Last_Error = $LAST_ERROR)")
fi

## Check if IO thread is running ##
if [[ "$IO_IS_RUNNING" != "Yes" ]]; then
    ERRORS=("${ERRORS[@]}" "I/O thread for reading the master's binary log is not running (Slave_IO_Running)")
fi

## Check for SQL thread ##
if [[ "$SQL_IS_RUNNING" != "Yes" ]]; then
    ERRORS=("${ERRORS[@]}" "SQL thread for executing events in the relay log is not running (Slave_SQL_Running)")
fi

## Check how slow the slave is ##
if [[ "$SECONDS_BEHIND_MASTER" == "NULL" ]]; then
    ERRORS=("${ERRORS[@]}" "The Slave is reporting 'NULL' (Seconds_Behind_Master)")
elif [[ "$SECONDS_BEHIND_MASTER" > 60 ]]; then
    ERRORS=("${ERRORS[@]}" "The Slave is at least 60 seconds behind the master (Seconds_Behind_Master)")
fi

totalErrors=${#ERRORS[@]}
### Send an Email if there is an error ###
if [[ "$totalErrors" -gt 0 ]]; then
    MESSAGE="$totalErrors errors detected on ${SERVER} involving the mysql replication. Below is a list of the reported errors:\n\n
    $(for i in $(seq 1 $totalErrors) ; do echo "\t$i: ${ERRORS[(($i-1))]}\n" ; done)
    \nPlease correct this ASAP!
    "
    # https://explainshell.com/explain?cmd=echo+-e+%24MESSAGE
    # -e     enable interpretation of backslash escapes
    echo -e ${MESSAGE}
#/usr/sbin/ssmtp -t << EOF
#to: emailaddress@email.co.uk
#from: emailaddress@email.co.uk
#subject: $SUBJECT
#$MESSAGE
#EOF
## If no error, then will do the below to grab a dump, upload it to S3 and then remove it from the local file. Follows the naming convention of the two variables $FILENAME and $DATE.
elif [[ "$totalErrors" -le 0 ]]; then
    echo "Everything is OK!"
#    mysqldump --login-path=mysql_login --single-transaction database_name | gzip > application.database_name.$DATE.sql.gz
#    aws s3 cp application.database_name.$DATE.sql.gz s3://aws-bucket-name/$FILENAME/application.database_name.$DATE.sql.gz
#    rm application.database_name.$DATE.sql.gz
#
#/usr/sbin/ssmtp -t << EOF
#to: emailaddress@email.co.uk
#from: emailaddress@email.co.uk
#subject: $BACKUP
#Backup now complete.
#EOF

fi
