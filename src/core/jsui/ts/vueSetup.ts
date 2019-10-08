import 'core-js/es/promise'
import Vue from 'vue'
import Vuex from 'vuex'
import Vuetify from 'vuetify'
import '../node_modules/vuetify/dist/vuetify.min.css'
Vue.use(Vuetify)
Vue.use(Vuex)

let mutationObservers = []
const opts = {}

export default {
    vuetify: new Vuetify(opts),
    store: new Vuex.Store({
        state: {
            miniMainNavBar: false,
            primaryNavBarWidth: 256,
        },
        mutations: {
            toggleMiniNavBar(state) {
                state.miniMainNavBar = !state.miniMainNavBar
            },
            changePrimaryNavBarWidth(state, width) {
                state.primaryNavBarWidth = width
            },
        },
        actions: {
            mountPrimaryNavBar(context, nav) {
                let observer = new MutationObserver(function(records : MutationRecord[], _: MutationObserver) {
                    for (let mutation of records) {
                        if (mutation.type == "attributes" && mutation.attributeName == "style") {
                            // For whatever reason, mutation.target.offsetWidth does not get updated
                            // immediately (even though console.log begs to differ) so the only
                            // reliable thing to use is the target's style. Assume that it'll always
                            // be in pixels...
                            //@ts-ignore
                            context.commit('changePrimaryNavBarWidth', parseInt(mutation.target.style.width, 10))
                            break
                        }
                    }
                })
                observer.observe(nav.$el, {
                    attributes: true,
                    subtree: true
                })
                mutationObservers.push(observer)
            },

        }
    })
}
