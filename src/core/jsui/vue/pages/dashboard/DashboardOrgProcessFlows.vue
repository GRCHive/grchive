<template>
    <section class="max-height">
        <dashboard-app-bar>
        </dashboard-app-bar>

        <dashboard-home-page-nav-bar :selected-page="1"></dashboard-home-page-nav-bar>
        <process-flows-nav-bar></process-flows-nav-bar>
        <v-content class="max-height" ref="sectionDiv">
            <process-flow-editor @on-change="recomputeProcessFlowHeaderHeight"></process-flow-editor>
            <v-divider></v-divider>
            <process-flow-toolbar></process-flow-toolbar>
            <v-divider ref="headerDivider"></v-divider>
            <process-flow-renderer></process-flow-renderer>

            <process-flow-attribute-editor :custom-clip-height="headerClipHeight" 
                                           ref="attrEditor"
                                           :show-hide="showHideAttrEditor"
            ></process-flow-attribute-editor>
        </v-content>

        <v-btn color="primary"
               id="attrPullButton"
               small
               :style="attributePullButtonStyle"
               class="no-transition"
               @click="clickAttributePullTab"
        >
            <v-icon v-if="!showHideAttrEditor">mdi-chevron-up</v-icon>
            <v-icon v-else>mdi-chevron-down</v-icon>
        </v-btn>
    </section>
</template>

<script lang="ts">

import DashboardAppBar from '../../components/dashboard/DashboardAppBar.vue'
import DashboardHomePageNavBar from '../../components/dashboard/DashboardHomePageNavBar.vue'
import ProcessFlowsNavBar from '../../components/dashboard/ProcessFlowsNavBar.vue'
import ProcessFlowEditor from '../../components/dashboard/ProcessFlowEditor.vue'
import ProcessFlowRenderer from '../../components/dashboard/ProcessFlowRenderer.vue'
import ProcessFlowToolbar from '../../components/dashboard/ProcessFlowToolbar.vue'
import ProcessFlowAttributeEditor from '../../components/dashboard/ProcessFlowAttributeEditor.vue'
import Vue from 'vue'
import VueSetup from '../../../ts/vueSetup'
import VueRouter from 'vue-router'

export default Vue.extend({
    components : {
        DashboardAppBar,
        DashboardHomePageNavBar,
        ProcessFlowsNavBar,
        ProcessFlowEditor,
        ProcessFlowRenderer,
        ProcessFlowToolbar,
        ProcessFlowAttributeEditor,
    },
    data : () => ({
        headerClipHeight : 0,
        attrEditorTop: 0,
        attrEditorBottom: 0,
        attrEditorLeft: 0,
        showHideAttrEditor: true
    }),
    router: new VueRouter({
        base : window.location.pathname,
        routes: [
            { path: '/:flowId' } 
        ]
    }),
    methods: {
        recomputeProcessFlowHeaderHeight() {
            Vue.nextTick(() => {
                //@ts-ignore
                const contentDiv = this.$refs.sectionDiv.$el.firstElementChild
                const sectionTop = contentDiv.getBoundingClientRect().top
                //@ts-ignore
                const dividerBottom = this.$refs.headerDivider.$el.getBoundingClientRect().bottom
                this.headerClipHeight = dividerBottom - sectionTop

                Vue.nextTick(() => {
                    //@ts-ignore
                    const attrEditorEl = this.$refs.attrEditor.$el
                    const attrEditorRect = attrEditorEl.getBoundingClientRect()
                    this.attrEditorTop = attrEditorRect.top
                    this.attrEditorBottom = attrEditorRect.bottom
                    this.attrEditorLeft = attrEditorRect.left
                })
            })
        },
        clickAttributePullTab() {
            this.showHideAttrEditor = !this.showHideAttrEditor

            // Need to track the thing through its transition
            // Probably more ideal to make the button CSS
            // transition in the same way but then we have the problem
            // of the initial jump on page load...so do this hack instead.
            let id = setInterval(() => {
                this.recomputeProcessFlowHeaderHeight()
            }, 16)
            setTimeout(() => {
                clearInterval(id) 
            }, 1000)
        }
    },
    computed: {
        attributePullButtonStyle() {
            let leftTranslate : string = this.attrEditorLeft.toString()
            let topTranslate : string = ((this.attrEditorTop + this.attrEditorBottom) / 2).toString()

            return {
                "height": "20px",
                "width": "80px",
                "position": "absolute",
                "left": "0px",
                "top": "0px",
                "transform-origin": "bottom center",
                // Transform top-left corner of the button to the left center point on the
                // attribute editor. Then make it so the bottom center of the button is there instead.
                // Finally rotate.
                "transform": `translate(${leftTranslate}px, ${topTranslate}px) translate(-50%, -100%) rotate(-90deg)`,
                "z-index": 5
            }
        }
    },
    created() {
        VueSetup.currentRouter = this.$router
    },
    mounted() {
        this.recomputeProcessFlowHeaderHeight()
    }
})
</script>

<style scoped>

#attrPullButton {
    border-radius: 3px 3px 0px 0px !important;
}

.no-transition {
    transition: none !important;
}

</style>
