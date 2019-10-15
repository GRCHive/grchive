<template>
    <g :id="node.Id.toString()"
       :transform="`translate(${tx}, ${ty})`"
       v-if="ready"
       ref="basegroup"
    >
        <rect :width="nodeLayout.boxWidth"
              :height="nodeLayout.boxHeight"
              :class="styleClass + ` ` + (isNodeSelected ? `node-selected-box` : `node-box`)"
              @mousedown="onMouseDown($event)"
              @mouseup="onMouseUp($event)"
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
                      @mousedown="onPlugMouseDown($event, input, true)"
                      @mouseup="onPlugMouseUp($event, input, true)"
                ></rect>

                <rect :width="plugWidth"
                      :height="plugHeight"
                      v-for="(output, iIndex) in nodeLayout.groupLayout[group].relevantOutputs"
                      :key="`output` + iIndex.toString()"
                      :transform="`translate(
                        ${nodeLayout.groupLayout[group].outputLayouts[output.Id].plugTransform.tx},
                        ${nodeLayout.groupLayout[group].outputLayouts[output.Id].plugTransform.ty})`"
                      @mousedown="onPlugMouseDown($event, output, false)"
                      @mouseup="onPlugMouseUp($event, output, false)"
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

export default Vue.extend({
    props: {
        node: {
            type: Object as () => ProcessFlowNode
        }
    },
    data: () => ({
        margins : {
            left: 5,
            right: 5,
            top: 5,
            bottom: 5
        },
        plugHeight: 20,
        plugWidth: 20
    }),
    methods : {
        onMouseDown(e: MouseEvent) {
            this.$emit("onmousedown", e, this.node.Id)
        },
        onMouseUp(e: MouseEvent) {
            this.$emit("onmouseup", e, this.node.Id)
        },
        onPlugMouseDown(e : MouseEvent, io : ProcessFlowInputOutput, isInput: boolean) {
            this.$emit("onplugmousedown", e, this.node.Id, io, isInput)
        },
        onPlugMouseUp(e : MouseEvent, io : ProcessFlowInputOutput, isInput: boolean) {
            this.$emit("onplugmouseup", e, this.node.Id, io, isInput)
        },
    },
    computed: {
        ready() : boolean {
            return RenderLayout.store.state.ready
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
            // TODO: How do we keep this in sync with the server?
            //       Maybe it should get queried along with the types.
            const typeId = (this.node as ProcessFlowNode).NodeTypeId
            switch(typeId){
                case 1:
                    return "activity-manual"
                case 2:
                    return "activity-automated"
                case 3:
                    return "decision"
                case 4:
                    return "start"
                case 5:
                    return "general-ledger-entry"
                case 6:
                    return "system"
                default:
                    break;
            }
            return ""
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

            Vue.nextTick(() => {
                RenderLayout.store.dispatch('associateNodeLayoutWithComponent', {
                    nodeId: this.node.Id,
                    component: this
                })
            })
        }
    }
})

</script>

<style scoped>

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

.activity-manual {
    fill: #001F3F;
}

.activity-automated {
    fill: #0074D9;
}

.decision {
    fill: #7FDBFF;
}

.start {
    fill: #39CCCC;
}

.general-ledger-entry {
    fill: #3D9970;
}

.system {
    fill: #2ECC40;
}

.activity-manual-text {
    fill: white;
}

.activity-automated-text {
    fill: white;
}

.decision-text {
    fill: black;
}

.start-text {
    fill: black;
}

.general-ledger-entry-text {
    fill: white;
}

.system-text {
    fill: black;
}

</style>
