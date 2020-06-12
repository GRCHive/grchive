import './vueInit'
import vueOpts from  './vueSetup'
import Vue from 'vue'
import { VApp } from 'vuetify/lib'
const LandingPage = () => import( /* webpackChunkName: "LandingPage" */ '../vue/pages/LandingPage.vue')
const ContactUsPage = () => import( /* webpackChunkName: "ContactUsPage" */ '../vue/pages/ContactUsPage.vue')
const LoginPage = () => import( /* webpackChunkName: "LoginPage" */ '../vue/pages/LoginPage.vue')
const RegistrationPage = () => import( /* webpackChunkName: "RegistrationPage" */ '../vue/pages/RegistrationPage.vue')
const GettingStartedPage = () => import( /* webpackChunkName: "GettingStartedPage" */ '../vue/pages/GettingStartedPage.vue')
const LearnMorePage = () => import( /* webpackChunkName: "LearnMorePage" */ '../vue/pages/LearnMorePage.vue')
const RedirectPage = () => import( /* webpackChunkName: "RedirectPage" */ '../vue/pages/RedirectPage.vue')
const ErrorPage = () => import( /* webpackChunkName: "ErrorPage" */ '../vue/pages/ErrorPage.vue')
const SnackBar = () => import( /* webpackChunkName: "SnackBar" */ '../vue/components/SnackBar.vue')
import '../sass/main.scss'
import { PageParamsStore, PageParamsStoreState  } from '../ts/pageParams'

function mountApp(inData : PageParamsStoreState) {
    PageParamsStore.commit('replaceState', inData)
    new Vue({
        el: '#app',
        components: {
            VApp,
            LandingPage,
            ContactUsPage,
            LoginPage,
            GettingStartedPage,
            LearnMorePage,
            RedirectPage,
            SnackBar,
            ErrorPage,
            RegistrationPage
        },
        vuetify: vueOpts.vuetify
    }).$mount('#app')
}

export default {
    mountApp
}
