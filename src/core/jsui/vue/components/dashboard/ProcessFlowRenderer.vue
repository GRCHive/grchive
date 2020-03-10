<template>
    <div id="flowRenderer" class="ma-0" :style="contentStyle" ref="renderer">
        <div class="max-height" v-if="hasProcessFlowToRender">
            <process-flow-svg-renderer
                :svg-width="svgWidth"
                :svg-height="svgHeight"
                v-if="renderReady"
            ></process-flow-svg-renderer>

            <div class="max-height" v-else>
                <v-row class="max-height ma-0" align="center" width="100%">
                    <v-col class="pa-0">
                        <v-row justify="center" class="ma-0">
                            <v-btn x-large icon :loading="true">
                            </v-btn>
                        </v-row>
                    </v-col>
                </v-row>
            </div>
        </div>

        <div class="max-height" v-else>
            <v-row class="max-height ma-0" align="center" width="100%">
                <v-col class="pa-0">
                    <v-row justify="center" class="ma-0">
                        <p class="display-1">This process flow is empty!</p>
                    </v-row>
                    <v-row justify="center" class="ma-0">
                        <p class="body-1">Get started by clicking the "Add Node" button.</p>
                    </v-row>
                    <v-row justify="center" class="ma-0">
                        <v-btn icon @click="refreshProcessFlow"
                                    :disabled="processFlowLoading"
                                    :loading="processFlowLoading"
                        >
                            <v-icon x-large>mdi-refresh</v-icon>
                        </v-btn>
                    </v-row>
                </v-col>
            </v-row>
        </div>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import VueSetup from '../../../ts/vueSetup'
import ProcessFlowSvgRenderer from './graph/ProcessFlowSvgRenderer.vue'
import RenderLayout from '../../../ts/render/renderLayout'

export default Vue.extend({
    components: {
        ProcessFlowSvgRenderer
    },
    props: {
        contentMaxHeightClip: Number,
        contentMaxWidthClip: Number,
        displayRect: {
            type: Object as () => IDOMRect
        }
    },
    computed: {
        svgHeight() : number {
            return (<IDOMRect>this.displayRect).height
        },
        svgWidth(): number {
            return (<IDOMRect>this.displayRect).width
        },
        hasProcessFlowToRender() : boolean {
            return !!VueSetup.store.state.currentProcessFlowFullData &&
                VueSetup.store.state.currentProcessFlowFullData!.NodeKeys.length > 0
        },
        renderReady() : boolean {
            return RenderLayout.store.state.ready && 
                this.svgHeight > 0 && 
                this.svgWidth > 0
        },
        contentStyle() : any {
            return {
                "height": "100%",
                "max-height": `calc(100% - ${this.contentMaxHeightClip.toString()}px)`,
                "width": "100%",
                "max-width": `calc(100% - ${this.contentMaxWidthClip.toString()}px)`
            }
        },
        processFlowLoading() : boolean {
            return VueSetup.store.getters.isFullRequestInProgress
        }
    },
    methods: {
        refreshProcessFlow() {
            VueSetup.store.dispatch('refreshCurrentProcessFlowFullData', 
                VueSetup.store.state.currentProcessFlowBasicData!.Id)
        }
    },
})
</script>
