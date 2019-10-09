package mysqlclusterinit

import (
	"fmt"
	"strings"

	"github.com/bingoohuang/gonet"
	"github.com/bingoohuang/sqlmore"
	"github.com/sirupsen/logrus"

	// support mysql
	_ "github.com/go-sql-driver/mysql"
)

// Settings 表示舒适化MySQL集群所需要的参数结构
type Settings struct {
	Master1IP    string   // Master1的IP
	Master2IP    string   // Master2的IP
	SlaveIps     []string // Slave的IP地址
	RootPassword string   // Root用户密码
	Port         int      // MySQL 端口号
	ReplUsr      string   // 复制用用户名
	ReplPassword string   // 复制用户密码
	Debug        bool     // 测试模式，只打印SQL和HAProxy配置, 不实际执行
	LocalIP      string   // 指定本机的IP地址，不指定则自动从网卡中获取
}

// Result 表示初始化结果
type Result struct {
	Error   error
	Sqls    []string
	HAProxy string
}

// InitMySQLCluster 初始化MySQL Master-Master集群
func (s Settings) InitMySQLCluster() (r Result) {
	if r.Sqls, r.Error = s.createMySQCluster(); r.Error != nil {
		return r
	}

	r.HAProxy = s.createHAProxyConfig()

	if s.Debug {
		logrus.Infof("SQL:%s", strings.Join(r.Sqls, ";\n"))
		logrus.Infof("HAProxy:%s", r.HAProxy)
	}

	return r
}

func (s Settings) createMySQCluster() (sqls []string, err error) {
	localIP := s.initLocalIPMap()

	sqls = s.createInitSqls(localIP)
	if len(sqls) == 0 {
		logrus.Infof("InitMySQLCluster bypassed, nor master or slave on %v", localIP)
	} else {
		err = s.execMultiSqls(sqls)
	}

	return
}

func (s Settings) initLocalIPMap() map[string]bool {
	if s.LocalIP == "" {
		return gonet.ListLocalIPMap()
	}

	return map[string]bool{s.LocalIP: true}
}

func (s Settings) execMultiSqls(sqls []string) error {
	if s.Debug {
		return nil
	}

	ds := fmt.Sprintf("root:%s@tcp(127.0.0.1:%d)/root", s.RootPassword, s.Port)
	db := sqlmore.NewSQLMore("mysql", ds).MustOpen()
	defer db.Close()

	for _, sqlStr := range sqls {
		if r := sqlmore.ExecSQL(db, sqlStr, 0, ""); r.Error != nil {
			return fmt.Errorf("exec sql %s error %w", sqlStr, r.Error)
		}

		logrus.Infof("execSQL %s completed", sqlStr)
	}

	logrus.Infof("createMySQCluster completed")
	return nil
}

func (s Settings) createInitSqls(localIP map[string]bool) []string {
	if _, ok := localIP[s.Master1IP]; ok {
		return s.initMasterSqls(1, s.Master2IP)
	}
	if _, ok := localIP[s.Master2IP]; ok {
		return s.initMasterSqls(2, s.Master1IP)
	}

	for seq, slaveIP := range s.SlaveIps {
		if _, ok := localIP[slaveIP]; ok {
			return s.initSlaveSqls(seq+3, s.Master2IP)
		}
	}

	return []string{}
}

func (s Settings) createHAProxyConfig() string {
	rwConfig := fmt.Sprintf(`
listen mysql-rw
  bind 0.0.0.0:13306
  mode tcp
  option tcpka
  server mysql-1 %s:%d check inter 1s
  server mysql-2 %s:%d check inter 1s backup
`, s.Master1IP, s.Port, s.Master2IP, s.Port)

	rConfig := fmt.Sprintf(`
listen mysql-ro
  bind 0.0.0.0:23306
  mode tcp
  option tcpka
  server mysql-1 %s:%d check inter 1s
  server mysql-2 %s:%d check inter 1s
`, s.Master1IP, s.Port, s.Master2IP, s.Port)

	for seq, slaveIP := range s.SlaveIps {
		rConfig += fmt.Sprintf("  server mysql-%d %s:%d check inter 1s\n", seq+3, slaveIP, s.Port)
	}

	return rwConfig + rConfig
}

func (s Settings) initMasterSqls(serverID int, masterTo string) []string {
	return []string{
		fmt.Sprintf("SET GLOBAL server_id=%d", serverID),
		fmt.Sprintf("CREATE USER '%s'@'%%'", s.ReplUsr),
		fmt.Sprintf("GRANT REPLICATION SLAVE ON *.* "+
			"TO '%s'@'%%' IDENTIFIED BY '%s'", s.ReplUsr, s.ReplPassword),
		"STOP SLAVE",
		fmt.Sprintf("CHANGE MASTER TO master_host='%s', master_port=%d, master_user='%s', "+
			"master_password='%s', master_auto_position = 1", masterTo, s.Port, s.ReplUsr, s.ReplPassword),
		"START SLAVE",
	}
}

func (s Settings) initSlaveSqls(serverID int, masterTo string) []string {
	return []string{
		fmt.Sprintf("SET GLOBAL server_id=%d", serverID),
		"STOP SLAVE",
		fmt.Sprintf("CHANGE MASTER TO master_host='%s', master_port=%d, master_user='%s', "+
			"master_password='%s', master_auto_position = 1", masterTo, s.Port, s.ReplUsr, s.ReplPassword),
		"START SLAVE",
	}
}
