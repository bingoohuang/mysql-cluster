package mci

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/bingoohuang/sqlmore"
	"github.com/gobars/cmd"
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

// CheckMySQLCluster 检查MySQL集群配置
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

	_, status := cmd.Bash(fmt.Sprintf(`netstat -tunlp | grep ":%d"`, s.Port), cmd.Timeout(1*time.Second))
	if status.Error != nil {
		fmt.Printf("netstat error %v\n", status.Error)
		os.Exit(1)
	}

	if len(status.Stdout) == 0 {
		fmt.Printf("netstat result empty\n")
		os.Exit(1)
	}

	// [root@BJCA-device ~]# netstat -tunlp | grep ":3306"
	// tcp6       0      0 :::3306                 :::*                    LISTEN      28132/mysqld
	re := regexp.MustCompile(`(?is)LISTEN\s+(\d+)/(\w+)`)
	stdout := strings.Join(status.Stdout, "\n")
	fmt.Printf("netstat stdout %s\n", stdout)

	subs := re.FindAllStringSubmatch(stdout, -1)
	if len(subs) != 1 {
		fmt.Printf("netstat too many results\n")
		os.Exit(1)
	}

	fmt.Printf("netstat found cmd %s with listen port %s\n", subs[0][2], subs[0][1])

	db := s.MustOpenDB()
	defer db.Close()

	result := sqlmore.ExecSQL(db, s.CheckSQL, 100, "NULL")
	if err := PrintSQLResult(os.Stdout, os.Stderr, s.CheckSQL, result); err != nil {
		fmt.Printf("PrintSQLResult error %v\n", err)
		os.Exit(1)
	}
}
