<template>
    <div class="max-height">
        <dashboard-app-bar ref="dashboardAppBar">
        </dashboard-app-bar>

        <dashboard-org-settings-nav-bar></dashboard-org-settings-nav-bar>

        <v-content>
            <v-tabs>
                <v-tab>
                    PBC
                </v-tab>
                
                <v-tab-item>
                    <pbc-notification-settings-component>
                    </pbc-notification-settings-component>
                </v-tab-item>
            </v-tabs>
        </v-content>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'

import DashboardAppBar from '../../components/dashboard/DashboardAppBar.vue'
import DashboardOrgSettingsNavBar from '../../components/dashboard/DashboardOrgSettingsNavBar.vue'
import PbcNotificationSettingsComponent from '../../components/dashboard/settings/PbcNotificationSettingsComponent.vue'

@Component({
    components: {
        DashboardAppBar,
        DashboardOrgSettingsNavBar,
        PbcNotificationSettingsComponent,
    }
})
export default class DashboardOrgSettingsNotifications extends Vue {
}

</script>
