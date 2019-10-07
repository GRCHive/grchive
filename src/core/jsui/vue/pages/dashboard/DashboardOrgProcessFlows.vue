<template>
    <section>
        <dashboard-app-bar>
        </dashboard-app-bar>

        <dashboard-home-page-nav-bar :selected-page="1" :mini="false"></dashboard-home-page-nav-bar>
        <v-content>
        </v-content>

    </section>
</template>

<script>

import DashboardAppBar from '../../components/dashboard/DashboardAppBar.vue'
import DashboardHomePageNavBar from '../../components/dashboard/DashboardHomePageNavBar.vue'

export default {
    components : {
        DashboardAppBar,
        DashboardHomePageNavBar
    }
}
</script>
