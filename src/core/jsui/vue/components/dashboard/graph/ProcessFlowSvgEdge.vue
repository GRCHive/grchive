<template>
    <path :d="d"
          class="flowEdge">
    </path>
</template>

<script lang="ts">

import Vue from 'vue'
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
        endIsInput: Boolean
    },
    computed: {
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
    }
})
</script>

<style scoped>

.flowEdge {
    stroke: black;
    fill: transparent;
}

</style>
