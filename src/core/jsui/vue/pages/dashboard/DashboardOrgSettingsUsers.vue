<template>
    <div class="max-height">
        <dashboard-app-bar ref="dashboardAppBar">
        </dashboard-app-bar>

        <dashboard-org-settings-nav-bar :selected-page="0"></dashboard-org-settings-nav-bar>

        <v-content>
            <dashboard-org-users-manage></dashboard-org-users-manage>
        </v-content>
    </div>
</template>

<script lang="ts">

import DashboardAppBar from '../../components/dashboard/DashboardAppBar.vue'
import DashboardOrgSettingsNavBar from '../../components/dashboard/DashboardOrgSettingsNavBar.vue'
import DashboardOrgUsersManage from '../../components/dashboard/DashboardOrgUsersManage.vue'
import Vue from 'vue'

export default Vue.extend({
    components : {
        DashboardAppBar,
        DashboardOrgSettingsNavBar,
        DashboardOrgUsersManage,
    },
})
</script>
