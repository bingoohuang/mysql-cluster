# mysqlclusterinit

a tool to setup MySQL cluster-cluster and haproxy config.

## Build

`env GOOS=linux GOARCH=amd64 go install ./...`

## Usage

```bash
➜  mysqlclusterinit -h
Usage of mysqlclusterinit:
      --Debug                        Debug
      --HAProxyCfg string            HAProxyCfg
      --HAProxyRestartShell string   HAProxyRestartShell
      --LocalAddr string             LocalAddr
      --Master1Addr string           Master1Addr
      --Master2Addr string           Master2Addr
      --MySQLCnf string              MySQLCnf
      --Port int                     Port
      --ReplPassword string          ReplPassword
      --ReplUsr string               ReplUsr
      --RootPassword string          RootPassword
      --SlaveAddrs string            SlaveAddrs
  -m, --checkmysql                   check mysql
  -c, --config string                config file path (default "./config.toml")
  -v, --version                      show version
pflag: help requested
```

## Demo

```bash
➜  MCI_DEBUG=true mysqlclusterinit --LocalAddr=mysqlmaster1 --Master1Addr=mysqlmaster1 --Master2Addr=mysqlmaster2 --RootPassword=123
INFO[0000] config: {
	"Master1Addr": "mysqlmaster1",
	"Master2Addr": "mysqlmaster2",
	"SlaveAddrs": null,
	"RootPassword": "123",
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
# ./mysqlclusterinit -m -c tool-config.toml
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
