package mci

import (
	"errors"
	"fmt"

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
`, s.replaceIPToLocalhost(s.Master1Addr), s.Port, s.replaceIPToLocalhost(s.Master2Addr), s.Port)

	replaceIP1, originalIP1 := s.replaceIPToLocalhost2(s.Master1Addr)
	replaceIP2, originalIP2 := s.replaceIPToLocalhost2(s.Master2Addr)
	rConfig := fmt.Sprintf(`
listen mysql-ro
  bind 127.0.0.1:23306
  mode tcp
  option tcpka
  server mysql-1 %s:%d check inter 1s # %s:%d
  server mysql-2 %s:%d check inter 1s # %s:%d
`, replaceIP1, s.Port, originalIP1, s.Port, replaceIP2, s.Port, originalIP2, s.Port)

	for seq, slaveIP := range s.SlaveAddrs {
		if slaveIP != "" {
			replaceIP, originalIP := s.replaceIPToLocalhost2(slaveIP)
			rConfig += fmt.Sprintf("  server mysql-%d %s:%d check inter 1s # %s:%d\n",
				seq+3, replaceIP, s.Port, originalIP, s.Port)
		}
	}

	return rwConfig + rConfig
}

const localhost = "127.0.0.1"

func (s Settings) replaceIPToLocalhost(ip string) string {
	if s.isLocalAddr(ip) {
		return localhost
	}

	return ip
}

func (s Settings) replaceIPToLocalhost2(ip string) (local, ip2 string) {
	if s.isLocalAddr(ip) {
		return localhost, ip
	}

	return ip, ip
}

func (s Settings) overwriteHAProxyCnf(r *Result) error {
	if s.HAProxyCfg == "" {
		return errors.New("HAProxyCfg required")
	}

	if err := FileExists(s.HAProxyCfg); err != nil {
		return err
	}

	logrus.Infof("prepare to overwriteHAProxyCnf %s", r.HAProxy)

	if err := ReplaceFileContent(s.HAProxyCfg,
		`(?is)#\s*MySQLClusterConfigStart(.+)#\s*MySQLClusterConfigEnd`, r.HAProxy); err != nil {
		logrus.Warnf("overwriteHAProxyCnf error: %v", err)
		return err
	}

	logrus.Infof("overwriteHAProxyCnf completed")

	return nil
}

func (s Settings) resetHAProxyCnf() error {
	if err := FileExists(s.HAProxyCfg); err != nil {
		return err
	}

	logrus.Infof("prepare to resetHAProxyCnf")

	cnf := fmt.Sprintf(`
listen mysql-rw
  bind 127.0.0.1:13306
  mode tcp
  option tcpka
  server mysql-1 127.0.0.1:%d check inter 1s
listen mysql-ro
  bind 127.0.0.1:23306
  mode tcp
  option tcpka
  server mysql-1 127.0.0.1:%d check inter 1s
`, s.Port, s.Port)

	if err := ReplaceFileContent(s.HAProxyCfg,
		`(?is)#\s*MySQLClusterConfigStart(.+)#\s*MySQLClusterConfigEnd`, cnf); err != nil {
		logrus.Warnf("resetHAProxyCnf error: %v", err)
		return err
	}

	logrus.Infof("resetHAProxyCnf completed")

	return nil
}
