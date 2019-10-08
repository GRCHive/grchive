import 'core-js/es/promise'
import Vue from 'vue'
import Vuex from 'vuex'
import Vuetify from 'vuetify'
import '../node_modules/vuetify/dist/vuetify.min.css'
Vue.use(Vuetify)
Vue.use(Vuex)

const opts = {}
export default {
    vuetify: new Vuetify(opts),
    store: new Vuex.Store({
    })
}
