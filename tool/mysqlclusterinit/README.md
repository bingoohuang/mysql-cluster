# mysqlclusterinit

a tool to setup MySQL cluster-cluster and haproxy config.

## build

`env GOOS=linux GOARCH=amd64 go install ./...`


```bash
root@630f67855bfe:/# /tmp/mysqlclusterinit -c /tmp/tool-config.toml -m
ShowSlaveStatus:{SlaveIOState:Waiting for master to send event MasterHost:mysqlmaster2 MasterUser:repl MasterPort:3306 SlaveSqlRunningState:Slave has read all relay log; waiting for more updates AutoPosition:true SlaveIoRunning:Yes SlaveSqlRunning:Yes MasterServerId:2}
Variables:{ServerId:1 LogBin:ON SqlLogBin:ON GtidMode:ON GtidNext:AUTOMATIC SlaveSkipErrors:ALL BinlogFormat:ROW}
```