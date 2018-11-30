package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/juju/errors"
	"github.com/ngaut/log"
	"github.com/zhouqiang-cl/hack/types"
	"github.com/zhouqiang-cl/hack/util"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type failpointCtl struct {
	toplogic *types.Topological
}

func newFailpointCtl(toplogic *types.Topological) *failpointCtl {
	return &failpointCtl{toplogic: toplogic}
}

func (f *failpointCtl) start(typ string) error {
	switch typ {
	case "random":
		f.cleanFailpoint()
		return errors.Trace(f.doRandomFailpoint())
	case "certain":
		f.cleanFailpoint()
		return errors.Trace(f.doCertainFailpoint())
	}
	return errors.NotSupportedf("typ %s", typ)
}

func (f *failpointCtl) doRandomFailpoint() error {
	kvs := f.getRandomTiKVs()
	for _, kv := range kvs {
		pathes := f.getRandomFailpointPath()
		err := doFailpoints(kv, pathes)
		if err != nil {
			log.Errorf("do failpoint error %v", err)
		}
	}
	return nil
}

func (f *failpointCtl) cleanFailpoint() error {
	kvs := f.getRandomTiKVs()
	for _, kv := range kvs {
		err := emptyFailpoints(kv)
		if err != nil {
			log.Errorf("clean failpoint error %v", err)
		}
	}
	return nil
}

func (f *failpointCtl) doCertainFailpoint() error {
	// rule := "pct(5)->delay(100)|pct(1)->timeout()"
	for _, kv := range f.toplogic.TiKV {
		err := doFailpoints(kv, getAllPath())
		if err != nil {
			return errors.Trace(err)
		}
	}
	return nil
}

func (f *failpointCtl) getRandomTiKVs() []string {
	var (
		kvs []string
	)

	leftTiKVs := make([]string, len(f.toplogic.TiKV))
	count := rand.Intn(len(f.toplogic.TiKV))
	copy(leftTiKVs, f.toplogic.TiKV)
	for i := 0; i < count+1; i++ {
		index := rand.Intn(len(leftTiKVs))
		kvs = append(kvs, leftTiKVs[index])
		leftTiKVs = append(leftTiKVs[:index], leftTiKVs[index+1:]...)
	}
	return kvs
}

func (f *failpointCtl) getRandomFailpointPath() []string {
	var (
		pathes []string
	)
	allPathes := getAllPath()
	count := rand.Intn(len(allPathes))
	for i := 0; i < count+1; i++ {
		index := rand.Intn(len(allPathes))
		pathes = append(pathes, allPathes[index])
		allPathes = append(allPathes[:index], allPathes[index+1:]...)
	}
	return pathes
}

func doFailpoints(kv string, pathes []string) error {
	rule := "pct(5)->delay(100)|pct(1)->timeout()"
	for _, path := range pathes {
		err := doFailpoint(kv, path, rule)
		if err != nil {
			return errors.Trace(err)
		}
	}
	return nil
}

func doFailpoint(kv string, path string, rule string) error {
	cfg := &types.FailpointConfig{
		Path:  path,
		Value: rule,
	}
	data, err := json.Marshal(cfg)
	if err != nil {
		return errors.Trace(err)
	}
	url := fmt.Sprintf("http://%s:10008/config/failpoint/add", kv)
	_, err = util.DoPost(url, data)
	return errors.Trace(err)
}

func emptyFailpoints(kv string) error {
	url := fmt.Sprintf("http://%s:10008/config/failpoint/clean", kv)
	_, err := util.DoPost(url, []byte{})
	return errors.Trace(err)
}
