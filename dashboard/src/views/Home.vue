<template>
    <!--<div :style="{ height: '300px'}" class="clusterChart" ref="clusterChart"></div>-->
    <div>
        <div class="clusterChart" id="clusterChart">

        </div>
        <!--<el-button @click="drawCluster" class="el-button">Cluster Info</el-button>-->
        <div id="metricChart" style="width: 600px;height:400px;">

        </div>
    </div>
</template>

<script>
    import ajax from '../request/index'

    export default {
        data() {
            return {
                clusterInfo: {
                    tidb: [],
                    tikv: [],
                    pd: [],
                },
                partition: {
                    kind: "",
                    groups: [],
                    real_groups: [],
                }
            }
        },
        created() {
            // this.drawCluster();
        },
        updated() {
        },
        mounted() {
            this.drawCluster();
            this.drawMetric();
        },
        methods: {
            drawCluster() {
                ajax.getClusterInfo().then(result => {
                    if (result.data != null) {
                        this.clusterInfo = result.data;

                        var myChart = this.$echarts.init(document.getElementById('clusterChart'));

                        var increment = 360 / this.clusterInfo.tikv.length;
                        var startAngle = 270;
                        var r = 250;
                        var datas = []

                        this.clusterInfo.tikv.forEach((tikv, index) => {
                            var angle = startAngle + increment * index;
                            var rads = angle * Math.PI / 180;
                            datas.push({
                                name: "tikv:" + tikv,
                                x: Math.trunc(500 + r * Math.cos(rads)),
                                y: Math.trunc(500 + r * Math.sin(rads)),
                            })
                        })

                        ajax.getPartitionInfo().then(result => {
                            if (result.data != null) {
                                this.partition = result.data
                            }

                            console.log(result.data);
                            var links = [];
                            if (this.partition.real_groups.length <= 1) {
                                this.clusterInfo.tikv.forEach((tikv, index) => {
                                    for (var i = 0; i < index; i++) {
                                        links.push({
                                            source: "tikv:" + tikv,
                                            target: "tikv:" + this.clusterInfo.tikv[i],
                                        })
                                    }
                                    for (var i = index + 1; i < this.clusterInfo.tikv.length; i++) {
                                        links.push({
                                            source: "tikv:" + tikv,
                                            target: "tikv:" + this.clusterInfo.tikv[i],
                                        })
                                    }
                                })
                            } else {
                                switch (this.partition.kind) {
                                    case "full":
                                        this.partition.real_groups.forEach(item => {
                                            item.forEach((h, index) => {
                                                for (var i = 0; i < index; i++) {
                                                    links.push({
                                                        source: "tikv:" + h,
                                                        target: "tikv:" + item[i],
                                                        symbolSize: [3, 15],
                                                    })
                                                }
                                                for (var i = index + 1; i < item.length; i++) {
                                                    links.push({
                                                        source: "tikv:" + h,
                                                        target: "tikv:" + item[i],
                                                        symbolSize: [3, 15],
                                                    })
                                                }
                                            })
                                        })
                                        break;
                                    case "partial":
                                        if (this.partition.real_groups.length >= 2) {
                                            for (var index = 0; index <= this.partition.real_groups[0].length; index++) {
                                                var itemt = this.partition.real_groups[0];
                                                var item = this.partition.real_groups[0];
                                                if (this.partition.real_groups.length > 2) {
                                                    item = item.concat(this.partition.real_groups[2])
                                                }
                                                for (var i = 0; i < index; i++) {
                                                    links.push({
                                                        source: "tikv:" + itemt[index],
                                                        target: "tikv:" + item[i],
                                                        symbolSize: [3, 15],
                                                    })
                                                }
                                                for (var i = index + 1; i < item.length; i++) {
                                                    links.push({
                                                        source: "tikv:" + itemt[index],
                                                        target: "tikv:" + item[i],
                                                        symbolSize: [3, 15],
                                                    })
                                                }

                                            }
                                            for (var index = 0; index <= this.partition.real_groups[1].length; index++) {
                                                var itemt = this.partition.real_groups[1];
                                                var item = this.partition.real_groups[1];
                                                if (this.partition.real_groups.length > 2) {
                                                    item = item.concat(this.partition.real_groups[2])
                                                }
                                                for (var i = 0; i < index; i++) {
                                                    links.push({
                                                        source: "tikv:" + itemt[index],
                                                        target: "tikv:" + item[i],
                                                        symbolSize: [3, 15],
                                                    })
                                                }
                                                for (var i = index + 1; i < item.length; i++) {
                                                    links.push({
                                                        source: "tikv:" + itemt[index],
                                                        target: "tikv:" + item[i],
                                                        symbolSize: [3, 15],
                                                    })
                                                }

                                            }

                                            if (this.partition.real_groups.length == 3) {
                                                for (var index = 0; index <= this.partition.real_groups[2].length; index++) {
                                                    var itemt = this.partition.real_groups[2];
                                                    var item = this.partition.real_groups[2];
                                                    item = item.concat(this.partition.real_groups[0])
                                                    item = item.concat(this.partition.real_groups[1])
                                                    for (var i = 0; i < index; i++) {
                                                        links.push({
                                                            source: "tikv:" + itemt[index],
                                                            target: "tikv:" + item[i],
                                                            symbolSize: [3, 15],
                                                        })
                                                    }
                                                    for (var i = index + 1; i < item.length; i++) {
                                                        links.push({
                                                            source: "tikv:" + itemt[index],
                                                            target: "tikv:" + item[i],
                                                            symbolSize: [3, 15],
                                                        })
                                                    }

                                                }
                                            }
                                        } else {
                                            this.$notify.error({
                                                title: 'Error',
                                                message: "partition is not supported !!"
                                            })
                                        }
                                        break;
                                    case "simplex":
                                        if (this.partition.real_groups.length >= 2) {
                                            for (var index = 0; index <= this.partition.real_groups[0].length; index++) {
                                                var itemt = this.partition.real_groups[0];
                                                var item = this.partition.real_groups[0];
                                                item = item.concat(this.partition.real_groups[1]);
                                                if (this.partition.real_groups.length > 2) {
                                                    item = item.concat(this.partition.real_groups[2])
                                                }
                                                for (var i = 0; i < index; i++) {
                                                    links.push({
                                                        source: "tikv:" + itemt[index],
                                                        target: "tikv:" + item[i],
                                                        symbolSize: [3, 15],
                                                    })
                                                }
                                                for (var i = index + 1; i < item.length; i++) {
                                                    links.push({
                                                        source: "tikv:" + itemt[index],
                                                        target: "tikv:" + item[i],
                                                        symbolSize: [3, 15],
                                                    })
                                                }

                                            }
                                            for (var index = 0; index <= this.partition.real_groups[1].length; index++) {
                                                var itemt = this.partition.real_groups[1];
                                                var item = this.partition.real_groups[1];
                                                if (this.partition.real_groups.length > 2) {
                                                    item = item.concat(this.partition.real_groups[2])
                                                }
                                                for (var i = 0; i < index; i++) {
                                                    links.push({
                                                        source: "tikv:" + itemt[index],
                                                        target: "tikv:" + item[i],
                                                        symbolSize: [3, 15],
                                                    })
                                                }
                                                for (var i = index + 1; i < item.length; i++) {
                                                    links.push({
                                                        source: "tikv:" + itemt[index],
                                                        target: "tikv:" + item[i],
                                                        symbolSize: [3, 15],
                                                    })
                                                }

                                            }
                                            if (this.partition.real_groups.length == 3) {
                                                for (var index = 0; index <= this.partition.real_groups[2].length; index++) {
                                                    var itemt = this.partition.real_groups[2];
                                                    var item = this.partition.real_groups[2];
                                                    item = item.concat(this.partition.real_groups[0])
                                                    item = item.concat(this.partition.real_groups[1])
                                                    for (var i = 0; i < index; i++) {
                                                        links.push({
                                                            source: "tikv:" + itemt[index],
                                                            target: "tikv:" + item[i],
                                                            symbolSize: [3, 15],
                                                        })
                                                    }
                                                    for (var i = index + 1; i < item.length; i++) {
                                                        links.push({
                                                            source: "tikv:" + itemt[index],
                                                            target: "tikv:" + item[i],
                                                            symbolSize: [3, 15],
                                                        })
                                                    }

                                                }
                                            }
                                        } else {
                                            this.$notify.error({
                                                title: 'Error',
                                                message: "partition is not supported !!"
                                            })
                                        }
                                        break;
                                    default:
                                        this.$notify.error({
                                            title: 'Error',
                                            message: "partition is not supported !!"
                                        })
                                }
                            }

                            var option = {
                                title: {
                                    text: 'TiDB Cluster Info'
                                },
                                tooltip: {},
                                animationDurationUpdate: 1500,
                                animationEasingUpdate: 'quinticInOut',
                                series: [
                                    {
                                        type: 'graph',
                                        layout: 'none',
                                        symbolSize: 50,
                                        roam: true,
                                        label: {
                                            normal: {
                                                show: true
                                            }
                                        },
                                        edgeSymbol: ['circle', 'arrow'],
                                        edgeSymbolSize: [4, 10],
                                        edgeLabel: {
                                            normal: {
                                                textStyle: {
                                                    fontSize: 20
                                                }
                                            }
                                        },
                                        data: datas,
                                        // links: [],
                                        links: links,
                                        lineStyle: {
                                            normal: {
                                                opacity: 0.9,
                                                width: 2,
                                                curveness: 0.05
                                            }
                                        }
                                    }
                                ]
                            };

                            myChart.setOption(option)
                        }).catch(resp => {
                            this.$notify.error({
                                title: 'Error',
                                message: resp.message
                            })
                        })
                    } else {
                        this.$notify.error({
                            title: 'Error',
                            message: "cluster info is empty"
                        })
                    }
                }).catch(resp => {
                    this.$notify.error({
                        title: 'Error',
                        message: resp.message
                    })
                })
            },

            getOption(title, time, value) {
                var option = {
                    title: {
                        text: title,
                    },
                    tooltip: {
                        trigger: 'axis'
                    },
                    legend: {
                        data: [title]
                    },
                    toolbox: {
                        show: true,
                        feature: {
                            dataZoom: {
                                yAxisIndex: 'none'
                            },
                            dataView: {readOnly: false},
                            magicType: {type: ['line', 'bar']},
                            restore: {},
                            saveAsImage: {}
                        }
                    },
                    xAxis: {
                        type: 'category',
                        boundaryGap: false,
                        data: time
                    },
                    yAxis: {
                        type: 'value',
                        axisLabel: {
                            formatter: '{value} s'
                        }
                    },
                    series: [
                        {
                            name: title,
                            type: 'line',
                            data: value,
                            markPoint: {
                                data: [
                                    {value: "partition", xAxis: time[20], yAxis: value[20]}
                                ]
                            },
                            markLine: {
                                data: [
                                    {type: 'average', name: '平均值'}
                                ]
                            }
                        }
                    ]
                };
                return option
            },

            drawData(title, metric, timeFrom, timeTo, id) {
                var myChart = this.$echarts.init(document.getElementById(id));
                ajax.getDuration(metric, timeFrom, timeTo).then(result => {
                    var time = []
                    var value = []
                    console.log(result.data.data.result[0])
                    result.data.data.result[0].values.forEach((e, index) => {
                        var date = new Date(e[0] * 1000);
                        // Hours part from the timestamp
                        var hours = date.getHours();
                        // Minutes part from the timestamp
                        var minutes = "0" + date.getMinutes();
                        // Seconds part from the timestamp
                        var seconds = "0" + date.getSeconds();
                        var formattedTime = hours + ':' + minutes.substr(-2) + ':' + seconds.substr(-2);
                        time.push(formattedTime)
                        value.push(parseFloat(e[1]).toFixed(4))
                    })
                    var option = this.getOption(title, time, value)
                    myChart.setOption(option)
                }).catch(resp => {
                    this.$notify.error({
                        title: 'ERROR',
                        message: resp.message,
                        duration: 0
                    })
                })
            },

            drawMetric() {
                console.log("test");
                const end = new Date();
                const start = new Date();
                start.setTime(start.getTime() - 3600 * 1000 * 24);
                var startTimestamp = start.getTime() / 1000
                var endTimestamp = end.getTime() / 1000
                var option = this.drawData("Duration", "tidb_server_handle_query_duration_seconds_bucket", startTimestamp, endTimestamp, "metricChart")
                console.log("test");
            }
        }
    }
</script>

<style scoped>
    .clusterChart {
        width: 100em;
        height: 60em;
    }
</style>