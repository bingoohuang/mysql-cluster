export MYSQL_PWD=$MYSQL_ROOT_PASSWORD;

mysql -h 127.0.0.1 -uroot -AN -vvv <<InputComesFromHERE
create database bjca;
create table bjca.t2(age int);
insert into bjca.t2 values(100);
InputComesFromHERE
