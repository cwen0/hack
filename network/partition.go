package network

import (
	"github.com/juju/errors"
	"github.com/zhouqiang-cl/hack/types"
	"github.com/zhouqiang-cl/hack/utils"
)

// GetProxyPartitionConfig gets proxy
func GetProxyPartitionConfig(node *types.Topological, partition *types.Partition) (map[string]*types.NetworkConfig, error) {
	switch partition.Kind {
	case types.FullPartition:
		if len(partition.Groups) < 2 {
			return nil, errors.NotValidf("%v", partition)
		}

		return fullPartition(node, partition.Groups), nil
	case types.PartialPartition:
		if len(partition.Groups) != 2 {
			return nil, errors.NotValidf("%v", partition)
		}

		return partialPartion(node, partition.Groups[0], partition.Groups[1]), nil
	case types.SimplexPartition:
		if len(partition.Groups) != 2 {
			return nil, errors.NotValidf("%v", partition)
		}

		return simplexPartition(node, partition.Groups[0], partition.Groups[1]), nil
	default:
		return nil, errors.NotSupportedf("partition kind %s", partition.Kind)
	}
	return nil, nil
}

func fullPartition(node *types.Topological, groups []types.Group) map[string]*types.NetworkConfig {
	cfgs := make(map[string]*types.NetworkConfig)
	groupsHosts := groups[0].Hosts
	for index, group := range groups {
		if index == 0 {
			continue
		}
		groupsHosts = append(groupsHosts, group.Hosts...)
	}

	for _, group := range groups {
		for _, host := range group.Hosts {
			cfg := &types.NetworkConfig{
				Ingress: group.Hosts,
				Egress:  group.Hosts,
			}

			addOtherHosts(node.PD, groupsHosts, cfg)
			addOtherHosts(node.TiDB, groupsHosts, cfg)
			addOtherHosts(node.TiKV, groupsHosts, cfg)

			cfgs[host] = cfg
		}
	}

	return cfgs
}

func partialPartion(node *types.Topological, group1, group2 types.Group) map[string]*types.NetworkConfig {
	cfgs := make(map[string]*types.NetworkConfig)
	groupsHosts := group1.Hosts
	groupsHosts = append(groupsHosts, group2.Hosts...)

	for _, host := range group1.Hosts {
		cfg := &types.NetworkConfig{
			Ingress: group1.Hosts,
			Egress:  group1.Hosts,
		}

		addOtherHosts(node.PD, groupsHosts, cfg)
		addOtherHosts(node.TiDB, groupsHosts, cfg)
		addOtherHosts(node.TiKV, groupsHosts, cfg)

		cfgs[host] = cfg
	}

	for _, host := range group2.Hosts {
		cfg := &types.NetworkConfig{
			Ingress: group2.Hosts,
			Egress:  group2.Hosts,
		}

		addOtherHosts(node.PD, groupsHosts, cfg)
		addOtherHosts(node.TiDB, groupsHosts, cfg)
		addOtherHosts(node.TiKV, groupsHosts, cfg)

		cfgs[host] = cfg
	}

	return cfgs
}

func simplexPartition(node *types.Topological, group1, group2 types.Group) map[string]*types.NetworkConfig {
	cfgs := make(map[string]*types.NetworkConfig)
	groupsHosts := group1.Hosts
	groupsHosts = append(groupsHosts, group2.Hosts...)

	for _, host := range group1.Hosts {
		cfg := &types.NetworkConfig{
			Ingress: group1.Hosts,
			Egress:  group1.Hosts,
		}
		cfg.Egress = append(cfg.Egress, group2.Hosts...)
		addOtherHosts(node.PD, groupsHosts, cfg)
		addOtherHosts(node.TiDB, groupsHosts, cfg)
		addOtherHosts(node.TiKV, groupsHosts, cfg)

		cfgs[host] = cfg
	}

	for _, host := range group2.Hosts {
		cfg := &types.NetworkConfig{
			Ingress: group2.Hosts,
			Egress:  group2.Hosts,
		}
		cfg.Ingress = append(cfg.Ingress, group1.Hosts...)
		addOtherHosts(node.PD, groupsHosts, cfg)
		addOtherHosts(node.TiDB, groupsHosts, cfg)
		addOtherHosts(node.TiKV, groupsHosts, cfg)

		cfgs[host] = cfg
	}

	return cfgs
}

func addOtherHosts(hosts, groups []string, cfg *types.NetworkConfig) {
	for _, host := range hosts {
		if !utils.MatchInArray(groups, host) {
			cfg.Ingress = append(cfg.Ingress, host)
			cfg.Egress = append(cfg.Egress, host)
		}
	}
}
