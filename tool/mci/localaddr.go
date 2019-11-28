package mci

import (
	"sync"

	"github.com/bingoohuang/gonet"
)

var localAddrMap sync.Map // nolint

// IsLocalAddr tells if an addr is local or not.
func IsLocalAddr(addr string) bool {
	if addr == "" {
		return false
	}

	if addr == localhost || addr == "localhost" {
		return true
	}

	if yes, ok := localAddrMap.Load(addr); ok {
		return yes.(bool)
	}

	yes, _ := gonet.IsLocalAddr(addr)
	localAddrMap.Store(addr, yes)

	return yes
}

const localhost = "127.0.0.1"

// ReplaceAddr2Local try to replace an local IP to localhost
func ReplaceAddr2Local(ip string) (replaced, original string) {
	if IsLocalAddr(ip) {
		return localhost, ip
	}

	return ip, ip
}

// TryReplaceAddr2Local try to replace an local IP to localhost
func TryReplaceAddr2Local(ip string) (replaced string) {
	replaced, _ = ReplaceAddr2Local(ip)

	return replaced
}
