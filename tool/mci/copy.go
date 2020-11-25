package mci

import (
	"fmt"
	"github.com/bingoohuang/now"
	"github.com/gobars/cmd"
	"github.com/sirupsen/logrus"
)

// nolint:goerr113
func (s Settings) copyMaster1Data(slaveServers []string) error {
	env := cmd.Env(`MYSQL_PWD=` + s.Password)
	dumpTime, err := s.syncMaster1(env)
	if err != nil {
		return err
	}

	if dumpTime == "" {
		return nil
	}

	for _, slaveServer := range slaveServers {
		importCmd := fmt.Sprintf(`%s -u root -P %d -h %s<%s_mm1.sql`, s.MySQLCmd, s.Port, slaveServer, dumpTime)
		logrus.Infof("%s", importCmd)
		_, status := cmd.Bash(importCmd, cmd.Timeout(s.shellTimeoutDuration), env)

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

func (s Settings) syncMaster1(env cmd.OptionFn) (string, error) {
	if !s.Backup {
		return "", nil
	}

	dumpTime := now.MakeNow().Format("yyyyMMddHHmmss")
	dumpCmd := fmt.Sprintf(`%s -u root -P %d -h %s --all-databases --set-gtid-purged=OFF %s>%s_mm1.sql`,
		s.MySQLDumpCmd, s.Port, s.Master1Addr, s.MySQLDumpOptions, dumpTime)
	logrus.Infof("%s", dumpCmd)
	_, status := cmd.Bash(dumpCmd, cmd.Timeout(s.shellTimeoutDuration), env)

	if status.Error != nil {
		return "", fmt.Errorf("exec %s fail error %w", dumpCmd, status.Error)
	}

	if status.Exit != 0 {
		// nolint:goerr113
		return "", fmt.Errorf("exec %s fail exiting code %d, stdout:%s, stderr:%s",
			dumpCmd, status.Exit, status.Stdout, status.Stderr)
	}

	return dumpTime, nil
}
