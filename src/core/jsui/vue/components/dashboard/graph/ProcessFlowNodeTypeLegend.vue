<template>
    <g>
        <rect :width="currentWidth"
              :height="currentHeight"
              class="bg-rect"
        ></rect>

        <g id="legendContent" ref="legendContent"
           :transform="`translate(${margins.left}, ${margins.top})`">
            <g v-for="(item, index) in rawTypeOptions"
               :key="index"
            >
                <rect width="16" height="16"
                      :class="`${nodeTypeToClass(item.Id)} legend-box`"
                      :transform="`translate(0, ${index * 30})`"
                ></rect>
                <text dominant-baseline="middle"
                      :transform="`translate(24, ${index * 30 + 8})`"
                      class="bg-text"
                >
                    {{ item.Name }}
                </text>
            </g>
        </g>
    </g>
</template>

<script lang="ts">

import Vue from 'vue'
import MetadataStore from '../../../../ts/metadata'
import { nodeTypeToClass } from '../../../../ts/render/nodeCssUtils'

export default Vue.extend({
    data : () => ({
        currentWidth: 200,
        currentHeight: 200,
        margins: {
            top: 5,
            left: 5,
            bottom: 5,
            right: 5
        }
    }),
    computed: {
        rawTypeOptions() : ProcessFlowNodeType[] {
            return MetadataStore.state.nodeTypes
        },
    },
    methods: {
        nodeTypeToClass: nodeTypeToClass,
        updateWidthHeight() {
            const contentSvg : SVGSVGElement = <SVGSVGElement>(this.$refs.legendContent)
            this.currentWidth =
                this.margins.left +
                this.margins.right +
                contentSvg.getBBox().width;

            this.currentHeight =
                this.margins.top +
                this.margins.bottom +
                contentSvg.getBBox().height;
        }
    },
    mounted() {
        this.updateWidthHeight()
    },
    watch : {
        rawTypeOptions() {
            this.updateWidthHeight()
        }
    }
})
</script>

<style scoped>

#legend-content {
    user-select: none;
}

.legend-box {
    stroke: white;
    stroke-width: 1px;
}

.bg-rect {
    fill: black;
    fill-opacity: 80%;
}

.bg-text {
    fill: white;
}

</style>
