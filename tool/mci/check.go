package mci

import (
	"fmt"
	"os"
	"strings"

	"github.com/bingoohuang/sqlmore"
)

// CheckMySQLCluster 检查MySQL集群配置
func (s Settings) CheckMySQLCluster() {
	if s.ValidateAndSetDefault(SetDefault) != nil {
		os.Exit(1)
	}

	db := s.MustOpenGormDB()
	defer db.Close()

	if status, err := ShowSlaveStatus(db); err == nil {
		fmt.Printf("ShowSlaveStatus:%s\n", JSONPretty(status))
	}

	if variables, err := ShowVariables(db); err == nil {
		fmt.Printf("Variables:%s\n", JSONPretty(variables))
	}
}

// CheckHAProxyServers 检查HAProxy中的MySQL集群配置
func (s Settings) CheckHAProxyServers() {
	if s.ValidateAndSetDefault(SetDefault) != nil {
		os.Exit(1)
	}

	linesInFile, err := SearchPatternLinesInFile(s.HAProxyCfg,
		`(?is)mysql-ro(.+)MySQLClusterConfigEnd`, `(?i)server\s+\S+\s(\d+(\.\d+){3})(:\d+)?`)
	if err != nil {
		fmt.Printf("SearchPatternLinesInFile error %v\n", err)
		return
	}

	fmt.Println(strings.Join(linesInFile, "\n"))
}

// CheckMySQL 检查MySQL连接
// refer https://github.com/zhishutech/mysqlha-keepalived-3node/blob/master/keepalived/checkMySQL.py
func (s Settings) CheckMySQL() {
	if s.ValidateAndSetDefault(SetDefault) != nil {
		os.Exit(1)
	}

	psLines, err := Ps([]string{"mysqld"}, []string{"mysqld_safe"})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ps error %v\n", err)
		os.Exit(1)
	}

	if len(psLines) == 0 {
		fmt.Fprintf(os.Stderr, "Ps result is empty\n")
		os.Exit(1)
	}

	fmt.Println(strings.Join(psLines, "\n"))

	pid, cmdName, err := NetstatListen(s.Port)
	if err != nil {
		fmt.Fprintf(os.Stderr, "NetstatListen error %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("netstat found cmd %s with pid %d\n", cmdName, pid)

	if !strings.HasPrefix(cmdName, "mysqld") {
		fmt.Printf("cmd %s is not msyqld\n", cmdName)
		os.Exit(1)
	}

	db := s.MustOpenDB()
	defer db.Close()

	result := sqlmore.ExecSQL(db, s.CheckSQL, 100, "NULL")
	if err := PrintSQLResult(os.Stdout, os.Stderr, s.CheckSQL, result); err != nil {
		fmt.Printf("PrintSQLResult error %v\n", err)
		os.Exit(1)
	}
}
