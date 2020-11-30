package mci_test

import (
	"testing"

	"github.com/bingoohuang/tool/mci"
	"github.com/elliotchance/testify-stats/assert"
)

func TestRename(t *testing.T) {
	renameSQL := mci.CreateRenameSQL([]mci.TableBean{
		{
			Schema: "a-b",
			Name:   "c-d",
		},
	}, false)

	assert.Equal(t, "rename table `a-b`.`c-d` to `a-b`.`c-d_mci1`", renameSQL)
}
