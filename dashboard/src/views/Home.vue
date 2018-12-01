<template>
    <!--<div :style="{ height: '300px'}" class="clusterChart" ref="clusterChart"></div>-->
    <div>
        <h1>Cluster Info</h1>
        <div id="clusterChart" style="width: 600px;height:400px;">

        </div>
        <el-button @click="drawCluster" class="el-button">Cluster Info</el-button>
    </div>
</template>

<script>
    export default {
        data() {
            return {}
        },
        created() {
            // this.drawCluster();
        },
        updated() {
        },
        methods: {
            drawCluster() {
                console.log("test");
                var myChart = this.$echarts.init(document.getElementById('clusterChart'));
                console.log("test");

                var option = {
                    title: {
                        text: 'Graph 简单示例'
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
                            data: [{
                                name: '节点1',
                                x: 300,
                                y: 300
                            }, {
                                name: '节点2',
                                x: 800,
                                y: 300
                            }, {
                                name: '节点3',
                                x: 550,
                                y: 100
                            }, {
                                name: '节点4',
                                x: 550,
                                y: 500
                            }],
                            // links: [],
                            links: [{
                                source: 0,
                                target: 1,
                                symbolSize: [5, 20],
                                label: {
                                    normal: {
                                        show: true
                                    }
                                },
                                lineStyle: {
                                    normal: {
                                        width: 5,
                                        curveness: 0.2
                                    }
                                }
                            }, {
                                source: '节点2',
                                target: '节点1',
                                label: {
                                    normal: {
                                        show: true
                                    }
                                },
                                lineStyle: {
                                    normal: {curveness: 0.2}
                                }
                            }, {
                                source: '节点1',
                                target: '节点3'
                            }, {
                                source: '节点2',
                                target: '节点3'
                            }, {
                                source: '节点2',
                                target: '节点4'
                            }, {
                                source: '节点1',
                                target: '节点4'
                            }],
                            lineStyle: {
                                normal: {
                                    opacity: 0.9,
                                    width: 2,
                                    curveness: 0
                                }
                            }
                        }
                    ]
                };

                console.log("test");

                myChart.setOption(option)
            }
        }
    }
</script>

<style scoped>

</style>