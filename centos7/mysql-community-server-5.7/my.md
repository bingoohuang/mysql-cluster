1. interactive_timeout

    含义：服务器关闭交互式连接前等待活动的秒数。交互式客户端定义为在mysql_real_connect()中使用CLIENT_INTERACTIVE选项的客户端。
    默认值：28800秒（8小时）

1. wait_timeout

    含义：服务器关闭非交互连接之前等待活动的秒数。
    在线程启动时，根据全局wait_timeout值或全局interactive_timeout值初始化会话wait_timeout值，取决于客户端类型(由mysql_real_connect()的连接选项CLIENT_INTERACTIVE定义)。
    默认值：28800秒（8小时）

1. innodb_flush_log_at_timeout
1. wait_timeout
1. innodb_flush_method
1. innodb_log_file_size
  
    innodb_log_file_size表示每个文件的大小。因此总的redo log大小为innodb_log_files_in_group * innodb_log_file_size。

