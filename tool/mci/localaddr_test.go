package mci_test

import (
	"testing"

	"github.com/bingoohuang/gonet"
	"github.com/bingoohuang/tool/mci"
	"github.com/stretchr/testify/assert"
)

func TestIsLocalAddr(t *testing.T) {
	assert.True(t, mci.IsLocalAddr("127.0.0.1"))
	assert.True(t, mci.IsLocalAddr("localhost"))
	assert.False(t, mci.IsLocalAddr(""))

	for _, ip := range gonet.ListIps() {
		assert.True(t, mci.IsLocalAddr(ip))
	}
}
