<template>
    <div>
       <v-content class="max-height" ref="sectionDiv">
        <div :style="contentContainerStyle">
            <process-flow-editor @on-change="recomputeProcessFlowHeaderHeight"></process-flow-editor>
            <v-divider></v-divider>
            <process-flow-toolbar></process-flow-toolbar>
            <v-divider ref="headerDivider"></v-divider>
            <process-flow-renderer :content-max-height-clip="headerClipHeight"
                                   :content-max-width-clip="attrEditorClipWidth"
                                   :display-rect="rendererClientRect"
                                   ref="rendererVue"
            ></process-flow-renderer>

            <process-flow-attribute-editor :custom-clip-height="headerClipHeight" 
                                           ref="attrEditor"
                                           :show-hide="showHideAttrEditor && isNodeSelected"
            ></process-flow-attribute-editor>
        </div>
    </v-content>

        <v-btn color="primary"
               id="attrPullButton"
               small
               :style="attributePullButtonStyle"
               class="no-transition"
               @click="clickAttributePullTab"
                v-if="isNodeSelected"
        >
            <v-icon v-if="!showHideAttrEditor">mdi-chevron-up</v-icon>
            <v-icon v-else>mdi-chevron-down</v-icon>
        </v-btn>
    </div>
</template>

<script lang="ts">
import DashboardAppBar from '../../components/dashboard/DashboardAppBar.vue'
import DashboardHomePageNavBar from '../../components/dashboard/DashboardHomePageNavBar.vue'
import ProcessFlowEditor from '../../components/dashboard/ProcessFlowEditor.vue'
import ProcessFlowRenderer from '../../components/dashboard/ProcessFlowRenderer.vue'
import ProcessFlowToolbar from '../../components/dashboard/ProcessFlowToolbar.vue'
import ProcessFlowAttributeEditor from '../../components/dashboard/ProcessFlowAttributeEditor.vue'
import Vue from 'vue'
import VueSetup from '../../../ts/vueSetup'
import VueRouter from 'vue-router'
import LocalSettings from '../../../ts/localSettings'
import RenderLayout from '../../../ts/render/renderLayout'

export default Vue.extend({
    components : {
        DashboardAppBar,
        DashboardHomePageNavBar,
        ProcessFlowEditor,
        ProcessFlowRenderer,
        ProcessFlowToolbar,
        ProcessFlowAttributeEditor,
    },
    data : () => ({
        appBarClipHeight: 0,
        headerClipHeight : 0,
        attrEditorTop: 0,
        attrEditorBottom: 0,
        attrEditorLeft: 0,
        attrEditorClipWidth: 256,
    }),
    router: new VueRouter({
        base : window.location.pathname,
        routes: [
            { 
                path: '/:flowId?',
                name: 'flow'
            } 
        ]
    }),
    methods: {
        updateClientRect() {
            //@ts-ignore
            const rect = this.$refs.rendererVue.$el.getBoundingClientRect()
            const rendererClientRect =  <IDOMRect>{
                top: rect.top,
                bottom: rect.bottom,
                left: rect.left,
                right: rect.right,
                width: rect.width,
                height: rect.height
            }
            RenderLayout.store.commit('setRendererRect', rendererClientRect)
        },
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
                    this.attrEditorClipWidth = this.$root.$el.clientWidth - this.attrEditorLeft
                })
            })

            // Spend the next second or so making sure we keep track of the total size
            // of the rendering section.
            let intervalId = setInterval(() => {
                this.updateClientRect()
            }, 16)

            setTimeout(function() {
                clearInterval(intervalId)
            }, 1000)
        },
        clickAttributePullTab() {
            LocalSettings.commit('setShowHideAttributeEditor', !this.showHideAttrEditor)
            this.trackAttributeEditor()
        },
        trackAttributeEditor() {
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
        rendererClientRect() : IDOMRect {
            return RenderLayout.store.state.rendererRect
        },
        showHideAttrEditor() {
            return LocalSettings.state.showHideAttributeEditor
        },
        isNodeSelected() : boolean {
            return VueSetup.store.getters.isNodeSelected
        },
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
        },
        contentContainerStyle() {
            return {
                "height": "100vh",
                "maxHeight": `calc(100vh - ${this.appBarClipHeight}px)`
            }
        }
    },
    created() {
        VueSetup.currentRouter = this.$router
    },
    mounted() {
        this.recomputeProcessFlowHeaderHeight()
        //@ts-ignore
        this.appBarClipHeight = this.$refs.dashboardAppBar.$el.offsetHeight

        window.addEventListener('resize', this.updateClientRect)
    },
    watch: {
        isNodeSelected() {
            this.trackAttributeEditor()
        }
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
