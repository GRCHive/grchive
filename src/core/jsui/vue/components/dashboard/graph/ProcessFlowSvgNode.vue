<template>
    <g :id="node.Id.toString()"
       :transform="`translate(${tx}, ${ty})`"
       v-if="ready"
       ref="basegroup"
       class="node"
    >
        <g>
            <process-flow-svg-risk-control-dropdown
                :node="node"
                :parent-width="nodeLayout.boxWidth"
                :parent-height="nodeLayout.boxHeight">
            </process-flow-svg-risk-control-dropdown>
        </g>

        <rect :width="nodeLayout.boxWidth"
              :height="nodeLayout.boxHeight"
              :class="styleClass + ` ` + 
                (isNodeSelected ? `node-selected-box` : 
                    (hoverNode ? `highlight` : `node-box`))"
              @mousedown="onMouseDownNode($event)"
              @mouseup="onMouseUpNode($event)"
              @mouseenter="hoverNode = true"
              @mouseleave="hoverNode = false"
        ></rect>
        <g ref="textgroup"
           class="no-pointer"
           :transform="`translate(${nodeLayout.titleTransform.tx}, ${nodeLayout.titleTransform.ty})`"
        >
            <text dominant-baseline="hanging"
                  :class="`title ` + styleClass + `-text`"
                  text-rendering="optimizeLegibility"
                  ref="title"
            >{{ node.Name }}</text> 

            <g v-for="(group, index) in nodeLayout.groupKeys" :key="index"
               :transform="`translate(
                ${nodeLayout.groupLayout[group].transform.tx},
                ${nodeLayout.groupLayout[group].transform.ty})`">
                <text dominant-baseline="hanging"
                      :class="`subtitle-1 font-weight-bold ` + styleClass + `-text`"
                      text-rendering="optimizeLegibility"
                      :transform="`translate(
                       ${nodeLayout.groupLayout[group].titleTransform.tx},
                       ${nodeLayout.groupLayout[group].titleTransform.ty})`">
                    {{ group }}
                </text>
                
                <text dominant-baseline="hanging"
                      :class="`body-2 ` + styleClass + `-text`"
                      text-rendering="optimizeLegibility"
                      v-for="(input, iIndex) in nodeLayout.groupLayout[group].relevantInputs"
                      :key="`input` + iIndex.toString()"
                      :transform="`translate(
                        ${nodeLayout.groupLayout[group].inputLayouts[input.Id].textTransform.tx},
                        ${nodeLayout.groupLayout[group].inputLayouts[input.Id].textTransform.ty})`">
                    {{ input.Name }}
                </text>

                <text dominant-baseline="hanging"
                      :class="`body-2 ` + styleClass + `-text`"
                      text-rendering="optimizeLegibility"
                      text-anchor="end"
                      v-for="(output, oIndex) in nodeLayout.groupLayout[group].relevantOutputs"
                      :key="`output` + oIndex.toString()"
                      :transform="`translate(
                        ${nodeLayout.groupLayout[group].outputLayouts[output.Id].textTransform.tx},
                        ${nodeLayout.groupLayout[group].outputLayouts[output.Id].textTransform.ty})`">
                    {{ output.Name }}
                </text>
            </g>
        </g>

        <g ref="ioPlugs">
            <g v-for="(group, index) in nodeLayout.groupKeys" :key="index"
               :transform="`translate(
                ${nodeLayout.groupLayout[group].transform.tx},
                ${nodeLayout.groupLayout[group].transform.ty})`">
                <rect :width="plugWidth"
                      :height="plugHeight"
                      v-for="(input, iIndex) in nodeLayout.groupLayout[group].relevantInputs"
                      :key="`input` + iIndex.toString()"
                      :transform="`translate(
                        ${nodeLayout.groupLayout[group].inputLayouts[input.Id].plugTransform.tx},
                        ${nodeLayout.groupLayout[group].inputLayouts[input.Id].plugTransform.ty})`"
                      @mousedown="onMouseDownPlug($event, input, true)"
                      @mouseup="onMouseUpPlug($event, input, true)"
                      @mouseenter="hoverInputPlugId = input.Id"
                      @mouseleave="hoverInputPlugId = -1"
                      :class="(hoverInputPlugId == input.Id ? `highlight` : ``)"
                ></rect>

                <rect :width="plugWidth"
                      :height="plugHeight"
                      v-for="(output, iIndex) in nodeLayout.groupLayout[group].relevantOutputs"
                      :key="`output` + iIndex.toString()"
                      :transform="`translate(
                        ${nodeLayout.groupLayout[group].outputLayouts[output.Id].plugTransform.tx},
                        ${nodeLayout.groupLayout[group].outputLayouts[output.Id].plugTransform.ty})`"
                      @mousedown="onMouseDownPlug($event, output, false)"
                      @mouseup="onMouseUpPlug($event, output, false)"
                      @mouseenter="hoverOutputPlugId = output.Id"
                      @mouseleave="hoverOutputPlugId = -1"
                      :class="(hoverOutputPlugId == output.Id ? `highlight` : ``)"
                ></rect>
            </g>
        </g>
    </g>
</template>

<script lang="ts">

interface GroupedInputOutputData {
    name : string,
    inputs: ProcessFlowInputOutput[],
    outputs: ProcessFlowInputOutput[]
}

interface GroupedInputOutputDisplayData {
    groupTransform: TransformData,
    // These transforms are relative to the group transform.
    inputTransforms: TransformData[],
    outputTransforms: TransformData[]
    bbox: BoundingBox
}

import Vue from 'vue'
import VueSetup from '../../../../ts/vueSetup'
import MetadataStore from '../../../../ts/metadata'
import RenderLayout from '../../../../ts/render/renderLayout'
import ProcessFlowSvgRiskControlDropdown from './ProcessFlowSvgRiskControlDropdown.vue'
import { nodeTypeToClass } from '../../../../ts/render/nodeCssUtils'

export default Vue.extend({
    props: {
        node: {
            type: Object as () => ProcessFlowNode
        }
    },
    components: {
        ProcessFlowSvgRiskControlDropdown
    },
    data: () => ({
        hoverNode: false,
        hoverInputPlugId: -1,
        hoverOutputPlugId: -1
    }),
    methods : {
        onMouseDownNode(e: MouseEvent) {
            this.$emit("onnodemousedown", e, this.node.Id)
            e.stopPropagation()
        },
        onMouseUpNode(e: MouseEvent) {
            this.$emit("onnodemouseup", e, this.node.Id)
            e.stopPropagation()
        },
        onMouseDownPlug(e : MouseEvent, io : ProcessFlowInputOutput, isInput: boolean) {
            this.$emit("onplugmousedown", e, this.node.Id, io, isInput)
            e.stopPropagation()
        },
        onMouseUpPlug(e : MouseEvent, io : ProcessFlowInputOutput, isInput: boolean) {
            this.$emit("onplugmouseup", e, this.node.Id, io, isInput)
            e.stopPropagation()
        },
        reassociateComponent() {
            if (!this.ready) {
                return
            }
            Vue.nextTick(() => {
                RenderLayout.store.dispatch('associateNodeLayoutWithComponent', {
                    nodeId: this.node.Id,
                    component: this
                })
            })
        }
    },
    computed: {
        plugHeight() : number {
            return RenderLayout.params.plugHeight
        },
        plugWidth() : number {
            return RenderLayout.params.plugWidth
        },
        ready() : boolean {
            return RenderLayout.store.getters.isReadyForNode(this.node.Id)
        },
        nodeLayout() : NodeLayout {
            return RenderLayout.store.getters.nodeLayout(this.node.Id)
        },
        tx() : number {
            return this.nodeLayout.transform.tx
        },
        ty() : number {
            return this.nodeLayout.transform.ty
        },
        styleClass() {
            return nodeTypeToClass(this.node.NodeTypeId)
        },
        isNodeSelected() : boolean {
            return VueSetup.store.state.selectedNodeId == this.node.Id
        },
    },
    watch: {
        ready(val: boolean) {
            if (!val) {
                return
            }
            this.reassociateComponent()
        },
        nodeLayout() {
            this.reassociateComponent()
        }
    },
    mounted() {
        if (this.ready) {
            this.reassociateComponent()
        }
    }
})

</script>

<style scoped>

.highlight {
    stroke-width: 3px;
    stroke: orange;
}

.no-pointer {
    pointer-events: none;
    user-select: none;
}

.node-selected-box {
    stroke-width: 3px;
    stroke: red;
}

.node-box {
    stroke-width: 1px;
    stroke: black;
}

</style>
