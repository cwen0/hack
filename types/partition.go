package types

// PartitionKind defines the kind of the network partition
type PartitionKind string

// // kinds of Partition
const (
	FullPartition    PartitionKind = "full"
	PartialPartition               = "partial"
	SimplexPartition               = "simplex"
)

// Partition defines the config for network partition
type Partition struct {
	Groups    []Group       `json:"groups"`
	Kind      PartitionKind `json:"kind"`
}

// Group define the a set of nodes for partition
type Group struct {
	Hosts []string `json:"hosts"`
}

