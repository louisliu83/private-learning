package prom

var (
	proxyBytesMap = map[string]int64{}
)

func AddProxyDataBytes(host string, n int64) {
	if m, ok := proxyBytesMap[host]; ok {
		proxyBytesMap[host] = m + n
	} else {
		proxyBytesMap[host] = n
	}
}
