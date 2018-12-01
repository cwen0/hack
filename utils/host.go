package utils

var (
	defaultHostMap = map[string]string{
		"tidb-cluster-tikv-0.tidb-cluster-tikv-peer.thix-cat-tidb-cluster.svc:20160":"10.128.31.56",
		"tidb-cluster-tikv-1.tidb-cluster-tikv-peer.thix-cat-tidb-cluster.svc:20160":"10.128.31.62",
		"tidb-cluster-tikv-2.tidb-cluster-tikv-peer.thix-cat-tidb-cluster.svc:20160":"10.128.31.51",
	}
)
// Resolve resolves ip. now just using map
func Resolve(hostname string) (string, bool) {
	ip, exist := defaultHostMap[hostname]
	return ip, exist
}