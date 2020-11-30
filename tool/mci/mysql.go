package mci

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/bingoohuang/gonet"
	"github.com/bingoohuang/gou/pbe"
	"github.com/bingoohuang/sqlx"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func (s Settings) resetMySQCluster() error {
	s.currentHost = localhostIPv4
	resetSqls := s.resetSlaveSqls()

	return s.execSqls(resetSqls)
}

func (s Settings) createMySQCluster() ([]MySQLNode, error) {
	nodes := s.createInitSqls()

	if err := s.fixMySQLConf(nodes); err != nil {
		return nodes, err
	}

	if !s.Debug { // 所有节点都做root向master1的root授权
		if err := s.prepareCluster(nodes); err != nil {
			return nodes, err
		}
	}

	if IsLocalAddr(s.Master1Addr) && !s.Debug {
		if err := s.master1LocalProcess(nodes); err != nil {
			return nodes, err
		}
	}

	return nodes, nil
}

func (s Settings) master1LocalProcess(nodes []MySQLNode) error {
	backupServers := []string{s.Master2Addr}
	backupServers = append(backupServers, s.SlaveAddrs...)

	if err := s.backupTables(backupServers); err != nil {
		return err
	}

	if err := s.createClusters(nodes); err != nil {
		return err
	}

	if err := s.copyMaster1Data(backupServers); err != nil {
		return err
	}

	if err := s.resetMaster(nodes); err != nil {
		return err
	}

	return s.startSlaves(nodes)
}

func (s Settings) fixMySQLConf(nodes []MySQLNode) error {
	processed := 0

	for _, node := range nodes {
		if !IsLocalAddr(node.Addr) {
			continue
		}

		processed++

		if err := s.fixMySQLConfServerID(node.ServerID); err != nil {
			return err
		}

		if err := s.fixAutoIncrementOffset(node.AutoIncrementOffset); err != nil {
			return err
		}

		if err := s.fixServerUUID(); err != nil {
			return err
		}

		if err := ExecuteBash("MySQLRestartShell", s.MySQLRestartShell, 0); err != nil {
			return err
		}
	}

	if processed == 0 {
		logrus.Infof("CreateMySQLCluster bypassed, neither master nor slave for host %v", gonet.ListIps())
	}

	return nil
}

func (s Settings) createClusters(nodes []MySQLNode) error {
	for _, node := range nodes {
		s.currentHost = node.Addr
		if err := s.execSqls(node.Sqls); err != nil {
			return err
		}
	}

	return nil
}

func (s Settings) startSlaves(nodes []MySQLNode) error {
	for _, node := range nodes {
		s.currentHost = node.Addr
		if err := s.execSqls([]string{"start slave"}); err != nil {
			return err
		}
	}

	return nil
}

func (s Settings) resetMaster(nodes []MySQLNode) error {
	for _, node := range nodes {
		s.currentHost = node.Addr
		if err := s.execSqls([]string{"reset master"}); err != nil {
			return err
		}
	}

	return nil
}

func (s Settings) backupTables(servers []string) error {
	for _, server := range servers {
		s.currentHost = server
		if _, err := s.renameTables(s.NoBackup); err != nil {
			return err
		}
	}

	return nil
}

func (s Settings) prepareCluster(nodes []MySQLNode) error {
	s.currentHost = localhostIPv4

	return s.execSqls([]string{
		fmt.Sprintf("SET GLOBAL server_id=%d", s.findLocalServerID(nodes)),
		"STOP SLAVE", "RESET SLAVE ALL",
		fmt.Sprintf(`DROP USER IF EXISTS '%s'@'%s'`, s.User, s.Master1Addr),
		fmt.Sprintf(`CREATE USER '%s'@'%s' IDENTIFIED BY '%s'`, s.User, s.Master1Addr, s.Password),
		fmt.Sprintf(`GRANT ALL PRIVILEGES ON *.* TO '%s'@'%s' WITH GRANT OPTION`, s.User, s.Master1Addr),
		// DROP USER IF EXISTS 'root'@'%'; CREATE USER 'root'@'%' IDENTIFIED BY 'C2D747DB89F6_a';
		// GRANT ALL PRIVILEGES ON *.* TO 'root'@'%' WITH GRANT OPTION;
	})
}

func (s Settings) findLocalServerID(nodes []MySQLNode) int {
	for _, node := range nodes {
		if IsLocalAddr(node.Addr) {
			return node.ServerID
		}
	}

	return 0
}

// MustOpenDB must open the db.
// nolint:gomnd
func (s Settings) MustOpenDB() *sql.DB {
	pwd, err := pbe.Ebp(s.Password)
	if err != nil {
		panic(err)
	}

	ds := ""

	host := s.currentHost
	net := "tcp"

	// 不是连接本机MySQL，设置net名称，指定IP出口
	if IsLocalAddr(s.Master1Addr) && !IsLocalAddr(host) {
		net = "master1"
	}

	if gonet.IsIPv6(host) {
		ds = fmt.Sprintf("%s:%s@%s([%s]:%d)/", s.User, pwd, net, host, s.Port)
	} else {
		ds = fmt.Sprintf("%s:%s@%s(%s:%d)/", s.User, pwd, net, host, s.Port)
	}

	ds += `?` + s.MySQLDSNParams

	logrus.Infof("s.Master1Addr %s, host:%s, ds: %s", s.Master1Addr, ds, host)

	// 只有主1连接其他MySQL时，才设置
	if IsLocalAddr(s.Master1Addr) {
		viper.Set("mysqlNet", "master1")
		viper.Set("bindAddress", s.Master1Addr)
		logrus.Infof("bindAddress %s", s.Master1Addr)
	}

	sqlMore := sqlx.NewSQLMore("mysql", ds)

	if !s.NoLog {
		logrus.Debugf("DSN: %s", sqlMore.EnhancedURI)
	}

	db := sqlMore.Open()
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(60 * time.Second)

	return db
}

// MustOpenGormDB must open the db.
func (s Settings) MustOpenGormDB() *gorm.DB {
	gdb, _ := gorm.Open("mysql", s.MustOpenDB())
	return gdb
}

func (s Settings) renameTables(noBackup bool) (int, error) {
	db := s.MustOpenGormDB()
	defer db.Close()

	return RenameTables(db, noBackup)
}

func (s Settings) execSqls(sqls []string) error {
	if s.Debug {
		fmt.Print(strings.Join(sqls, ";\n"))
		return nil
	}

	db := s.MustOpenDB()
	defer db.Close()

	for _, sqlStr := range sqls {
		r := sqlx.ExecSQL(db, sqlStr, 0, "")
		if r.Error != nil {
			return fmt.Errorf("exec sql %s error %w", sqlStr, r.Error)
		}

		logrus.Infof("SQL:%s, %+v", sqlStr, r)
	}

	return nil
}

// MySQLNode means the information of MySQLNode in the cluster.
type MySQLNode struct {
	Addr                string
	AutoIncrementOffset int
	ServerID            int
	Sqls                []string
}

func (s Settings) createInitSqls() []MySQLNode {
	replPwd, err := pbe.Ebp(s.ReplPassword)
	if err != nil {
		panic(err)
	}

	m := make([]MySQLNode, 2+len(s.SlaveAddrs)) // nolint:gomnd

	const offset = 10000 // 0-4294967295, https://dev.mysql.com/doc/refman/5.7/en/replication-options.html

	m[0] = MySQLNode{
		Addr: s.Master1Addr, AutoIncrementOffset: 1, ServerID: offset + 1,
		Sqls: s.initSlaveSqls(s.Master2Addr, replPwd),
	}

	m[1] = MySQLNode{
		Addr: s.Master2Addr, AutoIncrementOffset: 2, ServerID: offset + 2, // nolint:gomnd
		Sqls: s.initSlaveSqls(s.Master1Addr, replPwd),
	}

	for seq, slaveAddr := range s.SlaveAddrs {
		m[2+seq] = MySQLNode{
			Addr: slaveAddr, AutoIncrementOffset: seq + 3, ServerID: offset + seq + 3, // nolint:gomnd
			Sqls: s.initSlaveSqls(s.Master2Addr, replPwd),
		}
	}

	return m
}

const (
	deleteUsers  = "DELETE FROM mysql.user WHERE user='%s'"
	createUser   = "CREATE USER '%s'@'%s' IDENTIFIED BY '%s'"
	grantSlave   = "GRANT REPLICATION SLAVE ON *.* TO '%s'@'%s' IDENTIFIED BY '%s'"
	changeMaster = `CHANGE MASTER TO master_host='%s',master_port=%d,master_user='%s',` +
		`master_password='%s',master_auto_position=1`
)

// https://dev.mysql.com/doc/refman/5.7/en/reset-slave.html
// RESET SLAVE makes the slave forget its replication position in the master's binary log.
// This statement is meant to be used for a clean start: It clears the master info
// and relay log info repositories, deletes all the relay log files,
// and starts a new relay log file. It also resets to 0 the replication delay specified
// with the MASTER_DELAY option to CHANGE MASTER TO.

// https://stackoverflow.com/a/32148683
// RESET SLAVE will leave behind master.info file with "default" values in such a way
// that SHOW SLAVE STATUS will still give output. So if you have slave monitoring on this host,
//after it becomes the master, you would still get alarms that are checking for 'Slave_IO_Running: Yes'
//
// RESET SLAVE ALL wipes slave info clean away, deleting master.info and
// SHOW SLAVE STATUS will report "Empty Set (0.00)"

func (s Settings) initSlaveSqls(masterTo, replPwd string) []string {
	sqls := []string{fmt.Sprintf(deleteUsers, s.ReplUsr), "FLUSH PRIVILEGES"}

	args := []interface{}{s.ReplUsr, "%", replPwd}
	sqls = append(sqls, fmt.Sprintf(createUser, args...), fmt.Sprintf(grantSlave, args...))

	return append(sqls, fmt.Sprintf(changeMaster, masterTo, s.Port, s.ReplUsr, replPwd))
}

func (s Settings) resetSlaveSqls() []string {
	return []string{
		fmt.Sprintf(deleteUsers, s.ReplUsr),
		"STOP SLAVE", "RESET SLAVE ALL",
		fmt.Sprintf(`DROP USER IF EXISTS '%s'@'%s'`, s.User, s.Master1Addr),
		"FLUSH PRIVILEGES",
	}
}

func (s *Settings) fixServerUUID() error {
	if s.Debug {
		fmt.Println("fix fixServerUUID")
		return nil
	}

	if values, err := SearchFileContent(s.MySQLCnf,
		`(?i)datadir\s*=\s*(.+)`); err != nil {
		logrus.Warnf("SearchFileContent error: %v", err)

		return err
	} else if len(values) > 0 {
		autoCnf := filepath.Join(strings.TrimSpace(values[0]), "auto.cnf")
		logrus.Infof("remove auto.cnf %s", autoCnf)

		return os.Remove(autoCnf)
	}

	// nolint:goerr113
	return fmt.Errorf("unable to find datadir in %s", s.MySQLCnf)
}

func (s Settings) fixMySQLConfServerID(serverID int) error {
	if s.Debug {
		fmt.Println("fix server-id =", serverID)
		return nil
	}

	if err := ReplaceFileContent(s.MySQLCnf,
		`(?i)server[-_]id\s*=\s*(\d+)`, fmt.Sprintf("%d", serverID)); err != nil {
		return fmt.Errorf("fixMySQLConfServerID %s error %w", s.MySQLCnf, err)
	}

	return nil
}

// fixAutoIncrementOffset fix auto_increment_offset.
func (s Settings) fixAutoIncrementOffset(offset int) error {
	if s.Debug {
		fmt.Println("fix increment-offset =", offset)
		return nil
	}

	if err := ReplaceFileContent(s.MySQLCnf,
		`(?i)auto[-_]increment[-_]offset\s*=\s*(\d+)`, fmt.Sprintf("%d", offset)); err != nil {
		return fmt.Errorf("fixAutoIncrementOffset %s error %w", s.MySQLCnf, err)
	}

	return nil
}
