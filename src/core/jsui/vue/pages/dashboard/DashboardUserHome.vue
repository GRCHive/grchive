<template>
    <section class="max-height">
        <dashboard-app-bar>
        </dashboard-app-bar>

        <v-content class="max-height">
            <user-home-page-nav-bar :selected-page="0"></user-home-page-nav-bar>
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
