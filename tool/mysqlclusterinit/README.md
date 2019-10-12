# mysqlclusterinit

a tool to setup MySQL cluster-cluster and haproxy config.

## build

`env GOOS=linux GOARCH=amd64 go install ./...`

```bash
[root@BJCA-device ~]# ./mysqlclusterinit -m -c tool-config.toml
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
