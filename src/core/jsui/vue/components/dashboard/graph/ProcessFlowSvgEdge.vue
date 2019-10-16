<template>
    <g v-if="ready">
        <path v-if="hoverEdge"
              :d="d"
              class="highlightEdge"> 
        </path>

        <path v-if="isEdgeSelected"
              :d="d"
              class="selectedFlowEdge"> 
        </path>

       <path :d="d"
              class="flowEdge">
        </path>

        <path :d="d"
              class="clickEdge"
              @click="onClick"
              @mouseenter="hoverEdge = true"
              @mouseleave="hoverEdge = false">
        </path>
    </g>
</template>

<script lang="ts">

import Vue from 'vue'
import VueSetup from '../../../../ts/vueSetup.ts'
import RenderLayout from '../../../../ts/render/renderLayout'
export default Vue.extend({
    props: {
        usePropEnd: Boolean,
        propEndX: Number,
        propEndY: Number,
        startNodeId: Number,
        startIo: {
            type: Object as () => ProcessFlowInputOutput,
        },
        startIsInput: Boolean,
        endNodeId: Number,
        endIo: {
            type: Object as () => ProcessFlowInputOutput,
        },
        endIsInput: Boolean,
        edgeId: {
            type: Number,
            default: -1
        }
    },
    data: () => ({
        hoverEdge: false
    }),
    computed: {
        isEdgeSelected(): boolean {
            return VueSetup.store.state.selectedEdgeId == this.edgeId
        },
        ready() : boolean {
            const inputRdy = RenderLayout.store.getters.isReadyForNode(this.startNodeId)
            if (!this.usePropEnd) {
                return inputRdy
            }
            const outputRdy = RenderLayout.store.getters.isReadyForNode(this.endNodeId)
            return inputRdy && outputRdy
        },
        startPoint() : Point2D {
            return RenderLayout.store.getters.getPlugLocation(this.startNodeId, this.startIo, this.startIsInput)
        },
        endPoint() : Point2D {
            if (this.usePropEnd) {
                return <Point2D>{
                    x: this.propEndX,
                    y: this.propEndY
                }
            } else {
                return RenderLayout.store.getters.getPlugLocation(this.endNodeId, this.endIo, this.endIsInput)
            }
        },
        d() : string {
            // Compute bezier control points.
            // The controls points just point directly left/right for simplicity.
            let halfwayX : number = (this.startPoint.x + this.endPoint.x) / 2.0
            let cPoint1 = <Point2D>{
                x: halfwayX,
                y: this.startPoint.y
            }

            let cPoint2 = <Point2D>{
                x: halfwayX,
                y: this.endPoint.y
            }
            return `M ${this.startPoint.x} ${this.startPoint.y}
                    C ${cPoint1.x} ${cPoint1.y},
                      ${cPoint2.x} ${cPoint2.y},
                      ${this.endPoint.x} ${this.endPoint.y}`
        }
    },
    methods: {
        onClick(e : MouseEvent){ 
            this.$emit("onedgeclick", e)
        }
    }
})
</script>

<style scoped>

.flowEdge {
    stroke: black;
    stroke-width: 2px;
    fill: transparent;
}

.selectedFlowEdge {
    stroke: red;
    stroke-width: 4px;
    fill: transparent;
}

.highlightEdge {
    stroke-width: 4px;
    stroke: orange;
    fill: transparent;
}

.clickEdge {
    stroke-opacity: 0%;
    stroke-width: 6px;
    fill: transparent;
}

.tempFlowEdge {
    pointer-events: none;
    user-select: none;
}

</style>
