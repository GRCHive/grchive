<template>
    <div class="max-height">
        <dashboard-app-bar ref="dashboardAppBar">
        </dashboard-app-bar>

        <dashboard-org-settings-nav-bar :selected-page="1"></dashboard-org-settings-nav-bar>

        <v-content>
            <full-edit-role-component></full-edit-role-component>
        </v-content>
    </div>
</template>

<script lang="ts">

import DashboardAppBar from '../../components/dashboard/DashboardAppBar.vue'
import DashboardOrgSettingsNavBar from '../../components/dashboard/DashboardOrgSettingsNavBar.vue'
import FullEditRoleComponent from '../../components/dashboard/FullEditRoleComponent.vue'
import Vue from 'vue'

export default Vue.extend({
    components : {
        DashboardAppBar,
        DashboardOrgSettingsNavBar,
        FullEditRoleComponent
    },
})
</script>
