<template>
    <div class="max-height">
        <dashboard-app-bar ref="dashboardAppBar">
        </dashboard-app-bar>

        <dashboard-home-page-nav-bar></dashboard-home-page-nav-bar>

        <v-content>
            <v-tabs>
                <v-tab>Build</v-tab>
                <v-tab-item>
                    <build-log-list></build-log-list>
                </v-tab-item>

                <v-tab>Run</v-tab>
                <v-tab-item>
                    <run-log-list></run-log-list>
                </v-tab-item>
            </v-tabs>
        </v-content>
    </div>
</template>

<script lang="ts">

import DashboardAppBar from '../../components/dashboard/DashboardAppBar.vue'
import DashboardHomePageNavBar from '../../components/dashboard/DashboardHomePageNavBar.vue'
import RunLogList from '../../generic/logs/RunLogList.vue'
import BuildLogList from '../../generic/logs/BuildLogList.vue'
import Vue from 'vue'

export default Vue.extend({
    components : {
        DashboardAppBar,
        DashboardHomePageNavBar,
        RunLogList,
        BuildLogList,
    },
})

</script>
