package mysqlclusterinit_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/bingoohuang/tool/mysqlclusterinit"
	"github.com/stretchr/testify/assert"
)

func TestReplaceFileContent(t *testing.T) {
	file, err := ioutil.TempFile("/tmp", "myname.*.conf")
	assert.Nil(t, err)

	fmt.Println("tmp file created ", file.Name())
	assert.Nil(t, ioutil.WriteFile(file.Name(), []byte(
		"a=b\nserver_id=0\nc=d\nserver_id=3"), 0644))
	assert.Nil(t, mysqlclusterinit.ReplaceFileContent(file.Name(),
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
	s, err := mysqlclusterinit.ReplaceContent("START\na\nb\nEND", `(?is)START(.+)END`, "\nHAHA\n")
	assert.Nil(t, err)
	assert.Equal(t, "START\nHAHA\nEND", s)
}
