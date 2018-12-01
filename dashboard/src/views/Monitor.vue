<template>
    <!--<div :style="{ height: '300px'}" class="clusterChart" ref="clusterChart"></div>-->
    <div>
        <h1>Cluster Info</h1>
        <div id="clusterChart" style="width: 600px;height:400px;">

        </div>
        <!--<el-button @click="drawCluster" class="el-button">Cluster Info</el-button>-->
    </div>
</template>

<script>
    import ajax from '../request/index'
    export default {
        data() {
            return {}
        },
        created() {
            // this.drawCluster();
        },
        updated() {
        },
        mounted() {
            this.drawCluster();
        },
        methods: {
            getOption(title, time, value) {
                var option = {
                    title: {
                        text: title,
                    },
                    tooltip: {
                        trigger: 'axis'
                    },
                    legend: {
                        data:[title]
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
                    xAxis:  {
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
                            name:title,
                            type:'line',
                            data:value,
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
            drawData(title,metric, timeFrom, timeTo, id) {
                var myChart = this.$echarts.init(document.getElementById(id));
                ajax.getDuration(metric, timeFrom,timeTo).then(result => {
                    var time = []
                    var value = []
                    console.log(result.data.data.result[0])
                    result.data.data.result[0].values.forEach((e, index) => {
                        var date = new Date(e[0]*1000);
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
                }).catch(resp =>{
                    this.$notify.error({
                    title: 'ERROR',
                    message: resp.message,
                    duration: 0
                    })
                })
            },
            drawCluster() {
                console.log("test");
                const end = new Date();
                const start = new Date();
                start.setTime(start.getTime() - 3600 * 1000 * 24);
                var startTimestamp = start.getTime()/1000
                var endTimestamp = end.getTime()/1000
                var option = this.drawData("Duration","tidb_server_handle_query_duration_seconds_bucket",startTimestamp,endTimestamp, "clusterChart")
                console.log("test");
            }
        }
    }
</script>

<style scoped>

</style>