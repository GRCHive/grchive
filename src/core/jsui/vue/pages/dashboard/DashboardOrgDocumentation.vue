<template>
    <div class="max-height">
        <dashboard-app-bar ref="dashboardAppBar">
        </dashboard-app-bar>

        <dashboard-home-page-nav-bar></dashboard-home-page-nav-bar>

        <v-content>
            <dashboard-documentation-list></dashboard-documentation-list>
        </v-content>
    </div>
</template>

<script lang="ts">

import DashboardAppBar from '../../components/dashboard/DashboardAppBar.vue'
import DashboardHomePageNavBar from '../../components/dashboard/DashboardHomePageNavBar.vue'
import DashboardDocumentationList from '../../components/dashboard/DashboardDocumentationList.vue'
import Vue from 'vue'

export default Vue.extend({
    components : {
        DashboardAppBar,
        DashboardHomePageNavBar,
        DashboardDocumentationList,
    },
})
</script>
