<template>
    <div class="max-height">
        <dashboard-app-bar>
        </dashboard-app-bar>

        <dashboard-home-page-nav-bar></dashboard-home-page-nav-bar>
        <v-content>
            <v-list-item>
                <v-list-item-content>
                    <v-list-item-title class="title">
                        PBC Dashboard
                    </v-list-item-title>
                </v-list-item-content>
            </v-list-item>

            <v-divider></v-divider>

            <div class="ma-4">
                <pbc-dashboard></pbc-dashboard>
            </div>

        </v-content>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'

import DashboardAppBar from '../../components/dashboard/DashboardAppBar.vue'
import DashboardHomePageNavBar from '../../components/dashboard/DashboardHomePageNavBar.vue'
import PbcDashboard from '../../components/dashboard/dashboards/PbcDashboard.vue'

@Component({
    components : {
        DashboardAppBar,
        DashboardHomePageNavBar,
        PbcDashboard
    }
})
export default class DashboardOrgHome extends Vue {
}
</script>
