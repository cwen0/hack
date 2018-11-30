# Use goreman to run `go get github.com/mattn/goreman`

pd: ./bin/pd-server --client-urls="http://127.0.0.1:12379" --peer-urls="http://127.0.0.1:12380" --data-dir=./var/default.pd --log-file ./var/pd.log

tikv1: sleep 5 && ./bin/tikv-server --pd 127.0.0.1:12379 -A 127.0.0.1:21161 --advertise-addr 127.0.0.1:11111 --data-dir ./var/store1 --log-file ./var/tikv1.log
tikv2: sleep 5 && ./bin/tikv-server --pd 127.0.0.1:12379 -A 127.0.0.1:21162 --advertise-addr 127.0.0.1:22222 --data-dir ./var/store2 --log-file ./var/tikv2.log
tikv3: sleep 5 && ./bin/tikv-server --pd 127.0.0.1:12379 -A 127.0.0.1:21163 --advertise-addr 127.0.0.1:33333 --data-dir ./var/store3 --log-file ./var/tikv3.log
# proxy
# bridge3: ./bin/bridge --reorder=false 127.0.0.1:11111 127.0.0.1:21163
proxy1: ./bin/proxy --listen-addr=:11111 --upstream 127.0.0.1:21161 --config-listen-addr=:10018
proxy2: ./bin/proxy --listen-addr=:22222 --upstream 127.0.0.1:21162 --config-listen-addr=:10019
proxy3: ./bin/proxy --listen-addr=:33333 --upstream 127.0.0.1:21163 --config-listen-addr=:10020

tidb: ./bin/tidb-server -P 4001 -status 10081 -path="127.0.0.1:12379" -store=tikv --log-file ./var/tidb.log --lease 60