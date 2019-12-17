1. interactive_timeout

    含义：服务器关闭交互式连接前等待活动的秒数。交互式客户端定义为在mysql_real_connect()中使用CLIENT_INTERACTIVE选项的客户端。
    默认值：28800秒（8小时）

1. wait_timeout

    含义：服务器关闭非交互连接之前等待活动的秒数。
    在线程启动时，根据全局wait_timeout值或全局interactive_timeout值初始化会话wait_timeout值，取决于客户端类型(由mysql_real_connect()的连接选项CLIENT_INTERACTIVE定义)。
    默认值：28800秒（8小时）
    
1. innodb_flush_log_at_timeout
1. innodb_flush_method
1. innodb_log_file_size

    表示每个文件的大小。因此总的redo log大小为innodb_log_files_in_group * innodb_log_file_size。

1. [半同步](https://dev.mysql.com/doc/refman/8.0/en/replication-options-master.html#sysvar_rpl_semi_sync_master_wait_for_slave_count)
    
    ```
    rpl_semi_sync_master_enabled        = 1                             #    0
    rpl_semi_sync_slave_enabled         = 1                             #    0
    rpl_semi_sync_master_timeout        = 1000                          #    1000(1 second) 同步复制中由于网络原因导致复制时间超过1s后，增强半同步复制就变成了异步复制了
    rpl_semi_sync_master_wait_for_slave_count = 1
    plugin_load_add                     = semisync_master.so            #
    plugin_load_add                     = semisync_slave.so             #
    ```
