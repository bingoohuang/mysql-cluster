package mci

import (
	"time"

	"github.com/bingoohuang/goreflect"
	"github.com/creasty/defaults"
	"github.com/sirupsen/logrus"
	"gopkg.in/dealancer/validate.v2"
)

// Settings 表示初始化化MySQL集群所需要的参数结构.
type Settings struct {
	Master1Addr  string   `validate:"empty=false" usage:"Master1的地址(IP，域名)"`
	Master2Addr  string   `validate:"empty=false" usage:"Master2的地址(IP，域名)"`
	SlaveAddrs   []string `usage:"Slave的地址(IP，域名)"`
	User         string   `default:"root" usage:"root用户名"`
	Password     string   `validate:"empty=false" usage:"root用户密码"`
	Port         int      `default:"3306" usage:"MySQL端口"`
	ReplUsr      string   `default:"repl" usage:"复制用用户名，默认rep"`
	ReplPassword string   `usage:"复制用户密码，如果不指定，则使用uuid生成"`

	Debug       bool `usage:"测试模式，只打印SQL和HAProxy配置, 不实际执行"`
	NoLog       bool
	IPv6Enabled bool `usage:"是否支持IPv6"`
	NoBackup    bool `usage:"是否备份其他（主2及从1...n)的数据"`

	MySQLCnf   string `default:"/etc/my.cnf" usage:"MySQL 配置文件的地址， 例如：/etc/mysql/conf.d/my.cnf, /etc/my.cnf"`
	HAProxyCfg string `default:"/etc/haproxy.cfg" usage:"HAProxy配置文件地址，例如：/etc/haproxy/haproxy.cfg, /etc/opt/rh/rh-haproxy18/haproxy/haproxy.cfg"`

	MySQLRestartShell   string `default:"systemctl restart mysqld" usage:"MySQL重启命令"`
	HAProxyRestartShell string `default:"systemctl restart haproxy" usage:"HAProxy重启命令"`

	CheckSQL string `default:"select current_date()" usage:"检查MySQL可用性的SQL"`

	// nolint
	MySQLDumpOptions string `default:"--ignore-table=mysql.help_topic --ignore-table=mysql.help_keyword --ignore-table=mysql.help_relation --ignore-table=mysql.help_category" usage:"msyqldump命令选项"`
	MySQLDumpCmd     string `default:"mysqldump" usage:"mysqldump命令路径"`
	MySQLCmd         string `default:"mysql" usage:"mysql命令路径"`
	ShellTimeout     string `usage:"执行shell的超时时间"`

	// https://github.com/go-sql-driver/mysql#parameters
	MySQLDSNParams string `default:"timeout=120s&writeTimeout=120s&readTimeout=120s" usage:"MySQL DSN 参数"`
	LogLevel       string `usage:"logrus日志级别"`

	shellTimeoutDuration time.Duration
	currentHost          string
}

// SettingsOption  stands for option of settings.
type SettingsOption int

const (
	// Validate means validation required.
	Validate SettingsOption = iota + 1
	// SetDefault means SetDefault required.
	SetDefault
)

// ValidateAndSetDefault validates and set defaults to s.
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

	if s.LogLevel != "" {
		if level, err := logrus.ParseLevel(s.LogLevel); err == nil {
			logrus.SetLevel(level)
		}
	}

	if s.Debug {
		logrus.Infof("config: %+v\n", JSONPretty(s))
	}

	return nil
}

// Setup setups settings.
func (s *Settings) Setup() {
	if s.ReplPassword == "" {
		s.ReplPassword = GeneratePasswordBySet(16, UpperLetters, DigitsLetters, LowerLetters, "-#") // nolint:gomnd
	}

	if s.ShellTimeout != "" {
		shellTimeoutDuration, err := time.ParseDuration(s.ShellTimeout)
		if err != nil {
			logrus.Fatalf("error to parse ShellTimeout %s error %v", s.ShellTimeout, err)
		}

		s.shellTimeoutDuration = shellTimeoutDuration
	}
}

// Result 表示初始化结果.
type Result struct {
	Nodes   []MySQLNode
	HAProxy string
}
