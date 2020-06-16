package mci_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/bingoohuang/tool/mci"
	"github.com/stretchr/testify/assert"
)

func TestReplaceFileContent(t *testing.T) {
	file, err := ioutil.TempFile("/tmp", "myname.*.conf")
	assert.Nil(t, err)

	fmt.Println("tmp file created ", file.Name())
	assert.Nil(t, ioutil.WriteFile(file.Name(), []byte(
		"a=b\nserver_id=0\nc=d\nserver_id=3"), 0600))
	assert.Nil(t, mci.ReplaceFileContent(file.Name(),
		`(?i)server[-_]id\s*=\s*(\d+)`, "123456"))

	bytes, err := ioutil.ReadFile(file.Name())
	assert.Nil(t, err)
	assert.Equal(t, "a=b\nserver_id=123456\nc=d\nserver_id=123456", string(bytes))

	assert.Nil(t, os.Remove(file.Name()))

	// https://golang.org/pkg/regexp/syntax/
	// (?flags)       set flags within current group; non-capturing
	//
	// Flag syntax is xyz (set) or -xyz (clear) or xy-z (set xy, clear z). The flags are:
	//
	// i              case-insensitive (default false)
	// m              multi-line mode: ^ and $ match begin/end line in addition to begin/end text (default false)
	// s              let . match \n (default false)
	// U              ungreedy: swap meaning of x* and x*?, x+ and x+?, etc (default false)
	s, err := mci.ReplaceRegexGroup1("START\na\nb\nEND", `(?is)START(.+)END`, "\nHAHA\n")
	assert.Nil(t, err)
	assert.Equal(t, "START\nHAHA\nEND", s)
}

func TestSearchPatternLines(t *testing.T) {
	s, err := mci.SearchPatternLines(`START
listen mysql-rw
  bind 127.0.0.1:13306
  mode tcp
  option tcpka
  server mysql-1 10.0.0.1:3306 check inter 1s
  server mysql-2 10.0.0.2:3306 check inter 1s backup

listen mysql-ro
  bind 127.0.0.1:23306
  mode tcp
  option tcpka
  server mysql-1 10.0.0.1:3306 check inter 1s
  server mysql-2 10.0.0.2:3306 check inter 1s
  server mysql-3 10.0.0.3:3306 check inter 1s
END`, `(?is)mysql-ro(.+)END`, `(?i)server\s+\S+\s(\d+(\.\d+){3})(:\d+)?`)
	assert.Nil(t, err)
	assert.Equal(t, []string{"10.0.0.1", "10.0.0.2", "10.0.0.3"}, s)
}
