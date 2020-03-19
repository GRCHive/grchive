<template>
    <div class="max-height">
        <dashboard-app-bar>
        </dashboard-app-bar>

        <user-home-page-nav-bar></user-home-page-nav-bar>
        <v-content class="ma-4">
            <notification-viewer></notification-viewer>
        </v-content>
    </div>
</template>

<script>

import DashboardAppBar from '../../components/dashboard/DashboardAppBar.vue'
import UserHomePageNavBar from '../../components/dashboard/UserHomePageNavBar.vue'
import NotificationViewer from '../../generic/notifications/NotificationViewer.vue'
export default {
    components : {
        DashboardAppBar,
        UserHomePageNavBar,
        NotificationViewer
    }
}
</script>
