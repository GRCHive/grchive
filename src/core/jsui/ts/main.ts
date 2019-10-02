import Vue from 'vue'
import vuetify from  './vuetify'
import LandingPage from '../vue/pages/LandingPage.vue'
import ContactUsPage from '../vue/pages/ContactUsPage.vue'
import LoginPage from '../vue/pages/LoginPage.vue'
import GettingStartedPage from '../vue/pages/GettingStartedPage.vue'
import LearnMorePage from '../vue/pages/LearnMorePage.vue'
import RedirectPage from '../vue/pages/RedirectPage.vue'
import SnackBar from '../vue/components/SnackBar.vue'
import '../sass/main.scss'

function mountApp(inData : Object) {
    new Vue({
        el: '#app',
        components: {
            LandingPage,
            ContactUsPage,
            LoginPage,
            GettingStartedPage,
            LearnMorePage,
            RedirectPage,
            SnackBar
        },
        data: () => (inData),
        vuetify
    }).$mount('#app')
}

export default {
    mountApp
}
