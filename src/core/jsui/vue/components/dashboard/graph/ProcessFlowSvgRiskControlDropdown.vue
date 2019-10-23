<template>
    <g :transform="`translate(0, ${parentHeight})`"
       :visibility="hasRiskControl ? `visible` : `hidden`">
        <rect :width="currentWidth"
              :height="expandedHeight + 2 * dropdownButtonMargin + dropdownButtonHeight"
              class="dropdown-rect"
        ></rect>

        <g ref="riskControlTextGroup" 
           :visibility="isExpanded ? `visible` : `hidden`">

            <text dominant-baseline="hanging"
                  class="body-2 dropdown-text no-pointer"
                  text-rendering="optimizeLegibility"
                  text-anchor="middle"
                  :transform="`translate(${currentWidth / 2}, ${dropdownButtonMargin})`"
                  ref="riskTitle"
            >RISKS</text> 

            <g ref="riskText">
                <text dominant-baseline="hanging"
                      class="body-2 dropdown-text no-pointer"
                      text-rendering="optimizeLegibility"
                      v-for="(item, index) in risks"
                      :key="index"
                      :transform="`translate(
                        ${getRiskLayout(item.Id).tx},
                        ${getRiskLayout(item.Id).ty})`"
                >{{ item.Name }}</text> 
            </g>

            <g ref="controlText">
            </g>
        </g>

        <g ref="dropdown"
           :transform="`translate(0, ${expandedHeight})`"
           @mousedown.stop
           @click.stop="toggleExpand"
           cursor="pointer"
        >
            <text dominant-baseline="hanging"
                  class="body-2 dropdown-text no-pointer"
                  text-rendering="optimizeLegibility"
                  text-anchor="middle"
                  :transform="`translate(${currentWidth / 2}, ${dropdownButtonMargin})`"
            >{{ !isExpanded ? "RISKS AND CONTROLS" : "COLLAPSE" }}</text> 
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
        },
        parentWidth: Number,
        parentHeight: Number
    },
    data : () => ({
        dropdownButtonMargin: 5,
        dropdownButtonHeight: 20,
        riskControlTextHeight: 200,
        isExpanded: false,
        riskTransformLayout: new Map<number, TransformData>(),
        riskTextWidth: 200
    }),
    computed: {
        hasRiskControl() : boolean {
            return this.risks.length > 0
        },
        risks() : ProcessFlowRisk[] {
            return VueSetup.store.getters.risksForNode(this.node.Id)
        },
        expandedHeight() : number {
            if (!this.isExpanded) {
                return 0
            } else {
                return this.riskControlTextHeight
            }
        },
        currentWidth() : number {
            if (!this.isExpanded) {
                return this.parentWidth
            } else {
                return Math.max(this.parentWidth, this.riskTextWidth)
            }
        }
    },
    methods: {
        recomputeRiskControlTextLayout() {
            const itemSpacingMargin : number = 5

            //@ts-ignore
            const riskTitle: SVGGraphicsElement = this.$refs.riskTitle
            let riskStartY = riskTitle.getBBox().y + riskTitle.getBBox().height

            this.riskTransformLayout = new Map<number, TransformData>()
            for (let r of this.risks) {
                riskStartY += itemSpacingMargin
                this.riskTransformLayout.set(r.Id, <TransformData>{
                    tx: itemSpacingMargin,
                    ty: riskStartY
                })
                // This is OK since we don't actually change the font
                riskStartY += riskTitle.getBBox().height
            }

            Vue.nextTick(() => {
                //@ts-ignore
                const totalGroup: SVGGraphicsElement = this.$refs.riskControlTextGroup
                this.riskControlTextHeight = totalGroup.getBBox().height

                //@ts-ignore
                const riskText : SVGGraphicsElement = this.$refs.riskText
                this.riskTextWidth = riskText.getBBox().width + 2 * itemSpacingMargin
            })
        },
        toggleExpand() {
            this.isExpanded = !this.isExpanded
        },
        getRiskLayout(riskId : number) : TransformData {
            if (!this.riskTransformLayout.has(riskId)) {
                return <TransformData>{
                    tx: 0,
                    ty: 0
                }
            } else {
                return this.riskTransformLayout.get(riskId)!
            }
        }
    },
    mounted() {
        //@ts-ignore
        const group : SVGGraphicsElement = this.$refs.dropdown
        this.dropdownButtonHeight = group.getBBox().height
        this.recomputeRiskControlTextLayout()
    },
    watch : {
        risks() {
            this.recomputeRiskControlTextLayout()
        },
    }
})
</script>

<style scoped>

.dropdown-rect {
    fill: black;
    fill-opacity: 80%;
}

.dropdown-text {
    fill: white;
}

.no-pointer {
    user-select: none;
}

</style>
