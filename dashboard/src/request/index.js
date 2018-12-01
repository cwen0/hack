// import { REFUSED } from 'dns'
import axios from 'axios'
const Proxy = '/api/datasources/proxy/1/api/v1/query_range'
class Ajax {
    getDuration(metric, start, end) {
       var data = axios.get(`${Proxy}?query=histogram_quantile(0.95%2C%20sum(rate(${metric}%5B1m%5D))%20by%20(le))&start=${start}&end=${end}&step=30`)
        return data
    }
}
export default new Ajax()
