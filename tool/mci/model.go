package mci

import (
	"github.com/bingoohuang/goreflect"
	"github.com/creasty/defaults"
	"github.com/sirupsen/logrus"
	"gopkg.in/dealancer/validate.v2"
)

// Settings 表示初始化化MySQL集群所需要的参数结构
type Settings struct {
	Master1Addr  string   `validate:"empty=false"` // Master1的地址(IP，域名)
	Master2Addr  string   `validate:"empty=false"` // Master2的地址(IP，域名)
	SlaveAddrs   []string // Slave的地址(IP，域名)
	User         string   `default:"root"`              // Root用户名
	Password     string   `validate:"empty=false"`      // Root用户密码
	Host         string   `default:"127.0.0.1"`         // MySQL 端口号
	Port         int      `default:"3306"`              // MySQL 端口号
	ReplUsr      string   `default:"repl"`              // 复制用用户名
	ReplPassword string   `default:"984d-CE5679F93918"` // 复制用户密码
	Debug        bool     // 测试模式，只打印SQL和HAProxy配置, 不实际执行
	LocalAddr    string   // 指定本机的IP地址，不指定则自动从网卡中获取
	MySQLCnf     string   `default:"/etc/my.cnf"`      // MySQL 配置文件的地址， 例如：/etc/mysql/conf.d/my.cnf, /etc/my.cnf
	HAProxyCfg   string   `default:"/etc/haproxy.cfg"` // HAProxy配置文件地址，
	// 例如：/etc/haproxy/haproxy.cfg, /etc/opt/rh/rh-haproxy18/haproxy/haproxy.cfg
	HAProxyRestartShell string `default:"systemctl restart haproxy"` // HAProxy重启命令

	CheckSQL string `default:"select current_date()"` // 检查MySQL可用性的SQL
}

type SettingsOption int

const (
	Validate SettingsOption = iota
	SetDefault
)

func (s *Settings) ValidateAndSetDefault(options ...SettingsOption) error {
	if goreflect.SliceContains(options, Validate) {
		if err := validate.Validate(s); err != nil {
			logrus.Errorf("error %v", err)
			return err
		}
	}

	if goreflect.SliceContains(options, SetDefault) {
		if err := defaults.Set(s); err != nil {
			logrus.Errorf("defaults set %v", err)
			return err
		}
	}

	if s.Debug {
		logrus.Infof("config: %+v\n", JSONPretty(s))
	}

	return nil
}

// Result 表示初始化结果
type Result struct {
	Sqls    []string
	HAProxy string
}

// ShowSlaveStatusBean 表示MySQL Slave Status
type ShowSlaveStatusBean struct {
	SlaveIOState         string `gorm:"column:Slave_IO_State"`
	MasterHost           string `gorm:"column:Master_Host"`
	MasterUser           string `gorm:"column:Master_User"`
	MasterPort           int    `gorm:"column:Master_Port"`
	SlaveSQLRunningState string `gorm:"column:Slave_SQL_Running_State"`
	AutoPosition         bool   `gorm:"column:Auto_Position"`
	SlaveIoRunning       string `gorm:"column:Slave_IO_Running"`
	SlaveSQLRunning      string `gorm:"column:Slave_SQL_Running"`
	MasterServerID       string `gorm:"column:Master_Server_Id"`
	SecondsBehindMaster  string `gorm:"column:Seconds_Behind_Master"`
}

// ShowVariablesBean 表示MySQL 服务器参数结果行
type ShowVariablesBean struct {
	VariableName string `gorm:"column:Variable_name"`
	Value        string `gorm:"column:Value"`
}

// Variables 表示MySQL 服务器参数
type Variables struct {
	ServerID               string `var:"server_id"`
	LogBin                 string `var:"log_bin"`
	SQLLogBin              string `var:"sql_log_bin"`
	GtidMode               string `var:"gtid_mode"`
	GtidNext               string `var:"gtid_next"`
	SlaveSkipErrors        string `var:"slave_skip_errors"`
	BinlogFormat           string `var:"binlog_format"`
	MasterInfoRepository   string `var:"master_info_repository"`
	RelayLogInfoRepository string `var:"relay_log_info_repository"`
	InnodbVersion          string `var:"innodb_version"`
}
