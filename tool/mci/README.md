# mci(mysqlclusterinit)

a tool to setup MySQL cluster-cluster and haproxy config.

## Build

`env GOOS=linux GOARCH=amd64 go install ./...`

`go fmt ./...&&goimports -w .&&golint ./...&&golangci-lint run --enable-all&&env GOOS=linux GOARCH=amd64 go install ./...&&upx ~/go/bin/linux_amd64/mci`

## Usage

```bash
🕙[ 21:17:15 ] ❯ mci -h 
Usage of mci:
      --Backup                       Backup (default true)
      --CheckSQL string              CheckSQL (default "select current_date()")
      --Debug                        Debug
      --HAProxyCfg string            HAProxyCfg (default "/etc/haproxy.cfg")
      --HAProxyRestartShell string   HAProxyRestartShell (default "systemctl restart haproxy")
      --Host string                  Host (default "127.0.0.1")
      --IPv6Enabled                  IPv6Enabled
      --LogLevel string              LogLevel
      --Master1Addr string           Master1Addr
      --Master2Addr string           Master2Addr
      --MySQLCmd string              MySQLCmd (default "mysql")
      --MySQLCnf string              MySQLCnf (default "/etc/my.cnf")
      --MySQLDSNParams string        MySQLDSNParams (default "timeout=120s&writeTimeout=120s&readTimeout=120s")
      --MySQLDumpCmd string          MySQLDumpCmd (default "mysqldump")
      --MySQLDumpOptions string      MySQLDumpOptions (default "--ignore-table=mysql.help_topic --ignore-table=mysql.help_keyword --ignore-table=mysql.help_relation --ignore-table=mysql.help_category")
      --MySQLRestartShell string     MySQLRestartShell (default "systemctl restart mysqld")
      --NoLog                        NoLog
      --Password string              Password
      --Port int                     Port (default 3306)
      --ReplPassword string          ReplPassword
      --ReplUsr string               ReplUsr (default "repl")
      --ShellTimeout string          ShellTimeout
      --SlaveAddrs string            SlaveAddrs
      --User string                  User (default "root")
      --checkmc string               check mysql cluster, format checkmc/json/table
      --checkmysql                   check mysql connection
  -c, --config string                config file path (default "./config.toml")
      --ebp strings                  PrintDecrypt by pbe, string or @file
      --pbe strings                  PrintEncrypt by pbe, string or @file
      --pbechg string                file to be change with another pbes
      --pbepwd string                pbe password
      --pbepwdnew string             new pbe pwd
  -r, --readips                      read haproxy server ips
      --removeSlaves string          remove slave nodes from cluster, eg 192.168.1.1,192.168.1.2
      --reset                        reset MySQL cluster
  -v, --version                      show version
pflag: help requested
```

## Demo

```bash
➜  MCI_DEBUG=true mysqlclusterinit --LocalAddr=mysqlmaster1 --Master1Addr=mysqlmaster1 --Master2Addr=mysqlmaster2 --Password=123
INFO[0000] config: {
	"Master1Addr": "mysqlmaster1",
	"Master2Addr": "mysqlmaster2",
	"SlaveAddrs": null,
	"Password": "123",
	"Port": 3306,
	"ReplUsr": "repl",
	"ReplPassword": "984d-CE5679F93918",
	"Debug": true,
	"LocalAddr": "mysqlmaster1",
	"MySQLCnf": "/etc/my.cnf",
	"HAProxyCfg": "/etc/haproxy/haproxy.cfg",
	"HAProxyRestartShell": "systemctl restart haproxy"
}

INFO[0000] SQL:SET GLOBAL server_id=1;
DROP USER IF EXISTS 'repl'@'%';
CREATE USER 'repl'@'%' IDENTIFIED BY '984d-CE5679F93918';
GRANT REPLICATION SLAVE ON *.* TO 'repl'@'%' IDENTIFIED BY '984d-CE5679F93918';
STOP SLAVE;
CHANGE MASTER TO master_host='mysqlmaster2', master_port=3306, master_user='repl', master_password='984d-CE5679F93918', master_auto_position = 1;
START SLAVE
INFO[0000] HAProxy:
listen mysql-rw
  bind 127.0.0.1:13306
  mode tcp
  option tcpka
  server mysql-1 mysqlmaster1:3306 check inter 1s
  server mysql-2 mysqlmaster2:3306 check inter 1s backup

listen mysql-ro
  bind 127.0.0.1:23306
  mode tcp
  option tcpka
  server mysql-1 mysqlmaster1:3306 check inter 1s
  server mysql-2 mysqlmaster2:3306 check inter 1s
➜  mysqlclusterinit git:(master) ✗
```

```bash
# ./mci -m -c tool-config.toml
ShowSlaveStatus:{
	"SlaveIOState": "Connecting to master",
	"MasterHost": "mysqlmaster2",
	"MasterUser": "repl",
	"MasterPort": 3306,
	"SlaveSQLRunningState": "Slave has read all relay log; waiting for more updates",
	"AutoPosition": true,
	"SlaveIoRunning": "Connecting",
	"SlaveSQLRunning": "Yes",
	"MasterServerID": "0"
}

Variables:{
	"ServerID": "1",
	"LogBin": "ON",
	"SQLLogBin": "ON",
	"GtidMode": "ON",
	"GtidNext": "AUTOMATIC",
	"SlaveSkipErrors": "ALL",
	"BinlogFormat": "ROW",
	"MasterInfoRepository": "TABLE",
	"RelayLogInfoRepository": "TABLE"
}
```


check MySQL availability

```bash
# ./mci --checkmysql --User=root --Password="A1765527-61a0" --Host=127.0.0.1 --Port=3306
netstat found cmd mysqld with pid 680
INFO[0000] mysql ds:root:A1765527-61a0@tcp(127.0.0.1:3306)/
SQL: select current_date()
cost: 945.128µs
+---+---------------------------+
| # | CURRENT_DATE()            |
+---+---------------------------+
| 1 | 2019-10-22T00:00:00+08:00 |
+---+---------------------------+
# echo $?
0
# ./mci --checkmysql --User=root --Password="A1765527-61a0" --Host=127.0.0.1 --Port=3307
NetstatListen error netstat  netstat -tunlp | grep ":3307" result empty
# echo $?
1
```


一些验证脚本

```bash
# 启动docker集群
docker-compose -f mci.yml rm -fsv && docker-compose -f mci.yml up
# 登录主1
docker-compose -f mci.yml exec mm1 bash
# 登录主2
docker-compose -f mci.yml exec mm2 bash
# 登录MySQL服务
MYSQL_PWD=root mysql -u root -P 3306

# 主1锁定表，准备导出数据 
FLUSH TABLES WITH READ LOCK;
# 主1导出数据
MYSQL_PWD=root mysqldump -u root -P 3306 -h mm1 --all-databases > mm1.sql
# 主2，从1，从2，...，从n 导入数据
MYSQL_PWD=root mysql -u root -P 3306 -h mm2 <  mm1.sql
# 主1解除表锁定 
UNLOCK TABLES;

# 执行mci工具，顺序: 主2，从1，从2，...，从n, 主1
/tmp/mci -c /tmp/mci.toml

MYSQL_PWD=root mysql -u root -P 3306  -vvv -e "SHOW SLAVE STATUS \G"

```

```sql
-- 主1上制造已有数据
create database bjca; create table bjca.t1(age int); insert into bjca.t1 values(100); select * from bjca.t1;
-- 主2上制造已有数据
create database bjca; create table bjca.t1(age int); insert into bjca.t1 values(200); select * from bjca.t1;

-- 重置主1
SET GLOBAL server_id=10001; STOP SLAVE; RESET SLAVE ALL; 
DROP USER IF EXISTS 'root'@'mm1'; REATE USER 'root'@'mm1' IDENTIFIED BY 'root';
GRANT ALL PRIVILEGES ON *.* TO 'root'@'mm1' WITH GRANT OPTION;
DROP USER IF EXISTS 'repl'@'%';
-- 重置主2
SET GLOBAL server_id=10002; STOP SLAVE; RESET SLAVE ALL; 
DROP USER IF EXISTS 'root'@'mm1'; REATE USER 'root'@'mm1' IDENTIFIED BY 'root';
GRANT ALL PRIVILEGES ON *.* TO 'root'@'mm1' WITH GRANT OPTION;
DROP USER IF EXISTS 'repl'@'%';

-- 重置主1 Master信息
RESET MASTER;
-- 重置主2 Master信息
RESET MASTER;

-- 主1主2 创建复制用户
CREATE USER 'repl'@'%' IDENTIFIED BY 'repl'; GRANT REPLICATION SLAVE ON *.* TO 'repl'@'%' IDENTIFIED BY 'repl';

-- 主1指向主2
CHANGE MASTER TO master_host='mm2', master_port=3306, master_user='repl', master_password='repl', master_auto_position = 1;
-- 主2指向主1
CHANGE MASTER TO master_host='mm1', master_port=3306, master_user='repl', master_password='repl', master_auto_position = 1;

-- 主1主2 启动复制进程
START SLAVE;
-- 查看复制状态
SHOW SLAVE STATUS\G

-- 主1上
insert into bjca.t1 values(101);
create table bjca.t2(age int);
-- 主2上
select * from bjca.t1;
insert into bjca.t1 values(200);
insert into bjca.t2 values(200);
-- 主1上
select * from bjca.t2;

```

## Problems

1. [master and slave have equal MySQL server UUIDs](doc/serverUUID.md)