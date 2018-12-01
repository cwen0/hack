// import { REFUSED } from 'dns'
import axios from 'axios'

const Proxy = '/operation'

const GrafanaProxy = '/api/datasources/proxy/1/api/v1/query_range'

// get all test templates
// Mock.mock(`${Proxy}/topology`, 'get', {
//     "tidb": ["10.0.0.1"],
//     "pd": ["10.0.0.2"],
//     "tikv": ["10.0.0.3", "10.0.0.4", "10.0.0.5", "10.0.0.6", "10.0.0.7"]
// })
//
// Mock.mock(`${Proxy}/store/leaders`, 'get', 4)
// // Mock.mock(`${Proxy}/stores`, 'get', [
// //     {"ip": "10.0.0.3", "leader_count": 3},
// //     {"ip": "10.0.0.4", "leader_count": 4},
// //     {"ip": "10.0.0.5", "leader_count": 5},
// //     {"ip": "10.0.0.6", "leader_count": 6},
// //     {"ip": "10.0.0.7", "leader_count": 7},
// // ])
//
// Mock.mock(`${Proxy}/stores`, 'get', {
//     "count": 5,
//     "stores": [
//         {
//             "store": {
//                 "address": "10.0.0.3",
//             },
//             "status": {
//                 "leader_count": 3,
//             }
//         },
//         {
//             "store": {
//                 "address": "10.0.0.4",
//             },
//             "status": {
//                 "leader_count": 4,
//             }
//         },
//         {
//             "store": {
//                 "address": "10.0.0.5",
//             },
//             "status": {
//                 "leader_count": 5,
//             }
//         },
//         {
//             "store": {
//                 "address": "10.0.0.6",
//             },
//             "status": {
//                 "leader_count": 6,
//             }
//         },
//         {
//             "store": {
//                 "address": "10.0.0.7",
//             },
//             "status": {
//                 "leader_count": 7,
//             }
//         }
//     ]
// })
//
// Mock.mock(`${Proxy}/partition`, 'get', {
//     "kind": "full",
//     "groups": [["10.0.0.3", "10.0.0.4", "10.0.0.5"], ["10.0.0.6", "10.0.0.7"]],
//     "real_groups": [["10.0.0.3", "10.0.0.4", "10.0.0.5"], ["10.0.0.6", "10.0.0.7"]],
//     // "kind": "partial",
//     // "groups": [["10.0.0.3", "10.0.0.4"], ["10.0.0.6", "10.0.0.7"]],
//     // "real_groups": [["10.0.0.3", "10.0.0.4"], ["10.0.0.6", "10.0.0.7"], ["10.0.0.5"]],
//     // "kind": "simplex",
//     // "groups": [["10.0.0.3", "10.0.0.4"], ["10.0.0.6", "10.0.0.7"]],
//     // "real_groups": [["10.0.0.3", "10.0.0.4"], ["10.0.0.6", "10.0.0.7"], ["10.0.0.5"]],
// })

class Ajax {
    getClusterInfo() {
        return axios.get(`${Proxy}/topology`)
    }

    getPartitionInfo() {
        return axios.get(`${Proxy}/partition`)
    }

    //getStoreLeaderCount(ip) {
    //    return axios.get(`${Proxy}/store/leaders?ip=${ip}`)
    //}

    getStoresInfo() {
        return axios.get(`${Proxy}/store`)
    }

    getDuration(metric, start, end) {
        var data = axios.get(`${GrafanaProxy}?query=histogram_quantile(0.95%2C%20sum(rate(${metric}%5B1m%5D))%20by%20(le))&start=${start}&end=${end}&step=30`)
        return data
    }

    setevictTikvLeader(ip) {
        return axios.post(`${proxy}/evictleadere/${ip}`)
    }

    setNetworkPartition(data) {
        return axios.post(`${proxy}/partition`, data)
    }

    cleanNetworkPartition() {
        return axios.post(`${proxy}/partition/clean`)
    }

    setFailpoint(data) {
        return axios.post(`${proxy}/failpoint`, data)
    }
}

export default new Ajax()
