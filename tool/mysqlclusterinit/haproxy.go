package mysqlclusterinit

import (
	"errors"
	"fmt"
	"time"

	"github.com/gobars/cmd"
	"github.com/sirupsen/logrus"
)

func (s Settings) createHAProxyConfig() string {
	rwConfig := fmt.Sprintf(`
listen mysql-rw
  bind 127.0.0.1:13306
  mode tcp
  option tcpka
  server mysql-1 %s:%d check inter 1s
  server mysql-2 %s:%d check inter 1s backup
`, s.Master1Addr, s.Port, s.Master2Addr, s.Port)
	rConfig := fmt.Sprintf(`
listen mysql-ro
  bind 127.0.0.1:23306
  mode tcp
  option tcpka
  server mysql-1 %s:%d check inter 1s
  server mysql-2 %s:%d check inter 1s
`, s.Master1Addr, s.Port, s.Master2Addr, s.Port)

	for seq, slaveIP := range s.SlaveAddrs {
		if slaveIP != "" {
			rConfig += fmt.Sprintf("  server mysql-%d %s:%d check inter 1s\n", seq+3, slaveIP, s.Port)
		}
	}

	return rwConfig + rConfig
}

func (s Settings) restartHAProxy(r *Result) {
	if s.HAProxyRestartShell == "" {
		logrus.Warnf("HAProxyRestartShell is empty")
		return
	}

	_, status := cmd.Bash(s.HAProxyRestartShell, cmd.Timeout(5*time.Second), cmd.Buffered(false))
	if status.Error != nil {
		r.Error = status.Error
		return
	}

	if status.Exit != 0 {
		r.Error = fmt.Errorf("error exiting code %d, stdout:%s, stderr:%s",
			status.Exit, status.Stdout, status.Stderr)
		return
	}
}

func (s Settings) overwriteHAProxyCnf(r *Result) {
	if s.HAProxyCfg == "" {
		r.Error = errors.New("HAProxyCfg required")
		return
	}

	if r.Error = FileExists(s.HAProxyCfg); r.Error != nil {
		return
	}

	logrus.Infof("prepare to overwriteHAProxyCnf %s", r.HAProxy)

	if r.Error = ReplaceFileContent(s.HAProxyCfg,
		`(?is)#\s*MySQLClusterConfigStart(.+)#\s*MySQLClusterConfigEnd`, r.HAProxy); r.Error == nil {
		logrus.Infof("overwriteHAProxyCnf completed")
	} else {
		logrus.Warnf("overwriteHAProxyCnf error: %v", r.Error)
	}
}
