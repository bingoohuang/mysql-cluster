package mci_test

import (
	"testing"

	"github.com/bingoohuang/tool/mci"
	"github.com/stretchr/testify/assert"
)

func TestContainsSub(t *testing.T) {
	assert.True(t, mci.ContainsSub("abc", "a"))
	assert.False(t, mci.ContainsSub("abc", "d", "e"))
}
