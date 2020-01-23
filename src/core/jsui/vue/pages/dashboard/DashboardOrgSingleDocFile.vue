<template>
    <section class="max-height">
        <dashboard-app-bar @height-change="onHeightChange">
        </dashboard-app-bar>

        <dashboard-home-page-nav-bar></dashboard-home-page-nav-bar>

        <v-content class="max-height">
            <full-edit-documentation-component ref="edit"></full-edit-documentation-component>
        </v-content>
    </section>
</template>

<script lang="ts">

import DashboardAppBar from '../../components/dashboard/DashboardAppBar.vue'
import DashboardHomePageNavBar from '../../components/dashboard/DashboardHomePageNavBar.vue'
import FullEditDocumentationComponent from '../../components/dashboard/FullEditDocumentationComponent.vue'
import Vue from 'vue'

export default Vue.extend({
    components : {
        DashboardAppBar,
        DashboardHomePageNavBar,
        FullEditDocumentationComponent
    },
    methods: {
        // This is not ideal -- ideally information about the header changes goes through a vuex store 
        // and we do the logic inside the edit component.
        onHeightChange() {
            // Need to poll for when the element getBoundingClientRect actually changes.
            // Ideally it changes within a few seconds. Maybe this logic should go inside
            // the component too...
            let intId = setInterval(() => {
                // @ts-ignore
                this.$refs.edit.updateViewerRect()
                // @ts-ignore
                this.$refs.edit.updateMetadataRect()
            }, 100)

            setTimeout(() => {
                clearInterval(intId)
            }, 3000)
        }
    }
})
</script>
