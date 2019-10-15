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
                ref="item.Id"
                @onmousedown="onMouseDownNode"
                @onmouseup="onMouseUpNode"
                @onplugmousedown="onPlugMouseDown"
                @onplugmouseup="onPlugMouseUp"
            >
            </process-flow-svg-node>
        </g>

        <g id="edges">
            <!-- One temporary edge that the user sees when they click and drag from one plug to another -->
            <process-flow-svg-edge v-if="drawingEdge">
            </process-flow-svg-edge>
        </g>
    </svg>
</template>

<script lang="ts">

import Vue from 'vue'
import VueSetup from '../../../../ts/vueSetup'
import RenderLayout from '../../../../ts/render/renderLayout'
import ProcessFlowSvgNode from './ProcessFlowSvgNode.vue'
import ProcessFlowSvgEdge from './ProcessFlowSvgEdge.vue'

export default Vue.extend({
    components: {
        ProcessFlowSvgNode,
        ProcessFlowSvgEdge
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
        viewBoxY: 0,
        drawingEdge: false
    }),
    methods: {
        updateTemporaryEdge() {
        },
        saveTemporaryEdge() {
            this.drawingEdge = false
        },
        doMoveNode(e : MouseEvent) {
            if (!VueSetup.store.getters.isNodeSelected) {
                return
            }

            RenderLayout.store.commit('addNodeDisplayTranslation', {
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

            if (e.button == 0) {
                this.drawingEdge = false
            }
        },
        onMouseLeave(e : MouseEvent) {
            this.moveNodeActive = false
            this.moveViewBoxActive = false
            this.drawingEdge = false
        },
        onContextMenu(e : Event) {
            e.preventDefault()
        },
        onPlugMouseDown(e : MouseEvent, nodeId : number, io : ProcessFlowInputOutput, isInput: boolean) {
            if (e.button != 0) {
                return
            }

            e.stopPropagation()
            this.drawingEdge = true
        },
        onPlugMouseUp(e : MouseEvent, nodeId : number, io : ProcessFlowInputOutput, isInput: boolean) {
            if (e.button != 0) {
                return
            }

            e.stopPropagation()
            this.saveTemporaryEdge()
        },
    },
})

</script>

<style scoped>

#svgrenderer {
    display: block;
}

</style>
