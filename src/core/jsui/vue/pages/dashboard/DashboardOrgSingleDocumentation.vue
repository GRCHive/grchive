<template>
    <div class="max-height">
        <dashboard-app-bar>
        </dashboard-app-bar>

        <dashboard-home-page-nav-bar></dashboard-home-page-nav-bar>

        <v-content>
            <full-edit-documentation-category-component></full-edit-documentation-category-component>
        </v-content>
    </div>
</template>

<script lang="ts">

import DashboardAppBar from '../../components/dashboard/DashboardAppBar.vue'
import DashboardHomePageNavBar from '../../components/dashboard/DashboardHomePageNavBar.vue'
import FullEditDocumentationCategoryComponent from '../../components/dashboard/FullEditDocumentationCategoryComponent.vue'
import Vue from 'vue'

export default Vue.extend({
    components : {
        DashboardAppBar,
        DashboardHomePageNavBar,
        FullEditDocumentationCategoryComponent
    },
})
</script>
