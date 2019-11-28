package mci

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/bingoohuang/gou/str"

	"github.com/bingoohuang/goreflect"

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

	yes, _ := TellIsLocalAddr(addr)
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

var primaryIP, _, _ = HostIP("eth0", "en0") // nolint

// ReplaceLocalAddr2MainIP replace a single local address to main iface ip
func ReplaceLocalAddr2MainIP(address string) string {
	host, _ := str.Split2(address, ":", true, true)
	if IsLocalAddr(host) {
		return strings.ReplaceAll(address, host, primaryIP)
	}

	return address
}

// ReplaceLocalAddr2MainIPAll replace local addresses slice to main iface ips slice
func ReplaceLocalAddr2MainIPAll(addresses []string) []string {
	for i, addr := range addresses {
		addresses[i] = ReplaceLocalAddr2MainIP(addr)
	}

	return addresses
}

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

// TellIsLocalAddr 判断addr（ip，域名等）是否指向本机
// 由于IP可能经由iptable指向，或者可能是域名，或者其它，不能直接与本机IP做对比
// 本方法构建一个临时的HTTP服务，然后使用指定的addr去连接改HTTP服务，如果能连接上，说明addr是指向本机的地址
func TellIsLocalAddr(addr string) (bool, error) {
	if addr == "127.0.0.1" {
		return true, nil
	}

	localIPMap := gonet.ListLocalIPMap()
	if _, ok := localIPMap[addr]; ok {
		return true, nil
	}

	port, err := gonet.FreePort()
	if err != nil {
		return false, err
	}

	radStr := gonet.RandString(512)
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprint(w, radStr)
	})

	server := &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: mux}
	exitChan := make(chan bool)

	go func() {
		_ = server.ListenAndServe()
		exitChan <- true
	}()

	resp, err := gonet.HTTPGet(fmt.Sprintf("http://%s:%d", addr, port))
	_ = server.Close()

	if err != nil {
		return false, err
	}

	select {
	case <-time.After(10 * time.Second):
	case <-exitChan:
	}

	return string(resp) == radStr, nil
}
