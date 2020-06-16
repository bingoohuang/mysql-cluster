package mci_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/bingoohuang/tool/mci"

	"github.com/elliotchance/testify-stats/assert"
)

func TestReadMySQLServersFromHAProxyCfg(t *testing.T) {
	// Create our Temp File:  This will create a filename like /tmp/prefix-123456
	// We can use a pattern of "pre-*.txt" to get an extension like: /tmp/pre-123456.txt
	tmpFile, err := ioutil.TempFile(os.TempDir(), "haproxy-")
	assert.Nil(t, err)

	// Remember to clean up the file afterwards
	defer os.Remove(tmpFile.Name())

	fmt.Println("Created File: " + tmpFile.Name())

	// Example writing to the file
	text := []byte(`
# MySQLClusterConfigStart
listen mysql-ro
  bind 127.0.0.1:13306
  mode tcp
  option tcpka
  server mysql-1 127.0.0.1:9633 check inter 1s # fd15:4ba5:5a2b:1008::3:9633
  server mysql-2 fd15:4ba5:5a2b:1008::4:9633 check inter 1s # fd15:4ba5:5a2b:1008::4:9633
# MySQLClusterConfigEnd
`)
	_, err = tmpFile.Write(text)
	assert.Nil(t, err)

	// Close the file
	err = tmpFile.Close()
	assert.Nil(t, err)

	s := mci.Settings{HAProxyCfg: tmpFile.Name()}
	cfg, err := s.ReadMySQLServersFromHAProxyCfg()
	assert.Nil(t, err)

	fmt.Println(cfg)
}
