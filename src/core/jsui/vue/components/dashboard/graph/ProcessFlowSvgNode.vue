<template>
    <g :id="node.Id.toString()"
       :transform="`translate(${tx}, ${ty})`"
    >
        <rect :width="rectWidth"
              :height="rectHeight"
              :class="styleClass + ` ` + (isNodeSelected ? `node-selected-box` : `node-box`)"
              @mousedown="onMouseDown($event)"
              @mouseup="onMouseUp($event)">
        </rect>
        <g ref="textgroup"
           class="no-pointer"
           :transform="`translate(${margins.left}, ${margins.top})`"
        >
            <text dominant-baseline="hanging"
                  :class="`title ` + styleClass + `-text`"
                  text-rendering="optimizeLegibility"
            >{{ node.Name }}</text> 
        </g>
    </g>
</template>

<script lang="ts">

import Vue from 'vue'
import VueSetup from '../../../../ts/vueSetup'

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
    }),
    methods : {
        onMouseDown(e: MouseEvent) {
            this.$emit("onmousedown", e, this.node.Id)
            e.stopPropagation()
        },
        onMouseUp(e: MouseEvent) {
            this.$emit("onmouseup", e, this.node.Id)
            e.stopPropagation()
        }
    },
    computed: {
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
        }
    },
    mounted() {
        //@ts-ignore
        const textgroup = this.$refs.textgroup
        //@ts-ignore
        this.textWidth = textgroup.getBBox().width
        //@ts-ignore
        this.textHeight = textgroup.getBBox().height
    }
})

</script>

<style scoped>

.no-pointer {
    pointer-events: none;
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
