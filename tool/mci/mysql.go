package mci

import (
	"database/sql"
	"fmt"
	"io"
	"strings"

	"github.com/bingoohuang/gossh/pbe"

	"github.com/bingoohuang/now"

	"github.com/bingoohuang/gonet"
	"github.com/bingoohuang/sqlmore"
	"github.com/jedib0t/go-pretty/table"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"github.com/tkrajina/go-reflector/reflector"
)

func (s Settings) createMySQCluster() ([]MySQLNode, error) {
	nodes := s.createInitSqls()

	if !s.Debug { // 所有节点都做root向master1的root授权
		if err := s.prepareCluster(nodes); err != nil {
			return nodes, err
		}
	}

	if s.isLocalAddr(s.Master1Addr) && !s.Debug {
		mysqlServers := []string{s.Master1Addr, s.Master2Addr}
		mysqlServers = append(mysqlServers, s.SlaveAddrs...)

		if err := s.stopSlaves(mysqlServers); err != nil {
			return nodes, err
		}

		if err := s.backupTables(mysqlServers); err != nil {
			return nodes, err
		}

		if err := s.createClusters(nodes); err != nil {
			return nodes, err
		}
	}

	if err := s.fixMySQLConf(nodes); err != nil {
		return nodes, err
	}

	return nodes, nil
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
	for _, server := range servers[1:] {
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

	// 授权 GRANT ALL ON *.* TO root@'192.168.136.23' IDENTIFIED BY 'xx';
	// 回收：DROP USER root@'192.168.136.23';
	return s.execMultiSqls([]string{
		fmt.Sprintf("SET GLOBAL server_id=%d", s.findServerID(nodes)),
		fmt.Sprintf(`GRANT ALL ON *.* TO root@'%s' IDENTIFIED BY '%s'`, s.Master1Addr, s.Password),
	})
}

func (s Settings) findServerID(nodes []MySQLNode) int {
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
	return []string{
		fmt.Sprintf("DROP USER IF EXISTS '%s'@'%%'", s.ReplUsr),
		fmt.Sprintf("CREATE USER '%s'@'%%' IDENTIFIED BY '%s'", s.ReplUsr, replPwd),
		fmt.Sprintf("GRANT REPLICATION SLAVE ON *.* "+"TO '%s'@'%%' IDENTIFIED BY '%s'", s.ReplUsr, replPwd),
		"RESET SLAVE",
		fmt.Sprintf("CHANGE MASTER TO master_host='%s', master_port=%d, master_user='%s', "+
			"master_password='%s', master_auto_position = 1", masterTo, s.Port, s.ReplUsr, replPwd),
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

// ShowTables show all tables
func ShowTables(db *gorm.DB, postfix string) (beans []TableBean, err error) {
	sql := `select * from information_schema.tables
		where TABLE_SCHEMA not in ('performance_schema', 'information_schema', 'mysql', 'sys') 
		and TABLE_NAME not like '%` + postfix + `'`

	if s := db.Raw(sql).Scan(&beans); s.Error != nil {
		logrus.Warnf("show slave status error: %v", s.Error)
		return beans, s.Error
	}

	return beans, nil
}

// RenameTables rename the non-system databases' table to another name.
func RenameTables(db *gorm.DB, postfix string) error {
	tables, err := ShowTables(db, postfix)
	if err != nil {
		return err
	}

	if len(tables) == 0 {
		return nil
	}

	renameSqls := make([]string, len(tables))
	for i, t := range tables {
		renameSqls[i] = fmt.Sprintf("%s.%s to %s.%s%s",
			t.Schema, t.Name, t.Schema, t.Name, postfix)
	}

	// https://dev.mysql.com/doc/refman/5.7/en/rename-table.html
	// RENAME TABLE
	//    tbl_name TO new_tbl_name
	//    [, tbl_name2 TO new_tbl_name2] ...
	joined := "rename table " + strings.Join(renameSqls, ", ")
	logrus.Infof("sql:%s", joined)

	return db.Exec(joined).Error
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

// PrintSQLResult prints the result r of sqlStr execution
func PrintSQLResult(stdout, stderr io.Writer, sqlStr string, r sqlmore.ExecResult) error {
	if r.Error != nil {
		fmt.Fprintf(stderr, "error %v\n", r.Error)
		return r.Error
	}

	fmt.Fprintf(stdout, "SQL: %s\n", sqlStr)
	fmt.Fprintf(stdout, "Cost: %s\n", r.CostTime.String())

	if !r.IsQuerySQL {
		return nil
	}

	cols := len(r.Headers) + 1
	header := make(table.Row, cols)
	header[0] = "#"

	for i, h := range r.Headers {
		header[i+1] = h
	}

	t := table.NewWriter()
	t.SetOutputMirror(stdout)
	t.AppendHeader(header)

	for i, r := range r.Rows {
		row := make(table.Row, cols)
		row[0] = i + 1

		for j, c := range r {
			row[j+1] = c
		}

		t.AppendRow(row)
	}

	t.Render()

	return nil
}
