package mysqlclusterinit

// Settings 表示舒适化MySQL集群所需要的参数结构
type Settings struct {
	Master1Addr         string   // Master1的地址(IP，域名)
	Master2Addr         string   // Master2的地址(IP，域名)
	SlaveAddrs          []string // Slave的地址(IP，域名)
	RootPassword        string   // Root用户密码
	Port                int      // MySQL 端口号
	ReplUsr             string   // 复制用用户名
	ReplPassword        string   // 复制用户密码
	Debug               bool     // 测试模式，只打印SQL和HAProxy配置, 不实际执行
	LocalIP             string   // 指定本机的IP地址，不指定则自动从网卡中获取
	MySQLCnf            string   // MySQL 配置文件的地址， 例如：/etc/mysql/conf.d/my.cnf
	HAProxyCfg          string   // HAProxy配置文件地址，例如：/etc/haproxy/haproxy.cfg
	HAProxyRestartShell string   // HAProxy重启命令
}

// Result 表示初始化结果
type Result struct {
	Error   error
	Sqls    []string
	HAProxy string
}

type ShowSlaveStatusBean struct {
	SlaveIOState         string `gorm:"column:Slave_IO_State"`
	MasterHost           string `gorm:"column:Master_Host"`
	MasterUser           string `gorm:"column:Master_User"`
	MasterPort           int    `gorm:"column:Master_Port"`
	SlaveSqlRunningState string `gorm:"column:Slave_SQL_Running_State"`
	AutoPosition         bool   `gorm:"column:Auto_Position"`
	SlaveIoRunning       string `gorm:"column:Slave_IO_Running"`
	SlaveSqlRunning      string `gorm:"column:Slave_SQL_Running"`
	MasterServerId       string `gorm:"column:Master_Server_Id"`
}

type ShowVariablesBean struct {
	VariableName string `gorm:"column:Variable_name"`
	Value        string `gorm:"column:Value"`
}

type Variables struct {
	ServerId        string `var:"server_id"`
	LogBin          string `var:"log_bin"`
	SqlLogBin       string `var:"sql_log_bin"`
	GtidMode        string `var:"gtid_mode"`
	GtidNext        string `var:"gtid_next"`
	SlaveSkipErrors string `var:"slave_skip_errors"`
	BinlogFormat    string `var:"binlog_format"`
}
