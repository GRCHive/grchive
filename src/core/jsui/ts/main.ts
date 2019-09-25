import Vue from 'vue';
import Vuetify from 'vuetify';
Vue.use(Vuetify);

const opts = {};
export default new Vuetify(opts);

var app = new Vue({
    el: '#app',
    data: {
        message: 'Hello World!'
    }
});
