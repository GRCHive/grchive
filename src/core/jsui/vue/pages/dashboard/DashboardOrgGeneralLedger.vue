<template>
    <section class="max-height">
        <dashboard-app-bar ref="dashboardAppBar">
        </dashboard-app-bar>

        <dashboard-home-page-nav-bar :selected-page="4"
        ></dashboard-home-page-nav-bar>

        <v-content>
            <dashboard-general-ledger-list></dashboard-general-ledger-list>
        </v-content>
    </section>
</template>

<script lang="ts">

import DashboardAppBar from '../../components/dashboard/DashboardAppBar.vue'
import DashboardHomePageNavBar from '../../components/dashboard/DashboardHomePageNavBar.vue'
import DashboardGeneralLedgerList from '../../components/dashboard/DashboardGeneralLedgerList.vue'
import Vue from 'vue'

export default Vue.extend({
    components : {
        DashboardAppBar,
        DashboardHomePageNavBar,
        DashboardGeneralLedgerList
    },
})
</script>
