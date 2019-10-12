package mysqlclusterinit

import (
	"fmt"

	"github.com/bingoohuang/sqlmore"
	"github.com/sirupsen/logrus"
)

// CheckMySQL 检查MySQL集群配置
func (s Settings) CheckMySQL() {
	ds := fmt.Sprintf("root:%s@tcp(127.0.0.1:%d)/", s.RootPassword, s.Port)

	db, err := sqlmore.NewSQLMore("mysql", ds).GormOpen()
	if err != nil {
		logrus.Warnf("open db error: %v", err)
		return
	}

	defer db.Close()

	if status, err := ShowSlaveStatus(db); err == nil {
		fmt.Printf("ShowSlaveStatus:%s\n", PrettyJSONSilent(status))
	}

	if variables, err := ShowVariables(db); err == nil {
		fmt.Printf("Variables:%s\n", PrettyJSONSilent(variables))
	}
}
