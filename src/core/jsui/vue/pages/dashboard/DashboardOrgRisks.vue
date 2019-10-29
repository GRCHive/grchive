<template>
    <section class="max-height">
        <dashboard-app-bar ref="dashboardAppBar">
        </dashboard-app-bar>

        <dashboard-home-page-nav-bar :selected-page="2"
        ></dashboard-home-page-nav-bar>

        <v-content>
            <dashboard-risk-list></dashboard-risk-list>
        </v-content>
    </section>
</template>

<script lang="ts">

import DashboardAppBar from '../../components/dashboard/DashboardAppBar.vue'
import DashboardHomePageNavBar from '../../components/dashboard/DashboardHomePageNavBar.vue'
import DashboardRiskList from '../../components/dashboard/DashboardRiskList.vue'
import Vue from 'vue'

export default Vue.extend({
    components : {
        DashboardAppBar,
        DashboardHomePageNavBar,
        DashboardRiskList
    },
})
</script>
