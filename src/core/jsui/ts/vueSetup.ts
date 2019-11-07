import 'core-js/es/promise'
import Vue from 'vue'
import Vuex, { StoreOptions } from 'vuex'
import VueRouter from 'vue-router'
import Vuetify from 'vuetify/lib'

Vue.use(Vuex)
Vue.use(VueRouter)
Vue.use(Vuetify)

import { deleteProcessFlowEdge, TDeleteProcessFlowEdgeInput, TDeleteProcessFlowEdgeOutput } from './api/apiProcessFlowEdges'
import { deleteProcessFlowNode, TDeleteProcessFlowNodeInput, TDeleteProcessFlowNodeOutput } from './api/apiProcessFlowNodes'
import RelationshipMap from './relationship'
import { FullProcessFlowData } from './processFlow'
import VuexState from './processFlowState'
import { getFullProcessFlow, TGetFullProcessFlowInput, TGetFullProcessFlowOutput } from './api/apiProcessFlow'

let mutationObservers = []
const opts = {}

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
        deleteProcessFlow(state, flowId) {
            if (state.currentProcessFlowFullData.FlowId == flowId) {
                state.currentProcessFlowIndex = 0
                state.currentProcessFlowFullData = {} as FullProcessFlowData
                state.selectedNodeId = -1
                state.selectedEdgeId = -1
            }

            state.allProcessFlowBasicData.splice(
                state.allProcessFlowBasicData.findIndex(
                    (ele) => ele.Id == flowId),
                1)
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
            let currentNodeData = state.currentProcessFlowFullData.Nodes[nodeId]
            currentNodeData.Name = node.Name
            currentNodeData.Description = node.Description
            currentNodeData.NodeTypeId = node.NodeTypeId
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
                    (ele) => { return ele === edgeId }),
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
                    (ele) => { return ele == nodeId}),
                1)
            Vue.delete(
                state.currentProcessFlowFullData.Nodes,
                nodeId)
        },
        setRisk(state, risk : ProcessFlowRisk) {
            if (!(risk.Id in state.currentProcessFlowFullData.Risks)) {
                state.currentProcessFlowFullData.RiskKeys.push(risk.Id)
                Vue.set(state.currentProcessFlowFullData.Risks, risk.Id, risk)
            } else {
                state.currentProcessFlowFullData.Risks[risk.Id].Name = risk.Name
                state.currentProcessFlowFullData.Risks[risk.Id].Description = risk.Description
            }
        },
        deleteRiskFromNode(state, {nodeId, riskIds}) {
            for (let riskId of riskIds) {
                state.currentProcessFlowFullData.NodeRiskRelationships.delete(
                    state.currentProcessFlowFullData.Nodes[nodeId],
                    state.currentProcessFlowFullData.Risks[riskId])
            }
        },
        deleteRiskGlobal(state, riskIds) {
            for (let id of riskIds) {
                state.currentProcessFlowFullData.NodeRiskRelationships.deleteB(
                    state.currentProcessFlowFullData.Risks[id]
                )
                Vue.delete(state.currentProcessFlowFullData.Risks, id)
                state.currentProcessFlowFullData.RiskKeys.splice(
                    state.currentProcessFlowFullData.RiskKeys.findIndex((ele) => ele == id),
                    1)

            }
        },
        deleteNodeFromRisks(state, nodeId) {
            state.currentProcessFlowFullData.NodeRiskRelationships.deleteA(
                state.currentProcessFlowFullData.Nodes[nodeId]
            )
        },
        addRisksToNode(state, {nodeId, riskIds}) {
            for (let id of riskIds) {
                state.currentProcessFlowFullData.NodeRiskRelationships.add(
                    state.currentProcessFlowFullData.Nodes[nodeId],
                    state.currentProcessFlowFullData.Risks[id]
                )
            }
        },
        setControl(state, {control}) {
            if (control.Id in state.currentProcessFlowFullData.Controls) {
                state.currentProcessFlowFullData.Controls[control.Id].Id = control.Id
                state.currentProcessFlowFullData.Controls[control.Id].Name = control.Name
                state.currentProcessFlowFullData.Controls[control.Id].Description = control.Description
                state.currentProcessFlowFullData.Controls[control.Id].ControlTypeId = control.ControlTypeId
                state.currentProcessFlowFullData.Controls[control.Id].FrequencyType = control.FrequencyType
                state.currentProcessFlowFullData.Controls[control.Id].FrequencyInterval = control.FrequencyInterval
                state.currentProcessFlowFullData.Controls[control.Id].OwnerId = control.OwnerId
            } else {
                Vue.set(
                    state.currentProcessFlowFullData.Controls,
                    control.Id,
                    control)
                state.currentProcessFlowFullData.ControlKeys.push(control.Id)
            }
        },
        addControlToNode(state, {controlId, nodeId}) {
            state.currentProcessFlowFullData.NodeControlRelationships.add(
                state.currentProcessFlowFullData.Nodes[nodeId],
                state.currentProcessFlowFullData.Controls[controlId]
            )
        },
        addControlToRisk(state, {controlId, riskId}) {
            state.currentProcessFlowFullData.RiskControlRelationships.add(
                state.currentProcessFlowFullData.Risks[riskId],
                state.currentProcessFlowFullData.Controls[controlId]
            )
        },
        deleteControlFromRiskNode(state, {controlId, nodeId, riskId}) {
            state.currentProcessFlowFullData.NodeControlRelationships.delete(
                state.currentProcessFlowFullData.Nodes[nodeId],
                state.currentProcessFlowFullData.Controls[controlId]
            )

            state.currentProcessFlowFullData.RiskControlRelationships.delete(
                state.currentProcessFlowFullData.Risks[riskId],
                state.currentProcessFlowFullData.Controls[controlId]
            )
        },
        deleteControlGlobal(state, {controlId}) {
            state.currentProcessFlowFullData.NodeControlRelationships.deleteB(
                state.currentProcessFlowFullData.Controls[controlId]
            )

            state.currentProcessFlowFullData.RiskControlRelationships.deleteB(
                state.currentProcessFlowFullData.Controls[controlId]
            )

            Vue.delete(state.currentProcessFlowFullData.Controls, controlId)
            state.currentProcessFlowFullData.ControlKeys.splice(
                state.currentProcessFlowFullData.ControlKeys.findIndex((ele) => ele == controlId),
                1)
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
        requestSetCurrentProcessFlowIndex(context, {index}) {
            context.commit('setSelectedProcessFlowEdge', -1)
            context.commit('setSelectedProcessFlowNode', -1)
            context.commit('setCurrentProcessFlowIndex', index)
            context.dispatch('refreshCurrentProcessFlowFullData')
        },
        refreshCurrentProcessFlowFullData(context) {
            const id = context.getters.currentProcessFlowBasicData.Id
            context.commit('setFullProcessFlowRequestedId', id)
            getFullProcessFlow(id, <TGetFullProcessFlowInput>{
            }).then(
                (resp : TGetFullProcessFlowOutput) => {
                    let newData = <FullProcessFlowData>{
                        FlowId: id,
                        Nodes: Object(),
                        NodeKeys: [] as number[],
                        Edges: Object(),
                        EdgeKeys: [] as number[],
                        Inputs: Object(),
                        Outputs: Object(),
                        Risks: Object(),
                        RiskKeys: [] as number[],
                        Controls: Object(),
                        ControlKeys: [] as number[],
                        NodeRiskRelationships: Vue.observable(new RelationshipMap<ProcessFlowNode,ProcessFlowRisk>()),
                        NodeControlRelationships: Vue.observable(new RelationshipMap<ProcessFlowNode,ProcessFlowControl>()),
                        RiskControlRelationships: Vue.observable(new RelationshipMap<ProcessFlowRisk,ProcessFlowControl>()),
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

                    for (let data of resp.data.Risks) {
                        newData.Risks[data.Id] = data
                        newData.RiskKeys.push(data.Id)
                    }

                    for (let data of resp.data.Controls) {
                        newData.Controls[data.Id] = data
                        newData.ControlKeys.push(data.Id)
                    }

                    for (let data of resp.data.NodeRisk) {
                        newData.NodeRiskRelationships.add(
                            newData.Nodes[data.NodeId],
                            newData.Risks[data.RiskId])
                    }

                    for (let data of resp.data.NodeControl) {
                        newData.NodeControlRelationships.add(
                            newData.Nodes[data.NodeId],
                            newData.Controls[data.ControlId])
                    }

                    for (let data of resp.data.RiskControl) {
                        newData.RiskControlRelationships.add(
                            newData.Risks[data.RiskId],
                            newData.Controls[data.ControlId])
                    }
                    
                    context.commit('setCurrentProcessFlowFullData', newData)
                    context.commit('setFullProcessFlowRequestedId', -1)
                }
            ).catch(
                (err : any) => {
                    // TODO: Somehow display something went wrong??
                    console.log("Failed to obtain process flow.", err)
                    context.commit('setCurrentProcessFlowFullData', {} as FullProcessFlowData)
                    context.commit('setFullProcessFlowRequestedId', -1)
                }
            )
        },
        requestDeletionOfSelection(context) {
            if (context.state.selectedEdgeId != -1) {
                let edgeId = context.state.selectedEdgeId
                deleteProcessFlowEdge({
                    edgeId: edgeId
                }).then(() => {
                    context.commit('setSelectedProcessFlowEdge', -1)
                    context.commit('deleteEdgeById', edgeId)
                })
            }

            if (context.state.selectedNodeId != -1) {
                let nodeId = context.state.selectedNodeId
                deleteProcessFlowNode({
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
            context.commit('deleteNodeFromRisks', nodeId)

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
        deleteBatchRisks(context, {nodeId, riskIds, global}) {
            context.commit('deleteRiskFromNode', {nodeId, riskIds})
            if (global) {
                context.commit('deleteRiskGlobal', riskIds)
            }
        },
        deleteBatchControls(context, {nodeId, controlIds, riskIds, global}) {
            for (let i = 0; i < controlIds.length; ++i) {
                context.commit('deleteControlFromRiskNode', {
                    controlId: controlIds[i],
                    nodeId: nodeId,
                    riskIds: riskIds[i]
                })

                if (global) {
                    context.commit('deleteControlGlobal', {controlId: controlIds[i]})
                }
            }
        }
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
        },
        risksForNode: (state) => (nodeId : number) : ProcessFlowRisk[] => {
            return state.currentProcessFlowFullData.NodeRiskRelationships.changed &&
                state.currentProcessFlowFullData.NodeRiskRelationships.getB(
                    state.currentProcessFlowFullData.Nodes[nodeId]
                )
        },
        controlsForRiskNode: (state) => (riskId : number, nodeId : number) : RiskControl[] => {
            let controlsForNode = 
                state.currentProcessFlowFullData.NodeControlRelationships.changed && 
                state.currentProcessFlowFullData.NodeControlRelationships.getB(
                    state.currentProcessFlowFullData.Nodes[nodeId]
                )
            let controlsForRisk = 
                state.currentProcessFlowFullData.RiskControlRelationships.changed && 
                state.currentProcessFlowFullData.RiskControlRelationships.getB(
                    state.currentProcessFlowFullData.Risks[riskId]
                )

            return controlsForNode.filter(val => controlsForRisk.includes(val)).map(ele => ({
                risk: state.currentProcessFlowFullData.Risks[riskId],
                control: ele
            }))
        }
    }
}

export default {
    vuetify: new Vuetify(opts),
    store: new Vuex.Store<VuexState>(store),
    currentRouter: {} as VueRouter
}
