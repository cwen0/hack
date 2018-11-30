package types

// Topological is topological
type Topological struct {
	PD   []string `json:"pd"`
	TiDB []string `json:"tidb"`
	TiKV []string `json:"tikv"`
}
