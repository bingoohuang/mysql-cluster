package mci

import (
	"strings"
	"sync"

	"github.com/bingoohuang/goreflect"

	"github.com/bingoohuang/gonet"
)

var localAddrMap sync.Map // nolint

// IsLocalAddr tells if an addr is local or not.
func IsLocalAddr(addr string) bool {
	if addr == "" {
		return false
	}

	if yes, ok := localAddrMap.Load(addr); ok {
		return yes.(bool)
	}

	yes, _ := gonet.IsLocalAddr(addr)
	localAddrMap.Store(addr, yes)

	return yes
}

const (
	localhostIPv4 = "127.0.0.1"
	localhostIPv6 = "::1"
)

// ReplaceAddr2Local try to replace an local IP to localhost.
func ReplaceAddr2Local(ip string) (replaced, original string) {
	if IsLocalAddr(ip) {
		if gonet.IsIPv6(ip) {
			return localhostIPv6, ip
		}

		return localhostIPv4, ip
	}

	return ip, ip
}

// TryReplaceAddr2Local try to replace an local IP to localhost.
func TryReplaceAddr2Local(ip string) (replaced string) {
	replaced, _ = ReplaceAddr2Local(ip)

	return replaced
}

var primaryIP, _, _ = HostIP("eth0", "en0") // nolint

// ReplaceLocalAddr2MainIP replace a single local address to main iface ip.
func ReplaceLocalAddr2MainIP(address string) string {
	sepPos := strings.LastIndex(address, ":")

	host := address
	if sepPos > 0 {
		host = address[:sepPos]
	}

	if IsLocalAddr(host) {
		return strings.ReplaceAll(address, host, primaryIP)
	}

	return address
}

// ReplaceLocalAddr2MainIPAll replace local addresses slice to main iface ips slice.
func ReplaceLocalAddr2MainIPAll(addresses []string) []string {
	for i, addr := range addresses {
		addresses[i] = ReplaceLocalAddr2MainIP(addr)
	}

	return addresses
}

// HostIP 根据 primaryIfaceName 确定的名字，返回主IP primaryIP，以及以空格分隔的本机IP列表 ipList.
// PrimaryIfaceName 表示主网卡的名称，用于获取主IP(v4)，不设置时，从eth0(linux), en0(darwin)，或者第一个ip v4的地址.
func HostIP(primaryIfaceNames ...string) (primaryIP string, ipList []string, err error) {
	ips, err := gonet.ListIfaces()
	if err != nil {
		return
	}

	ipList = make([]string, 0)

	for _, addr := range ips {
		if goreflect.SliceContains(primaryIfaceNames, addr.IfaceName) {
			primaryIP = addr.IP.String()
		}

		ipList = append(ipList, addr.IP.String())
	}

	if primaryIP == "" && len(ipList) > 0 {
		primaryIP = ipList[0]
	}

	return
}
