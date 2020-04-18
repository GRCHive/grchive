<template>
    <div class="max-height">
        <dashboard-app-bar ref="dashboardAppBar">
        </dashboard-app-bar>

        <dashboard-home-page-nav-bar></dashboard-home-page-nav-bar>

        <v-content>
            <v-row>
                <v-col cols="6">
                    <log-viewer :run-id="runId" run-log></log-viewer>
                </v-col>

                <v-col cols="6">
                    <log-viewer :run-id="runId"></log-viewer>
                </v-col>
            </v-row>
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
export default class DashboardOrgSingleRunLog extends Vue {
    get runId() : number {
        return parseInt(PageParamsStore.state.resource!.Id, 10)
    }
}

</script>
