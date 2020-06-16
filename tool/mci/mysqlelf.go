package mci

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
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

// ShowTables show all tables.
func ShowTables(db *gorm.DB, excludedDbs ...string) (beans []TableBean, err error) {
	sql := `select * from information_schema.tables
		where TABLE_SCHEMA not in ('performance_schema', 'information_schema', 'mysql', 'sys'`

	if len(excludedDbs) > 0 {
		sql += `,'` + strings.Join(excludedDbs, "','") + `'`
	}

	sql += `)`

	if s := db.Raw(sql).Scan(&beans); s.Error != nil {
		logrus.Warnf("show slave status error: %v", s.Error)
		return beans, s.Error
	}

	return beans, nil
}

// RenameTables rename the non-system databases' table to another name.
// Returns the number of table renamed.
func RenameTables(db *gorm.DB) (int, error) {
	tables, err := ShowTables(db)
	if err != nil {
		return 0, err
	}

	renameSQL := createRenameSQL(tables)
	if renameSQL == "" {
		return 0, nil
	}

	logrus.Infof("sql:%s", renameSQL)

	if err := db.Exec(renameSQL).Error; err != nil {
		return 0, err
	}

	return len(tables), nil
}

func createRenameSQL(tables []TableBean) string {
	oldBackMap := map[string]bool{}
	newBackMap := map[string]TableBean{}
	reg := regexp.MustCompile(`.+_mci\d*`)

	for _, t := range tables {
		if reg.MatchString(t.Name) {
			oldBackMap[t.Schema+"."+t.Name] = true
		} else {
			newBackMap[t.Schema+"."+t.Name] = t
		}
	}

	if len(newBackMap) == 0 {
		logrus.Info("there is no tables to backup")
		return ""
	}

	needBackups := make(map[string]string)

	for k, t := range newBackMap {
		for i := 1; i < 9999999; i++ {
			k2 := fmt.Sprintf("%s.%s_mci%d", t.Schema, t.Name, i)
			if _, ok := oldBackMap[k2]; !ok {
				needBackups[k2] = k
				break
			}
		}
	}

	parts := make([]string, 0, len(needBackups))
	for nk, k := range needBackups {
		parts = append(parts, fmt.Sprintf("%s to %s", k, nk))
	}

	// https://dev.mysql.com/doc/refman/5.7/en/rename-table.html
	// RENAME TABLE
	//    tbl_name TO new_tbl_name
	//    [, tbl_name2 TO new_tbl_name2] ...
	return "rename table " + strings.Join(parts, ", ")
}

// ShowSlaveStatusBean 表示MySQL Slave Status.
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
	LastIOError          string `gorm:"column:Last_IO_Error"`
}

// ShowSlaveStatus show slave status to bean.
func ShowSlaveStatus(db *gorm.DB) (bean ShowSlaveStatusBean, err error) {
	if s := db.Raw("show slave status").Scan(&bean); s.Error != nil {
		logrus.Warnf("show slave status error: %v", s.Error)
		return bean, s.Error
	}

	return bean, nil
}

// ShowVariablesBean 表示MySQL 服务器参数结果行.
type ShowVariablesBean struct {
	VariableName string `gorm:"column:Variable_name" var:"field"`
	Value        string `gorm:"column:Value" var:"value"`
}

// Variables 表示MySQL 服务器参数.
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
	ServerUUID             string `var:"server_uuid"`
}

// ShowVariables shows variables to variables bean.
func ShowVariables(db *gorm.DB) (variables Variables, err error) {
	var beans []ShowVariablesBean

	if s := db.Raw("show variables").Scan(&beans); s.Error != nil {
		logrus.Warnf("show variables error: %v", s.Error)
		return Variables{}, s.Error
	}

	err = FlattenBeans(beans, &variables, "var")

	return variables, err
}
