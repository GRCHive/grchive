import vueOpts from  './vueSetup'
import MetadataStore from './metadata'
import Vue from 'vue'
import DashboardOrgHome from '../vue/pages/dashboard/DashboardOrgHome.vue'
import DashboardOrgProcessFlows from '../vue/pages/dashboard/DashboardOrgProcessFlows.vue'
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
            DashboardUserHome,
            SnackBar
        },
        data: () => (inData),
        vuetify: vueOpts.vuetify,
        mounted() {
            //@ts-ignore
            MetadataStore.dispatch('initialize', {csrf: inData.csrf})
        }
    })
}

export default {
    mountApp
}
