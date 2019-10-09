package mysqlclusterinit

import (
	"fmt"

	"github.com/bingoohuang/sqlmore"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"github.com/tkrajina/go-reflector/reflector"
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

	ShowSlaveStatus(db)
	ShowVariables(db)
}

func ShowSlaveStatus(db *gorm.DB) {
	var showSlaveStatus ShowSlaveStatusBean
	scan := db.Raw("show slave status").Scan(&showSlaveStatus)
	if scan.Error != nil {
		logrus.Warnf("show slave status error: %v", scan.Error)
		return
	}
	fmt.Printf("ShowSlaveStatus:%+v\n", showSlaveStatus)
}

func ShowVariables(db *gorm.DB) {
	var variables []ShowVariablesBean
	scan := db.Raw("show variables").Scan(&variables)
	if scan.Error != nil {
		logrus.Warnf("show variables error: %v", scan.Error)
		return
	}
	vs := Variables{}
	obj := reflector.New(&vs)
	fields := obj.Fields()

	for _, v := range variables {
		for _, f := range fields {
			tag, _ := f.Tag("var")
			if v.VariableName == tag {
				if err := f.Set(v.Value); err != nil {
					logrus.Warnf("Set error: %v", err)
				}
				break
			}
		}
	}

	fmt.Printf("Variables:%+v\n", vs)
}
