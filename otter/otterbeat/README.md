# otterbeat

## Build

1. `make install` and the executable is installed in `$GOPATH/bin`.
1. `otterbeat --init` to generate the demo cnf.toml and ctl files at the current directory.
1. edit the cnf.toml file.
1. `./ctl start` to startup the otterbeat.
1. `./ctl tail` to trace the log of `~/logs/otterbeat/otterbeat.log`.

## 从页面采集指标

<details><summary>Pipeline管理页面</summary>

![image](https://user-images.githubusercontent.com/1940588/79715445-50641480-8306-11ea-9a5a-cb5322cff428.png)

</details>

## 从数据库采集指标

表名|含义|更新频率|使用方式
---|---|---|---
DELAY_STAT | 延时统计表 | 60秒 |  按GMT_MODIFIED取变化值
LOG_RECORD | 日志记录表 | NA | 按ID取
TABLE_HISTORY_STAT | 同步明细统计表 |  NA | 按GMT_MODIFIED取变化值
TABLE_STAT |表同步汇总统计表 |  NA | 按GMT_MODIFIED取变化值
THROUGHPUT_STAT | 同步流量统计表 |  60秒 | 按GMT_MODIFIED取变化值

<details><summary>DELAY_STAT</summary>

```sql
CREATE TABLE `DELAY_STAT` (
  `ID` bigint(20) NOT NULL AUTO_INCREMENT,
  `DELAY_TIME` bigint(20) NOT NULL,
  `DELAY_NUMBER` bigint(20) NOT NULL,
  `PIPELINE_ID` bigint(20) NOT NULL,
  `GMT_CREATE` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
  `GMT_MODIFIED` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`ID`),
  KEY `idx_PipelineID_GmtModified_ID` (`PIPELINE_ID`,`GMT_MODIFIED`,`ID`),
  KEY `idx_Pipeline_GmtCreate` (`PIPELINE_ID`,`GMT_CREATE`),
  KEY `idx_GmtCreate_id` (`GMT_CREATE`,`ID`)
) ENGINE=InnoDB AUTO_INCREMENT=21 DEFAULT CHARSET=utf8
```

ID|DELAY_TIME|DELAY_NUMBER|PIPELINE_ID|GMT_CREATE           |GMT_MODIFIED         |
--|----------|------------|-----------|---------------------|---------------------|
 1|      1772|           0|          1|2020-04-20 02:49:57.0|2020-04-20 02:49:57.0|
 3|      6204|           0|          3|2020-04-20 02:49:57.0|2020-04-20 02:49:57.0|
 5|     14246|           0|          1|2020-04-20 02:50:57.0|2020-04-20 02:50:57.0|
 7|     13376|           0|          3|2020-04-20 02:50:57.0|2020-04-20 02:50:57.0|
 9|       296|           0|          1|2020-04-20 02:51:57.0|2020-04-20 02:51:57.0|
11|       548|           0|          1|2020-04-20 04:13:11.0|2020-04-20 04:13:11.0|
13|       713|           0|          1|2020-04-20 04:20:11.0|2020-04-20 04:20:11.0|

</details>

<details><summary>LOG_RECORD</summary>

```sql
CREATE TABLE `LOG_RECORD` (
  `ID` bigint(20) NOT NULL AUTO_INCREMENT,
  `NID` varchar(200) DEFAULT NULL,
  `CHANNEL_ID` varchar(200) NOT NULL,
  `PIPELINE_ID` varchar(200) NOT NULL,
  `TITLE` varchar(1000) DEFAULT NULL,
  `MESSAGE` text,
  `GMT_CREATE` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
  `GMT_MODIFIED` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`ID`),
  KEY `logRecord_pipelineId` (`PIPELINE_ID`)
) ENGINE=InnoDB AUTO_INCREMENT=33 DEFAULT CHARSET=utf8
```

ID|NID|CHANNEL_ID|PIPELINE_ID|TITLE          |MESSAGE                                                                                                                                                                                                                                                        |GMT_CREATE           |GMT_MODIFIED         |
--|---|----------|-----------|---------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|---------------------|---------------------|
 1|-1 |1         |3          |POSITIONTIMEOUT|pid:3 position 671 seconds no update                                                                                                                                                                                                                           |2020-04-20 04:19:11.0|2020-04-20 04:19:11.0|
 3|   |-1        |-1         |EXCEPTION      |pid:-1 nid:null exception:cid:1 stop recovery successful for rid:3                                                                                                                                                                                             |2020-04-20 04:19:24.0|2020-04-20 04:19:24.0|
 5|-1 |1         |1          |POSITIONTIMEOUT|pid:1 position 659 seconds no update                                                                                                                                                                                                                           |2020-04-20 04:31:11.0|2020-04-20 04:31:11.0|
 7|-1 |1         |3          |POSITIONTIMEOUT|pid:3 position 705 seconds no update                                                                                                                                                                                                                           |2020-04-20 04:31:11.0|2020-04-20 04:31:11.0|
 9|   |-1        |-1         |EXCEPTION      |pid:-1 nid:null exception:cid:1 stop recovery successful for rid:3                                                                                                                                                                                             |2020-04-20 04:31:24.0|2020-04-20 04:31:24.0|
11|   |-1        |-1         |EXCEPTION      |pid:-1 nid:null exception:cid:1 stop recovery successful for rid:11                                                                                                                                                                                            |2020-04-20 04:31:26.0|2020-04-20 04:31:26.0|
13|-1 |1         |1          |POSITIONTIMEOUT|pid:1 position 660 seconds no update                                                                                                                                                                                                                           |2020-04-20 04:43:11.0|2020-04-20 04:43:11.0|
15|-1 |1         |3          |POSITIONTIMEOUT|pid:3 position 703 seconds no update                                                                                                                                                                                                                           |2020-04-20 04:43:11.0|2020-04-20 04:43:11.0|
17|   |-1        |-1         |EXCEPTION      |pid:-1 nid:null exception:cid:1 stop recovery successful for rid:3                                                                                                                                                                                             |2020-04-20 04:43:24.0|2020-04-20 04:43:24.0|
19|3  |1         |3          |EXCEPTION      |pid:3 nid:3 exception:canal:canalb:com.alibaba.otter.canal.parse.exception.CanalParseException: java.net.SocketException: Broken pipe (Write failed)¶Caused by: java.net.SocketException: Broken pipe (Write failed)¶ at java.net.SocketOutputStream.socketWrit|2020-04-20 04:43:24.0|2020-04-20 04:43:24.0|
21|   |-1        |-1         |EXCEPTION      |pid:-1 nid:null exception:cid:1 stop recovery successful for rid:11                                                                                                                                                                                            |2020-04-20 04:43:26.0|2020-04-20 04:43:26.0|
23|3  |1         |3          |EXCEPTION      |pid:3 nid:3 exception:canal:canalb:java.net.SocketTimeoutException: Timeout occurred, failed to read total 4 bytes in 25000 milliseconds, actual read only 0 bytes¶ at com.alibaba.otter.canal.parse.driver.mysql.socket.BioSocketChannel.read(BioSocketChannel|2020-04-20 04:43:34.0|2020-04-20 04:43:34.0|
25|-1 |1         |1          |POSITIONTIMEOUT|pid:1 position 660 seconds no update                                                                                                                                                                                                                           |2020-04-20 04:55:11.0|2020-04-20 04:55:11.0|
27|-1 |1         |3          |POSITIONTIMEOUT|pid:3 position 703 seconds no update                                                                                                                                                                                                                           |2020-04-20 04:55:11.0|2020-04-20 04:55:11.0|
29|   |-1        |-1         |EXCEPTION      |pid:-1 nid:null exception:cid:1 stop recovery successful for rid:3                                                                                                                                                                                             |2020-04-20 04:55:24.0|2020-04-20 04:55:24.0|
31|   |-1        |-1         |EXCEPTION      |pid:-1 nid:null exception:cid:1 stop recovery successful for rid:11                                                                                                                                                                                            |2020-04-20 04:55:26.0|2020-04-20 04:55:26.0|

</details>

<details><summary>TABLE_HISTORY_STAT</summary>

```sql
CREATE TABLE `TABLE_HISTORY_STAT` (
  `ID` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `FILE_SIZE` bigint(20) DEFAULT NULL,
  `FILE_COUNT` bigint(20) DEFAULT NULL,
  `INSERT_COUNT` bigint(20) DEFAULT NULL,
  `UPDATE_COUNT` bigint(20) DEFAULT NULL,
  `DELETE_COUNT` bigint(20) DEFAULT NULL,
  `DATA_MEDIA_PAIR_ID` bigint(20) DEFAULT NULL,
  `PIPELINE_ID` bigint(20) DEFAULT NULL,
  `START_TIME` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
  `END_TIME` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
  `GMT_CREATE` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
  `GMT_MODIFIED` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`ID`),
  KEY `idx_DATA_MEDIA_PAIR_ID_END_TIME` (`DATA_MEDIA_PAIR_ID`,`END_TIME`),
  KEY `idx_GmtCreate_id` (`GMT_CREATE`,`ID`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8
```

ID|FILE_SIZE|FILE_COUNT|INSERT_COUNT|UPDATE_COUNT|DELETE_COUNT|DATA_MEDIA_PAIR_ID|PIPELINE_ID|START_TIME           |END_TIME             |GMT_CREATE           |GMT_MODIFIED         |
--|---------|----------|------------|------------|------------|------------------|-----------|---------------------|---------------------|---------------------|---------------------|
 1|        0|         0|        1000|           0|           0|                 3|          3|2020-04-20 02:49:50.0|2020-04-20 02:49:55.0|2020-04-20 02:49:57.0|2020-04-20 02:49:57.0|
 3|        0|         0|       10000|           0|           0|                 1|          1|2020-04-20 02:50:01.0|2020-04-20 02:50:21.0|2020-04-20 02:50:57.0|2020-04-20 02:50:57.0|
 5|        0|         0|        9000|           0|           0|                 3|          3|2020-04-20 02:49:51.0|2020-04-20 02:50:11.0|2020-04-20 02:50:57.0|2020-04-20 02:50:57.0|
 
 </details>

<details><summary>TABLE_STAT</summary>

```sql
CREATE TABLE `TABLE_STAT` (
  `ID` bigint(20) NOT NULL AUTO_INCREMENT,
  `FILE_SIZE` bigint(20) NOT NULL,
  `FILE_COUNT` bigint(20) NOT NULL,
  `INSERT_COUNT` bigint(20) NOT NULL,
  `UPDATE_COUNT` bigint(20) NOT NULL,
  `DELETE_COUNT` bigint(20) NOT NULL,
  `DATA_MEDIA_PAIR_ID` bigint(20) NOT NULL,
  `PIPELINE_ID` bigint(20) NOT NULL,
  `GMT_CREATE` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
  `GMT_MODIFIED` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`ID`),
  KEY `idx_PipelineID_DataMediaPairID` (`PIPELINE_ID`,`DATA_MEDIA_PAIR_ID`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8
```

ID|FILE_SIZE|FILE_COUNT|INSERT_COUNT|UPDATE_COUNT|DELETE_COUNT|DATA_MEDIA_PAIR_ID|PIPELINE_ID|GMT_CREATE           |GMT_MODIFIED         |
--|---------|----------|------------|------------|------------|------------------|-----------|---------------------|---------------------|
 1|        0|         0|       10000|           0|           0|                 3|          3|2020-04-20 02:49:55.0|2020-04-20 02:50:11.0|
 3|        0|         0|       10000|           0|           0|                 1|          1|2020-04-20 02:50:07.0|2020-04-20 02:50:21.0|
 
 </details>
 
 <details><summary>THROUGHPUT_STAT</summary>
 
 ```sql
 CREATE TABLE `THROUGHPUT_STAT` (
   `ID` bigint(20) NOT NULL AUTO_INCREMENT,
   `TYPE` varchar(20) NOT NULL,
   `NUMBER` bigint(20) NOT NULL,
   `SIZE` bigint(20) NOT NULL,
   `PIPELINE_ID` bigint(20) NOT NULL,
   `START_TIME` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
   `END_TIME` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
   `GMT_CREATE` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
   `GMT_MODIFIED` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   PRIMARY KEY (`ID`),
   KEY `idx_PipelineID_Type_GmtCreate_ID` (`PIPELINE_ID`,`TYPE`,`GMT_CREATE`,`ID`),
   KEY `idx_PipelineID_Type_EndTime_ID` (`PIPELINE_ID`,`TYPE`,`END_TIME`,`ID`),
   KEY `idx_GmtCreate_id` (`GMT_CREATE`,`ID`)
 ) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8
```

ID|TYPE|NUMBER|SIZE   |PIPELINE_ID|START_TIME           |END_TIME             |GMT_CREATE           |GMT_MODIFIED         |
--|----|------|-------|-----------|---------------------|---------------------|---------------------|---------------------|
 1|ROW |  1000| 511328|          3|2020-04-20 02:49:50.0|2020-04-20 02:49:55.0|2020-04-20 02:49:57.0|2020-04-20 02:49:57.0|
 3|ROW | 10000|5111872|          1|2020-04-20 02:50:01.0|2020-04-20 02:50:21.0|2020-04-20 02:50:57.0|2020-04-20 02:50:57.0|
 5|ROW |  9000|4601240|          3|2020-04-20 02:49:51.0|2020-04-20 02:50:11.0|2020-04-20 02:50:57.0|2020-04-20 02:50:57.0|
 
 </details>
 