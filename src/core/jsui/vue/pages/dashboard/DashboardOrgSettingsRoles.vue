<template>
    <div class="max-height">
        <dashboard-app-bar ref="dashboardAppBar">
        </dashboard-app-bar>

        <dashboard-org-settings-nav-bar :selected-page="1"></dashboard-org-settings-nav-bar>

        <v-content>
        </v-content>
    </div>
</template>

<script lang="ts">

import DashboardAppBar from '../../components/dashboard/DashboardAppBar.vue'
import DashboardOrgSettingsNavBar from '../../components/dashboard/DashboardOrgSettingsNavBar.vue'
import Vue from 'vue'

export default Vue.extend({
    components : {
        DashboardAppBar,
        DashboardOrgSettingsNavBar
    },
})
</script>
