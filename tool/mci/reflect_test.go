package mci_test

import (
	"testing"

	"github.com/bingoohuang/tool/mci"
	"github.com/stretchr/testify/assert"
)

type myBean struct {
	VariableName string `var:"field"`
	Value        string `var:"value"`
}

type myVar struct {
	ServerID string `var:"server_id"`
	LogBin   string `var:"log_bin"`
}

func TestFlattenBeans(t *testing.T) {
	beans := []myBean{
		{VariableName: "server_id", Value: "100"},
		{VariableName: "log_bin", Value: "abc"},
	}

	var v myVar

	assert.Nil(t, mci.FlattenBeans(beans, &v, "var"))
	assert.Equal(t, myVar{ServerID: "100", LogBin: "abc"}, v)
}
