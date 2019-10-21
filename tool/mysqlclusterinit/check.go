package mysqlclusterinit

import (
	"fmt"
	"os"
)

// CheckMySQL 检查MySQL集群配置
func (s Settings) CheckMySQL() {
	db := s.MustOpenGormDB()
	defer db.Close()

	if s.ValidateAndSetDefault() != nil {
		os.Exit(1)
	}

	if status, err := ShowSlaveStatus(db); err == nil {
		fmt.Printf("ShowSlaveStatus:%s\n", PrettyJSONSilent(status))
	}

	if variables, err := ShowVariables(db); err == nil {
		fmt.Printf("Variables:%s\n", PrettyJSONSilent(variables))
	}
}
