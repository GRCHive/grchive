<template>
    <div class="max-height">
        <dashboard-app-bar ref="dashboardAppBar">
        </dashboard-app-bar>

        <dashboard-home-page-nav-bar></dashboard-home-page-nav-bar>

        <v-content class="max-height">
            <schedule-viewer class="max-height"></schedule-viewer>
        </v-content>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import DashboardAppBar from '../../components/dashboard/DashboardAppBar.vue'
import DashboardHomePageNavBar from '../../components/dashboard/DashboardHomePageNavBar.vue'
import ScheduleViewer from '../../generic/tasks/ScheduleViewer.vue'

@Component({
    components: {
        DashboardAppBar,
        DashboardHomePageNavBar,
        ScheduleViewer
    }
})
export default class DashboardOrgScriptSchedule extends Vue {
}

</script>
