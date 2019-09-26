import Vue from 'vue'
import vuetify from  './vuetify'

import HomePage from '../vue/HomePage.vue'

new Vue({
    el: '#homepage',
    components: {HomePage},
    vuetify
}).$mount('#homepage')
