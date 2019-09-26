import Vue from 'vue'
import vuetify from  './vuetify'
import LandingPage from '../vue/pages/LandingPage.vue'
import ContactUsPage from '../vue/pages/ContactUsPage.vue'
import LoginPage from '../vue/pages/LoginPage.vue'

new Vue({
    el: '#app',
    components: {
        LandingPage,
        ContactUsPage,
        LoginPage,
    },
    vuetify
}).$mount('#app')
