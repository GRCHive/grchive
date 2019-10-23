import Vue from 'vue'
import Vuex, { StoreOptions } from 'vuex'
import { getProcessFlowIOTypes } from './api/apiProcessFlowIO'
import { getProcessFlowNodeTypes } from './api/apiProcessFlowNodes'
import { getControlTypes } from './api/apiControls'

interface MetadataStoreState {
    ioTypes: ProcessFlowIOType[],
    idToIoTypes: Record<number, ProcessFlowIOType>,
    nodeTypes: ProcessFlowNodeType[],
    idToNodeTypes: Record<number, ProcessFlowNodeType>,
    controlTypes: ProcessFlowControlType[],
    idToControlTypes: Record<number, ProcessFlowControlType>
}

const metaDataStore: StoreOptions<MetadataStoreState> = {
    state: {
        ioTypes: [] as ProcessFlowIOType[],
        idToIoTypes: Object() as Record<number, ProcessFlowIOType>,
        nodeTypes: [] as ProcessFlowNodeType[],
        idToNodeTypes: Object() as Record<number, ProcessFlowNodeType>,
        controlTypes: [] as ProcessFlowControlType[],
        idToControlTypes: Object() as Record<number, ProcessFlowControlType>
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
        },
        setControlTypes(state, inTypes: ProcessFlowControlType[]) {
            state.controlTypes = inTypes
            state.idToControlTypes = Object() as Record<number, ProcessFlowControlType>
            for (let typ of inTypes) {
                Vue.set(state.idToControlTypes, typ.Id, typ)
            }
        }
    },
    actions: {
        // Initialization functions
        initialize(context, data) {
            context.dispatch('initializeProcessFlowIOTypes', data)
            context.dispatch('initializeProcessFlowNodeTypes', data)
            context.dispatch('initializeProcessFlowControlTypes', data)
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
        },
        initializeProcessFlowControlTypes(context, {csrf}) {
            getControlTypes({csrf}).then((resp : TGetControlTypesOutput) => {
                context.commit('setControlTypes', resp.data)
            })
        }
    }
}

let store = new Vuex.Store<MetadataStoreState>(metaDataStore)
export default store
