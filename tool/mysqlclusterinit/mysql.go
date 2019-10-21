package mysqlclusterinit

import (
	"database/sql"
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/tkrajina/go-reflector/reflector"

	"github.com/bingoohuang/gonet"
	"github.com/bingoohuang/sqlmore"
	"github.com/sirupsen/logrus"
)

func (s Settings) createMySQCluster() (sqls []string, err error) {
	seq := 0
	seq, sqls = s.createInitSqls()

	if len(sqls) == 0 {
		logrus.Infof("InitMySQLCluster bypassed, nor master or slave for host %v", gonet.ListLocalIps())
	}

	if err := s.execMultiSqls(sqls); err != nil {
		return sqls, err
	}

	if err := s.fixMySQLConfServerID(seq); err != nil {
		return sqls, err
	}

	if err := s.fixAutoIncrementOffset(seq); err != nil {
		return sqls, err
	}

	return sqls, nil
}

func (s Settings) MustOpenDB() *sql.DB {
	ds := fmt.Sprintf("root:%s@tcp(127.0.0.1:%d)/", s.RootPassword, s.Port)
	logrus.Infof("mysql ds:%s", ds)

	return sqlmore.NewSQLMore("mysql", ds).MustOpen()
}

func (s Settings) MustOpenGormDB() *gorm.DB {
	gdb, _ := gorm.Open("mysql", s.MustOpenDB())
	return gdb
}

func (s Settings) execMultiSqls(sqls []string) error {
	if s.Debug {
		return nil
	}

	db := s.MustOpenDB()
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

func (s Settings) isLocalAddr(addr string) bool {
	if s.LocalAddr == addr {
		return true
	}

	if s.LocalAddr != "" {
		return false
	}

	if yes, _ := gonet.IsLocalAddr(addr); yes {
		return yes
	}

	return false
}

func (s Settings) createInitSqls() (int, []string) {
	if s.isLocalAddr(s.Master1Addr) {
		return 1, s.initMasterSqls(1, s.Master2Addr)
	}

	if s.isLocalAddr(s.Master2Addr) {
		return 2, s.initMasterSqls(2, s.Master1Addr)
	}

	for seq, slaveIP := range s.SlaveAddrs {
		if s.isLocalAddr(slaveIP) {
			return seq + 3, s.initSlaveSqls(seq+3, s.Master2Addr)
		}
	}

	return 0, []string{}
}

func (s Settings) initMasterSqls(serverID int, masterTo string) []string {
	return []string{
		fmt.Sprintf("SET GLOBAL server_id=%d", serverID),
		fmt.Sprintf("DROP USER IF EXISTS '%s'@'%%'", s.ReplUsr),
		fmt.Sprintf("CREATE USER '%s'@'%%' IDENTIFIED BY '%s'", s.ReplUsr, s.ReplPassword),
		fmt.Sprintf("GRANT REPLICATION SLAVE ON *.* "+
			"TO '%s'@'%%' IDENTIFIED BY '%s'", s.ReplUsr, s.ReplPassword),
		"STOP SLAVE",
		"RESET SLAVE",
		fmt.Sprintf("CHANGE MASTER TO master_host='%s', master_port=%d, master_user='%s', "+
			"master_password='%s', master_auto_position = 1", masterTo, s.Port, s.ReplUsr, s.ReplPassword),
		"START SLAVE",
	}
}

func (s Settings) initSlaveSqls(serverID int, masterTo string) []string {
	return []string{
		fmt.Sprintf("SET GLOBAL server_id=%d", serverID),
		"STOP SLAVE",
		"RESET SLAVE",
		fmt.Sprintf("CHANGE MASTER TO master_host='%s', master_port=%d, master_user='%s', "+
			"master_password='%s', master_auto_position = 1", masterTo, s.Port, s.ReplUsr, s.ReplPassword),
		"START SLAVE",
	}
}

func (s Settings) fixMySQLConfServerID(serverID int) error {
	if s.Debug {
		return nil
	}

	if err := ReplaceFileContent(s.MySQLCnf,
		`(?i)server[-_]id\s*=\s*(\d+)`, fmt.Sprintf("%d", serverID)); err != nil {
		return fmt.Errorf("fixMySQLConfServerID %s error %w", s.MySQLCnf, err)
	}

	return nil
}

// auto_increment_offset
func (s Settings) fixAutoIncrementOffset(offset int) error {
	if s.Debug {
		return nil
	}

	if err := ReplaceFileContent(s.MySQLCnf,
		`(?i)auto[-_]increment[-_]offset\s*=\s*(\d+)`, fmt.Sprintf("%d", offset)); err != nil {
		return fmt.Errorf("fixAutoIncrementOffset %s error %w", s.MySQLCnf, err)
	}

	return nil
}

// ShowSlaveStatus show slave status to bean
func ShowSlaveStatus(db *gorm.DB) (bean ShowSlaveStatusBean, err error) {
	if s := db.Raw("show slave status").Scan(&bean); s.Error != nil {
		logrus.Warnf("show slave status error: %v", s.Error)
		return bean, s.Error
	}

	return bean, nil
}

// ShowVariables shows variables to variables bean
func ShowVariables(db *gorm.DB) (variables Variables, err error) {
	fieldsMap := make(map[string]reflector.ObjField)

	for _, f := range reflector.New(&variables).Fields() {
		if tag, _ := f.Tag("var"); tag != "" {
			fieldsMap[tag] = f
		}
	}

	var beans []ShowVariablesBean
	if s := db.Raw("show variables").Scan(&beans); s.Error != nil {
		logrus.Warnf("show variables error: %v", s.Error)
		return Variables{}, s.Error
	}

	for _, b := range beans {
		if f, ok := fieldsMap[b.VariableName]; !ok {
			continue
		} else if err := f.Set(b.Value); err != nil {
			logrus.Warnf("Set error: %v", err)
		}
	}

	return variables, nil
}
