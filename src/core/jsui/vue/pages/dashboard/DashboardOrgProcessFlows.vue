<template>
    <section>
        <dashboard-app-bar>
        </dashboard-app-bar>

        <dashboard-home-page-nav-bar :selected-page="1"></dashboard-home-page-nav-bar>
        <process-flows-nav-bar></process-flows-nav-bar>
        <v-content>
            <process-flow-editor></process-flow-editor>
            <v-divider></v-divider>
            <process-flow-toolbar></process-flow-toolbar>
            <v-divider></v-divider>
            <process-flow-renderer></process-flow-renderer>
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
    },
    router: new VueRouter({
        base : window.location.pathname,
        routes: [
            { path: '/:flowId' } 
        ]
    }),
    created() {
        VueSetup.currentRouter = this.$router
    },
}
</script>
