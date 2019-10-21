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
import { deleteProcessFlowEdge } from './api/apiProcessFlowEdges'
import { deleteProcessFlowNode } from './api/apiProcessFlowNodes'

let mutationObservers = []
const opts = {}

interface FullProcessFlowDataResponse {
    data : FullProcessFlowResponseData
}

const store : StoreOptions<VuexState> = {
    state: {
        primaryNavBarWidth: 256,
        allProcessFlowBasicData: [],
        currentProcessFlowIndex : 0,
        currentProcessFlowFullData: {} as FullProcessFlowData,
        fullProcessFlowRequestedId: -1,
        selectedNodeId: -1,
        selectedEdgeId: -1
    },
    mutations: {
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
        },
        setFullProcessFlowRequestedId(state, data) {
            state.fullProcessFlowRequestedId = data
        },
        setSelectedProcessFlowNode(state, id) {
            state.selectedNodeId = id
        },
        setSelectedProcessFlowEdge(state, id) {
            state.selectedEdgeId = id
        },
        addNodeInput(state, {nodeId, input}) {
            state.currentProcessFlowFullData.Nodes[nodeId].Inputs.push(input)
            Vue.set(
                state.currentProcessFlowFullData.Inputs,
                input.Id,
                input)
        },
        addNodeOutput(state, {nodeId, output}) {
            state.currentProcessFlowFullData.Nodes[nodeId].Outputs.push(output)
            Vue.set(
                state.currentProcessFlowFullData.Outputs,
                output.Id,
                output)
        },
        removeNodeInput(state, {nodeId, inputId}) {
            let idx : number = state.currentProcessFlowFullData.Nodes[nodeId].Inputs.findIndex(
                (ele : ProcessFlowInputOutput) => {
                    return ele.Id == inputId
                })
            if (idx != -1) {
                state.currentProcessFlowFullData.Nodes[nodeId].Inputs.splice(idx, 1)
                Vue.delete(
                    state.currentProcessFlowFullData.Inputs,
                    inputId)
            }
        },
        removeNodeOutput(state, {nodeId, outputId}) {
            let idx : number = state.currentProcessFlowFullData.Nodes[nodeId].Outputs.findIndex(
                (ele : ProcessFlowInputOutput) => {
                    return ele.Id == outputId
                })
            if (idx != -1) {
                state.currentProcessFlowFullData.Nodes[nodeId].Outputs.splice(idx, 1)
                Vue.delete(
                    state.currentProcessFlowFullData.Outputs,
                    outputId)
            }
        },
        updateNodeInputOutput(state, {nodeId, io, isInput}) {
            let relevantArr : ProcessFlowInputOutput[] = 
                isInput ? 
                    state.currentProcessFlowFullData.Nodes[nodeId].Inputs :
                    state.currentProcessFlowFullData.Nodes[nodeId].Outputs;

            let idx : number = relevantArr.findIndex(
                (ele : ProcessFlowInputOutput) => {
                    return ele.Id == io.Id
                })
            if (idx != -1) {
                Vue.set(relevantArr, idx, io)
            }
        },
        updateNodePartial(state, {nodeId, node}) {
            let currentNodeData = {...state.currentProcessFlowFullData.Nodes[nodeId]}
            currentNodeData.Name = node.Name
            currentNodeData.Description = node.Description
            currentNodeData.NodeTypeId = node.NodeTypeId
            Vue.set(state.currentProcessFlowFullData.Nodes, nodeId, currentNodeData)
        },
        addNewEdge(state, {edge}) {
            if (edge.Id in state.currentProcessFlowFullData.Edges) {
                return
            }
            state.currentProcessFlowFullData.EdgeKeys.push(edge.Id)
            Vue.set(
                state.currentProcessFlowFullData.Edges,
                edge.Id,
                edge)
        },
        deleteEdgeById(state, edgeId) {
            if (!(edgeId in state.currentProcessFlowFullData.Edges)) {
                return
            }
            state.currentProcessFlowFullData.EdgeKeys.splice(
                state.currentProcessFlowFullData.EdgeKeys.findIndex(
                    (ele) => { ele == edgeId }),
                1)
            Vue.delete(
                state.currentProcessFlowFullData.Edges,
                edgeId)
        },
        deleteNodeById(state, nodeId) {
            if (!(nodeId in state.currentProcessFlowFullData.Nodes)) {
                return
            }

            state.currentProcessFlowFullData.NodeKeys.splice(
                state.currentProcessFlowFullData.NodeKeys.findIndex(
                    (ele) => { ele == nodeId}),
                1)
            Vue.delete(
                state.currentProcessFlowFullData.Nodes,
                nodeId)
        }
    },
    actions: {
        mountPrimaryNavBar(context, nav) {
            context.commit('changePrimaryNavBarWidth', parseInt(nav.$el.style.width, 10))
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
            const id = context.getters.currentProcessFlowBasicData.Id
            const baseUrl = createGetProcessFlowFullDataUrl(id)
            const queryParams = qs.stringify({
                csrf
            })
            context.commit('setFullProcessFlowRequestedId', id)
            axios.get(baseUrl + "?" + queryParams).then(
                (resp : FullProcessFlowDataResponse) => {
                    let newData = <FullProcessFlowData>{
                        FlowId: id,
                        Nodes: Object(),
                        NodeKeys: [] as number[],
                        Edges: Object(),
                        EdgeKeys: [] as number[],
                        Inputs: Object(),
                        Outputs: Object()
                    }
                    for (let data of resp.data.Nodes) {
                        newData.Nodes[data.Id] = data
                        newData.NodeKeys.push(data.Id)
                        for (let inp of data.Inputs) {
                            newData.Inputs[inp.Id] = inp
                        }

                        for (let inp of data.Outputs) {
                            newData.Outputs[inp.Id] = inp
                        }
                    }
                    for (let data of resp.data.Edges) {
                        newData.Edges[data.Id] = data
                        newData.EdgeKeys.push(data.Id)
                    }
                    context.commit('setCurrentProcessFlowFullData', newData)
                    context.commit('setFullProcessFlowRequestedId', -1)
                }
            ).catch(
                (_) => {
                    // TODO: Somehow display something went wrong??
                    console.log("Failed to obtain process flow.")
                    context.commit('setCurrentProcessFlowFullData', {} as FullProcessFlowData)
                    context.commit('setFullProcessFlowRequestedId', -1)
                }
            )
        },
        requestDeletionOfSelection(context, {csrf}) {
            if (context.state.selectedEdgeId != -1) {
                let edgeId = context.state.selectedEdgeId
                deleteProcessFlowEdge({
                    csrf: csrf,
                    edgeId: edgeId
                }).then(() => {
                    context.commit('setSelectedProcessFlowEdge', -1)
                    context.commit('deleteEdgeById', edgeId)
                })
            }

            if (context.state.selectedNodeId != -1) {
                let nodeId = context.state.selectedNodeId
                deleteProcessFlowNode({
                    csrf: csrf,
                    nodeId: nodeId
                }).then(() => {
                    context.commit('setSelectedProcessFlowNode', -1)
                    context.dispatch('deleteNodeById', nodeId)
                })
            }
        },
        deleteNodeById(context, nodeId) {
            if (!(nodeId in context.state.currentProcessFlowFullData.Nodes)) {
                return
            }
        
            let node : ProcessFlowNode = context.state.currentProcessFlowFullData.Nodes[nodeId]

            let inputsToRemove = []
            for (let inp of node.Inputs) {
                inputsToRemove.push(inp.Id)
            }

            let outputsToRemove = []
            for (let out of node.Outputs) {
                outputsToRemove.push(out.Id)
            }
            context.dispatch('deleteBatchNodeInput', {nodeId: nodeId, inputs: inputsToRemove})
            context.dispatch('deleteBatchNodeOutput', {nodeId: nodeId, outputs: outputsToRemove})
            context.commit('deleteNodeById', nodeId)

        },
        deleteBatchNodeInput(context, {nodeId, inputs}) {
            const set = new Set(inputs)
            for (let i of inputs) {
                context.commit('removeNodeInput', {nodeId: nodeId, inputId: i})
            }

            for (let i = context.state.currentProcessFlowFullData.EdgeKeys.length - 1; i >= 0; --i) {
                const edgeKey = context.state.currentProcessFlowFullData.EdgeKeys[i]
                const edge = context.state.currentProcessFlowFullData.Edges[edgeKey]
                if (set.has(edge.InputIoId)) {
                    context.commit('deleteEdgeById', edge.Id)
                }
            }
        },
        deleteBatchNodeOutput(context, {nodeId, outputs}) {
            const set = new Set(outputs)
            for (let i of outputs) {
                context.commit('removeNodeOutput', {nodeId: nodeId, outputId: i})
            }

            for (let i = context.state.currentProcessFlowFullData.EdgeKeys.length - 1; i >= 0; --i) {
                const edgeKey = context.state.currentProcessFlowFullData.EdgeKeys[i]
                const edge = context.state.currentProcessFlowFullData.Edges[edgeKey]
                if (set.has(edge.OutputIoId)) {
                    context.commit('deleteEdgeById', edge.Id)
                }
            }
        },
    },
    getters: {
        currentProcessFlowBasicData: (state) => {
            return state.allProcessFlowBasicData[state.currentProcessFlowIndex]
        },
        isFullRequestInProgress: (state) => {
            return state.fullProcessFlowRequestedId != -1
        },
        isNodeSelected: (state) => {
            return state.selectedNodeId != -1
        },
        nodeInfo: (state) => (nodeId : number) : ProcessFlowNode => {
            return state.currentProcessFlowFullData.Nodes[nodeId]
        },
        currentNodeInfo: (state, getters) : ProcessFlowNode  => {
            if (!getters.isNodeSelected) {
                return {} as ProcessFlowNode
            }
            return getters.nodeInfo(state.selectedNodeId)
        }
    }
}

export default {
    vuetify: new Vuetify(opts),
    store: new Vuex.Store<VuexState>(store),
    currentRouter: {} as VueRouter
}
