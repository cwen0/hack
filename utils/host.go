package utils

var (
	defaultHostMap = map[string]string{
		"a":"1234",
	}
)
// Resolve resolves ip. now just using map
func Resolve(hostname string) (string, bool) {
	ip, exist := defaultHostMap[hostname]
	return ip, exist
}