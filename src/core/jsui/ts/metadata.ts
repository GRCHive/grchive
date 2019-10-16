import Vue from 'vue'
import Vuex, { StoreOptions } from 'vuex'
import { getProcessFlowIOTypes } from './api/apiProcessFlowIO'
import { getProcessFlowNodeTypes } from './api/apiProcessFlowNodes'

interface MetadataStoreState {
    ioTypes: ProcessFlowIOType[],
    idToIoTypes: Record<number, ProcessFlowIOType>,
    nodeTypes: ProcessFlowNodeType[],
    idToNodeTypes: Record<number, ProcessFlowNodeType>
}

const metaDataStore: StoreOptions<MetadataStoreState> = {
    state: {
        ioTypes: [] as ProcessFlowIOType[],
        idToIoTypes: Object() as Record<number, ProcessFlowIOType>,
        nodeTypes: [] as ProcessFlowNodeType[],
        idToNodeTypes: Object() as Record<number, ProcessFlowNodeType>
    },
    mutations: {
        setIoTypes(state, inTypes : ProcessFlowIOType[]) {
            state.ioTypes = inTypes
            state.idToIoTypes = Object() as Record<number, ProcessFlowIOType>
            for (let typ of inTypes) {
                Vue.set(state.idToIoTypes, typ.Id, typ)
            }
        },
        setNodeTypes(state, inTypes: ProcessFlowNodeType[]) {
            state.nodeTypes = inTypes
            state.idToNodeTypes = Object() as Record<number, ProcessFlowNodeType>
            for (let typ of inTypes) {
                Vue.set(state.idToNodeTypes, typ.Id, typ)
            }
        }
    },
    actions: {
        // Initialization functions
        initialize(context, data) {
            context.dispatch('initializeProcessFlowIOTypes', data)
            context.dispatch('initializeProcessFlowNodeTypes', data)
        },
        initializeProcessFlowIOTypes(context, {csrf}) {
            getProcessFlowIOTypes({csrf}).then((resp : TGetProcessFlowIOTypesOutput) => {
                context.commit('setIoTypes', resp.data)
            })
        },
        initializeProcessFlowNodeTypes(context, {csrf}) {
            getProcessFlowNodeTypes({csrf}).then((resp : TGetProcessFlowNodeTypesOutput) => {
                context.commit('setNodeTypes', resp.data)
            })
        }
    }
}

let store = new Vuex.Store<MetadataStoreState>(metaDataStore)
export default store
