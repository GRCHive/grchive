import vueOpts from  './vueSetup'
import MetadataStore from './metadata'
import RenderLayout from './render/renderLayout'
import Vue from 'vue'
import DashboardOrgHome from '../vue/pages/dashboard/DashboardOrgHome.vue'
import DashboardOrgProcessFlows from '../vue/pages/dashboard/DashboardOrgProcessFlows.vue'
import DashboardOrgRisks from '../vue/pages/dashboard/DashboardOrgRisks.vue'
import DashboardUserHome from '../vue/pages/dashboard/DashboardUserHome.vue'
import SnackBar from '../vue/components/SnackBar.vue'

import '../sass/main.scss'
import '@mdi/font/scss/materialdesignicons.scss'

function mountApp(inData : Object) {
    new Vue({
        el: '#app',
        components: {
            DashboardOrgHome,
            DashboardOrgProcessFlows,
            DashboardOrgRisks,
            DashboardUserHome,
            SnackBar
        },
        data: () => (inData),
        vuetify: vueOpts.vuetify,
        mounted() {
            //@ts-ignore
            MetadataStore.dispatch('initialize', inData)

            //@ts-ignore
            RenderLayout.store.dispatch('initialize', {
                //@ts-ignore
                host: inData.host,
                //@ts-ignore
                csrf: inData.csrf,
                processFlowStore: vueOpts.store})
        }
    })
}

export default {
    mountApp
}
