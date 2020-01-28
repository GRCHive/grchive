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
import VueSetup from '../../../../ts/vueSetup'
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
            if (this.usePropEnd) {
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
            let cPoint1x : number = this.computeControlPointX(this.startPoint.x, halfwayX, this.startIsInput)
            let cPoint2x : number = this.computeControlPointX(this.endPoint.x, halfwayX, this.endIsInput)

            let cPoint1 = <Point2D>{
                x: cPoint1x,
                y: this.startPoint.y
            }

            let cPoint2 = <Point2D>{
                x: cPoint2x,
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
        },
        // Computes the proper control point X so that the edge is visible connecting into the 
        // plug even when the input is to the left of the output. I don't particularly like
        // the math in here since it will create extra long curves if the two nodes are very far apart...
        computeControlPointX(plugX : number, desiredX : number, isInput : boolean) : number {
            if (isInput) {
                if (plugX < desiredX) {
                    return plugX - (desiredX - plugX)
                } else {
                    return desiredX
                }
            } else {
                if (plugX > desiredX) {
                    return plugX + (plugX - desiredX)
                } else {
                    return desiredX
                }
            }
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
    stroke: transparent;
    stroke-opacity: 0%;
    stroke-width: 10px;
    fill: transparent;
}

.tempFlowEdge {
    pointer-events: none;
    user-select: none;
}

</style>
