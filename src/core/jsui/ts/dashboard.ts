import vueOpts from  './vueSetup'
import MetadataStore from './metadata'
import RenderLayout from './render/renderLayout'
import Vue from 'vue'
import { VApp } from 'vuetify/lib'
import DashboardOrgHome from '../vue/pages/dashboard/DashboardOrgHome.vue'
import DashboardOrgProcessFlows from '../vue/pages/dashboard/DashboardOrgProcessFlows.vue'
import DashboardOrgSingleProcessFlow from '../vue/pages/dashboard/DashboardOrgSingleProcessFlow.vue'
import DashboardOrgRisks from '../vue/pages/dashboard/DashboardOrgRisks.vue'
import DashboardOrgSingleRisk from '../vue/pages/dashboard/DashboardOrgSingleRisk.vue'
import DashboardOrgControls from '../vue/pages/dashboard/DashboardOrgControls.vue'
import DashboardOrgSingleControl from '../vue/pages/dashboard/DashboardOrgSingleControl.vue'
import DashboardUserHome from '../vue/pages/dashboard/DashboardUserHome.vue'
import SnackBar from '../vue/components/SnackBar.vue'
import { getCurrentCSRF } from './csrf'

import '../sass/main.scss'

function mountApp(inData : Object) {
    new Vue({
        el: '#app',
        components: {
            VApp,
            DashboardOrgHome,
            DashboardOrgProcessFlows,
            DashboardOrgSingleProcessFlow,
            DashboardOrgRisks,
            DashboardOrgSingleRisk,
            DashboardOrgControls,
            DashboardOrgSingleControl,
            DashboardUserHome,
            SnackBar
        },
        data: () => (inData),
        vuetify: vueOpts.vuetify,
        mounted() {
            //@ts-ignore
            MetadataStore.dispatch('initialize', {
                ...inData,
                csrf: getCurrentCSRF()
            })

            //@ts-ignore
            RenderLayout.store.dispatch('initialize', {
                //@ts-ignore
                host: inData.host,
                csrf: getCurrentCSRF(),
                processFlowStore: vueOpts.store})
        }
    })
}

export default {
    mountApp
}
