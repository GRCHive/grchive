<template>
    <div class="max-height">
        <dashboard-app-bar ref="dashboardAppBar">
        </dashboard-app-bar>

        <dashboard-home-page-nav-bar></dashboard-home-page-nav-bar>

        <v-content>
            <log-viewer :commit="commit"></log-viewer>
        </v-content>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import DashboardAppBar from '../../components/dashboard/DashboardAppBar.vue'
import DashboardHomePageNavBar from '../../components/dashboard/DashboardHomePageNavBar.vue'
import LogViewer from '../../generic/logs/LogViewer.vue'
import { PageParamsStore } from '../../../ts/pageParams'

@Component({
    components: {
        DashboardAppBar,
        DashboardHomePageNavBar,
        LogViewer,
    }
})
export default class DashboardOrgSingleBuildLog extends Vue {
    get commit() : string {
        return PageParamsStore.state.resource!.Id
    }
}

</script>
