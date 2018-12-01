package main

import (
	"encoding/json"
	"fmt"
	"github.com/unrolled/render"
	"math/rand"
	"net/http"
	"time"


	"github.com/ngaut/log"
	"github.com/juju/errors"
	"github.com/zhouqiang-cl/hack/types"
	"github.com/zhouqiang-cl/hack/utils"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type failpointHandler struct {
	c  *Manager
	rd *render.Render
}

func newFailpointHandler(c *Manager, rd *render.Render) *failpointHandler {
	return &failpointHandler{
		c:  c,
		rd: rd,
	}
}

func (f *failpointHandler) CreateFailpoint(w http.ResponseWriter, r *http.Request) {
	fpType := r.URL.Query()["type"]
	if len(fpType) == 0 {
		f.rd.JSON(w, http.StatusBadRequest, "miss parameter ip")
		return
	}
	switch fpType[0] {
	case "random":
		log.Debugf("clean failpoint")
		cleanFailpoint(f.c.pdAddr)
		log.Debugf("clean failpoint finished")
		err := doRandomFailpoint(f.c.pdAddr)
		log.Debugf("do random failpoint")
		if err != nil {
			f.rd.JSON(w, http.StatusInternalServerError, err.Error())
			return
		}
	case "certain":
		cleanFailpoint(f.c.pdAddr)
		err := doCertainFailpoint(f.c.pdAddr)
		if err != nil {
			f.rd.JSON(w, http.StatusInternalServerError, err.Error())
			return
		}
	case "clean":
		err := cleanFailpoint(f.c.pdAddr)
		if err != nil {
			f.rd.JSON(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	state = State{
		operation: StateFailpoint,
		partition: types.Partition{},
	}

	logs.items = append(logs.items, Log{
		operation: OperationFailpoint,
		parameter: fpType[0],
		timeStamp: time.Now().Unix(),
	})

	f.rd.JSON(w, http.StatusOK, nil)
}

func doRandomFailpoint(pdAddr string) error {
	kvs, err := getRandomTiKVs(pdAddr)
	if err != nil {
		return err
	}
	for _, kv := range kvs {
		pathes := getRandomFailpointPath()
		err := doFailpoints(kv, pathes)
		if err != nil {
			return err
		}
	}
	return nil
}

func cleanFailpoint(pdAddr string) error {
	kvs, err := getRandomTiKVs(pdAddr)
	if err != nil {
		return err
	}
	for _, kv := range kvs {
		err := emptyFailpoints(kv)
		if err != nil {
			return err
		}
	}
	return nil
}

func doCertainFailpoint(pdAddr string) error {
	// rule := "pct(5)->delay(100)|pct(1)->timeout()"
	topology, err := getTopologyInfo(pdAddr)
	if err != nil {
		return errors.Trace(err)
	}
	for _, kv := range topology.TiKV {
		err := doFailpoints(kv, getAllPath())
		if err != nil {
			return errors.Trace(err)
		}
	}
	return nil
}

func getRandomTiKVs(pdAddr string) ([]string, error) {
	var (
		kvs []string
	)

	topology, err := getTopologyInfo(pdAddr)
	if err != nil {
		return nil, errors.Trace(err)
	}

	leftTiKVs := make([]string, len(topology.TiKV))
	count := rand.Intn(len(topology.TiKV))
	copy(leftTiKVs, topology.TiKV)
	for i := 0; i < count+1; i++ {
		index := rand.Intn(len(leftTiKVs))
		kvs = append(kvs, leftTiKVs[index])
		leftTiKVs = append(leftTiKVs[:index], leftTiKVs[index+1:]...)
	}
	return kvs, nil
}

func getRandomFailpointPath() []string {
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
	url := fmt.Sprintf("http://%s:10008/operation/failpoint/add", kv)
	_, err = utils.DoPost(url, data)
	return errors.Trace(err)
}

func emptyFailpoints(kv string) error {
	url := fmt.Sprintf("http://%s:10008/operation/failpoint/clean", kv)
	_, err := utils.DoPost(url, []byte{})
	return errors.Trace(err)
}
