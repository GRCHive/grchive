import Vue from 'vue'
import Vuex, { StoreOptions } from 'vuex'
import { getProcessFlowIOTypes, TGetProcessFlowIOTypesInput, TGetProcessFlowIOTypesOutput } from './api/apiProcessFlowIO'
import { getProcessFlowNodeTypes, TGetProcessFlowNodeTypesInput, TGetProcessFlowNodeTypesOutput } from './api/apiProcessFlowNodes'
import { getControlTypes, TGetControlTypesInput, TGetControlTypesOutput } from './api/apiControls'
import { getAllOrgUsers, TGetAllOrgUsersInput, TGetAllOrgUsersOutput } from './api/apiUsers'

interface MetadataStoreState {
    ioTypes: ProcessFlowIOType[]
    idToIoTypes: Record<number, ProcessFlowIOType>
    nodeTypes: ProcessFlowNodeType[]
    idToNodeTypes: Record<number, ProcessFlowNodeType>
    controlTypes: ProcessFlowControlType[]
    idToControlTypes: Record<number, ProcessFlowControlType>
    controlTypeInitialized: boolean
    availableUsers: User[]
    idToUsers: Record<number, User>
    usersInitialized: boolean
}

const metaDataStore: StoreOptions<MetadataStoreState> = {
    state: {
        ioTypes: [] as ProcessFlowIOType[],
        idToIoTypes: Object() as Record<number, ProcessFlowIOType>,
        nodeTypes: [] as ProcessFlowNodeType[],
        idToNodeTypes: Object() as Record<number, ProcessFlowNodeType>,
        controlTypes: [] as ProcessFlowControlType[],
        idToControlTypes: Object() as Record<number, ProcessFlowControlType>,
        controlTypeInitialized: false,
        availableUsers: [] as User[],
        idToUsers: Object() as Record<number, User>,
        usersInitialized: false
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
            state.controlTypeInitialized = true
        },
        setUsers(state, inUsers : User[]) {
            state.availableUsers = inUsers
            state.idToUsers = Object() as Record<number, User>
            for (let user of inUsers) {
                Vue.set(state.idToUsers, user.Id, user)
            }
            state.usersInitialized = true
        }
    },
    actions: {
        // Initialization functions
        initialize(context, data) {
            context.dispatch('initializeProcessFlowIOTypes')
            context.dispatch('initializeProcessFlowNodeTypes')
            context.dispatch('initializeProcessFlowControlTypes')
            context.dispatch('initializeUsers', data)
        },
        initializeProcessFlowIOTypes(context) {
            getProcessFlowIOTypes(<TGetProcessFlowIOTypesInput>{
            }).then((resp : TGetProcessFlowIOTypesOutput) => {
                context.commit('setIoTypes', resp.data)
            })
        },
        initializeProcessFlowNodeTypes(context) {
            getProcessFlowNodeTypes(<TGetProcessFlowNodeTypesInput>{
            }).then((resp : TGetProcessFlowNodeTypesOutput) => {
                context.commit('setNodeTypes', resp.data)
            })
        },
        initializeProcessFlowControlTypes(context) {
            getControlTypes(<TGetControlTypesInput>{
            }).then((resp : TGetControlTypesOutput) => {
                context.commit('setControlTypes', resp.data)
            })
        },
        initializeUsers(context, {orgGroupId}) {
            if (!orgGroupId && orgGroupId != 0) {
                return
            }

            getAllOrgUsers(<TGetAllOrgUsersInput>{
                org: orgGroupId
            }).then((resp : TGetAllOrgUsersOutput) => {
                context.commit('setUsers', resp.data)
            })
        }
    },
}

let store = new Vuex.Store<MetadataStoreState>(metaDataStore)
export default store
