<template>
    <section class="max-height">
        <dashboard-app-bar>
        </dashboard-app-bar>

        <v-content class="max-height">
            <user-home-page-nav-bar></user-home-page-nav-bar>
        </v-content>

    </section>
</template>

<script>

import DashboardAppBar from '../../components/dashboard/DashboardAppBar.vue'
import UserHomePageNavBar from '../../components/dashboard/UserHomePageNavBar.vue'
export default {
    components : {
        DashboardAppBar,
        UserHomePageNavBar
    }
}
</script>
