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
import DashboardUserProfile from '../vue/pages/dashboard/DashboardUserProfile.vue'
import DashboardUserOrgs from '../vue/pages/dashboard/DashboardUserOrgs.vue'
import DashboardOrgSettingsUsers from '../vue/pages/dashboard/DashboardOrgSettingsUsers.vue'
import DashboardOrgSettingsRoles from '../vue/pages/dashboard/DashboardOrgSettingsRoles.vue'
import DashboardOrgSettingsSingleRole from '../vue/pages/dashboard/DashboardOrgSettingsSingleRole.vue'
import SnackBar from '../vue/components/SnackBar.vue'
import { getCurrentCSRF } from './csrf'
import { PageParamsStore, PageParamsStoreState  } from '../ts/pageParams'

import '../sass/main.scss'

function mountApp(inData : PageParamsStoreState) {
    PageParamsStore.commit('replaceState', inData)

    if (!!PageParamsStore.state.organization!.Name) {
        document.title = `${PageParamsStore.state.organization!.Name} :: ${PageParamsStore.state.site!.CompanyName}`
    } else {
        document.title = `${PageParamsStore.state.user!.FirstName}  ${PageParamsStore.state.user!.LastName} :: ${PageParamsStore.state.site!.CompanyName}`
    }

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
            DashboardUserProfile,
            DashboardUserOrgs,
            DashboardOrgSettingsUsers,
            DashboardOrgSettingsRoles,
            DashboardOrgSettingsSingleRole,
            SnackBar,
        },
        vuetify: vueOpts.vuetify,
        mounted() {
            MetadataStore.dispatch('initialize', {
                csrf: getCurrentCSRF(),
                orgGroupId: PageParamsStore.state.organization!.OktaGroupName
            })

            RenderLayout.store.dispatch('initialize', {
                host: PageParamsStore.state.site!.Host,
                csrf: getCurrentCSRF(),
                processFlowStore: vueOpts.store})
        }
    })
}

export default {
    mountApp
}
