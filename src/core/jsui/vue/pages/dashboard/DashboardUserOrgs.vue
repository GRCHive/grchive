<template>
    <div class="max-height">
        <dashboard-app-bar>
        </dashboard-app-bar>

        <user-home-page-nav-bar :selected-page="0"></user-home-page-nav-bar>
        <v-content class="ma-4">
            <user-dashboard></user-dashboard>
        </v-content>
    </div>
</template>

<script>

import DashboardAppBar from '../../components/dashboard/DashboardAppBar.vue'
import UserHomePageNavBar from '../../components/dashboard/UserHomePageNavBar.vue'
import UserDashboard from '../../components/dashboard/UserDashboard.vue'
export default {
    components : {
        DashboardAppBar,
        UserHomePageNavBar,
        UserDashboard,
    }
}

</script>
