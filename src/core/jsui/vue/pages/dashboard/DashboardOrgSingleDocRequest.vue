<template>
    <div class="max-height">
        <dashboard-app-bar>
        </dashboard-app-bar>

        <dashboard-home-page-nav-bar></dashboard-home-page-nav-bar>

        <v-content>
            <full-edit-doc-request-component></full-edit-doc-request-component>
        </v-content>
    </div>
</template>

<script lang="ts">

import DashboardAppBar from '../../components/dashboard/DashboardAppBar.vue'
import DashboardHomePageNavBar from '../../components/dashboard/DashboardHomePageNavBar.vue'
import FullEditDocRequestComponent from '../../components/dashboard/FullEditDocRequestComponent.vue'
import Vue from 'vue'

export default Vue.extend({
    components : {
        DashboardAppBar,
        DashboardHomePageNavBar,
        FullEditDocRequestComponent
    },
})
</script>
