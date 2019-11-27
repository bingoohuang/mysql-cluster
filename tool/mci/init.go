package mci

import (
	"os"

	"github.com/sirupsen/logrus"

	// support mysql
	_ "github.com/go-sql-driver/mysql"
)

// CreateMySQLCluster 初始化MySQL Master-Master集群
func (s Settings) CreateMySQLCluster() (r Result, err error) {
	if s.ValidateAndSetDefault(Validate, SetDefault) != nil {
		os.Exit(1)
	}

	if r.Nodes, err = s.createMySQCluster(); err != nil {
		return r, err
	}

	r.HAProxy = s.createHAProxyConfig()

	if s.Debug {
		logrus.Infof("HAProxy:%s", r.HAProxy)
		return r, err
	}

	if err := s.overwriteHAProxyCnf(&r); err != nil {
		return r, err
	}

	if err := executeBash("HAProxyRestartShell", s.shellTimeoutDuration, s.HAProxyRestartShell); err != nil {
		return r, err
	}

	return r, err
}

// ResetMySQLCluster 重置MySQL集群
func (s Settings) ResetMySQLCluster() (err error) {
	if s.ValidateAndSetDefault(Validate, SetDefault) != nil {
		os.Exit(1)
	}

	if err = s.resetMySQCluster(); err != nil {
		return err
	}

	if err := s.resetHAProxyCnf(); err != nil {
		return err
	}

	return err
}
