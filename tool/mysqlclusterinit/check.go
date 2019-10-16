package mysqlclusterinit

import (
	"fmt"
)

// CheckMySQL 检查MySQL集群配置
func (s Settings) CheckMySQL() {
	db := s.MustOpenGormDB()
	defer db.Close()

	if status, err := ShowSlaveStatus(db); err == nil {
		fmt.Printf("ShowSlaveStatus:%s\n", PrettyJSONSilent(status))
	}

	if variables, err := ShowVariables(db); err == nil {
		fmt.Printf("Variables:%s\n", PrettyJSONSilent(variables))
	}
}
