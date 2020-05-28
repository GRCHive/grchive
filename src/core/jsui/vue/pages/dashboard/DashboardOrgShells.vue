<template>
    <div class="max-height">
        <dashboard-app-bar ref="dashboardAppBar">
        </dashboard-app-bar>

        <dashboard-home-page-nav-bar></dashboard-home-page-nav-bar>

        <v-content>
            <v-tabs>
                <v-tab>
                    Bash
                </v-tab>
                <v-tab-item>
                    <shell-table-with-controls
                        :shell-type="ShellTypes.Bash"
                    >
                    </shell-table-with-controls>
                </v-tab-item>

                <v-tab>
                    PowerShell
                </v-tab>
                <v-tab-item>
                    <shell-table-with-controls
                        :shell-type="ShellTypes.PowerShell"
                    >
                    </shell-table-with-controls>
                </v-tab-item>
            </v-tabs>
        </v-content>
    </div>
</template>

<script lang="ts">

import DashboardAppBar from '../../components/dashboard/DashboardAppBar.vue'
import DashboardHomePageNavBar from '../../components/dashboard/DashboardHomePageNavBar.vue'
import ShellTableWithControls from '../../generic/resources/ShellTableWithControls.vue'
import { ShellTypes } from '../../../ts/shell'
import Vue from 'vue'

export default Vue.extend({
    components : {
        DashboardAppBar,
        DashboardHomePageNavBar,
        ShellTableWithControls,
    },
    data : () => ({
        ShellTypes
    })
})
</script>
