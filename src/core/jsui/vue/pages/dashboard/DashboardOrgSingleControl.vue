<template>
    <section class="max-height">
        <dashboard-app-bar ref="dashboardAppBar">
        </dashboard-app-bar>

        <dashboard-home-page-nav-bar :selected-page="3"
        ></dashboard-home-page-nav-bar>

        <v-content>
            <full-edit-control-component></full-edit-control-component>
        </v-content>
    </section>
</template>

<script lang="ts">

import DashboardAppBar from '../../components/dashboard/DashboardAppBar.vue'
import DashboardHomePageNavBar from '../../components/dashboard/DashboardHomePageNavBar.vue'
import FullEditControlComponent from '../../components/dashboard/FullEditControlComponent.vue'
import Vue from 'vue'

export default Vue.extend({
    components : {
        DashboardAppBar,
        DashboardHomePageNavBar,
        FullEditControlComponent
    },
})
</script>
