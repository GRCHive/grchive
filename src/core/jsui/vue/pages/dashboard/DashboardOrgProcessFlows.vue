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
            <process-flow-attribute-editor :custom-clip-height="headerClipHeight"></process-flow-attribute-editor>
        </v-content>

    </section>
</template>

<script>

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

export default {
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
        headerClipHeight : 0
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
                const contentDiv = this.$refs.sectionDiv.$el.firstElementChild
                const sectionTop = contentDiv.getBoundingClientRect().top
                const dividerBottom = this.$refs.headerDivider.$el.getBoundingClientRect().bottom
                this.headerClipHeight = dividerBottom - sectionTop
            })
        }
    },
    created() {
        VueSetup.currentRouter = this.$router
    },
    mounted() {
        this.recomputeProcessFlowHeaderHeight()
    }
}
</script>
