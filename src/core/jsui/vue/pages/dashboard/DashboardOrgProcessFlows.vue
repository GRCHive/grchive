<template>
    <section>
        <dashboard-app-bar>
        </dashboard-app-bar>

        <dashboard-home-page-nav-bar :selected-page="1"></dashboard-home-page-nav-bar>
        <process-flows-nav-bar></process-flows-nav-bar>
        <v-content>
        </v-content>

    </section>
</template>

<script>

import DashboardAppBar from '../../components/dashboard/DashboardAppBar.vue'
import DashboardHomePageNavBar from '../../components/dashboard/DashboardHomePageNavBar.vue'
import ProcessFlowsNavBar from '../../components/dashboard/ProcessFlowsNavBar.vue'

export default {
    components : {
        DashboardAppBar,
        DashboardHomePageNavBar,
        ProcessFlowsNavBar
    }
}
</script>
