<template>
    <g :id="node.Id.toString()"
       :transform="`translate(${tx}, ${ty})`"
    >
        <rect :width="rectWidth"
              :height="rectHeight"
              :class="styleClass + ` ` + (isNodeSelected ? `node-selected-box` : `node-box`)"
              @mousedown="onMouseDown($event)"
              @mouseup="onMouseUp($event)"
        ></rect>
        <g ref="textgroup"
           class="no-pointer"
           :transform="`translate(${margins.left}, ${margins.top})`"
        >
            <text dominant-baseline="hanging"
                  :class="`title ` + styleClass + `-text`"
                  text-rendering="optimizeLegibility"
                  ref="title"
            >{{ node.Name }}</text> 

            <g v-for="(group, index) in groupedInputOutputs" :key="index"
               :transform="`translate(${inputOutputDisplay[index].groupTransform.tx}, ${inputOutputDisplay[index].groupTransform.ty})`">
                <text dominant-baseline="hanging"
                      :class="`subtitle-1 font-weight-bold ` + styleClass + `-text`"
                      text-rendering="optimizeLegibility">
                    {{ group.name }}
                </text>
                
                <text dominant-baseline="hanging"
                      :class="`body-2 ` + styleClass + `-text`"
                      text-rendering="optimizeLegibility"
                      v-for="(input, iIndex) in group.inputs"
                      :key="`input` + iIndex.toString()"
                      :transform="`translate(
                        ${inputOutputDisplay[index].inputTransforms[iIndex].tx},
                        ${inputOutputDisplay[index].inputTransforms[iIndex].ty})`">
                    {{ input.Name }}
                </text>

                <text dominant-baseline="hanging"
                      :class="`body-2 ` + styleClass + `-text`"
                      text-rendering="optimizeLegibility"
                      v-for="(output, oIndex) in group.outputs"
                      :key="`output` + oIndex.toString()"
                      :transform="`translate(
                        ${inputOutputDisplay[index].outputTransforms[oIndex].tx},
                        ${inputOutputDisplay[index].outputTransforms[oIndex].ty})`">
                    {{ output.Name }}
                </text>
            </g>
        </g>

        <g ref="ioPlugs">
            <g v-for="(group, index) in groupedInputOutputs" :key="index"
               :transform="`translate(${inputOutputDisplay[index].groupTransform.tx}, ${inputOutputDisplay[index].groupTransform.ty})`">
                <rect :width="plugWidth"
                      :height="plugHeight"
                      v-for="(input, iIndex) in group.inputs"
                      :key="`input` + iIndex.toString()"
                      :transform="`translate(
                        ${inputOutputDisplay[index].inputTransforms[iIndex].tx - plugWidth},
                        ${inputOutputDisplay[index].inputTransforms[iIndex].ty + 5}) `">
                </rect>

                <rect :width="plugWidth"
                      :height="plugHeight"
                      v-for="(output, oIndex) in group.outputs"
                      :key="`output` + oIndex.toString()"
                      :transform="`translate(
                        ${rectWidth},
                        ${inputOutputDisplay[index].outputTransforms[oIndex].ty + 5}) `">
                </rect>
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
        textHeight: 200,
        textWidth: 200,
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
        updateTextHeightWidth() {
            //@ts-ignore
            const textgroup = this.$refs.textgroup

            //@ts-ignore
            this.textWidth = textgroup.getBBox().width
            //@ts-ignore
            this.textHeight = textgroup.getBBox().height
        }
    },
    computed: {
        groupedInputOutputs() : GroupedInputOutputData[] {
            let inputs = VueSetup.store.getters.nodeInfo(this.node.Id).Inputs
            let outputs = VueSetup.store.getters.nodeInfo(this.node.Id).Outputs

            // Sort the inputs and outputs by their type.
            let groupMap = new Map()

            function addToGroupMap(io : ProcessFlowInputOutput, isInput : boolean) {
                const typeId : number = io.TypeId
                const key : string = MetadataStore.state.idToIoTypes[typeId].Name

                if (!groupMap.has(key)) {
                    groupMap.set(key, <GroupedInputOutputData>{
                        name: key,
                        inputs: [] as ProcessFlowInputOutput[],
                        outputs: [] as ProcessFlowInputOutput[]
                    })
                }

                if (isInput) {
                    groupMap.get(key).inputs.push(io)
                } else {
                    groupMap.get(key).outputs.push(io)
                }
            }

            for (let i = 0; i < inputs.length; ++i) {
                addToGroupMap(inputs[i], true)
            }

            for (let i = 0; i < outputs.length; ++i) {
                addToGroupMap(outputs[i], false)
            }

            // Each object in the 'grouped' array will have the form: 
            // {
            //     name: TYPE_NAME (e.g. Execution/Data)
            //     inputs: ProcessFlowInputOutput[]
            //     outputs: ProcessFlowInputOutput[]
            // }
            let grouped = [] as GroupedInputOutputData[]
            groupMap.forEach((value) => {
                grouped.push(value)
            })

            return grouped
        },
        inputOutputDisplay() : GroupedInputOutputDisplayData[] {
            let displayData = [] as GroupedInputOutputDisplayData[]
        
            const titleHeight: number = 26
            const subtitleHeight : number = 21
            const bodyHeight: number = 19
            const ioMargins = {
                betweenGroups: 5,
                betweenPlugs: 10
            }
            const inputOutputGap : number = 200

            let startY : number = titleHeight + ioMargins.betweenGroups

            // We make assumptions about the height of text for simplicity.
            // Can we have the height of the text be reactive? It seems hard/annoying.
            // We can probably use DOM selectors but eh. Go for hard-coded for now.
            this.groupedInputOutputs.forEach((group, i) => {
                let groupDisplay = <GroupedInputOutputDisplayData>{
                    groupTransform: <TransformData>{
                        tx: 0,
                        ty: startY
                    },
                    inputTransforms: [] as TransformData[],
                    outputTransforms: [] as TransformData[],
                    bbox: <BoundingBox>{
                        x: 0,
                        y: 0,
                        width: 0,
                        height: 0
                    }
                }

                let inputStartY = subtitleHeight + ioMargins.betweenPlugs
                for (let i = 0; i < group.inputs.length; ++i) {
                    groupDisplay.inputTransforms.push(<TransformData>{
                        tx: 0,
                        ty: inputStartY
                    })
                    inputStartY += bodyHeight + ioMargins.betweenPlugs
                }

                let outputStartY = subtitleHeight + ioMargins.betweenPlugs
                for (let i = 0; i < group.outputs.length; ++i) {
                    groupDisplay.outputTransforms.push(<TransformData>{
                        tx: inputOutputGap,
                        ty: outputStartY
                    })
                    outputStartY += bodyHeight + ioMargins.betweenPlugs
                }

                startY += Math.max(inputStartY, outputStartY)

                displayData.push(groupDisplay)

                startY += ioMargins.betweenGroups
            })

            Vue.nextTick(() => {
                this.updateTextHeightWidth()
            })
            return displayData
        },
        nodeDisplaySettings() : ProcessFlowNodeDisplay {
            return VueSetup.store.getters.findNodeDisplayData(this.node.Id)
        },
        tx() : number {
            return this.nodeDisplaySettings.Tx
        },
        ty() : number {
            return this.nodeDisplaySettings.Ty
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
        rectWidth() {
            //@ts-ignore
            return this.margins.left + 
            //@ts-ignore
                this.margins.right +
            //@ts-ignore
                this.textWidth
        },
        rectHeight() {
            //@ts-ignore
            return this.margins.top + 
            //@ts-ignore
                this.margins.bottom +
            //@ts-ignore
                this.textHeight
        },
        isNodeSelected() : boolean {
            return VueSetup.store.state.selectedNodeId == this.node.Id
        },
    },
    mounted() {
        this.updateTextHeightWidth()
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
