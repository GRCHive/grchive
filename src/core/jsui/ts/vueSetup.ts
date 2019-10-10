import 'core-js/es/promise'
import Vue from 'vue'
import Vuex, { StoreOptions } from 'vuex'
import VueRouter from 'vue-router'
import Vuetify from 'vuetify'
import '../node_modules/vuetify/dist/vuetify.min.css'

Vue.use(Vuex)
Vue.use(VueRouter)
Vue.use(Vuetify)

import axios from 'axios'
import {createGetProcessFlowFullDataUrl} from './url'
import * as qs from 'query-string'

let mutationObservers = []
const opts = {}

interface FullProcessFlowDataResponse {
    data : FullProcessFlowData
}

interface VuexState {
    miniMainNavBar : boolean,
    primaryNavBarWidth : number,
    allProcessFlowBasicData : ProcessFlowBasicData[],
    currentProcessFlowIndex: number,
    currentProcessFlowFullData: FullProcessFlowData
}

const store : StoreOptions<VuexState> = {
    state: {
        miniMainNavBar: false,
        primaryNavBarWidth: 256,
        allProcessFlowBasicData: [],
        currentProcessFlowIndex : 0,
        currentProcessFlowFullData: {} as FullProcessFlowData
    },
    mutations: {
        toggleMiniNavBar(state) {
            state.miniMainNavBar = !state.miniMainNavBar
        },
        changePrimaryNavBarWidth(state, width) {
            state.primaryNavBarWidth = width
        },
        setProcessFlowBasicData(state, data) {
            for (let d of data) {
                d.CreationTime = new Date(d.CreationTime)
                d.LastUpdatedTime = new Date(d.LastUpdatedTime)
            }
            state.allProcessFlowBasicData = data
        },
        setCurrentProcessFlowIndex(state, index) {
            state.currentProcessFlowIndex = index
        },
        setIndividualProcessFlowBasicData(state, {index, data}) {
            data.CreationTime = new Date(data.CreationTime)
            data.LastUpdatedTime = new Date(data.LastUpdatedTime)
            state.allProcessFlowBasicData.splice(index, 1, data)
        },
        setCurrentProcessFlowFullData(state, data) {
            state.currentProcessFlowFullData = data
        }
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
        requestSetCurrentProcessFlowIndex(context, {index, csrf}) {
            context.commit('setCurrentProcessFlowIndex', index)
            context.dispatch('refreshCurrentProcessFlowFullData', csrf)
        },
        refreshCurrentProcessFlowFullData(context, csrf) {
            const baseUrl = createGetProcessFlowFullDataUrl(context.getters.currentProcessFlowBasicData.Id)
            const queryParams = qs.stringify({
                csrf
            })
            axios.get(baseUrl + "?" + queryParams).then(
                (resp : FullProcessFlowDataResponse) => {
                    context.commit('setCurrentProcessFlowFullData', resp.data)
                }
            ).catch(
                (_) => {
                    context.commit('setCurrentProcessFlowFullData', {} as FullProcessFlowData)
                }
            )
        }
    },
    getters: {
        currentProcessFlowBasicData: (state) => {
            return state.allProcessFlowBasicData[state.currentProcessFlowIndex]
        }
    }
}

export default {
    vuetify: new Vuetify(opts),
    store: new Vuex.Store<VuexState>(store),
    currentRouter: {} as VueRouter
}
