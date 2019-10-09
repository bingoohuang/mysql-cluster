package mysqlclusterinit

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/tkrajina/go-reflector/reflector"

	"github.com/bingoohuang/gonet"
	"github.com/bingoohuang/sqlmore"
	"github.com/sirupsen/logrus"
)

func (s Settings) createMySQCluster() (sqls []string, err error) {
	serverID := 0
	serverID, sqls = s.createInitSqls()

	if len(sqls) == 0 {
		logrus.Infof("InitMySQLCluster bypassed, nor master or slave for host %v", gonet.ListLocalIps())
	} else {
		err = s.execMultiSqls(sqls)
		if err != nil {
			return
		}
		err = s.fixMySQLConfServerID(serverID)
	}

	return
}

func (s Settings) execMultiSqls(sqls []string) error {
	if s.Debug {
		return nil
	}

	ds := fmt.Sprintf("root:%s@tcp(127.0.0.1:%d)/", s.RootPassword, s.Port)
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

func (s Settings) createInitSqls() (int, []string) {
	if yes, _ := gonet.IsLocalAddr(s.Master1Addr); yes {
		return 1, s.initMasterSqls(1, s.Master2Addr)
	}

	if yes, _ := gonet.IsLocalAddr(s.Master2Addr); yes {
		return 2, s.initMasterSqls(2, s.Master1Addr)
	}

	for seq, slaveIP := range s.SlaveAddrs {
		if yes, _ := gonet.IsLocalAddr(slaveIP); yes {
			return seq + 3, s.initSlaveSqls(seq+3, s.Master2Addr)
		}
	}

	return 0, []string{}
}

func (s Settings) initMasterSqls(serverID int, masterTo string) []string {
	return []string{
		fmt.Sprintf("SET GLOBAL server_id=%d", serverID),
		fmt.Sprintf("DROP USER IF EXISTS '%s'@'%%'", s.ReplUsr),
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

func (s Settings) fixMySQLConfServerID(serverID int) error {
	if err := ReplaceFileContent(s.MySQLCnf,
		`(?i)server[-_]id\s*=\s*(\d+)`, fmt.Sprintf("%d", serverID)); err != nil {
		return fmt.Errorf("fixMySQLConfServerID %s error %w", s.MySQLCnf, err)
	}

	return nil
}

func ShowSlaveStatus(db *gorm.DB) (bean ShowSlaveStatusBean, err error) {
	if s := db.Raw("show slave status").Scan(&bean); s.Error != nil {
		logrus.Warnf("show slave status error: %v", s.Error)
		return bean, s.Error
	}

	return bean, nil
}

func ShowVariables(db *gorm.DB) (variables Variables, err error) {
	var beans []ShowVariablesBean
	if s := db.Raw("show variables").Scan(&beans); s.Error != nil {
		logrus.Warnf("show variables error: %v", s.Error)
		return Variables{}, s.Error
	}

	variablesMap := make(map[string]string)
	for _, v := range beans {
		variablesMap[v.VariableName] = v.Value
	}

	for _, f := range reflector.New(&variables).Fields() {
		tag, _ := f.Tag("var")
		if v, ok := variablesMap[tag]; !ok {
			continue
		} else if err := f.Set(v); err != nil {
			logrus.Warnf("Set error: %v", err)
		}
	}

	return variables, nil
}
