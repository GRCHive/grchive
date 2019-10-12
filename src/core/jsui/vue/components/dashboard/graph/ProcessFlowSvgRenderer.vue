<template>
    <svg id="svgrenderer"
         width="100%"
         height="100%"
         preserveAspectRatio="none"
         :viewBox="`${viewBox.x} ${viewBox.y} ${viewBox.width} ${viewBox.height}`"
         @mousemove="onMouseMove"
         @mousedown="onMouseDown"
         @mouseup="onMouseUp"
         @mouseleave="onMouseLeave"
         @contextmenu="onContextMenu"
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
    props: {
        svgWidth: Number,
        svgHeight: Number
    },
    computed: {
        nodes() {
            return VueSetup.store.state.currentProcessFlowFullData.Nodes
        },
        viewBox() {
            return {
                x: this.viewBoxX,
                y: this.viewBoxY,
                width: this.svgWidth,
                height: this.svgHeight
            }
        }
    },
    data: () => ({
        moveNodeActive: false,
        moveViewBoxActive: false,
        viewBoxX: 0,
        viewBoxY: 0
    }),
    methods: {
        doMoveNode(e : MouseEvent) {
            if (!VueSetup.store.getters.isNodeSelected) {
                return
            }

            VueSetup.store.commit('addNodeDisplayTranslation', {
                nodeId: VueSetup.store.state.selectedNodeId,
                tx: e.movementX,
                ty: e.movementY
            })
        },
        doMoveViewBox(e : MouseEvent) {
            this.viewBoxX -= e.movementX
            this.viewBoxY -= e.movementY
        },
        onMouseMove(e : MouseEvent) {
            if (this.moveNodeActive) {
                this.doMoveNode(e)
            } else if (this.moveViewBoxActive) {
                this.doMoveViewBox(e)
            }
        },
        onMouseDownNode(e : MouseEvent, nodeId : number) {
            if (e.button != 0) {
                return
            }

            VueSetup.store.commit('setSelectedProcessFlowNode', nodeId)
            this.moveNodeActive = true
            e.stopPropagation()
        },
        onMouseUpNode(e : MouseEvent, nodeId : number) {
            if (e.button != 0) {
                return
            }

            this.moveNodeActive = false
            e.stopPropagation()
        },
        onMouseDown(e : MouseEvent) {
            if (e.button == 0) {
                VueSetup.store.commit('setSelectedProcessFlowNode', -1)
            }

            if (e.button == 0 || e.button == 1) {
                this.moveViewBoxActive = true
            }
        },
        onMouseUp(e : MouseEvent) {
            if (e.button == 0 || e.button == 1) {
                this.moveViewBoxActive = false
            }
        },
        onMouseLeave(e : MouseEvent) {
            this.moveNodeActive = false
            this.moveViewBoxActive = false
        },
        onContextMenu(e : Event) {
            e.preventDefault()
        }
    },
})

</script>

<style scoped>

#svgrenderer {
    display: block;
}

</style>
