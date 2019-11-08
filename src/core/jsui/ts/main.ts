import vueOpts from  './vueSetup'
import Vue from 'vue'
import { VApp } from 'vuetify/lib'
import LandingPage from '../vue/pages/LandingPage.vue'
import ContactUsPage from '../vue/pages/ContactUsPage.vue'
import LoginPage from '../vue/pages/LoginPage.vue'
import GettingStartedPage from '../vue/pages/GettingStartedPage.vue'
import LearnMorePage from '../vue/pages/LearnMorePage.vue'
import RedirectPage from '../vue/pages/RedirectPage.vue'
import ErrorPage from '../vue/pages/ErrorPage.vue'
import SnackBar from '../vue/components/SnackBar.vue'
import '../sass/main.scss'
import { PageParamsStore, PageParamsStoreState  } from '../ts/pageParams'

function mountApp(inData : PageParamsStoreState) {
    PageParamsStore.commit('replaceState', inData)
    document.title = `${PageParamsStore.state.site!.CompanyName}`

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
        },
        vuetify: vueOpts.vuetify
    }).$mount('#app')
}

export default {
    mountApp
}
