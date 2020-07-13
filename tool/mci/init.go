package mci

import (
	"fmt"
	"os"
	"strings"

	"github.com/bingoohuang/gou/str"

	"github.com/sirupsen/logrus"

	// support mysql
	_ "github.com/go-sql-driver/mysql"
)

// CreateMySQLCluster 初始化MySQL Master-Master集群.
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

	if err := s.restartHAProxy(); err != nil {
		return r, err
	}

	return r, err
}

// ResetLocalMySQLClusterNode 重置MySQL集群.
func (s Settings) ResetLocalMySQLClusterNode() error {
	if s.ValidateAndSetDefault(Validate, SetDefault) != nil {
		os.Exit(1)
	}

	if err := s.resetMySQCluster(); err != nil {
		return err
	}

	if err := s.resetHAProxyCnf(); err != nil {
		return err
	}

	return s.restartHAProxy()
}

const reMySQLClusterConfig = `(?is)#\s*MySQLClusterConfigStart(.+)#\s*MySQLClusterConfigEnd`

// RemoveSlavesFromCluster 从集群中，删除从节点
// nolint:goerr113
func (s Settings) RemoveSlavesFromCluster(removeSlaves string) error {
	slavesToRemove := str.SplitTrim(removeSlaves, ",")
	if len(slavesToRemove) == 0 {
		return nil
	}

	if s.ValidateAndSetDefault(SetDefault) != nil {
		os.Exit(1)
	}

	mySQLClusterConfig, err := SearchFileContent(s.HAProxyCfg, reMySQLClusterConfig)
	if err != nil {
		return fmt.Errorf("SearchFileContent error %w", err)
	}

	if len(mySQLClusterConfig) == 0 {
		return fmt.Errorf("RemoveSlavesFromCluster error : no MySQLClusterConfig found in %s", s.HAProxyCfg)
	}

	if len(mySQLClusterConfig) > 1 {
		return fmt.Errorf("RemoveSlavesFromCluster error : more than one MySQLClusterConfig found in %s", s.HAProxyCfg)
	}

	lines := strings.Split(mySQLClusterConfig[0], "\n")
	changes := 0

	for i, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "#") {
			continue
		}

		if ContainsSub(line, slavesToRemove...) {
			lines[i] = "# " + line
			changes++
		}
	}

	if changes == 0 {
		fmt.Println("nothing to do")
		return nil
	}

	newConfig := strings.Join(lines, "\n")

	if err := ReplaceFileContent(s.HAProxyCfg, reMySQLClusterConfig, newConfig); err != nil {
		return fmt.Errorf("ReplaceFileContent HAProxyCfg error %w", err)
	}

	return s.restartHAProxy()
}

func (s Settings) restartHAProxy() error {
	return ExecuteBash("HAProxyRestartShell", s.HAProxyRestartShell, s.shellTimeoutDuration)
}
