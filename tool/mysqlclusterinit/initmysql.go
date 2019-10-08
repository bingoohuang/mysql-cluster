package mysqlclusterinit

import (
	"database/sql"
	"fmt"

	"github.com/bingoohuang/gonet"
	"github.com/bingoohuang/sqlmore"
	"github.com/sirupsen/logrus"
)

// MySQLClusterSettings 表示舒适化MySQL集群所需要的参数结构
type MySQLClusterSettings struct {
	MasterIP1    string
	MasterIP2    string
	SlaveIps     []string
	RootPassword string
	Port         int
	ReplUsr      string
	ReplPassword string
	Debug        bool // 测试模式，只打印SQL和HAProxy配置, 不实际执行
	LocalIP      string
}

type Result struct {
	Error   error
	Sqls    []string
	HAProxy string
}

// InitMySQLCluster 初始化MySQL Master-Master集群
func (s MySQLClusterSettings) InitMySQLCluster() Result {
	var result Result
	var err error
	if result.Sqls, err = s.createMySQCluster(); err != nil {
		result.Error = err
		return result
	}

	result.HAProxy = s.createHAProxyConfig()

	return result
}

func (s MySQLClusterSettings) createMySQCluster() ([]string, error) {
	localIP := s.initLocalIPMap()

	sqls := s.createInitSqls(localIP)
	if len(sqls) == 0 {
		logrus.Infof("InitMySQLCluster bypassed, nor master or slave on %v", localIP)
		return sqls, nil
	}

	if err := s.execMultiSqls(sqls); err != nil {
		return sqls, err
	}

	return sqls, nil
}

func (s MySQLClusterSettings) initLocalIPMap() map[string]bool {
	var localIP map[string]bool
	if s.LocalIP == "" {
		localIP = gonet.ListLocalIPMap()
	} else {
		localIP = make(map[string]bool)
		localIP[s.LocalIP] = true
	}
	return localIP
}

func (s MySQLClusterSettings) execMultiSqls(sqls []string) error {
	if s.Debug {
		return nil
	}

	ds := fmt.Sprintf("root:%s@tcp(127.0.0.1:%d)/root", s.RootPassword, s.Port)
	more := sqlmore.NewSQLMore("mysql", ds)
	db := more.MustOpen()
	defer db.Close()
	for _, sqlStr := range sqls {
		if err := s.execSQL(db, sqlStr); err != nil {
			return err
		}
	}
	logrus.Infof("createMySQCluster completed")
	return nil
}

func (s MySQLClusterSettings) createInitSqls(localIP map[string]bool) []string {
	if _, ok := localIP[s.MasterIP1]; ok {
		return s.initMasterSqls(1, s.MasterIP2)
	}
	if _, ok := localIP[s.MasterIP2]; ok {
		return s.initMasterSqls(2, s.MasterIP1)
	}

	for seq, slaveIP := range s.SlaveIps {
		if _, ok := localIP[slaveIP]; ok {
			return s.initSlaveSqls(seq+3, s.MasterIP2)
		}
	}

	return nil
}

func (s MySQLClusterSettings) createHAProxyConfig() string {
	rwConfig := fmt.Sprintf(`
listen mysql-rw
  bind 0.0.0.0:13306
  mode tcp
  option tcpka
  server mysql-1 %s:%d check inter 1s
  server mysql-2 %s:%d check inter 1s backup
`, s.MasterIP1, s.Port, s.MasterIP2, s.Port)

	rConfig := fmt.Sprintf(`
listen mysql-ro
  bind 0.0.0.0:23306
  mode tcp
  option tcpka
  server mysql-1 %s:%d check inter 1s
  server mysql-2 %s:%d check inter 1s
`, s.MasterIP1, s.Port, s.MasterIP2, s.Port)

	for seq, slaveIP := range s.SlaveIps {
		rConfig += fmt.Sprintf("  server mysql-%d %s:%d check inter 1s\n", seq+3, slaveIP, s.Port)
	}

	return rwConfig + rConfig
}

func (s MySQLClusterSettings) initMasterSqls(serverID int, masterTo string) []string {
	sqls := make([]string, 0)
	sqls = append(sqls, fmt.Sprintf("SET GLOBAL server_id=%d", serverID))
	sqls = append(sqls, fmt.Sprintf("CREATE USER '%s'@'%%'", s.ReplUsr))
	sqls = append(sqls, fmt.Sprintf("GRANT REPLICATION SLAVE ON *.* "+
		"TO '%s'@'%%' IDENTIFIED BY '%s'", s.ReplUsr, s.ReplPassword))
	sqls = append(sqls, "STOP SLAVE")
	sqls = append(sqls, fmt.Sprintf("CHANGE MASTER TO master_host='%s', master_port=%d, master_user='%s', "+
		"master_password='%s', master_auto_position = 1", masterTo, s.Port, s.ReplUsr, s.ReplPassword))
	sqls = append(sqls, "START SLAVE")

	return sqls
}

func (s MySQLClusterSettings) initSlaveSqls(serverID int, masterTo string) []string {
	sqls := make([]string, 0)
	sqls = append(sqls, fmt.Sprintf("SET GLOBAL server_id=%d", serverID))
	sqls = append(sqls, "STOP SLAVE")
	sqls = append(sqls, fmt.Sprintf("CHANGE MASTER TO master_host='%s', master_port=%d, master_user='%s', "+
		"master_password='%s', master_auto_position = 1", masterTo, s.Port, s.ReplUsr, s.ReplPassword))
	sqls = append(sqls, "START SLAVE")

	return sqls
}

func (s MySQLClusterSettings) execSQL(db *sql.DB, sqlStr string) error {
	if result := sqlmore.ExecSQL(db, sqlStr, 0, ""); result.Error != nil {
		return result.Error
	}

	logrus.Infof("execSQL %s completed", sqlStr)
	return nil
}
