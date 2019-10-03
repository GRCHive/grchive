import Vue from 'vue'
import DashboardOrgHome from '../vue/pages/dashboard/DashboardOrgHome.vue'
import SnackBar from '../vue/components/SnackBar.vue'
import vuetify from  './vuetify'

import '../sass/main.scss'

function mountApp(inData : Object) {
    new Vue({
        el: '#app',
        components: {
            DashboardOrgHome,
            SnackBar
        },
        data: () => (inData),
        vuetify
    }).$mount('#app')
}

export default {
    mountApp
}
