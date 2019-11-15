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
	User         string   `default:"root"`         // Root用户名
	Password     string   `validate:"empty=false"` // Root用户密码
	Host         string   `default:"127.0.0.1"`    // MySQL 端口号
	Port         int      `default:"3306"`         // MySQL 端口号
	ReplUsr      string   `default:"repl"`         // 复制用用户名，默认repl
	ReplPassword string   // 复制用户密码，如果不指定，则使用uuid生成
	Debug        bool     // 测试模式，只打印SQL和HAProxy配置, 不实际执行
	LocalAddr    string   // 指定本机的IP地址，不指定则自动从网卡中获取
	MySQLCnf     string   `default:"/etc/my.cnf"`      // MySQL 配置文件的地址， 例如：/etc/mysql/conf.d/my.cnf, /etc/my.cnf
	HAProxyCfg   string   `default:"/etc/haproxy.cfg"` // HAProxy配置文件地址，
	// 例如：/etc/haproxy/haproxy.cfg, /etc/opt/rh/rh-haproxy18/haproxy/haproxy.cfg
	HAProxyRestartShell string `default:"systemctl restart haproxy"` // HAProxy重启命令

	CheckSQL string `default:"select current_date()"` // 检查MySQL可用性的SQL

	MySQLDumpCmd string `default:"mysqldump"` // mysqldump命令路径
	MySQLCmd     string `default:"mysql"`     // mysql命令路径

	LocalAddrMap map[string]bool //  本地地址
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

	if s.Debug {
		logrus.Infof("config: %+v\n", JSONPretty(s))
	}

	return nil
}

// Setup setups settings.
func (s *Settings) Setup() {
	s.LocalAddrMap = make(map[string]bool)
	if s.LocalAddr != "" {
		s.LocalAddrMap[s.LocalAddr] = true
	}

	s.LocalAddrMap[s.Master1Addr] = s.isLocalAddr(s.Master1Addr)
	s.LocalAddrMap[s.Master2Addr] = s.isLocalAddr(s.Master2Addr)

	for _, slaveAddr := range s.SlaveAddrs {
		s.LocalAddrMap[slaveAddr] = s.isLocalAddr(slaveAddr)
	}

	if s.ReplPassword == "" {
		s.ReplPassword = GeneratePasswordBySet(16, UpperLetters, DigitsLetters, LowerLetters, "-#")
	}
}

// Result 表示初始化结果
type Result struct {
	Nodes   []MySQLNode
	HAProxy string
}
