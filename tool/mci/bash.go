package mci

import (
	"fmt"
	"time"

	"github.com/gobars/cmd"
	"github.com/sirupsen/logrus"
)

// ExecuteBash executes bash.
func ExecuteBash(name string, bash string, shellTimeout time.Duration) error {
	if bash == "" {
		logrus.Warnf("%s is empty", name)
		return nil
	}

	logrus.Infof("start execute %s %s", name, bash)

	start := time.Now()

	_, status := cmd.Bash(bash, cmd.Timeout(shellTimeout), cmd.Buffered(false))
	if status.Error != nil {
		logrus.Infof("start execute %s %s error %v", name, bash, status.Error)
		return fmt.Errorf("execute %s %s error %w", name, bash, status.Error)
	}

	if status.Exit != 0 {
		logrus.Infof("start execute %s %s exiting code %d, stdout:%s, stderr:%s",
			name, bash, status.Exit, status.Stdout, status.Stderr)

		// nolint:goerr113
		return fmt.Errorf("execute %s %s exiting code %d, stdout:%s, stderr:%s",
			name, bash, status.Exit, status.Stdout, status.Stderr)
	}

	logrus.Infof("completed execute %s %s cost %v", name, bash, time.Since(start))

	return nil
}
