package network

import (
	"github.com/juju/errors"
	"github.com/zhouqiang-cl/hack/types"
	"github.com/zhouqiang-cl/hack/utils"
	"math/rand"
)

// GetProxyPartitionConfig gets proxy
func GetProxyPartitionConfig(node *types.Topological, partition *types.Partition) (map[string]*types.NetworkConfig, error) {
	switch partition.Kind {
	case types.FullPartition:
		if len(partition.Groups) == 0 {
			partition.Groups = randomFullPartitionGroups(node)
		}

		if len(partition.Groups) < 2 {
			return nil, errors.NotValidf("%v", partition)
		}

		configs, realGroups := fullPartition(node, partition.Groups)
		partition.RealGroups = realGroups
		return configs, nil
	case types.PartialPartition:
		if len(partition.Groups) == 0 {
			partition.Groups = randomPartialPartition(node)
		}

		if len(partition.Groups) != 2 {
			return nil, errors.NotValidf("%v", partition)
		}

		configs, realGroups := partialPartition(node, partition.Groups[0], partition.Groups[1])
		partition.RealGroups = realGroups
		return configs, nil
	case types.SimplexPartition:
		if len(partition.Groups) == 0 {
			partition.Groups = randomSimplexPartition(node)
		}

		if len(partition.Groups) != 2 {
			return nil, errors.NotValidf("%v", partition)
		}

		configs, realGroups := simplexPartition(node, partition.Groups[0], partition.Groups[1])
		partition.RealGroups = realGroups
		return configs, nil
	default:
		return nil, errors.NotSupportedf("partition kind %s", partition.Kind)
	}
	return nil, nil
}

func fullPartition(node *types.Topological, groups []types.Group) (map[string]*types.NetworkConfig, []types.Group) {
	cfgs := make(map[string]*types.NetworkConfig)
	groupsHosts := groups[0].Hosts
	for index, group := range groups {
		if index == 0 {
			continue
		}
		groupsHosts = append(groupsHosts, group.Hosts...)
	}
	var readGroups []types.Group

	for _, group := range groups {
		// readGroup := types.Group{}
		for _, host := range group.Hosts {
			cfg := &types.NetworkConfig{
				Ingress: group.Hosts,
				Egress:  group.Hosts,
			}

			addOtherHosts(node.PD, groupsHosts, cfg)
			addOtherHosts(node.TiDB, groupsHosts, cfg)
			addOtherHosts(node.TiKV, groupsHosts, cfg)
			//if index == 0 {
			//	readGroup.Hosts = cfg.Ingress
			//}

			cfgs[host] = cfg
		}

		readGroups = append(readGroups, group)
	}

	return cfgs, readGroups
}

func partialPartition(node *types.Topological, group1, group2 types.Group) (map[string]*types.NetworkConfig, []types.Group) {
	cfgs := make(map[string]*types.NetworkConfig)
	groupsHosts := group1.Hosts
	groupsHosts = append(groupsHosts, group2.Hosts...)
	var realGroups []types.Group

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
	realGroups = append(realGroups, group1)

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
	realGroups = append(realGroups, group2)

	var group3Host []string
	for _, h := range node.TiKV {
		if !utils.MatchInArray(groupsHosts, h) {
			group3Host = append(group3Host, h)
		}
	}

	if len(group3Host) > 0 {
		realGroups = append(realGroups, types.Group{Hosts: group3Host})
	}

	return cfgs, realGroups
}

func simplexPartition(node *types.Topological, group1, group2 types.Group) (map[string]*types.NetworkConfig, []types.Group) {
	cfgs := make(map[string]*types.NetworkConfig)
	groupsHosts := group1.Hosts
	groupsHosts = append(groupsHosts, group2.Hosts...)
	var realGroups []types.Group

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
	realGroups = append(realGroups, group1)

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
	realGroups = append(realGroups, group2)

	var group3Host []string
	for _, h := range node.TiKV {
		if !utils.MatchInArray(groupsHosts, h) {
			group3Host = append(group3Host, h)
		}
	}

	if len(group3Host) > 0 {
		realGroups = append(realGroups, types.Group{Hosts: group3Host})
	}

	return cfgs, realGroups
}

func addOtherHosts(hosts, groups []string, cfg *types.NetworkConfig) {
	for _, host := range hosts {
		if !utils.MatchInArray(groups, host) {
			cfg.Ingress = append(cfg.Ingress, host)
			cfg.Egress = append(cfg.Egress, host)
		}
	}
}

func randomFullPartitionGroups(nodes *types.Topological) []types.Group {
	var groups []types.Group
	randC := rand.Intn(len(nodes.TiKV))

	group1 := types.Group{}
	for i := 0; i < randC; i++ {
		group1.Hosts = append(group1.Hosts, nodes.TiKV[i])
	}
	groups = append(groups, group1)

	group2 := types.Group{}
	for i := randC; i < len(nodes.TiKV); i++ {
		group2.Hosts = append(group2.Hosts, nodes.TiKV[i])
	}

	groups = append(groups, group2)
	return groups
}

func randomPartialPartition(nodes *types.Topological) []types.Group {
	var groups []types.Group

	if len(nodes.TiKV) < 3 {
		return groups
	}

	group1 := types.Group{}
	group1.Hosts = append(group1.Hosts, nodes.TiKV[0])
	groups = append(groups, group1)

	group2 := types.Group{}
	group2.Hosts = append(group2.Hosts, nodes.TiKV[1])

	groups = append(groups, group2)
	return groups
}

func randomSimplexPartition(nodes *types.Topological) []types.Group {
	var groups []types.Group

	if len(nodes.TiKV) < 3 {
		return groups
	}

	group1 := types.Group{}
	group1.Hosts = append(group1.Hosts, nodes.TiKV[0])
	groups = append(groups, group1)

	group2 := types.Group{}
	group2.Hosts = append(group2.Hosts, nodes.TiKV[1])

	groups = append(groups, group2)
	return groups
}
