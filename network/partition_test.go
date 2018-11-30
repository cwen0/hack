package network

import (
	"testing"

	. "github.com/pingcap/check"
	"github.com/zhouqiang-cl/hack/types"
)

func TestGetPartitionConfig(t *testing.T) {
	TestingT(t)
}

var _ = Suite(&testPartitionConfig{})

type testPartitionConfig struct{}

func (t *testPartitionConfig) TestFullPartitionConfig(c *C) {
	node := &types.Topological{
		PD:   []string{"10.0.0.1"},
		TiDB: []string{"10.0.0.2"},
		TiKV: []string{"10.0.0.3", "10.0.0.4", "10.0.0.5"},
	}

	var groups1 []types.Group
	groups1 = append(groups1, types.Group{
		Hosts: []string{"10.0.0.3", "10.0.0.4"},
	})
	groups1 = append(groups1, types.Group{
		Hosts: []string{"10.0.0.5"},
	})

	cfg1 := fullPartition(node, groups1)
	expectCfg1 := make(map[string]*types.NetworkConfig)
	expectCfg1["10.0.0.3"] = &types.NetworkConfig{
		Ingress: []string{"10.0.0.3", "10.0.0.4", "10.0.0.1", "10.0.0.2"},
		Egress:  []string{"10.0.0.3", "10.0.0.4", "10.0.0.1", "10.0.0.2"},
	}
	expectCfg1["10.0.0.4"] = &types.NetworkConfig{
		Ingress: []string{"10.0.0.3", "10.0.0.4", "10.0.0.1", "10.0.0.2"},
		Egress:  []string{"10.0.0.3", "10.0.0.4", "10.0.0.1", "10.0.0.2"},
	}
	expectCfg1["10.0.0.5"] = &types.NetworkConfig{
		Ingress: []string{"10.0.0.5", "10.0.0.1", "10.0.0.2"},
		Egress:  []string{"10.0.0.5", "10.0.0.1", "10.0.0.2"},
	}
	c.Assert(cfg1, DeepEquals, expectCfg1)

	var groups2 []types.Group
	groups2 = append(groups2, types.Group{
		Hosts: []string{"10.0.0.1", "10.0.0.3", "10.0.0.4"},
	})
	groups2 = append(groups2, types.Group{
		Hosts: []string{"10.0.0.5"},
	})

	cfg2 := fullPartition(node, groups2)

	expectCfg2 := make(map[string]*types.NetworkConfig)
	expectCfg2["10.0.0.1"] = &types.NetworkConfig{
		Ingress: []string{"10.0.0.1", "10.0.0.3", "10.0.0.4", "10.0.0.2"},
		Egress:  []string{"10.0.0.1", "10.0.0.3", "10.0.0.4", "10.0.0.2"},
	}
	expectCfg2["10.0.0.3"] = &types.NetworkConfig{
		Ingress: []string{"10.0.0.1", "10.0.0.3", "10.0.0.4", "10.0.0.2"},
		Egress:  []string{"10.0.0.1", "10.0.0.3", "10.0.0.4", "10.0.0.2"},
	}
	expectCfg2["10.0.0.4"] = &types.NetworkConfig{
		Ingress: []string{"10.0.0.1", "10.0.0.3", "10.0.0.4", "10.0.0.2"},
		Egress:  []string{"10.0.0.1", "10.0.0.3", "10.0.0.4", "10.0.0.2"},
	}
	expectCfg2["10.0.0.5"] = &types.NetworkConfig{
		Ingress: []string{"10.0.0.5", "10.0.0.2"},
		Egress:  []string{"10.0.0.5", "10.0.0.2"},
	}
	c.Assert(cfg2, DeepEquals, expectCfg2)
}

func (t *testPartitionConfig) TestPartialPartitionConfig(c *C) {
	node := &types.Topological{
		PD:   []string{"10.0.0.1"},
		TiDB: []string{"10.0.0.2"},
		TiKV: []string{"10.0.0.3", "10.0.0.4", "10.0.0.5"},
	}

	group1 := types.Group{
		Hosts: []string{"10.0.0.3"},
	}
	group2 := types.Group{
		Hosts: []string{"10.0.0.5"},
	}

	cfg1 := partialPartion(node, group1, group2)
	expectCfg1 := make(map[string]*types.NetworkConfig)
	expectCfg1["10.0.0.3"] = &types.NetworkConfig{
		Ingress: []string{"10.0.0.3", "10.0.0.1", "10.0.0.2", "10.0.0.4"},
		Egress:  []string{"10.0.0.3", "10.0.0.1", "10.0.0.2", "10.0.0.4"},
	}
	expectCfg1["10.0.0.5"] = &types.NetworkConfig{
		Ingress: []string{"10.0.0.5", "10.0.0.1", "10.0.0.2", "10.0.0.4"},
		Egress:  []string{"10.0.0.5", "10.0.0.1", "10.0.0.2", "10.0.0.4"},
	}
	c.Assert(cfg1, DeepEquals, expectCfg1)
}

func (t *testPartitionConfig) TestSimplexPartitionConfig(c *C) {
	node := &types.Topological{
		PD:   []string{"10.0.0.1"},
		TiDB: []string{"10.0.0.2"},
		TiKV: []string{"10.0.0.3", "10.0.0.4", "10.0.0.5"},
	}

	group1 := types.Group{
		Hosts: []string{"10.0.0.3", "10.0.0.4"},
	}
	group2 := types.Group{
		Hosts: []string{"10.0.0.5"},
	}

	cfg1 := simplexPartition(node, group1, group2)
	expectCfg1 := make(map[string]*types.NetworkConfig)
	expectCfg1["10.0.0.3"] = &types.NetworkConfig{
		Ingress: []string{"10.0.0.3", "10.0.0.4", "10.0.0.1", "10.0.0.2"},
		Egress:  []string{"10.0.0.3", "10.0.0.4", "10.0.0.5", "10.0.0.1", "10.0.0.2"},
	}
	expectCfg1["10.0.0.4"] = &types.NetworkConfig{
		Ingress: []string{"10.0.0.3", "10.0.0.4", "10.0.0.1", "10.0.0.2"},
		Egress:  []string{"10.0.0.3", "10.0.0.4", "10.0.0.5", "10.0.0.1", "10.0.0.2"},
	}
	expectCfg1["10.0.0.5"] = &types.NetworkConfig{
		Ingress: []string{"10.0.0.5", "10.0.0.3", "10.0.0.4", "10.0.0.1", "10.0.0.2"},
		Egress:  []string{"10.0.0.5", "10.0.0.1", "10.0.0.2"},
	}
	c.Assert(cfg1, DeepEquals, expectCfg1)
}
