<template>
    <div class="max-height">
        <dashboard-app-bar>
        </dashboard-app-bar>

        <user-home-page-nav-bar :selected-page="0"></user-home-page-nav-bar>
        <v-content class="ma-4">
            <user-org-dashboard></user-org-dashboard>
        </v-content>
    </div>
</template>

<script>

import DashboardAppBar from '../../components/dashboard/DashboardAppBar.vue'
import UserHomePageNavBar from '../../components/dashboard/UserHomePageNavBar.vue'
import UserOrgDashboard from '../../components/dashboard/UserOrgDashboard.vue'
export default {
    components : {
        DashboardAppBar,
        UserHomePageNavBar,
        UserOrgDashboard,
    }
}

</script>
