package mysqlclusterinit_test

import (
	"testing"

	"github.com/bingoohuang/tool/mysqlclusterinit"
	"github.com/stretchr/testify/assert"
)

const ha = `
listen mysql-rw
  bind 0.0.0.0:13306
  mode tcp
  option tcpka
  server mysql-1 10.0.0.1:3306 check inter 1s
  server mysql-2 10.0.0.2:3306 check inter 1s backup

listen mysql-ro
  bind 0.0.0.0:23306
  mode tcp
  option tcpka
  server mysql-1 10.0.0.1:3306 check inter 1s
  server mysql-2 10.0.0.2:3306 check inter 1s
  server mysql-3 10.0.0.3:3306 check inter 1s
`

func TestMaster1(t *testing.T) {
	settings := &mysqlclusterinit.Settings{
		Master1IP:    "10.0.0.1",
		Master2IP:    "10.0.0.2",
		RootPassword: "123456",
		Port:         3306,
		ReplUsr:      "repl",
		ReplPassword: "repl_pwd",
		SlaveIps:     []string{"10.0.0.3"},
		Debug:        true,
		LocalIP:      "10.0.0.1",
	}

	result := settings.InitMySQLCluster()
	assert.Nil(t, result.Error)
	assert.Equal(t, []string{
		"SET GLOBAL server_id=1",
		"CREATE USER 'repl'@'%'",
		"GRANT REPLICATION SLAVE ON *.* TO 'repl'@'%' IDENTIFIED BY 'repl_pwd'",
		"STOP SLAVE",
		"CHANGE MASTER TO master_host='10.0.0.2', master_port=3306, " +
			"master_user='repl', master_password='repl_pwd', master_auto_position = 1",
		"START SLAVE",
	}, result.Sqls)
	assert.Equal(t, ha, result.HAProxy)
}

func TestMaster2(t *testing.T) {
	settings := &mysqlclusterinit.Settings{
		Master1IP:    "10.0.0.1",
		Master2IP:    "10.0.0.2",
		RootPassword: "123456",
		Port:         3306,
		ReplUsr:      "repl",
		ReplPassword: "repl_pwd",
		SlaveIps:     []string{"10.0.0.3"},
		Debug:        true,
		LocalIP:      "10.0.0.2",
	}

	result := settings.InitMySQLCluster()
	assert.Nil(t, result.Error)
	assert.Equal(t, []string{
		"SET GLOBAL server_id=2",
		"CREATE USER 'repl'@'%'",
		"GRANT REPLICATION SLAVE ON *.* TO 'repl'@'%' IDENTIFIED BY 'repl_pwd'",
		"STOP SLAVE",
		"CHANGE MASTER TO master_host='10.0.0.1', master_port=3306, " +
			"master_user='repl', master_password='repl_pwd', master_auto_position = 1",
		"START SLAVE",
	}, result.Sqls)
	assert.Equal(t, ha, result.HAProxy)
}

func TestSlave(t *testing.T) {
	settings := &mysqlclusterinit.Settings{
		Master1IP:    "10.0.0.1",
		Master2IP:    "10.0.0.2",
		RootPassword: "123456",
		Port:         3306,
		ReplUsr:      "repl",
		ReplPassword: "repl_pwd",
		SlaveIps:     []string{"10.0.0.3"},
		Debug:        true,
		LocalIP:      "10.0.0.3",
	}

	result := settings.InitMySQLCluster()
	assert.Nil(t, result.Error)
	assert.Equal(t, []string{
		"SET GLOBAL server_id=3",
		"STOP SLAVE",
		"CHANGE MASTER TO master_host='10.0.0.2', master_port=3306, " +
			"master_user='repl', master_password='repl_pwd', master_auto_position = 1",
		"START SLAVE",
	}, result.Sqls)
}

func TestNone(t *testing.T) {
	settings := &mysqlclusterinit.Settings{
		Master1IP:    "10.0.0.1",
		Master2IP:    "10.0.0.2",
		RootPassword: "123456",
		Port:         3306,
		ReplUsr:      "repl",
		ReplPassword: "repl_pwd",
		SlaveIps:     []string{"10.0.0.3"},
		Debug:        true,
		LocalIP:      "",
	}

	result := settings.InitMySQLCluster()
	assert.Nil(t, result.Error)
	assert.Equal(t, []string{}, result.Sqls)
	assert.Equal(t, ha, result.HAProxy)
}
