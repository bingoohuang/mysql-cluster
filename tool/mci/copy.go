package mci

import (
	"fmt"

	"github.com/bingoohuang/now"
	"github.com/gobars/cmd"
	"github.com/sirupsen/logrus"
)

func (s *Settings) copyMaster1Data(slaveServers []string) error {
	dumpTime := now.MakeNow().Format("yyyyMMddHHmmss")
	env := cmd.Env(`MYSQL_PWD=` + s.Password)
	dumpCmd := fmt.Sprintf(`%s -u root -P %d -h %s --all-databases --set-gtid-purged=OFF %s>%s_mm1.sql`,
		s.MySQLDumpCmd, s.Port, s.Master1Addr, s.MySQLDumpOptions, dumpTime)
	logrus.Infof("%s", dumpCmd)
	_, status := cmd.Bash(dumpCmd, env)

	if status.Error != nil {
		return fmt.Errorf("exec %s fail error %w", dumpCmd, status.Error)
	}

	if status.Exit != 0 {
		return fmt.Errorf("exec %s fail exiting code %d, stdout:%s, stderr:%s",
			dumpCmd, status.Exit, status.Stdout, status.Stderr)
	}

	for _, slaveServer := range slaveServers {
		importCmd := fmt.Sprintf(`%s -u root -P %d -h %s<%s_mm1.sql`, s.MySQLCmd, s.Port, slaveServer, dumpTime)
		logrus.Infof("%s", importCmd)
		_, status := cmd.Bash(importCmd, env)

		if status.Error != nil {
			return fmt.Errorf("exec %s fail error %w", importCmd, status.Error)
		}

		if status.Exit != 0 {
			return fmt.Errorf("exec %s fail exiting code %d, stdout:%s, stderr:%s",
				importCmd, status.Exit, status.Stdout, status.Stderr)
		}
	}

	return nil
}
