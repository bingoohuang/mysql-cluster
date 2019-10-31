package mci

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/bingoohuang/gossh/pbe"

	"github.com/bingoohuang/now"

	"github.com/bingoohuang/gonet"
	"github.com/bingoohuang/sqlmore"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

func (s Settings) createMySQCluster() ([]MySQLNode, error) {
	nodes := s.createInitSqls()

	if !s.Debug { // 所有节点都做root向master1的root授权
		if err := s.prepareCluster(nodes); err != nil {
			return nodes, err
		}
	}

	if s.isLocalAddr(s.Master1Addr) && !s.Debug {
		if err := s.master1LocalProcess(nodes); err != nil {
			return nodes, err
		}
	}

	if err := s.fixMySQLConf(nodes); err != nil {
		return nodes, err
	}

	return nodes, nil
}

func (s Settings) master1LocalProcess(nodes []MySQLNode) error {
	mysqlServers := []string{s.Master1Addr, s.Master2Addr}
	mysqlServers = append(mysqlServers, s.SlaveAddrs...)

	if err := s.stopSlaves(mysqlServers); err != nil {
		return err
	}

	backupServers := mysqlServers[1:]
	if err := s.backupTables(backupServers); err != nil {
		return err
	}

	if err := s.createClusters(nodes); err != nil {
		return err
	}

	return nil
}

func (s Settings) fixMySQLConf(nodes []MySQLNode) error {
	for _, node := range nodes {
		if !s.isLocalAddr(node.Addr) {
			continue
		}

		if err := s.fixMySQLConfServerID(node.ServerID); err != nil {
			return err
		}

		if err := s.fixAutoIncrementOffset(node.AutoIncrementOffset); err != nil {
			return err
		}

		return nil
	}

	logrus.Infof("InitMySQLCluster bypassed, neither master nor slave for host %v", gonet.ListLocalIps())

	return nil
}

func (s Settings) createClusters(nodes []MySQLNode) error {
	for _, node := range nodes {
		s.Host = node.Addr
		if err := s.execMultiSqls(node.Sqls); err != nil {
			return err
		}
	}

	return nil
}

func (s Settings) backupTables(servers []string) error {
	for _, server := range servers {
		s.Host = server
		postfix := "_mci" + now.MakeNow().Format("yyyyMMdd")

		if err := s.renameTables(postfix); err != nil {
			return err
		}
	}

	return nil
}

func (s Settings) stopSlaves(servers []string) error {
	for _, server := range servers {
		s.Host = server
		if err := s.execSQL("stop slave"); err != nil {
			return err
		}
	}

	return nil
}

func (s Settings) prepareCluster(nodes []MySQLNode) error {
	s.Host = "127.0.0.1"
	fs := fmt.Sprintf

	return s.execMultiSqls([]string{
		fs("SET GLOBAL server_id=%d", s.findLocalServerID(nodes)),
		fs(`DROP USER IF EXISTS '%s'@'%s'`, s.User, s.Master1Addr),
		fs(`CREATE USER '%s'@'%s' IDENTIFIED BY '%s'`, s.User, s.Master1Addr, s.Password),
		fs(`GRANT ALL PRIVILEGES ON *.* TO '%s'@'%s' WITH GRANT OPTION`, s.User, s.Master1Addr),
	})
}

func (s Settings) findLocalServerID(nodes []MySQLNode) int {
	for _, node := range nodes {
		if s.isLocalAddr(node.Addr) {
			return node.ServerID
		}
	}

	return 0
}

// MustOpenDB must open the db.
func (s Settings) MustOpenDB() *sql.DB {
	pwd, err := ebpFix(s.Password)
	if err != nil {
		panic(err)
	}

	ds := fmt.Sprintf("%s:%s@tcp(%s:%d)/", s.User, pwd, s.Host, s.Port)
	logrus.Infof("mysql ds:%s", ds)

	return sqlmore.NewSQLMore("mysql", ds).MustOpen()
}

// MustOpenGormDB must open the db.
func (s Settings) MustOpenGormDB() *gorm.DB {
	gdb, _ := gorm.Open("mysql", s.MustOpenDB())
	return gdb
}

func (s Settings) renameTables(postfix string) error {
	db := s.MustOpenGormDB()
	defer db.Close()

	return RenameTables(db, postfix)
}

func (s Settings) execSQL(sqlStr string) error {
	if s.Debug {
		fmt.Println(sqlStr + ";")
		return nil
	}

	db := s.MustOpenDB()
	defer db.Close()

	if r := sqlmore.ExecSQL(db, sqlStr, 0, ""); r.Error != nil {
		return fmt.Errorf("exec sql %s error %w", sqlStr, r.Error)
	}

	logrus.Infof("execSQL %s completed", sqlStr)

	return nil
}

func (s Settings) execMultiSqls(sqls []string) error {
	if s.Debug {
		fmt.Print(strings.Join(sqls, ";\n"))
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

// MySQLNode means the information of MySQLNode in the cluster.
type MySQLNode struct {
	Addr                string
	AutoIncrementOffset int
	ServerID            int
	Sqls                []string
}

// ebpFix will be removed later after pbe upgraded.
func ebpFix(p string) (string, error) {
	if strings.HasPrefix(`{PBE}`, p) {
		return pbe.Ebp(p)
	}

	return p, nil
}

func (s Settings) createInitSqls() []MySQLNode {
	replPwd, err := ebpFix(s.ReplPassword)
	if err != nil {
		panic(err)
	}

	m := make([]MySQLNode, 0)

	const offset = 10000 // 0-4294967295, https://dev.mysql.com/doc/refman/5.7/en/replication-options.html

	m = append(m, MySQLNode{
		Addr:                s.Master1Addr,
		AutoIncrementOffset: 1,
		ServerID:            offset + 1,
		Sqls:                s.initMasterSqls(s.Master2Addr, replPwd),
	})

	m = append(m, MySQLNode{
		Addr:                s.Master2Addr,
		AutoIncrementOffset: 2,
		ServerID:            offset + 2,
		Sqls:                s.initMasterSqls(s.Master1Addr, replPwd),
	})

	for seq, slaveAddr := range s.SlaveAddrs {
		m = append(m, MySQLNode{
			Addr:                slaveAddr,
			AutoIncrementOffset: seq + 3,
			ServerID:            offset + seq + 3,
			Sqls:                s.initSlaveSqls(s.Master2Addr, replPwd),
		})
	}

	return m
}

// https://dev.mysql.com/doc/refman/5.7/en/reset-slave.html
// RESET SLAVE makes the slave forget its replication position in the master's binary log.
// This statement is meant to be used for a clean start: It clears the master info
// and relay log info repositories, deletes all the relay log files,
// and starts a new relay log file. It also resets to 0 the replication delay specified
// with the MASTER_DELAY option to CHANGE MASTER TO.
func (s Settings) initMasterSqls(masterTo, replPwd string) []string {
	fs := fmt.Sprintf

	return []string{
		fs("DROP USER IF EXISTS '%s'@'%%'", s.ReplUsr),
		fs("CREATE USER '%s'@'%%' IDENTIFIED BY '%s'", s.ReplUsr, replPwd),
		fs("GRANT REPLICATION SLAVE ON *.* TO '%s'@'%%' IDENTIFIED BY '%s'", s.ReplUsr, replPwd),
		"RESET SLAVE",
		fs(`CHANGE MASTER TO master_host='%s', master_port=%d, master_user='%s', 
			master_password='%s', master_auto_position = 1`, masterTo, s.Port, s.ReplUsr, replPwd),
		"START SLAVE",
	}
}

func (s Settings) initSlaveSqls(masterTo, replPwd string) []string {
	return []string{
		"RESET SLAVE",
		fmt.Sprintf("CHANGE MASTER TO master_host='%s', master_port=%d, master_user='%s', "+
			"master_password='%s', master_auto_position = 1", masterTo, s.Port, s.ReplUsr, replPwd),
		"START SLAVE",
	}
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

// fixAutoIncrementOffset fix auto_increment_offset
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
