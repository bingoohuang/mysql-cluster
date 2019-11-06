MYSQL_PWD=$MYSQL_ROOT_PASSWORD mysql -u root -AN -vvv <<EOF
create database bjca;
create table bjca.t1(name varchar(100));
insert into bjca.t1 values('bingoo');
EOF
