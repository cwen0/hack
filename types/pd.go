package types

import "time"

type StoreState int32

// Case insensitive key/value for replica constraints.
type StoreLabel struct {
	Key   string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value string `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
}

type Store struct {
	Id      uint64        `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Address string        `protobuf:"bytes,2,opt,name=address,proto3" json:"address,omitempty"`
	State   StoreState    `protobuf:"varint,3,opt,name=state,proto3,enum=metapb.StoreState" json:"state,omitempty"`
	Labels  []*StoreLabel `protobuf:"bytes,4,rep,name=labels" json:"labels,omitempty"`
}

// MetaStore is TiKV store status defined in protobuf
type MetaStore struct {
	*Store
	StateName string `json:"state_name"`
}

// ByteSize is a retype uint64 for TOML and JSON.
type ByteSize uint64
// Duration is a wrapper of time.Duration for TOML and JSON.
type Duration struct {
	time.Duration
}

// StoreStatus is TiKV store status returned from PD RESTful interface
type StoreStatus struct {
	Capacity           ByteSize `json:"capacity"`
	Available          ByteSize `json:"available"`
	LeaderCount        int               `json:"leader_count"`
	RegionCount        int               `json:"region_count"`
	SendingSnapCount   uint32            `json:"sending_snap_count"`
	ReceivingSnapCount uint32            `json:"receiving_snap_count"`
	ApplyingSnapCount  uint32            `json:"applying_snap_count"`
	IsBusy             bool              `json:"is_busy"`

	StartTS         time.Time         `json:"start_ts"`
	LastHeartbeatTS time.Time         `json:"last_heartbeat_ts"`
	Uptime          Duration `json:"uptime"`
}

// StoreInfo is a single store info returned from PD RESTful interface
type StoreInfo struct {
	Store  *MetaStore   `json:"store"`
	Status *StoreStatus `json:"status"`
}

// StoresInfo is stores info returned from PD RESTful interface
type StoresInfo struct {
	Count  int          `json:"count"`
	Stores []*StoreInfo `json:"stores"`
}

type SchedulerInfo struct {
	Name    string `json:"name"`
	StoreID uint64 `json:"store_id"`
}
