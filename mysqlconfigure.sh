#!/usr/bin/env bash

echo "Waiting for mysql to get up"
# Give 30 seconds for master and slave to come up
sleep 15

echo "Create MySQL Servers (master / master repl)"


# https://dev.mysql.com/doc/refman/8.0/en/mysql-command-options.html#option_mysql_execute
# --skip-column-names, -N

# --no-auto-rehash, -A
# Enable automatic rehashing. This option is on by default, which enables database, table, and column name completion.

# --execute=statement, -e statement

# --verbose, -v

#Verbose mode. Produce more output about what the program does. This option can be given multiple times to produce more and more output.
#(For example, -v -v -v produces table output format even in batch mode.)

# Suppress warning messages from MySQL in shell script
# Warning: Using a password on the command line interface can be insecure.
# https://stackoverflow.com/a/24188878 https://unix.stackexchange.com/a/334971
export MYSQL_PWD=$MYSQL_PWD;

echo "* Create replication user"

mysql -h mysqlmaster1 -uroot -AN -vvv <<InputComesFromHERE
set GLOBAL max_connections=2000;
show variables like "%log_bin%";
CREATE USER '$MYSQL_REPL_USR'@'%';
GRANT REPLICATION SLAVE ON *.* TO '$MYSQL_REPL_USR'@'%' IDENTIFIED BY '$MYSQL_REPL_PWD';
InputComesFromHERE

mysql -h mysqlmaster2 -uroot -AN -vvv <<InputComesFromHERE
set GLOBAL max_connections=2000;
show variables like "%log_bin%";
CREATE USER '$MYSQL_REPL_USR'@'%';
GRANT REPLICATION SLAVE ON *.* TO '$MYSQL_REPL_USR'@'%' IDENTIFIED BY '$MYSQL_REPL_PWD';
InputComesFromHERE

echo "* Set MySQL2 as master on MySQL1"

mysql -h mysqlmaster1 -uroot -AN -vvv <<InputComesFromHERE
STOP SLAVE;
CHANGE MASTER TO master_host='mysqlmaster2', master_port=3306, master_user='$MYSQL_REPL_USR', master_password='$MYSQL_REPL_PWD', MASTER_AUTO_POSITION = 1;
START SLAVE;
InputComesFromHERE

echo "* Set MySQL1 as master on MySQL2"

mysql -h mysqlmaster2 -uroot -AN -vvv  <<InputComesFromHERE
STOP SLAVE;
CHANGE MASTER TO master_host='mysqlmaster1', master_port=3306, master_user='$MYSQL_REPL_USR', master_password='$MYSQL_REPL_PWD', MASTER_AUTO_POSITION = 1;
START SLAVE;
InputComesFromHERE

sleep 3

mysql -h mysqlmaster1 -uroot  -vvv -e "SHOW SLAVE STATUS \G"
mysql -h mysqlmaster2 -uroot  -vvv -e "SHOW SLAVE STATUS \G"

echo "MySQL servers created!"

MASTER1_IP=$(eval "getent hosts mysqlmaster1|awk '{print \$1}'")
MASTER2_IP=$(eval "getent hosts mysqlmaster2|awk '{print \$1}'")

echo $MASTER1_IP       : mysqlmaster1
echo $MASTER2_IP       : mysqlmaster2


