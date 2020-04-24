<template>
    <div class="max-height">
        <dashboard-app-bar ref="dashboardAppBar">
        </dashboard-app-bar>

        <dashboard-home-page-nav-bar></dashboard-home-page-nav-bar>

        <v-content>
            <v-tabs>
                <v-tab>Documents</v-tab>
                <v-tab-item>
                    <dashboard-doc-request-list></dashboard-doc-request-list>
                </v-tab-item>

                <v-tab>SQL Queries</v-tab>
                <v-tab-item>
                    <dashboard-sql-request-list class="ma-4"></dashboard-sql-request-list>
                </v-tab-item>

                <v-tab>Script Runs</v-tab>
                <v-tab-item>
                    <dashboard-script-run-request-list class="ma-4"></dashboard-script-run-request-list>
                </v-tab-item>
            </v-tabs>
        </v-content>
    </div>
</template>

<script lang="ts">

import DashboardAppBar from '../../components/dashboard/DashboardAppBar.vue'
import DashboardHomePageNavBar from '../../components/dashboard/DashboardHomePageNavBar.vue'
import DashboardDocRequestList from '../../components/dashboard/DashboardDocRequestList.vue'
import DashboardSqlRequestList from  '../../components/dashboard/DashboardSqlRequestList.vue'
import DashboardScriptRunRequestList from '../../components/dashboard/DashboardScriptRunRequestList.vue'
import Vue from 'vue'

export default Vue.extend({
    components : {
        DashboardAppBar,
        DashboardHomePageNavBar,
        DashboardDocRequestList,
        DashboardSqlRequestList,
        DashboardScriptRunRequestList
    },
})
</script>
