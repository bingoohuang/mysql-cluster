CREATE DATABASE /*!32312 IF NOT EXISTS*/ `otter` /*!40100 DEFAULT CHARACTER SET utf8 COLLATE utf8_bin */;

USE `otter`;

SET sql_mode='ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION';

CREATE TABLE `ALARM_RULE` (
  `ID` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `MONITOR_NAME` varchar(1024) DEFAULT NULL,
  `RECEIVER_KEY` varchar(1024) DEFAULT NULL,
  `STATUS` varchar(32) DEFAULT NULL,
  `PIPELINE_ID` bigint(20) NOT NULL,
  `DESCRIPTION` varchar(256) DEFAULT NULL,
  `GMT_CREATE` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
  `GMT_MODIFIED` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `MATCH_VALUE` varchar(1024) DEFAULT NULL,
  `PARAMETERS` text DEFAULT NULL,
  PRIMARY KEY (`ID`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

CREATE TABLE `AUTOKEEPER_CLUSTER` (
  `ID` bigint(20) NOT NULL AUTO_INCREMENT,
  `CLUSTER_NAME` varchar(200) NOT NULL,
  `SERVER_LIST` varchar(1024) NOT NULL,
  `DESCRIPTION` varchar(200) DEFAULT NULL,
  `GMT_CREATE` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
  `GMT_MODIFIED` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`ID`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

CREATE TABLE `CANAL` (
  `ID` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `NAME` varchar(200) DEFAULT NULL,
  `DESCRIPTION` varchar(200) DEFAULT NULL,
  `PARAMETERS` text DEFAULT NULL,
  `GMT_CREATE` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
  `GMT_MODIFIED` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`ID`),
  UNIQUE KEY `CANALUNIQUE` (`NAME`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

CREATE TABLE `CHANNEL` (
  `ID` bigint(20) NOT NULL AUTO_INCREMENT,
  `NAME` varchar(200) NOT NULL,
  `DESCRIPTION` varchar(200) DEFAULT NULL,
  `PARAMETERS` text DEFAULT NULL,
  `GMT_CREATE` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
  `GMT_MODIFIED` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`ID`),
  UNIQUE KEY `CHANNELUNIQUE` (`NAME`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

CREATE TABLE `COLUMN_PAIR` (
  `ID` bigint(20) NOT NULL AUTO_INCREMENT,
  `SOURCE_COLUMN` varchar(200) DEFAULT NULL,
  `TARGET_COLUMN` varchar(200) DEFAULT NULL,
  `DATA_MEDIA_PAIR_ID` bigint(20) NOT NULL,
  `GMT_CREATE` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
  `GMT_MODIFIED` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`ID`),
  KEY `idx_DATA_MEDIA_PAIR_ID` (`DATA_MEDIA_PAIR_ID`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

CREATE TABLE `COLUMN_PAIR_GROUP` (
  `ID` bigint(20) NOT NULL AUTO_INCREMENT,
  `DATA_MEDIA_PAIR_ID` bigint(20) NOT NULL,
  `COLUMN_PAIR_CONTENT` text DEFAULT NULL,
  `GMT_CREATE` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
  `GMT_MODIFIED` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`ID`),
  KEY `idx_DATA_MEDIA_PAIR_ID` (`DATA_MEDIA_PAIR_ID`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

CREATE TABLE `DATA_MEDIA` (
  `ID` bigint(20) NOT NULL AUTO_INCREMENT,
  `NAME` varchar(200) NOT NULL,
  `NAMESPACE` varchar(200) NOT NULL,
  `PROPERTIES` varchar(1000) NOT NULL,
  `DATA_MEDIA_SOURCE_ID` bigint(20) NOT NULL,
  `GMT_CREATE` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
  `GMT_MODIFIED` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`ID`),
  UNIQUE KEY `DATAMEDIAUNIQUE` (`NAME`,`NAMESPACE`,`DATA_MEDIA_SOURCE_ID`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

CREATE TABLE `DATA_MEDIA_PAIR` (
  `ID` bigint(20) NOT NULL AUTO_INCREMENT,
  `PULLWEIGHT` bigint(20) DEFAULT NULL,
  `PUSHWEIGHT` bigint(20) DEFAULT NULL,
  `RESOLVER` text DEFAULT NULL,
  `FILTER` text DEFAULT NULL,
  `SOURCE_DATA_MEDIA_ID` bigint(20) DEFAULT NULL,
  `TARGET_DATA_MEDIA_ID` bigint(20) DEFAULT NULL,
  `PIPELINE_ID` bigint(20) NOT NULL,
  `COLUMN_PAIR_MODE` varchar(20) DEFAULT NULL,
  `GMT_CREATE` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
  `GMT_MODIFIED` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`ID`),
  KEY `idx_PipelineID` (`PIPELINE_ID`,`ID`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

CREATE TABLE `DATA_MEDIA_SOURCE` (
  `ID` bigint(20) NOT NULL AUTO_INCREMENT,
  `NAME` varchar(200) NOT NULL,
  `TYPE` varchar(20) NOT NULL,
  `PROPERTIES` varchar(1000) NOT NULL,
  `GMT_CREATE` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
  `GMT_MODIFIED` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`ID`),
  UNIQUE KEY `DATAMEDIASOURCEUNIQUE` (`NAME`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

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
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

CREATE TABLE `LOG_RECORD` (
  `ID` bigint(20) NOT NULL AUTO_INCREMENT,
  `NID` varchar(200) DEFAULT NULL,
  `CHANNEL_ID` varchar(200) NOT NULL,
  `PIPELINE_ID` varchar(200) NOT NULL,
  `TITLE` varchar(1000) DEFAULT NULL,
  `MESSAGE` text DEFAULT NULL,
  `GMT_CREATE` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
  `GMT_MODIFIED` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`ID`),
  KEY `logRecord_pipelineId` (`PIPELINE_ID`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

CREATE TABLE `NODE` (
  `ID` bigint(20) NOT NULL AUTO_INCREMENT,
  `NAME` varchar(200) NOT NULL,
  `IP` varchar(200) NOT NULL,
  `PORT` bigint(20) NOT NULL,
  `DESCRIPTION` varchar(200) DEFAULT NULL,
  `PARAMETERS` text DEFAULT NULL,
  `GMT_CREATE` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
  `GMT_MODIFIED` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`ID`),
  UNIQUE KEY `NODEUNIQUE` (`NAME`,`IP`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

CREATE TABLE `PIPELINE` (
  `ID` bigint(20) NOT NULL AUTO_INCREMENT,
  `NAME` varchar(200) NOT NULL,
  `DESCRIPTION` varchar(200) DEFAULT NULL,
  `PARAMETERS` text DEFAULT NULL,
  `CHANNEL_ID` bigint(20) NOT NULL,
  `GMT_CREATE` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
  `GMT_MODIFIED` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`ID`),
  UNIQUE KEY `PIPELINEUNIQUE` (`NAME`,`CHANNEL_ID`),
  KEY `idx_ChannelID` (`CHANNEL_ID`,`ID`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

CREATE TABLE `PIPELINE_NODE_RELATION` (
  `ID` bigint(20) NOT NULL AUTO_INCREMENT,
  `NODE_ID` bigint(20) NOT NULL,
  `PIPELINE_ID` bigint(20) NOT NULL,
  `LOCATION` varchar(20) NOT NULL,
  `GMT_CREATE` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
  `GMT_MODIFIED` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`ID`),
  KEY `idx_PipelineID` (`PIPELINE_ID`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

CREATE TABLE `SYSTEM_PARAMETER` (
  `ID` bigint(20) unsigned NOT NULL,
  `VALUE` text DEFAULT NULL,
  `GMT_CREATE` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
  `GMT_MODIFIED` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

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
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

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
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

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
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

CREATE TABLE `USER` (
  `ID` bigint(20) NOT NULL AUTO_INCREMENT,
  `USERNAME` varchar(20) NOT NULL,
  `PASSWORD` varchar(20) NOT NULL,
  `AUTHORIZETYPE` varchar(20) NOT NULL,
  `DEPARTMENT` varchar(20) NOT NULL,
  `REALNAME` varchar(20) NOT NULL,
  `GMT_CREATE` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
  `GMT_MODIFIED` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`ID`),
  UNIQUE KEY `USERUNIQUE` (`USERNAME`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

CREATE TABLE  `DATA_MATRIX` (
  `ID` bigint(20) NOT NULL AUTO_INCREMENT,
  `GROUP_KEY` varchar(200) DEFAULT NULL,
  `MASTER` varchar(200) DEFAULT NULL,
  `SLAVE` varchar(200) DEFAULT NULL,
  `DESCRIPTION` varchar(200) DEFAULT NULL,
  `GMT_CREATE` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
  `GMT_MODIFIED` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`ID`),
  KEY `GROUPKEY` (`GROUP_KEY`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `meta_history` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `gmt_create` datetime NOT NULL COMMENT '创建时间',
  `gmt_modified` datetime NOT NULL COMMENT '修改时间',
  `destination` varchar(128) DEFAULT NULL COMMENT '通道名称',
  `binlog_file` varchar(64) DEFAULT NULL COMMENT 'binlog文件名',
  `binlog_offest` bigint(20) DEFAULT NULL COMMENT 'binlog偏移量',
  `binlog_master_id` varchar(64) DEFAULT NULL COMMENT 'binlog节点id',
  `binlog_timestamp` bigint(20) DEFAULT NULL COMMENT 'binlog应用的时间戳',
  `use_schema` varchar(1024) DEFAULT NULL COMMENT '执行sql时对应的schema',
  `sql_schema` varchar(1024) DEFAULT NULL COMMENT '对应的schema',
  `sql_table` varchar(1024) DEFAULT NULL COMMENT '对应的table',
  `sql_text` longtext DEFAULT NULL COMMENT '执行的sql',
  `sql_type` varchar(256) DEFAULT NULL COMMENT 'sql类型',
  `extra` text DEFAULT NULL COMMENT '额外的扩展信息',
  PRIMARY KEY (`id`),
  UNIQUE KEY binlog_file_offest(`destination`,`binlog_master_id`,`binlog_file`,`binlog_offest`),
  KEY `destination` (`destination`),
  KEY `destination_timestamp` (`destination`,`binlog_timestamp`),
  KEY `gmt_modified` (`gmt_modified`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COMMENT='表结构变化明细表';

CREATE TABLE IF NOT EXISTS `meta_snapshot` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `gmt_create` datetime NOT NULL COMMENT '创建时间',
  `gmt_modified` datetime NOT NULL COMMENT '修改时间',
  `destination` varchar(128) DEFAULT NULL COMMENT '通道名称',
  `binlog_file` varchar(64) DEFAULT NULL COMMENT 'binlog文件名',
  `binlog_offest` bigint(20) DEFAULT NULL COMMENT 'binlog偏移量',
  `binlog_master_id` varchar(64) DEFAULT NULL COMMENT 'binlog节点id',
  `binlog_timestamp` bigint(20) DEFAULT NULL COMMENT 'binlog应用的时间戳',
  `data` longtext DEFAULT NULL COMMENT '表结构数据',
  `extra` text DEFAULT NULL COMMENT '额外的扩展信息',
  PRIMARY KEY (`id`),
  UNIQUE KEY binlog_file_offest(`destination`,`binlog_master_id`,`binlog_file`,`binlog_offest`),
  KEY `destination` (`destination`),
  KEY `destination_timestamp` (`destination`,`binlog_timestamp`),
  KEY `gmt_modified` (`gmt_modified`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COMMENT='表结构记录表快照表';


insert into USER(ID,USERNAME,PASSWORD,AUTHORIZETYPE,DEPARTMENT,REALNAME,GMT_CREATE,GMT_MODIFIED) values(null,'admin','801fc357a5a74743894a','ADMIN','admin','admin',now(),now());
insert into USER(ID,USERNAME,PASSWORD,AUTHORIZETYPE,DEPARTMENT,REALNAME,GMT_CREATE,GMT_MODIFIED) values(null,'guest','471e02a154a2121dc577','OPERATOR','guest','guest',now(),now());

;INSERT INTO NODE (ID,NAME,IP,PORT,DESCRIPTION,PARAMETERS,GMT_CREATE,GMT_MODIFIED) VALUES
(1,'node1','nd1',2088,NULL,'{"downloadPort":9090,"mbeanPort":2090,"useExternalIp":false,"zkCluster":{"clusterName":"zka","description":"zka","gmtCreate":1586773022000,"gmtModified":1586773022000,"id":1,"serverList":["zk1:2181","zk2:2181","zk3:2181"]}}','2020-04-13 10:18:20.0','2020-04-13 10:18:20.0')
,(2,'node2','nd2',2088,NULL,'{"downloadPort":9090,"mbeanPort":2090,"useExternalIp":false,"zkCluster":{"clusterName":"zkb","description":"zkb","gmtCreate":1586773046000,"gmtModified":1586773046000,"id":2,"serverList":["zk4:2181","zk5:2181"]}}','2020-04-13 10:18:31.0','2020-04-13 10:18:31.0')

;INSERT INTO AUTOKEEPER_CLUSTER (ID,CLUSTER_NAME,SERVER_LIST,DESCRIPTION,GMT_CREATE,GMT_MODIFIED) VALUES
(1,'zka','["zk1:2181","zk2:2181","zk3:2181"]','zka','2020-04-13 10:17:02.0','2020-04-13 10:17:02.0')
,(2,'zkb','["zk4:2181","zk5:2181"]','zkb','2020-04-13 10:17:26.0','2020-04-13 10:17:26.0')
;
/*

;INSERT INTO CANAL (ID,NAME,DESCRIPTION,PARAMETERS,GMT_CREATE,GMT_MODIFIED) VALUES
(1,'canala',NULL,'{"connectionCharset":"UTF-8","connectionCharsetNumber":33,"dataDir":"../conf","dbAddresses":[],"dbPassword":"root","dbUsername":"root","ddlIsolation":false,"defaultConnectionTimeoutInSeconds":30,"detectingEnable":false,"detectingIntervalInSeconds":5,"detectingRetryTimes":3,"detectingSQL":"insert into retl.xdual values(1,now()) on duplicate key update x=now()","detectingTimeoutThresholdInSeconds":30,"fallbackIntervalInSeconds":60,"filterTableError":false,"groupDbAddresses":[[{"dbAddress":{"address":"ma","port":3306},"type":"MYSQL"}]],"gtidEnable":false,"haMode":"HEARTBEAT","heartbeatHaEnable":false,"indexMode":"MEMORY_META_FAILBACK","memoryStorageBufferMemUnit":1024,"memoryStorageBufferSize":32768,"memoryStorageRawEntry":true,"metaFileFlushPeriod":1000,"metaMode":"MIXED","parallel":false,"port":11111,"positions":[],"receiveBufferSize":16384,"runMode":"EMBEDDED","sendBufferSize":16384,"slaveId":10001,"sourcingType":"MYSQL","storageBatchMode":"MEMSIZE","storageMode":"MEMORY","storageScavengeMode":"ON_ACK","transactionSize":1024,"tsdbEnable":false,"tsdbSnapshotExpire":360,"tsdbSnapshotInterval":24,"zkClusterId":1,"zkClusters":[]}','2020-04-13 10:31:08.0','2020-04-13 10:31:08.0')
,(2,'canalb',NULL,'{"connectionCharset":"UTF-8","connectionCharsetNumber":33,"dataDir":"../conf","dbAddresses":[],"dbPassword":"root","dbUsername":"root","ddlIsolation":false,"defaultConnectionTimeoutInSeconds":30,"detectingEnable":false,"detectingIntervalInSeconds":5,"detectingRetryTimes":3,"detectingSQL":"insert into retl.xdual values(1,now()) on duplicate key update x=now()","detectingTimeoutThresholdInSeconds":30,"fallbackIntervalInSeconds":60,"filterTableError":false,"groupDbAddresses":[[{"dbAddress":{"address":"mb","port":3306},"type":"MYSQL"}]],"gtidEnable":false,"haMode":"HEARTBEAT","heartbeatHaEnable":false,"indexMode":"MEMORY_META_FAILBACK","memoryStorageBufferMemUnit":1024,"memoryStorageBufferSize":32768,"memoryStorageRawEntry":true,"metaFileFlushPeriod":1000,"metaMode":"MIXED","parallel":false,"port":11111,"positions":[],"receiveBufferSize":16384,"runMode":"EMBEDDED","sendBufferSize":16384,"slaveId":10002,"sourcingType":"MYSQL","storageBatchMode":"MEMSIZE","storageMode":"MEMORY","storageScavengeMode":"ON_ACK","transactionSize":1024,"tsdbEnable":false,"tsdbSnapshotExpire":360,"tsdbSnapshotInterval":24,"zkClusterId":1,"zkClusters":[]}','2020-04-13 10:31:23.0','2020-04-13 10:31:23.0')
;INSERT INTO CHANNEL (ID,NAME,DESCRIPTION,PARAMETERS,GMT_CREATE,GMT_MODIFIED) VALUES
(1,'aabb',NULL,'{"channelId":1,"enableRemedy":false,"remedyAlgorithm":"LOOPBACK","remedyDelayThresoldForMedia":60,"syncConsistency":"BASE","syncMode":"FIELD"}','2020-04-13 10:31:41.0','2020-04-13 10:31:46.0')
;INSERT INTO PIPELINE (ID,NAME,DESCRIPTION,PARAMETERS,CHANNEL_ID,GMT_CREATE,GMT_MODIFIED) VALUES
(1,'pipea',NULL,'{"arbitrateMode":"AUTOMATIC","batchTimeout":-1,"ddlSync":true,"destinationName":"canala","dryRun":false,"dumpEvent":false,"dumpSelector":true,"dumpSelectorDetail":false,"enableCompatibleMissColumn":true,"extractPoolSize":10,"fileDetect":false,"fileLoadPoolSize":15,"home":false,"lbAlgorithm":"Stick","loadPoolSize":15,"mainstemBatchsize":6000,"parallelism":5,"pipeChooseType":"AUTOMATIC","selectorMode":"Canal","skipDdlException":false,"skipFreedom":false,"skipLoadException":false,"skipNoRow":false,"skipSelectException":false,"useBatch":true,"useExternalIp":false,"useFileEncrypt":false,"useLocalFileMutliThread":false,"useTableTransform":false}',1,'2020-04-13 10:33:17.0','2020-04-13 10:33:17.0')
,(2,'pipeb',NULL,'{"arbitrateMode":"AUTOMATIC","batchTimeout":-1,"ddlSync":false,"destinationName":"canalb","dryRun":false,"dumpEvent":true,"dumpSelector":true,"dumpSelectorDetail":true,"enableCompatibleMissColumn":true,"extractPoolSize":10,"fileDetect":false,"fileLoadPoolSize":15,"home":false,"lbAlgorithm":"Stick","loadPoolSize":15,"mainstemBatchsize":6000,"parallelism":5,"pipeChooseType":"AUTOMATIC","selectorMode":"Canal","skipDdlException":true,"skipFreedom":false,"skipLoadException":true,"skipNoRow":false,"skipSelectException":true,"useBatch":true,"useExternalIp":false,"useFileEncrypt":false,"useLocalFileMutliThread":false,"useTableTransform":false}',1,'2020-04-13 10:33:57.0','2020-04-13 10:33:57.0')
;INSERT INTO PIPELINE_NODE_RELATION (ID,NODE_ID,PIPELINE_ID,LOCATION,GMT_CREATE,GMT_MODIFIED) VALUES
(1,1,1,'SELECT','2020-04-13 10:33:17.0','2020-04-13 10:33:17.0')
,(2,1,1,'EXTRACT','2020-04-13 10:33:17.0','2020-04-13 10:33:17.0')
,(3,2,1,'LOAD','2020-04-13 10:33:17.0','2020-04-13 10:33:17.0')
,(4,2,2,'SELECT','2020-04-13 10:33:57.0','2020-04-13 10:33:57.0')
,(5,2,2,'EXTRACT','2020-04-13 10:33:57.0','2020-04-13 10:33:57.0')
,(6,1,2,'LOAD','2020-04-13 10:33:57.0','2020-04-13 10:33:57.0')
;
*/

/*
;INSERT INTO DATA_MEDIA (ID,NAME,NAMESPACE,PROPERTIES,DATA_MEDIA_SOURCE_ID,GMT_CREATE,GMT_MODIFIED) VALUES
(1,'.*','aa','{"mode":"SINGLE","name":".*","namespace":"aa","source":{"driver":"com.mysql.jdbc.Driver","encode":"UTF8","gmtCreate":1586773705000,"gmtModified":1586773705000,"id":1,"name":"dsa","password":"root","type":"MYSQL","url":"jdbc:mysql://ma:3306","username":"root"}}',1,'2020-04-13 10:28:59.0','2020-04-13 10:28:59.0')
,(2,'.*','bb','{"mode":"SINGLE","name":".*","namespace":"bb","source":{"driver":"com.mysql.jdbc.Driver","encode":"UTF8","gmtCreate":1586773717000,"gmtModified":1586773717000,"id":2,"name":"dsb","password":"root","type":"MYSQL","url":"jdbc:mysql://mb:3306","username":"root"}}',2,'2020-04-13 10:29:11.0','2020-04-13 10:29:11.0')
;INSERT INTO DATA_MEDIA_PAIR (ID,PULLWEIGHT,PUSHWEIGHT,RESOLVER,`FILTER`,SOURCE_DATA_MEDIA_ID,TARGET_DATA_MEDIA_ID,PIPELINE_ID,COLUMN_PAIR_MODE,GMT_CREATE,GMT_MODIFIED) VALUES
(1,NULL,5,'{"blank":true,"clazzPath":"","extensionDataType":"CLAZZ","notBlank":false,"timestamp":1586774167238}','{"blank":true,"clazzPath":"","extensionDataType":"CLAZZ","notBlank":false,"timestamp":1586774167238}',1,2,1,'INCLUDE','2020-04-13 10:36:07.0','2020-04-13 10:36:07.0')
,(2,NULL,5,'{"blank":true,"clazzPath":"","extensionDataType":"CLAZZ","notBlank":false,"timestamp":1586774188220}','{"blank":true,"clazzPath":"","extensionDataType":"CLAZZ","notBlank":false,"timestamp":1586774188220}',2,1,2,'INCLUDE','2020-04-13 10:36:28.0','2020-04-13 10:36:28.0')
;INSERT INTO DATA_MEDIA_SOURCE (ID,NAME,`TYPE`,PROPERTIES,GMT_CREATE,GMT_MODIFIED) VALUES
(1,'dsa','MYSQL','{"driver":"com.mysql.jdbc.Driver","encode":"UTF8","name":"dsa","password":"root","type":"MYSQL","url":"jdbc:mysql://ma:3306","username":"root"}','2020-04-13 10:28:25.0','2020-04-13 10:28:25.0')
,(2,'dsb','MYSQL','{"driver":"com.mysql.jdbc.Driver","encode":"UTF8","name":"dsb","password":"root","type":"MYSQL","url":"jdbc:mysql://mb:3306","username":"root"}','2020-04-13 10:28:37.0','2020-04-13 10:28:37.0')

*/