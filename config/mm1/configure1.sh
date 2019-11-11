MYSQL_PWD=$MYSQL_ROOT_PASSWORD mysql -u root -AN -vvv <<EOF
create database mm1;
create table mm1.t1(name varchar(100));
insert into mm1.t1 values('bingoo');

create database mm2;
create table mm2.t1(name varchar(100));
insert into mm2.t1 values('bingoo');

create database ms1;
create table ms1.t1(name varchar(100));
insert into ms1.t1 values('bingoo');
EOF
