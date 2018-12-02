package utils

var (
	defaultHostMap = map[string]string{
		"tidb-cluster-tikv-0.tidb-cluster-tikv-peer.hackday-tidb-cluster.svc:20160":"10.128.24.193",
		"tidb-cluster-tikv-1.tidb-cluster-tikv-peer.hackday-tidb-cluster.svc:20160":"10.128.24.226",
		"tidb-cluster-tikv-2.tidb-cluster-tikv-peer.hackday-tidb-cluster.svc:20160":"10.128.28.139",
		"tidb-cluster-tikv-3.tidb-cluster-tikv-peer.hackday-tidb-cluster.svc:20160":"10.128.20.140",
		"tidb-cluster-tikv-4.tidb-cluster-tikv-peer.hackday-tidb-cluster.svc:20160":"10.128.28.174",
	}
)
// Resolve resolves ip. now just using map
func Resolve(hostname string) (string, bool) {
	ip, exist := defaultHostMap[hostname]
	return ip, exist
}