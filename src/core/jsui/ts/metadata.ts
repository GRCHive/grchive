import Vue from 'vue'
import Vuex, { StoreOptions } from 'vuex'
import { getProcessFlowIOTypes } from './api/apiProcessFlowIO'

interface MetadataStoreState {
    ioTypes: ProcessFlowIOType[],
    idToIoTypes: Record<number, ProcessFlowIOType>
}

const metaDataStore: StoreOptions<MetadataStoreState> = {
    state: {
        ioTypes: [] as ProcessFlowIOType[],
        idToIoTypes: Object() as Record<number, ProcessFlowIOType>
    },
    mutations: {
        setIoTypes(state, inTypes : ProcessFlowIOType[]) {
            state.ioTypes = inTypes
            state.idToIoTypes = Object() as Record<number, ProcessFlowIOType>
            for (let typ of inTypes) {
                Vue.set(state.idToIoTypes, typ.Id, typ)
            }
        }
    },
    actions: {
        // Initialization functions
        initialize(context, data) {
            context.dispatch('initializeProcessFlowIOTypes', data)
        },
        initializeProcessFlowIOTypes(context, {csrf}) {
            getProcessFlowIOTypes({csrf}).then((resp : TGetProcessFlowIOTypesOutput) => {
                context.commit('setIoTypes', resp.data)
            })
        }
    }
}

let store = new Vuex.Store<MetadataStoreState>(metaDataStore)
export default store
