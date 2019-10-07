<template>
    <section class="max-height">
        <dashboard-app-bar>
        </dashboard-app-bar>

        <user-home-page-nav-bar :selected-page="0" :mini="false"></user-home-page-nav-bar>
        <v-content class="ma-4">
            <user-profile-view-editor></user-profile-view-editor>
        </v-content>

    </section>
</template>

<script>

import DashboardAppBar from '../../components/dashboard/DashboardAppBar.vue'
import UserHomePageNavBar from '../../components/dashboard/UserHomePageNavBar.vue'
import UserProfileViewEditor from '../../components/dashboard/UserProfileViewEditor.vue'
export default {
    components : {
        DashboardAppBar,
        UserHomePageNavBar,
        UserProfileViewEditor
    }
}
</script>
