package mci

import (
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"github.com/tkrajina/go-reflector/reflector"
)

// TableBean means the table information in MySQL.
type TableBean struct {
	Schema       string `gorm:"column:TABLE_SCHEMA"`
	Name         string `gorm:"column:TABLE_NAME"`
	TableRows    int    `gorm:"column:TABLE_ROWS"`
	CreateTime   string `gorm:"column:CREATE_TIME"`
	UpdateTime   string `gorm:"column:UPDATE_TIME"`
	TableComment string `gorm:"column:TABLE_COMMENT"`
}

// ShowTables show all tables
func ShowTables(db *gorm.DB, postfix string, excludedDbs ...string) (beans []TableBean, err error) {
	sql := `select * from information_schema.tables
		where TABLE_SCHEMA not in ('performance_schema', 'information_schema', 'mysql', 'sys'`

	if len(excludedDbs) > 0 {
		sql += `,'` + strings.Join(excludedDbs, "','") + `'`
	}

	sql += `) and TABLE_NAME not like '%` + postfix + `'`

	if s := db.Raw(sql).Scan(&beans); s.Error != nil {
		logrus.Warnf("show slave status error: %v", s.Error)
		return beans, s.Error
	}

	return beans, nil
}

// RenameTables rename the non-system databases' table to another name.
// Returns the number of table renamed.
func RenameTables(db *gorm.DB, postfix string) (int, error) {
	tables, err := ShowTables(db, postfix)
	if err != nil {
		return 0, err
	}

	if len(tables) == 0 {
		logrus.Info("there is no tables to backup")
		return 0, nil
	}

	renameSqls := make([]string, len(tables))
	for i, t := range tables {
		renameSqls[i] = fmt.Sprintf("%s.%s to %s.%s%s", t.Schema, t.Name, t.Schema, t.Name, postfix)
	}

	// https://dev.mysql.com/doc/refman/5.7/en/rename-table.html
	// RENAME TABLE
	//    tbl_name TO new_tbl_name
	//    [, tbl_name2 TO new_tbl_name2] ...
	joined := "rename table " + strings.Join(renameSqls, ", ")
	logrus.Infof("sql:%s", joined)

	if err := db.Exec(joined).Error; err != nil {
		return 0, err
	}

	//time.Sleep(1 * time.Second) // 确保经过1秒
	//
	//timeNow := now.MakeNow().Format("yyyy-MM-dd HH:mm:ss")
	//purgeSQL := fmt.Sprintf("PURGE BINARY LOGS BEFORE '%s'", timeNow)
	//logrus.Infof("purgeSQL:%s", purgeSQL)
	//
	//// https://dev.mysql.com/doc/refman/5.7/en/purge-binary-logs.html
	//if err := db.Exec(purgeSQL).Error; err != nil {
	//	return err
	//}

	return len(tables), nil
}

// ShowSlaveStatusBean 表示MySQL Slave Status
type ShowSlaveStatusBean struct {
	SlaveIOState         string `gorm:"column:Slave_IO_State"`
	MasterHost           string `gorm:"column:Master_Host"`
	MasterUser           string `gorm:"column:Master_User"`
	MasterPort           int    `gorm:"column:Master_Port"`
	SlaveSQLRunningState string `gorm:"column:Slave_SQL_Running_State"`
	AutoPosition         bool   `gorm:"column:Auto_Position"`
	SlaveIoRunning       string `gorm:"column:Slave_IO_Running"`
	SlaveSQLRunning      string `gorm:"column:Slave_SQL_Running"`
	MasterServerID       string `gorm:"column:Master_Server_Id"`
	SecondsBehindMaster  string `gorm:"column:Seconds_Behind_Master"`
}

// ShowSlaveStatus show slave status to bean
func ShowSlaveStatus(db *gorm.DB) (bean ShowSlaveStatusBean, err error) {
	if s := db.Raw("show slave status").Scan(&bean); s.Error != nil {
		logrus.Warnf("show slave status error: %v", s.Error)
		return bean, s.Error
	}

	return bean, nil
}

// ShowVariablesBean 表示MySQL 服务器参数结果行
type ShowVariablesBean struct {
	VariableName string `gorm:"column:Variable_name"`
	Value        string `gorm:"column:Value"`
}

// Variables 表示MySQL 服务器参数
type Variables struct {
	ServerID               string `var:"server_id"`
	LogBin                 string `var:"log_bin"`
	SQLLogBin              string `var:"sql_log_bin"`
	GtidMode               string `var:"gtid_mode"`
	GtidNext               string `var:"gtid_next"`
	SlaveSkipErrors        string `var:"slave_skip_errors"`
	BinlogFormat           string `var:"binlog_format"`
	MasterInfoRepository   string `var:"master_info_repository"`
	RelayLogInfoRepository string `var:"relay_log_info_repository"`
	InnodbVersion          string `var:"innodb_version"`
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
