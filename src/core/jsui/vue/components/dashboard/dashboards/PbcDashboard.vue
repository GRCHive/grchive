<template>
    <div>
        <advanced-doc-request-filters
            v-model="filterData"
        >
        </advanced-doc-request-filters>

        <v-row>
            <v-col cols="6">
                <v-card ref="overallContainer">
                    <v-card-title>
                        Overall Progress
                    </v-card-title>
                    <v-divider></v-divider>

                    <v-row justify="center" v-if="loadingOverall">
                        <v-progress-circular size="64" indeterminate></v-progress-circular>
                    </v-row>
                    <div id="overall" ref="overallDiv" style="width: 100%; height: 600px;"></div>
                </v-card>
            </v-col>

            <v-col cols="6">
                <v-card ref="overallContainer">
                    <v-card-title>
                        By Category
                    </v-card-title>
                    <v-divider></v-divider>

                    <v-select
                        v-model="selectedCategory"
                        label="Category"
                        dense
                        hide-details
                        :items="categoryItems"
                        filled
                        class="mb-4"
                    >
                    </v-select>

                    <v-row justify="center" v-if="loadingCategory">
                        <v-progress-circular size="64" indeterminate></v-progress-circular>
                    </v-row>
                    <div id="category" ref="categoryDiv" style="width: 100%; height: 600px;"></div>
                </v-card>
            </v-col>
        </v-row>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import { Watch } from 'vue-property-decorator'

import * as echarts from 'echarts/lib/echarts.js'
import 'echarts/lib/chart/bar'
import 'echarts/lib/chart/pie'
import 'echarts/lib/component/tooltip'
import 'echarts/lib/component/legend'
import 'echarts/lib/component/dataZoom'
import 'zrender/lib/svg/svg'

import {
    DocRequestFilterData,
    NullDocRequestFilterData,
    copyDocRequestFilterData,
    getDocRequestStatusString,
    allDocRequestStatus, DocRequestStatus
} from '../../../../ts/docRequests'
import AdvancedDocRequestFilters from '../../../generic/filters/AdvancedDocRequestFilters.vue'

import {
    getPbcOverallProgress, TGetPbcOverallProgressOutputs,
    getPbcCategoryProgress, TGetPbcCategoryProgressOutputs,
} from '../../../../ts/api/analytics/apiPbcAnalytics'
import { PageParamsStore } from '../../../../ts/pageParams'
import { contactUsUrl } from '../../../../ts/url'

const statusToColor  = new Map<DocRequestStatus, string>()
statusToColor.set(DocRequestStatus.Open, '#424242')
statusToColor.set(DocRequestStatus.InProgress, '#1976d2')
statusToColor.set(DocRequestStatus.Feedback, '#fb8c00')
statusToColor.set(DocRequestStatus.Complete, '#82b1ff')
statusToColor.set(DocRequestStatus.Overdue, '#ff5252')
statusToColor.set(DocRequestStatus.Approved, '#4caf50')

@Component({
    components: {
        AdvancedDocRequestFilters,
    }
})
export default class PbcDashboard extends Vue {
    filterData: DocRequestFilterData = copyDocRequestFilterData(NullDocRequestFilterData)

    $refs!: {
        overallContainer: HTMLElement
        overallDiv: HTMLElement

        categoryContainer: HTMLElement
        categoryDiv: HTMLElement
    }

    overallGraph : any | null = null
    loadingOverall: boolean = false

    categoryGraph : any | null = null
    loadingCategory: boolean = false
    selectedCategory : string = 'assignee'
    numCategories : number = 0

    readonly categoryItems = [
        {
            text: 'Requester',
            value: 'requester',
        },
        {
            text: 'Assignee',
            value: 'assignee',
        },
        {
            text: 'Document Category',
            value: 'cat',
        },
        {
            text: 'Process Flow',
            value: 'flow',
        },
        {
            text: 'Control',
            value: 'control',
        },
        {
            text: 'Risk',
            value: 'risk',
        },
        {
            text: 'GL account',
            value: 'gl',
        },
        {
            text: 'System',
            value: 'system',
        },
    ]

    @Watch('filterData', {deep:true})
    refreshOverallGraph() {
        this.loadingOverall = true
        getPbcOverallProgress({
            orgId: PageParamsStore.state.organization!.Id,
            filter: this.filterData,
        }).then((resp : TGetPbcOverallProgressOutputs) => {
            let data : any[] = allDocRequestStatus.map((ele : DocRequestStatus) => {
                return {
                    name: getDocRequestStatusString(ele),
                    value: resp.data[ele],
                    itemStyle: {
                        color: statusToColor.get(ele)
                    }
                }
            })

            this.overallGraph.setOption({
                tooltip: {
                    formatter: '{a} <br/>{b} : {c} ({d}%)'
                },
                legend: {
                    orient: 'vertical',
                    left: 'left',
                    data: allDocRequestStatus.map((ele : DocRequestStatus) => getDocRequestStatusString(ele))
                },
                series: [
                    {
                        name: 'PBC Status',
                        type: 'pie',
                        radius: '80%',
                        center: ['50%', '50%'],
                        data: data,
                    }
                ]
            });
            Vue.nextTick(this.onResize)
        }).catch((err : any) => {
            // @ts-ignore
            this.$root.$refs.snackbar.showSnackBar(
                "Oops! Something went wrong. Try again.",
                true,
                "Contact Us",
                contactUsUrl,
                true);
        }).finally(() => {
            this.loadingOverall = false
        })
    }

    @Watch('filterData', {deep:true})
    @Watch('selectedCategory')
    refreshCategoryGraph() {
        this.loadingCategory = true
        getPbcCategoryProgress({
            orgId: PageParamsStore.state.organization!.Id,
            category: this.selectedCategory,
            filter: this.filterData,
        }).then((resp : TGetPbcCategoryProgressOutputs) => {
            this.numCategories = resp.data.length
            let dataSeries : any[] = allDocRequestStatus.map((s : DocRequestStatus) => {
                return {
                    name: getDocRequestStatusString(s),
                    type: 'bar',
                    stack: 'Stack',
                    label: {
                        show: true,
                        position: 'insideRight',
                    },
                    itemStyle: {
                        color: statusToColor.get(s)
                    },
                    data: resp.data.map((ele : any) => ele.Data[s]),
                }
            })

            this.categoryGraph.setOption({
                tooltip: {
                    trigger: 'axis',
                    axisPointer: {
                        type: 'shadow'
                    }
                },
                legend: {
                    data: allDocRequestStatus.map((ele : DocRequestStatus) => getDocRequestStatusString(ele))
                },
                xAxis: {
                    type: 'value',
                    minInterval: 1,
                },
                yAxis: {
                    type: 'category',
                    data: resp.data.map((ele: any) => ele.Name),
                    axisLabel: {
                        rotate: 90,
                    }
                },
                grid: {
                    containLabel: true,
                },
                dataZoom: [
                    {
                        show: true,
                        start: 0,
                        end: 100,
                        yAxisIndex: 0,
                    },
                    {
                        type: 'inside',
                        start: 0,
                        end: 100,
                        yAxisIndex: 0,
                    },
                ],
                series: dataSeries,
            });

            Vue.nextTick(this.onResize)
        }).catch((err : any) => {
            // @ts-ignore
            this.$root.$refs.snackbar.showSnackBar(
                "Oops! Something went wrong. Try again.",
                true,
                "Contact Us",
                contactUsUrl,
                true);
        }).finally(() => {
            this.loadingCategory = false
        })
    }

    mounted() {
        this.overallGraph = echarts.init(this.$refs.overallDiv, null, { renderer: 'canvas' })
        this.categoryGraph = echarts.init(this.$refs.categoryDiv, null, { renderer : 'canvas' })

        this.refreshOverallGraph()
        this.refreshCategoryGraph()

        window.addEventListener('resize', this.onResize)
    }

    onResize() {
        if (!!this.overallGraph) {
            this.overallGraph.resize()
        }

        if (!!this.categoryGraph) {
            this.categoryGraph.resize()
        }
    }
}

</script>
