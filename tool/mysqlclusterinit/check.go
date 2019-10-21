package mysqlclusterinit

import (
	"fmt"
	"os"
	"strings"
)

// CheckMySQL 检查MySQL集群配置
func (s Settings) CheckMySQL() {
	db := s.MustOpenGormDB()
	defer db.Close()

	if s.ValidateAndSetDefault(SetDefault) != nil {
		os.Exit(1)
	}

	if status, err := ShowSlaveStatus(db); err == nil {
		fmt.Printf("ShowSlaveStatus:%s\n", PrettyJSONSilent(status))
	}

	if variables, err := ShowVariables(db); err == nil {
		fmt.Printf("Variables:%s\n", PrettyJSONSilent(variables))
	}
}

// CheckMySQL 检查MySQL集群配置
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
