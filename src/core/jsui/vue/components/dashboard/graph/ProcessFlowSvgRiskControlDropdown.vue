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
                  class="body-2 dropdown-text no-pointer font-weight-bold"
                  text-rendering="optimizeLegibility"
                  text-anchor="start"
                  :transform="`translate(${dropdownButtonMargin}, ${dropdownButtonMargin})`"
                  ref="riskTitle"
            >RISKS</text> 

            <text dominant-baseline="hanging"
                  class="body-2 dropdown-text no-pointer font-weight-bold"
                  text-rendering="optimizeLegibility"
                  text-anchor="end"
                  :transform="`translate(${currentWidth - dropdownButtonMargin}, ${dropdownButtonMargin})`"
                  ref="controlTitle"
            >CONTROLS</text> 

            <g ref="riskText">
                <a v-for="(item, index) in risks"
                   :key="index"
                   :href="generateRiskUrl(item)"
                   target="_blank"
                >
                    <text dominant-baseline="hanging"
                          class="body-2 dropdown-text no-pointer"
                          text-rendering="optimizeLegibility"
                          :transform="`translate(
                            ${getRiskLayout(item.Id).tx},
                            ${getRiskLayout(item.Id).ty})`"
                    >{{ item.Name }}</text> 
                </a>
            </g>

            <g ref="controlText">
                <g v-for="risk in risks"
                   :key="risk.Id"
                >
                    <a v-for="(control, cindex) in controls[risk.Id]"
                       :key="cindex"
                       :href="generateControlUrl(control.control)"
                       target="_blank"
                    >
                        <text dominant-baseline="hanging"
                              class="body-2 dropdown-text no-pointer"
                              text-rendering="optimizeLegibility"
                              text-anchor="end"
                              :transform="`translate(
                                ${currentWidth - getControlLayout(control.control.Id).tx},
                                ${getControlLayout(control.control.Id).ty})`"
                        >{{ control.control.Name }}</text> 
                    </a>
                </g>
            </g>
        </g>

        <g :transform="`translate(0, ${expandedHeight})`"
           @mousedown.stop
           @click.stop="toggleExpand"
           cursor="pointer"
        >
            <rect :width="currentWidth"
                  :height="dropdownButtonHeight + 2 * dropdownButtonMargin"
                  class="button-rect"
            ></rect>

            <text dominant-baseline="hanging"
                  class="body-2 dropdown-text no-pointer font-weight-bold"
                  text-rendering="optimizeLegibility"
                  text-anchor="middle"
                  :transform="`translate(${currentWidth / 2}, ${dropdownButtonMargin})`"
                  ref="dropdown"
            >{{ !isExpanded ? "RISKS AND CONTROLS" : "COLLAPSE" }}</text> 
        </g>
    </g>

</template>

<script lang="ts">
import Vue from 'vue'
import VueSetup from '../../../../ts/vueSetup'
import { createRiskUrl, createControlUrl } from '../../../../ts/url'
import { PageParamsStore } from '../../../../ts/pageParams'

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
        riskTextWidth: 200,
        controlTextWidth: 200,
        controlTransformLayout: new Map<number, TransformData>()
    }),
    computed: {
        hasRiskControl() : boolean {
            return this.risks.length > 0
        },
        risks() : ProcessFlowRisk[] {
            return VueSetup.store.getters.risksForNode(this.node.Id)
        },
        controls() : Record<number, RiskControl[]> {
            let groupedControls = Object() as Record<number, RiskControl[]>
            for (let risk of this.risks) {
                groupedControls[risk.Id] = VueSetup.store.getters.controlsForRiskNode(
                    risk.Id,
                    this.node.Id)
            }
            return groupedControls
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
                return Math.max(this.parentWidth, this.riskTextWidth + this.controlTextWidth + 50)
            }
        }
    },
    methods: {
        recomputeRiskControlTextLayout() {
            const itemSpacingMargin : number = 9

            //@ts-ignore
            const riskTitle: SVGGraphicsElement = this.$refs.riskTitle
            let currentY = riskTitle.getBBox().y + riskTitle.getBBox().height

            this.riskTransformLayout = new Map<number, TransformData>()
            for (let r of this.risks) {
                currentY += itemSpacingMargin
                this.riskTransformLayout.set(r.Id, <TransformData>{
                    tx: itemSpacingMargin,
                    ty: currentY
                })

                if (r.Id in this.controls && this.controls[r.Id].length > 0) {
                    for (let c of this.controls[r.Id]) {
                        this.controlTransformLayout.set(c.control.Id, <TransformData>{
                            tx: itemSpacingMargin,
                            ty: currentY
                        })

                        currentY += itemSpacingMargin + riskTitle.getBBox().height
                    }
                    currentY -= riskTitle.getBBox().height
                }
                // This is OK since we don't actually change the font
                currentY += riskTitle.getBBox().height
            }

            Vue.nextTick(() => {
                //@ts-ignore
                const riskText : SVGGraphicsElement = this.$refs.riskText
                this.riskTextWidth = riskText.getBBox().width + 2 * itemSpacingMargin

                //@ts-ignore
                const controlText : SVGGraphicsElement = this.$refs.controlText
                this.controlTextWidth = controlText.getBBox().width + 2 * itemSpacingMargin

                Vue.nextTick(() => {
                    //@ts-ignore
                    const totalGroup: SVGGraphicsElement = this.$refs.riskControlTextGroup
                    this.riskControlTextHeight = totalGroup.getBBox().height
                })
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
        },
        getControlLayout(controlId: number) : TransformData {
            if (!this.controlTransformLayout.has(controlId)) {
                return <TransformData>{
                    tx: 0,
                    ty: 0
                }
            } else {
                return this.controlTransformLayout.get(controlId)!
            }
        },
        generateRiskUrl(risk : ProcessFlowRisk) {
            return createRiskUrl(
                PageParamsStore.state.organization!.OktaGroupName,
                risk.Id)
        },
        generateControlUrl(control : ProcessFlowControl) {
            return createControlUrl(
                PageParamsStore.state.organization!.OktaGroupName,
                control.Id)
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
        controls() {
            this.recomputeRiskControlTextLayout()
        }
    }
})
</script>

<style scoped>

.button-rect {
    fill: transparent;
}

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
