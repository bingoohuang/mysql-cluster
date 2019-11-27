package mci

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gobars/cmd"

	"github.com/bingoohuang/gossh/pbe"

	"github.com/bingoohuang/gonet"
	"github.com/bingoohuang/sqlmore"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

func (s Settings) resetMySQCluster() error {
	s.Host = localhost
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

	if s.isLocalAddr(s.Master1Addr) && !s.Debug {
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

	if err := s.startSlaves(nodes); err != nil {
		return err
	}

	return nil
}

func (s Settings) fixMySQLConf(nodes []MySQLNode) error {
	processed := 0

	for _, node := range nodes {
		if !s.isLocalAddr(node.Addr) {
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

		if s.MySQLUUIDClear {
			if err := executeBash("MySQLRestartShell", 0, s.MySQLRestartShell); err != nil {
				return err
			}
		}
	}

	if processed == 0 {
		logrus.Infof("CreateMySQLCluster bypassed, neither master nor slave for host %v", gonet.ListLocalIps())
	}

	return nil
}

func executeBash(name string, shellTimeout time.Duration, bash string) error {
	if bash == "" {
		logrus.Warnf("%s is empty", name)
		return nil
	}

	logrus.Infof("start execute %s %s", name, bash)

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		}
	}()

	start := time.Now()

	_, status := cmd.Bash(bash, cmd.Timeout(shellTimeout), cmd.Buffered(false))
	if status.Error != nil {
		logrus.Infof("start execute %s %s error %v", name, bash, status.Error)
		return fmt.Errorf("execute %s %s error %w", name, bash, status.Error)
	}

	if status.Exit != 0 {
		logrus.Infof("start execute %s %s exiting code %d, stdout:%s, stderr:%s",
			name, bash, status.Exit, status.Stdout, status.Stderr)

		return fmt.Errorf("execute %s %s exiting code %d, stdout:%s, stderr:%s",
			name, bash, status.Exit, status.Stdout, status.Stderr)
	}

	end := time.Now()

	logrus.Infof("completed execute %s %s cost %v", name, bash, end.Sub(start))

	return nil
}

func (s Settings) createClusters(nodes []MySQLNode) error {
	for _, node := range nodes {
		s.Host = node.Addr
		if err := s.execSqls(node.Sqls); err != nil {
			return err
		}
	}

	return nil
}

func (s Settings) startSlaves(nodes []MySQLNode) error {
	for _, node := range nodes {
		s.Host = node.Addr
		if err := s.execSqls([]string{"start slave"}); err != nil {
			return err
		}
	}

	return nil
}

func (s Settings) resetMaster(nodes []MySQLNode) error {
	for _, node := range nodes {
		s.Host = node.Addr
		if err := s.execSqls([]string{"reset master"}); err != nil {
			return err
		}
	}

	return nil
}

func (s Settings) backupTables(servers []string) error {
	for _, server := range servers {
		s.Host = server
		if _, err := s.renameTables(); err != nil {
			return err
		}
	}

	return nil
}

func (s Settings) prepareCluster(nodes []MySQLNode) error {
	s.Host = localhost

	return s.execSqls([]string{
		fmt.Sprintf("SET GLOBAL server_id=%d", s.findLocalServerID(nodes)),
		"STOP SLAVE", "RESET SLAVE ALL",
		fmt.Sprintf(`DROP USER IF EXISTS '%s'@'%s'`, s.User, s.Master1Addr),
		fmt.Sprintf(`CREATE USER '%s'@'%s' IDENTIFIED BY '%s'`, s.User, s.Master1Addr, s.Password),
		fmt.Sprintf(`GRANT ALL PRIVILEGES ON *.* TO '%s'@'%s' WITH GRANT OPTION`, s.User, s.Master1Addr),
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
	pwd, err := pbe.Ebp(s.Password)
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

func (s Settings) renameTables() (int, error) {
	db := s.MustOpenGormDB()
	defer db.Close()

	return RenameTables(db)
}

func (s Settings) execSqls(sqls []string) error {
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

		logrus.Infof("%s", sqlStr)
	}

	return nil
}

func (s Settings) isLocalAddr(addr string) bool {
	if yes, ok := s.localAddrMap[addr]; ok {
		return yes
	}

	if s.LocalAddr == addr {
		logrus.Infof("%s is local addr", addr)
		return true
	}

	if s.LocalAddr != "" {
		logrus.Infof("%s is not local addr", addr)
		return false
	}

	if yes, _ := gonet.IsLocalAddr(addr); yes {
		logrus.Infof("%s is local addr", addr)
		return yes
	}

	logrus.Infof("%s is not local addr", addr)

	return false
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

	m := make([]MySQLNode, 2+len(s.SlaveAddrs))

	const offset = 10000 // 0-4294967295, https://dev.mysql.com/doc/refman/5.7/en/replication-options.html

	m[0] = MySQLNode{Addr: s.Master1Addr, AutoIncrementOffset: 1, ServerID: offset + 1,
		Sqls: s.initSlaveSqls(s.Master2Addr, replPwd)}

	m[1] = MySQLNode{Addr: s.Master2Addr, AutoIncrementOffset: 2, ServerID: offset + 2,
		Sqls: s.initSlaveSqls(s.Master1Addr, replPwd)}

	for seq, slaveAddr := range s.SlaveAddrs {
		m[2+seq] = MySQLNode{Addr: slaveAddr, AutoIncrementOffset: seq + 3, ServerID: offset + seq + 3,
			Sqls: s.initSlaveSqls(s.Master2Addr, replPwd)}
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
	if !s.MySQLUUIDClear {
		return nil
	}

	if s.Debug {
		fmt.Println("fix fixServerUUID")
		return nil
	}

	if values, err := SearchFileContent(s.MySQLCnf,
		`(?i)datadir\s*=\s*(.+)`); err != nil {
		logrus.Warnf("overwriteHAProxyCnf error: %v", err)

		return err
	} else if len(values) > 0 {
		autoCnf := filepath.Join(strings.TrimSpace(values[0]), "auto.cnf")
		logrus.Infof("remove auto.cnf %s", autoCnf)

		return os.Remove(autoCnf)
	}

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
