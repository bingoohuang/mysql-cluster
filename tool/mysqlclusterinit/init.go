package mysqlclusterinit

import (
	"strings"

	"github.com/sirupsen/logrus"

	// support mysql
	_ "github.com/go-sql-driver/mysql"
)

// InitMySQLCluster 初始化MySQL Master-Master集群
func (s Settings) InitMySQLCluster() (r Result) {
	if r.Sqls, r.Error = s.createMySQCluster(); r.Error != nil {
		return r
	}

	r.HAProxy = s.createHAProxyConfig()

	if s.Debug {
		logrus.Infof("SQL:%s", strings.Join(r.Sqls, ";\n"))
		logrus.Infof("HAProxy:%s", r.HAProxy)

		return r
	}

	if s.overwriteHAProxyCnf(&r); r.Error != nil {
		return r
	}

	if s.restartHAProxy(&r); r.Error != nil {
		return r
	}

	return r
}
