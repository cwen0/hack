// import { REFUSED } from 'dns'
import axios from 'axios'
import Mock from 'mockjs'

const Proxy = '/api'

// get all test templates
Mock.mock(`${Proxy}/clusterInfo`, 'get', {
    "tidb": ["10.0.0.1"],
    "pd": ["10.0.0.2"],
    "tikv": ["10.0.0.3", "10.0.0.4", "10.0.0.5", "10.0.0.6", "10.0.0.7"]
})

Mock.mock(`${Proxy}/partition`, 'get', {
    "kind": "full",
    "groups": [["10.0.0.3", "10.0.0.4", "10.0.0.5"], ["10.0.0.6", "10.0.0.7"]],
    "real_groups": [["10.0.0.3", "10.0.0.4", "10.0.0.5"], ["10.0.0.6", "10.0.0.7"]],
    //"kind": "partial",
    //"groups": [["10.0.0.3", "10.0.0.4"], ["10.0.0.6", "10.0.0.7"]],
    //"real_groups": [["10.0.0.3", "10.0.0.4"], ["10.0.0.6", "10.0.0.7"], ["10.0.0.5"]],
    // "kind": "simplex",
    // "groups": [["10.0.0.3", "10.0.0.4"], ["10.0.0.6", "10.0.0.7"]],
    // "real_groups": [["10.0.0.3", "10.0.0.4"], ["10.0.0.6", "10.0.0.7"], ["10.0.0.5"]],
})

class Ajax {
    getClusterInfo() {
        return axios.get(`${Proxy}/clusterInfo`)
    }

    getPartitionInfo() {
        return axios.get(`${Proxy}/partition`)
    }
}
export default new Ajax()
