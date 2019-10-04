import Vue from 'vue'
import DashboardOrgHome from '../vue/pages/dashboard/DashboardOrgHome.vue'
import DashboardUserHome from '../vue/pages/dashboard/DashboardUserHome.vue'
import SnackBar from '../vue/components/SnackBar.vue'
import vuetify from  './vuetify'

import '../sass/main.scss'
import '@mdi/font/scss/materialdesignicons.scss'

function mountApp(inData : Object) {
    new Vue({
        el: '#app',
        components: {
            DashboardOrgHome,
            DashboardUserHome,
            SnackBar
        },
        data: () => (inData),
        vuetify
    }).$mount('#app')
}

export default {
    mountApp
}
