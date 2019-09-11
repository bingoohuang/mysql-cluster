# docker-compose-mysql-master-master
A docker-compose example for a mysql master master setup

## docker-compose scripts

1. 停止并强制删除相关容器实例 `docker-compose rm -fsv`
1. 启动容器 `docker-compose up`
1. 登录master1的MySQL服务 `docker-compose exec mysqlmaster1 mysql -uroot -proot`
1. 登录master2的MySQL服务 `docker-compose exec mysqlmaster2 mysql -uroot -proot`


## test SQL scripts

```sql
create database bjca;
user bjca;

CREATE TABLE `t1` (
  `id` int(11) NOT NULL,
  `a` int(11) DEFAULT NULL,
  `b` int(11) DEFAULT NULL,
  `c` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

select name from mysql.proc where name like 't1%';

DROP PROCEDURE IF EXISTS batch_t1;
DELIMITER $
CREATE PROCEDURE batch_t1()
BEGIN
    DECLARE i INT DEFAULT 1;
    DELETE FROM t1
    WHILE i<=10000 DO
        INSERT INTO t1(id,a,b,c) VALUES(i,i*2,i*3,i*4);
        SET i = i+1;
    END WHILE;
END $

CALL batch_test();

select count(*) from t1;

create table t1(id int auto_increment, `a` int(11) DEFAULT NULL,primary key(id));
insert into t1(a) values(3);
select * from t2;
```


## thanks

1. [MySQL master slave using docker](https://tarunlalwani.com/post/mysql-master-slave-using-docker/) and its related [github rep](https://github.com/tarunlalwani/docker-compose-mysql-master-slave)
1. [MySQL Master Slave Docker部署例子](https://chanjarster.github.io/post/mysql-master-slave-docker-example/) and its related [github rep](https://github.com/chanjarster/mysql-master-slave-docker-example)
1. [玩转一下MySQL双主集群](https://github.com/bingoohuang/blog/issues/118)

## tips

### 为什么需要单独建立docker的network

> What was happening was that the default docker network doesn't allow name >> DNS mapping.
> Containers on the default bridge network can only access each other by IP addresses, unless you use the --link option, which is considered legacy. On a user-defined bridge network, containers can resolve each other by name or alias.
> --[How to allow docker containers to see each other by their name?](https://serverfault.com/a/913075)