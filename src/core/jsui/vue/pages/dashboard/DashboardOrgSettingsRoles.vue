<template>
    <div class="max-height">
        <dashboard-app-bar ref="dashboardAppBar">
        </dashboard-app-bar>

        <dashboard-org-settings-nav-bar></dashboard-org-settings-nav-bar>

        <v-content>
            <dashboard-org-roles-manage></dashboard-org-roles-manage>
        </v-content>
    </div>
</template>

<script lang="ts">

import DashboardAppBar from '../../components/dashboard/DashboardAppBar.vue'
import DashboardOrgSettingsNavBar from '../../components/dashboard/DashboardOrgSettingsNavBar.vue'
import DashboardOrgRolesManage from '../../components/dashboard/DashboardOrgRolesManage.vue'
import Vue from 'vue'

export default Vue.extend({
    components : {
        DashboardAppBar,
        DashboardOrgSettingsNavBar,
        DashboardOrgRolesManage,
    },
})
</script>
