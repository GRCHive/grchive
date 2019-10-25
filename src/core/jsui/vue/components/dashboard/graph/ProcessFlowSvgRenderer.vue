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
        <g id="fullgraph" ref="fullgraph">
            <g id="edges">
                <!-- One temporary edge that the user sees when they click and drag from one plug to another -->
                <process-flow-svg-edge v-if="drawingEdge"
                                       :use-prop-end="true"
                                       :prop-end-x="tempEdgeEnd.x"
                                       :prop-end-y="tempEdgeEnd.y"
                                       :start-node-id="tempEdgeStart.nodeId"
                                       :start-io="tempEdgeStart.io"
                                       :start-is-input="tempEdgeStart.isInput"
                ></process-flow-svg-edge>

                <process-flow-svg-edge 
                    :use-prop-end="false"
                    v-for="key in edgeKeys"
                    :key="key"
                    :start-node-id="getInputOutputFromId(edges[key].InputIoId, true).ParentNodeId"
                    :start-io="getInputOutputFromId(edges[key].InputIoId, true)"
                    :start-is-input="true"
                    :end-node-id="getInputOutputFromId(edges[key].OutputIoId, false).ParentNodeId"
                    :end-io="getInputOutputFromId(edges[key].OutputIoId, false)"
                    :end-is-input="false"
                    :edge-id="key"
                    @onedgeclick="onEdgeClick(arguments[0], key)"
                ></process-flow-svg-edge>
            </g>

            <g id="nodes">
                <process-flow-svg-node
                    v-for="key in nodeKeys"
                    :key="nodes[key].Id"
                    :node="nodes[key]"
                    ref="nodes[key].Id"
                    @onnodemousedown="onMouseDownNode"
                    @onnodemouseup="onMouseUpNode"
                    @onplugmousedown="onMouseDownPlug"
                    @onplugmouseup="onMouseUpPlug"
                >
                </process-flow-svg-node>
            </g>
        </g>
    </svg>
</template>

<script lang="ts">

import Vue from 'vue'
import VueSetup from '../../../../ts/vueSetup'
import RenderLayout from '../../../../ts/render/renderLayout'
import ProcessFlowSvgNode from './ProcessFlowSvgNode.vue'
import ProcessFlowSvgEdge from './ProcessFlowSvgEdge.vue'
import { newProcessFlowEdge } from '../../../../ts/api/apiProcessFlowEdges'
import { contactUsUrl } from '../../../../ts/url'
import LocalSettings from '../../../../ts/localSettings'
import { convertClientPointToSvg, convertClientDeltaToSvg } from '../../../../ts/svg'

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
        nodes() : Record<number,ProcessFlowNode> {
            return VueSetup.store.state.currentProcessFlowFullData.Nodes
        },
        edges() : Record<number, ProcessFlowEdge> {
            return VueSetup.store.state.currentProcessFlowFullData.Edges
        },
        nodeKeys() : number[] {
            return VueSetup.store.state.currentProcessFlowFullData.NodeKeys
        },
        edgeKeys() : number[] {
            return VueSetup.store.state.currentProcessFlowFullData.EdgeKeys
        },
        viewBoxX() : number {
            return LocalSettings.state.viewBoxTransform.tx
        },
        viewBoxY() : number {
            return LocalSettings.state.viewBoxTransform.ty
        },
        viewBoxZoom() : number {
            return LocalSettings.state.viewBoxZoom
        },
        viewBox() : BoundingBox {
            return {
                x: this.viewBoxX,
                y: this.viewBoxY,
                width: this.svgWidth / this.viewBoxZoom,
                height: this.svgHeight / this.viewBoxZoom
            }
        }
    },
    data: () => ({
        moveNodeActive: false,
        moveViewBoxActive: false,
        // Edge drawing properties
        drawingEdge: false,
        tempEdgeStart: {
            nodeId: -1,
            io: {} as ProcessFlowInputOutput,
            isInput: false
        },
        tempEdgeEnd: {
            x: 0,
            y: 0
        }
    }),
    methods: {
        getInputOutputFromId(ioId : number, isInput: boolean): ProcessFlowInputOutput {
            if (isInput) {
                return VueSetup.store.state.currentProcessFlowFullData.Inputs[ioId]
            } else {
                return VueSetup.store.state.currentProcessFlowFullData.Outputs[ioId]
            }
        },
        saveTemporaryEdge(endIo: ProcessFlowInputOutput, endIsInput: boolean) {
            this.drawingEdge = false
            // Need to make sure we connect an input to an output
            if (endIsInput == this.tempEdgeStart.isInput) {
                //@ts-ignore
                this.$root.$refs.snackbar.showSnackBar(
                    "You must connect an input to an output.",
                    false,
                    "",
                    "",
                    true);
                return
            }

            // Can't connect to the same node
            if (endIo.ParentNodeId == this.tempEdgeStart.nodeId) {
                //@ts-ignore
                this.$root.$refs.snackbar.showSnackBar(
                    "You can not connect an edge from a node to itself.",
                    false,
                    "",
                    "",
                    true);
                return
            }

            newProcessFlowEdge(<TNewProcessFlowEdgeInput>{
                //@ts-ignore
                csrf: this.$root.csrf,
                inputIoId: endIsInput ? endIo.Id : this.tempEdgeStart.io.Id,
                outputIoId: endIsInput ? this.tempEdgeStart.io.Id : endIo.Id
            }).then((resp : TNewProcessFlowEdgeOutput) => {
                this.drawingEdge = false
                VueSetup.store.commit('addNewEdge', {edge: resp.data})
            }).catch((err) => {
                console.log(err)
                //@ts-ignore
                this.$root.$refs.snackbar.showSnackBar(
                    "Oops! Something went wrong, please reload the page and try again.",
                    true,
                    "Contact Us",
                    contactUsUrl,
                    true);
            })
        },
        endMoveNode(nodeId: number) {
            if (this.moveNodeActive && nodeId != -1) {
                RenderLayout.store.dispatch('syncNodeTransform', {nodeId})
            }
            this.moveNodeActive = false
        },
        doMoveNode(e : MouseEvent) {
            if (!VueSetup.store.getters.isNodeSelected) {
                return
            }

            let svg : SVGSVGElement = <SVGSVGElement>this.$refs.svgrenderer
            let delta : SVGPoint = convertClientDeltaToSvg(svg, e.movementX, e.movementY)

            RenderLayout.store.commit('addNodeDisplayTranslation', {
                nodeId: VueSetup.store.state.selectedNodeId,
                tx: delta.x,
                ty: delta.y
            })
        },
        doMoveViewBox(e : MouseEvent) {
            let svg : SVGSVGElement = <SVGSVGElement>this.$refs.svgrenderer
            let delta : SVGPoint = convertClientDeltaToSvg(svg, e.movementX, e.movementY)
            LocalSettings.commit('setViewBoxTransform', {
                tx: this.viewBoxX - delta.x,
                ty: this.viewBoxY - delta.y,
            })
        },
        doMoveTempEdgeEnd(e: MouseEvent) {
            let svg : SVGSVGElement = <SVGSVGElement>this.$refs.svgrenderer
            let realPt = convertClientPointToSvg(svg, e.clientX, e.clientY)
            this.tempEdgeEnd.x = realPt.x
            this.tempEdgeEnd.y = realPt.y
        },
        onMouseMove(e : MouseEvent) {
            if (this.moveNodeActive) {
                this.doMoveNode(e)
            } else if (this.moveViewBoxActive) {
                this.doMoveViewBox(e)
            } else if (this.drawingEdge) {
                this.doMoveTempEdgeEnd(e)
            }
        },
        onMouseDownNode(e : MouseEvent, nodeId : number) {
            if (e.button != 0) {
                return
            }

            VueSetup.store.commit('setSelectedProcessFlowNode', nodeId)
            VueSetup.store.commit('setSelectedProcessFlowEdge', -1)
            this.moveNodeActive = true
        },
        onMouseUpNode(e : MouseEvent, nodeId : number) {
            if (e.button != 0) {
                return
            }

            this.endMoveNode(nodeId)
            this.drawingEdge = false
        },
        onMouseDown(e : MouseEvent) {
            if (e.button == 0) {
                VueSetup.store.commit('setSelectedProcessFlowNode', -1)
                VueSetup.store.commit('setSelectedProcessFlowEdge', -1)
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
            this.endMoveNode(VueSetup.store.state.selectedNodeId) 
            this.moveViewBoxActive = false
            this.drawingEdge = false
        },
        onContextMenu(e : Event) {
            e.preventDefault()
        },
        onMouseDownPlug(e : MouseEvent, nodeId : number, io : ProcessFlowInputOutput, isInput: boolean) {
            if (e.button != 0) {
                return
            }

            VueSetup.store.commit('setSelectedProcessFlowEdge', -1)
            this.drawingEdge = true
            this.tempEdgeStart.nodeId = nodeId
            this.tempEdgeStart.io = io
            this.tempEdgeStart.isInput = isInput
            this.doMoveTempEdgeEnd(e)
        },
        onMouseUpPlug(e : MouseEvent, nodeId : number, io : ProcessFlowInputOutput, isInput: boolean) {
            if (e.button != 0) {
                return
            }

            this.saveTemporaryEdge(io, isInput)
        },
        onEdgeClick(e : MouseEvent, edgeId: number) {
            VueSetup.store.commit('setSelectedProcessFlowEdge', edgeId)
            VueSetup.store.commit('setSelectedProcessFlowNode', -1)
            e.stopPropagation()
        }
    },
    mounted() {
        RenderLayout.store.commit('updateFullGraphComponent', this.$refs.fullgraph)
    }
})

</script>

<style scoped>

#svgrenderer {
    display: block;
}

</style>
