package mci

import (
	"github.com/bingoohuang/gonet"
	"github.com/bingoohuang/goreflect"
)

// HostIP 根据 primaryIfaceName 确定的名字，返回主IP primaryIP，以及以空格分隔的本机IP列表 ipList
// PrimaryIfaceName 表示主网卡的名称，用于获取主IP(v4)，不设置时，从eth0(linux), en0(darwin)，或者第一个ip v4的地址
func HostIP(primaryIfaceNames ...string) (primaryIP string, ipList []string, err error) {
	ips, err := gonet.ListLocalIfaceAddrs()
	if err != nil {
		return
	}

	ipList = make([]string, 0)

	for _, addr := range ips {
		if goreflect.SliceContains(primaryIfaceNames, addr.IfaceName) {
			primaryIP = addr.IP
		}

		ipList = append(ipList, addr.IP)
	}

	if primaryIP == "" && len(ipList) > 0 {
		primaryIP = ipList[0]
	}

	return
}
