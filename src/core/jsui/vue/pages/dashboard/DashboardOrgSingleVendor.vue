<template>
    <div class="max-height">
        <dashboard-app-bar>
        </dashboard-app-bar>

        <dashboard-home-page-nav-bar></dashboard-home-page-nav-bar>

        <v-content>
            <full-edit-vendor-component></full-edit-vendor-component>
        </v-content>
    </div>
</template>

<script lang="ts">

import DashboardAppBar from '../../components/dashboard/DashboardAppBar.vue'
import DashboardHomePageNavBar from '../../components/dashboard/DashboardHomePageNavBar.vue'
import Vue from 'vue'
import FullEditVendorComponent from '../../components/dashboard/FullEditVendorComponent.vue'

export default Vue.extend({
    components : {
        DashboardAppBar,
        DashboardHomePageNavBar,
        FullEditVendorComponent
    },
})

</script>
