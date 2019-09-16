# docker-compose-mysql-master-master

A docker-compose example for a mysql master master setup

## docker-compose scripts

1. 查看状态 `docker-compose ps`
1. 强删集群 `docker-compose rm -fsv`
1. 启动集群 `docker-compose up`
1. 登录1号 `docker-compose exec mysqlmaster1 mysql -uroot -proot`
1. 登录2号 `docker-compose exec mysqlmaster2 mysql -uroot -proot`
1. 停止1号 `docker-compose stop mysqlmaster1`
1. 启动1号 `docker-compose start mysqlmaster1`
1. 检查集群 `./MySQLReplicationCheck.sh`

## 测试场景

1. 自增字段场景

    1. MySQLServer1插入t1表10条数据
    1. 检查MySQLServer2中t1表10条数据是否同步，自增字段取值情况
    1. MySQLServer2插入t1表10条数据
    1. 检查MySQLServer1中t1表10条数据是否同步，自增字段取值情况

1. 双向同步场景

    1. MySQLServer1插入t1表100条数据
    1. MySQLServer2插入t1表100条数据
    1. 各自检查同步状态，是否都是200条

1. 节点重启场景

    1. MySQLServer2停
    1. MySQLServer1插入t1表100条数据
    1. MySQLServer2启动
    1. 查看MySQLServer2中t1表同步状态

## test SQL scripts

```sql
select name from mysql.proc where name like 't1%';

call bjca.batch_t1(100);
select count(*) from bjca.t1;
select * from bjca.t1;
insert into bjca.t1(a) values(3);
```

## thanks

1. [mysql refman 5.7 Chapter 16 Replication](https://dev.mysql.com/doc/refman/5.7/en/replication.html)
1. [MySQL master slave using docker](https://tarunlalwani.com/post/mysql-master-slave-using-docker/) and its related [github rep](https://github.com/tarunlalwani/docker-compose-mysql-master-slave)
1. [MySQL Master Slave Docker部署例子](https://chanjarster.github.io/post/mysql-master-slave-docker-example/) and its related [github rep](https://github.com/chanjarster/mysql-master-slave-docker-example)
1. [玩转一下MySQL双主集群](https://github.com/bingoohuang/blog/issues/118)
1. [High-Availability MySQL cluster with load balancing using HAProxy and Heartbeat.](https://github.com/bingoohuang/docker-compose-mysql-master-master)

## tips

### 为什么需要单独建立docker的network

> What was happening was that the default docker network doesn't allow name >> DNS mapping.
> Containers on the default bridge network can only access each other by IP addresses, unless you use the --link option, which is considered legacy. On a user-defined bridge network, containers can resolve each other by name or alias.
>
> --[How to allow docker containers to see each other by their name?](https://serverfault.com/a/913075)
