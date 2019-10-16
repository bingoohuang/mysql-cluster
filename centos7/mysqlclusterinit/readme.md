# 初始化MySQL集群

1. 检查config.toml配置，调整参数
1. 初始化命令 `./mysqlclusterinit -c config.toml`
1. 检查命令 `./mysqlclusterinit -c config.toml -m`

示例输出：

```bash
[root@BJCA-device ~]# ./mysqlclusterinit -c config.toml
INFO[0000] execSQL SET GLOBAL server_id=1 completed
INFO[0000] execSQL DROP USER IF EXISTS 'repl'@'%' completed
INFO[0000] execSQL CREATE USER 'repl'@'%' IDENTIFIED BY 'BE30FD30-5e8f' completed
INFO[0000] execSQL GRANT REPLICATION SLAVE ON *.* TO 'repl'@'%' IDENTIFIED BY 'BE30FD30-5e8f' completed
INFO[0000] execSQL STOP SLAVE completed
INFO[0000] execSQL CHANGE MASTER TO master_host='192.168.136.23', master_port=3306, master_user='repl', master_password='BE30FD30-5e8f', master_auto_position = 1 completed
INFO[0000] execSQL START SLAVE completed
INFO[0000] createMySQCluster completed
INFO[0000] prepare to overwriteHAProxyCnf
listen mysql-rw
  bind 127.0.0.1:13306
  mode tcp
  option tcpka
  server mysql-1 192.168.136.22:3306 check inter 1s
  server mysql-2 192.168.136.23:3306 check inter 1s backup

listen mysql-ro
  bind 127.0.0.1:23306
  mode tcp
  option tcpka
  server mysql-1 192.168.136.22:3306 check inter 1s
  server mysql-2 192.168.136.23:3306 check inter 1s
INFO[0000] overwriteHAProxyCnf completed
[root@BJCA-device ~]# ./mysqlclusterinit -c config.toml -m
ShowSlaveStatus:{
	"SlaveIOState": "",
	"MasterHost": "192.168.136.23",
	"MasterUser": "repl",
	"MasterPort": 3306,
	"SlaveSQLRunningState": "Slave has read all relay log; waiting for more updates",
	"AutoPosition": true,
	"SlaveIoRunning": "No",
	"SlaveSQLRunning": "Yes",
	"MasterServerID": "1"
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
	"RelayLogInfoRepository": "TABLE",
	"InnodbVersion": "5.7.27"
}
```
