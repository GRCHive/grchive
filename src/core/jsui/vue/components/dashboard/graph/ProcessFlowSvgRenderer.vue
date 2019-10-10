<template>
    <svg id="svgrenderer"
         width="100%"
         height="100%"
         @mousemove="onMouseMove"
         @mousedown="onMouseDown"
         @mouseup="onMouseUp"
         ref="svgrenderer"
    >
        <g id="nodes">
            <process-flow-svg-node
                v-for="item in nodes"
                :key="item.Id"
                :node="item"
                @onmousedown="onMouseDownNode"
                @onmouseup="onMouseUpNode"
            >
            </process-flow-svg-node>
        </g>
    </svg>
</template>

<script lang="ts">

import Vue from 'vue'
import VueSetup from '../../../../ts/vueSetup'
import ProcessFlowSvgNode from './ProcessFlowSvgNode.vue'

export default Vue.extend({
    components: {
        ProcessFlowSvgNode
    },
    computed: {
        nodes() {
            return VueSetup.store.state.currentProcessFlowFullData.Nodes
        }
    },
    data: () => ({
        moveNodeActive: false
    }),
    methods: {
        onMouseMove(e : MouseEvent) {
            if (!VueSetup.store.getters.isNodeSelected || !this.moveNodeActive) {
                return
            }

            VueSetup.store.commit('addNodeDisplayTranslation', {
                nodeId: VueSetup.store.state.selectedNodeId,
                tx: e.movementX,
                ty: e.movementY
            })
        },
        onMouseDownNode(e : MouseEvent, nodeId : number) {
            VueSetup.store.commit('setSelectedProcessFlowNode', nodeId)
            this.moveNodeActive = true
        },
        onMouseUpNode(e : MouseEvent, nodeId : number) {
            this.moveNodeActive = false
        },
        onMouseDown(e : MouseEvent) {
            VueSetup.store.commit('setSelectedProcessFlowNode', -1)
        },
        onMouseUp(e : MouseEvent) {
            VueSetup.store.commit('setSelectedProcessFlowNode', -1)
        }
    }
})

</script>

<style scoped>

#svgrenderer {
    display: block;
}

</style>
