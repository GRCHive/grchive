import Vue from 'vue'
import DashboardHome from '../vue/pages/dashboard/DashboardHome.vue'
import vuetify from  './vuetify'

import '../sass/main.scss'

function mountApp(inData : Object) {
    new Vue({
        el: '#app',
        components: {
            DashboardHome
        },
        data: () => (inData),
        vuetify
    }).$mount('#app')
}

export default {
    mountApp
}
