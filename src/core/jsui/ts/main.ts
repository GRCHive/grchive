import Vue from 'vue'
import vuetify from  './vuetify'
import LandingPage from '../vue/LandingPage.vue'

new Vue({
    el: '#landingpage',
    components: {
        LandingPage
    },
    vuetify
}).$mount('#landingpage')
