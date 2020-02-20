import 'core-js/stable'
import 'regenerator-runtime/runtime'
import './highlight'
import './vueInit'

import vueOpts from  './vueSetup'
import MetadataStore from './metadata'
import RenderLayout from './render/renderLayout'
import Vue from 'vue'
import { VApp } from 'vuetify/lib'
const DashboardOrgHome = () => import( /* webpackChunkName: "DashboardOrgHome" */ '../vue/pages/dashboard/DashboardOrgHome.vue')
const DashboardOrgProcessFlows = () => import( /* webpackChunkName: "DashboardOrgProcessFlows" */ '../vue/pages/dashboard/DashboardOrgProcessFlows.vue')
const DashboardOrgSingleProcessFlow = () => import( /* webpackChunkName: "DashboardOrgSingleProcessFlow" */ '../vue/pages/dashboard/DashboardOrgSingleProcessFlow.vue')
const DashboardOrgRisks = () => import( /* webpackChunkName: "DashboardOrgRisks" */ '../vue/pages/dashboard/DashboardOrgRisks.vue')
const DashboardOrgSingleRisk = () => import( /* webpackChunkName: "DashboardOrgSingleRisk" */ '../vue/pages/dashboard/DashboardOrgSingleRisk.vue')
const DashboardOrgControls = () => import( /* webpackChunkName: "DashboardOrgControls" */ '../vue/pages/dashboard/DashboardOrgControls.vue')
const DashboardOrgSingleControl = () => import( /* webpackChunkName: "DashboardOrgSingleControl" */ '../vue/pages/dashboard/DashboardOrgSingleControl.vue')
const DashboardUserProfile = () => import( /* webpackChunkName: "DashboardUserProfile" */ '../vue/pages/dashboard/DashboardUserProfile.vue')
const DashboardUserOrgs = () => import( /* webpackChunkName: "DashboardUserOrgs" */ '../vue/pages/dashboard/DashboardUserOrgs.vue')
const DashboardOrgSettingsUsers = () => import( /* webpackChunkName: "DashboardOrgSettingsUsers" */ '../vue/pages/dashboard/DashboardOrgSettingsUsers.vue')
const DashboardOrgSettingsRoles = () => import( /* webpackChunkName: "DashboardOrgSettingsRoles" */ '../vue/pages/dashboard/DashboardOrgSettingsRoles.vue')
const DashboardOrgSettingsSingleRole = () => import( /* webpackChunkName: "DashboardOrgSettingsSingleRole" */ '../vue/pages/dashboard/DashboardOrgSettingsSingleRole.vue')
const DashboardOrgGeneralLedger = () => import( /* webpackChunkName: "DashboardOrgGeneralLedger" */ '../vue/pages/dashboard/DashboardOrgGeneralLedger.vue')
const DashboardOrgSingleGeneralLedgerAccount = () => import( /* webpackChunkName: "DashboardOrgSingleGeneralLedgerAccount" */ '../vue/pages/dashboard/DashboardOrgSingleGeneralLedgerAccount.vue')
const DashboardOrgSystems = () => import( /* webpackChunkName: "DashboardOrgSystems" */ '../vue/pages/dashboard/DashboardOrgSystems.vue')
const DashboardOrgDatabases = () => import( /* webpackChunkName: "DashboardOrgDatabases" */ '../vue/pages/dashboard/DashboardOrgDatabases.vue')
const DashboardOrgServers = () => import( /* webpackChunkName: "DashboardOrgServers" */ '../vue/pages/dashboard/DashboardOrgServers.vue')
const DashboardOrgSingleSystem = () => import( /* webpackChunkName: "DashboardOrgSingleSystem" */ '../vue/pages/dashboard/DashboardOrgSingleSystem.vue')
const DashboardOrgSingleDb = () => import( /* webpackChunkName: "DashboardOrgSingleDb" */ '../vue/pages/dashboard/DashboardOrgSingleDb.vue')
const DashboardOrgDocumentation = () => import( /* webpackChunkName: "DashboardOrgDocumentation" */ '../vue/pages/dashboard/DashboardOrgDocumentation.vue')
const DashboardOrgSingleDocumentation = () => import( /* webpackChunkName: "DashboardOrgSingleDocumentation" */ '../vue/pages/dashboard/DashboardOrgSingleDocumentation.vue')
const DashboardOrgDocRequests = () => import( /* webpackChunkName: "DashboardOrgDocRequests" */ '../vue/pages/dashboard/DashboardOrgDocRequests.vue')
const DashboardOrgSingleDocRequest = () => import( /* webpackChunkName: "DashboardOrgSingleDocRequest" */ '../vue/pages/dashboard/DashboardOrgSingleDocRequest.vue')
const DashboardOrgSingleSqlRequest = () => import( /* webpackChunkName: "DashboardOrgSingleSqlRequest" */ '../vue/pages/dashboard/DashboardOrgSingleSqlRequest.vue')
const DashboardOrgSingleServer = () => import( /* webpackChunkName: "DashboardOrgSingleServer" */ '../vue/pages/dashboard/DashboardOrgSingleServer.vue')
const DashboardOrgVendors = () => import( /* webpackChunkName: "DashboardOrgVendors" */ '../vue/pages/dashboard/DashboardOrgVendors.vue')
const DashboardOrgSingleVendor = () => import( /* webpackChunkName: "DashboardOrgSingleVendor" */ '../vue/pages/dashboard/DashboardOrgSingleVendor.vue')
const DashboardOrgSingleDocFile = () => import( /* webpackChunkName: "DashboardOrgSingleDocFile" */ '../vue/pages/dashboard/DashboardOrgSingleDocFile.vue')
const SnackBar = () => import( /* webpackChunkName: "SnackBar" */ '../vue/components/SnackBar.vue')
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
            DashboardOrgGeneralLedger,
            DashboardOrgSingleGeneralLedgerAccount,
            DashboardOrgSystems,
            DashboardOrgDatabases,
            DashboardOrgServers,
            DashboardOrgSingleSystem,
            DashboardOrgSingleDb,
            DashboardOrgDocumentation,
            DashboardOrgSingleDocumentation,
            DashboardOrgDocRequests,
            DashboardOrgSingleDocRequest,
            DashboardOrgSingleServer,
            DashboardOrgVendors,
            DashboardOrgSingleVendor,
            DashboardOrgSingleDocFile,
            DashboardOrgSingleSqlRequest,
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
                processFlowStore: vueOpts.store
            })
        }
    })
}


export default {
    mountApp
}
