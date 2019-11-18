# 即主从架构中使用了相同的UUID

最近在部署MySQL主从复制架构的时候，碰到了"Last_IO_Error: Fatal error: The slave I/O thread stops because master and slave have equal MySQL server UUIDs;  these UUIDs must be different for replication to work." 这个错误提示。即主从架构中使用了相同的UUID。检查server_id系统变量，已经是不同的设置，那原因是？接下来为具体描述。 

## 错误消息

```bash
mysql> show slave staus;
 
Last_IO_Error: Fatal error: The slave I/O thread stops because master and slave have equal MySQL server UUIDs; 
these UUIDs must be different for replication to work.
```
     
## 查看主从的server_id变量

```bash
master_mysql> show variables like 'server_id';
+---------------+-------+
| Variable_name | Value |
+---------------+-------+
| server_id     | 33    |
+---------------+-------+
 
slave_mysql> show variables like 'server_id';
+---------------+-------+
| Variable_name | Value |
+---------------+-------+
| server_id     | 11    |
+---------------+-------+
```

从上面的情形可知，主从mysql已经使用了不同的server_id
 
## 解决故障


###查看auto.cnf文件

```bash
[root@dbsrv1 ~] cat /data/mysqldata/auto.cnf  ### 主上的uuid
[auto]
server-uuid=62ee10aa-b1f7-11e4-90ae-080027615026
 
[root@dbsrv2 ~]# more /data/mysqldata/auto.cnf ###从上的uuid，果然出现了重复，原因是克隆了虚拟机，只改server_id不行
[auto]
server-uuid=62ee10aa-b1f7-11e4-90ae-080027615026
 
[root@dbsrv2 ~]# mv /data/mysqldata/auto.cnf  /data/mysqldata/auto.cnf.bk  ###重命名该文件
[root@dbsrv2 ~]# service mysql restart          ###重启mysql
Shutting down MySQL.[  OK  ]
Starting MySQL.[  OK  ]
[root@dbsrv2 ~]# more /data/mysqldata/auto.cnf  ###重启后自动生成新的auto.cnf文件，即新的UUID
[auto]
server-uuid=6ac0fdae-b5d7-11e4-a9f3-0800278ce5c9
```

1. 老的重启方式：`/etc/init.d/mysql restart`
1. 新的重启方式: `systemctl restart mysqld`

###再次查看slave的状态已经正常

```bash
[root@dbsrv1 ~]# mysql -uroot -pxxx -e "show slave status\G"|grep Running
Warning: Using a password on the command line interface can be insecure.
             Slave_IO_Running: Yes
            Slave_SQL_Running: Yes
      Slave_SQL_Running_State: Slave has read all relay log; waiting for the slave I/O thread to update it
```

###主库端查看自身的uuid

```bash
master_mysql> show variables like 'server_uuid';
+---------------+--------------------------------------+
| Variable_name | Value                                |
+---------------+--------------------------------------+
| server_uuid   | 62ee10aa-b1f7-11e4-90ae-080027615026 |
+---------------+--------------------------------------+
1 row in set (0.00 sec)
```

###主库端查看从库的uuid

```bash
master_mysql> show slave hosts;
+-----------+------+------+-----------+--------------------------------------+
| Server_id | Host | Port | Master_id | Slave_UUID                           |
+-----------+------+------+-----------+--------------------------------------+
|        33 |      | 3306 |        11 | 62ee10aa-b1f7-11e4-90ae-080027615030 |
|        22 |      | 3306 |        11 | 6ac0fdae-b5d7-11e4-a9f3-0800278ce5c9 |
+-----------+------+------+-----------+--------------------------------------+
```

## 延生参考

### 有关server_id的描述

> The server ID, used in replication to give each master and slave a unique identity. This variable is set
> by the --server-id option. For each server participating in replication, you should pick a
> positive integer in the range from 1 to 232– 1(2的32次方减1) to act as that server's ID.

 

### 有关 server_uuid的描述

> Beginning with MySQL 5.6, the server generates a true UUID in addition to the --server-id
> supplied by the user. This is available as the global, read-only variable server_uuid(全局只读变量)

> When starting, the MySQL server automatically obtains a UUID as follows:
> a).  Attempt to read and use the UUID written in the file data_dir/auto.cnf (where data_dir is
> the server's data directory); exit on success.
> b). Otherwise, generate a new UUID and save it to this file, creating the file if necessary.
> The auto.cnf file has a format similar to that used for my.cnf or my.ini files. In MySQL 5.6,
> auto.cnf has only a single [auto] section containing a single server_uuid [1992] setting and
> value;
> 
> 
> `Important
> The auto.cnf file is automatically generated; you should not attempt to write
> or modify this file`
> 
> 
> Also beginning with MySQL 5.6, when using MySQL replication, masters and slaves know one
> another's UUIDs. The value of a slave's UUID can be seen in the output of SHOW SLAVE HOSTS. Once
> START SLAVE has been executed (but not before), the value of the master's UUID is available on the
> slave in the output of SHOW SLAVE STATUS.
> 
> In MySQL 5.6.5 and later, a server's server_uuid is also used in GTIDs for transactions originating
> on that server. For more information, see Section 16.1.3, “Replication with Global Transaction
————————————————

[原文链接](https://blog.csdn.net/leshami/article/details/43854505)


