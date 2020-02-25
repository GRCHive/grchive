import Vue from 'vue'
import Vuex, { StoreOptions } from 'vuex'
import Vuetify from 'vuetify/lib'
import { deleteProcessFlowEdge, TDeleteProcessFlowEdgeInput, TDeleteProcessFlowEdgeOutput } from './api/apiProcessFlowEdges'
import { deleteProcessFlowNode, TDeleteProcessFlowNodeInput, TDeleteProcessFlowNodeOutput } from './api/apiProcessFlowNodes'
import RelationshipMap from './relationship'
import { FullProcessFlowData } from './processFlow'
import VuexState from './processFlowState'
import { getFullProcessFlow, TGetFullProcessFlowInput, TGetFullProcessFlowOutput } from './api/apiProcessFlow'
import {
    TAllNodeSystemLinkOutput,
    allNodeSystemLink
} from './api/apiNodeSystemLinks'
import {
    TAllNodeGLLinkOutput,
    allNodeGLLink
} from './api/apiNodeGLLinks'

import { System } from './systems'
import { GeneralLedger } from './generalLedger'
import { PageParamsStore } from './pageParams'

let mutationObservers = []
const opts = {}

function sortIo(a : ProcessFlowInputOutput, b : ProcessFlowInputOutput) : number {
    return (a.TypeId - b.TypeId) || (a.IoOrder - b.IoOrder)
}

const store : StoreOptions<VuexState> = {
    state: {
        primaryNavBarWidth: 256,
        currentProcessFlowBasicData: null,
        currentProcessFlowFullData: null,
        selectedNodeId: -1,
        selectedEdgeId: -1
    },
    mutations: {
        setProcessFlowBasicData(state, data) {
            state.currentProcessFlowBasicData = data
        },
        changePrimaryNavBarWidth(state, width) {
            state.primaryNavBarWidth = width
        },
        setCurrentProcessFlowFullData(state, data) {
            state.currentProcessFlowFullData = data
        },
        setSelectedProcessFlowNode(state, id) {
            state.selectedNodeId = id
        },
        setSelectedProcessFlowEdge(state, id) {
            state.selectedEdgeId = id
        },
        addNodeInput(state, {nodeId, input}) {
            if (!state.currentProcessFlowFullData) {
                return
            }

            state.currentProcessFlowFullData.Nodes[nodeId].Inputs.push(input)
            state.currentProcessFlowFullData.Nodes[nodeId].Inputs.sort(sortIo)
            Vue.set(
                state.currentProcessFlowFullData.Inputs,
                input.Id,
                input)
        },
        addNodeOutput(state, {nodeId, output}) {
            if (!state.currentProcessFlowFullData) {
                return
            }
            state.currentProcessFlowFullData.Nodes[nodeId].Outputs.push(output)
            state.currentProcessFlowFullData.Nodes[nodeId].Outputs.sort(sortIo)
            Vue.set(
                state.currentProcessFlowFullData.Outputs,
                output.Id,
                output)
        },
        removeNodeInput(state, {nodeId, inputId}) {
            if (!state.currentProcessFlowFullData) {
                return
            }
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
            if (!state.currentProcessFlowFullData) {
                return
            }
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
            if (!state.currentProcessFlowFullData) {
                return
            }

            let relevantNodeArr : ProcessFlowInputOutput[] = 
                isInput ? 
                    state.currentProcessFlowFullData.Nodes[nodeId].Inputs :
                    state.currentProcessFlowFullData.Nodes[nodeId].Outputs;

            let nodeIdx : number = relevantNodeArr.findIndex(
                (ele : ProcessFlowInputOutput) => {
                    return ele.Id == io.Id
                })
            if (nodeIdx != -1) {
                Vue.set(relevantNodeArr, nodeIdx, io)
                relevantNodeArr.sort(sortIo)
            }

            let relevantArr : Record<number, ProcessFlowInputOutput> = 
                isInput ?
                    state.currentProcessFlowFullData.Inputs :
                    state.currentProcessFlowFullData.Outputs;

            if (!!relevantArr[io.Id]) {
                Vue.set(relevantArr, io.Id, io)
            }
        },
        updateNodePartial(state, {nodeId, node}) {
            if (!state.currentProcessFlowFullData) {
                return
            }
            let currentNodeData = state.currentProcessFlowFullData.Nodes[nodeId]
            currentNodeData.Name = node.Name
            currentNodeData.Description = node.Description
            currentNodeData.NodeTypeId = node.NodeTypeId
        },
        addNewEdge(state, {edge}) {
            if (!state.currentProcessFlowFullData) {
                return
            }
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
            if (!state.currentProcessFlowFullData) {
                return
            }
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
            if (!state.currentProcessFlowFullData) {
                return
            }
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
            if (!state.currentProcessFlowFullData) {
                return
            }
            if (!(risk.Id in state.currentProcessFlowFullData.Risks)) {
                state.currentProcessFlowFullData.RiskKeys.push(risk.Id)
                Vue.set(state.currentProcessFlowFullData.Risks, risk.Id, risk)
            } else {
                state.currentProcessFlowFullData.Risks[risk.Id].Name = risk.Name
                state.currentProcessFlowFullData.Risks[risk.Id].Description = risk.Description
            }
        },
        deleteRiskFromNode(state, {nodeId, riskIds}) {
            if (!state.currentProcessFlowFullData) {
                return
            }
            for (let riskId of riskIds) {
                state.currentProcessFlowFullData.NodeRiskRelationships.delete(
                    state.currentProcessFlowFullData.Nodes[nodeId],
                    state.currentProcessFlowFullData.Risks[riskId])
            }
        },
        deleteRiskGlobal(state, riskIds) {
            if (!state.currentProcessFlowFullData) {
                return
            }
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
            if (!state.currentProcessFlowFullData) {
                return
            }
            state.currentProcessFlowFullData.NodeRiskRelationships.deleteA(
                state.currentProcessFlowFullData.Nodes[nodeId]
            )
        },
        addRisksToNode(state, {nodeId, riskIds}) {
            if (!state.currentProcessFlowFullData) {
                return
            }
            for (let id of riskIds) {
                state.currentProcessFlowFullData.NodeRiskRelationships.add(
                    state.currentProcessFlowFullData.Nodes[nodeId],
                    state.currentProcessFlowFullData.Risks[id]
                )
            }
        },
        setControl(state, {control}) {
            if (!state.currentProcessFlowFullData) {
                return
            }
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
            if (!state.currentProcessFlowFullData) {
                return
            }
            state.currentProcessFlowFullData.NodeControlRelationships.add(
                state.currentProcessFlowFullData.Nodes[nodeId],
                state.currentProcessFlowFullData.Controls[controlId]
            )
        },
        addControlToRisk(state, {controlId, riskId}) {
            if (!state.currentProcessFlowFullData) {
                return
            }
            state.currentProcessFlowFullData.RiskControlRelationships.add(
                state.currentProcessFlowFullData.Risks[riskId],
                state.currentProcessFlowFullData.Controls[controlId]
            )
        },
        deleteControlFromRiskNode(state, {controlId, nodeId, riskId}) {
            if (!state.currentProcessFlowFullData) {
                return
            }

            if (nodeId != -1) {
                state.currentProcessFlowFullData.NodeControlRelationships.delete(
                    state.currentProcessFlowFullData.Nodes[nodeId],
                    state.currentProcessFlowFullData.Controls[controlId]
                )
            }

            if (riskId != -1) {
                state.currentProcessFlowFullData.RiskControlRelationships.delete(
                    state.currentProcessFlowFullData.Risks[riskId],
                    state.currentProcessFlowFullData.Controls[controlId]
                )
            }
        },
        deleteControlGlobal(state, {controlId}) {
            if (!state.currentProcessFlowFullData) {
                return
            }
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
        },
        addNodeSystemLink(state, {nodeId, system}) {
            if (!state.currentProcessFlowFullData) {
                return
            }

            let arr : System[] | null = state.currentProcessFlowFullData.NodeSystemLinks[nodeId]
            if (arr === null) {
                return
            }

            arr.push(system)
        },
        deleteNodeSystemLink(state, {nodeId, systemId}) {
            if (!state.currentProcessFlowFullData) {
                return
            }

            let arr : System[] | null = state.currentProcessFlowFullData.NodeSystemLinks[nodeId]
            if (arr === null) {
                return
            }

            let idx : number = arr.findIndex((ele : System) => ele.Id == systemId)
            if (idx == -1) {
                return
            }

            arr.splice(idx, 1)
        },
        addNodeGLLink(state, {nodeId, account}) {
            if (!state.currentProcessFlowFullData) {
                return
            }

            let gl: GeneralLedger | null = state.currentProcessFlowFullData.NodeGLLinks[nodeId]
            if (gl === null) {
                return
            }

            gl.addRawAccount(account)
        },
        deleteNodeGLLink(state, {nodeId, accountId}) {
            if (!state.currentProcessFlowFullData) {
                return
            }

            let gl: GeneralLedger | null = state.currentProcessFlowFullData.NodeGLLinks[nodeId]
            if (gl === null) {
                return
            }

            gl.removeAccount(accountId)
        },

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
        refreshCurrentProcessFlowFullData(context, id) {
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
                        NodeSystemLinks: Object(),
                        NodeGLLinks: Object(),
                    }
        
                    for (let data of resp.data.Graph.Nodes) {
                        newData.Nodes[data.Id] = data
                        newData.NodeKeys.push(data.Id)
                        for (let inp of data.Inputs) {
                            newData.Inputs[inp.Id] = inp
                        }

                        for (let inp of data.Outputs) {
                            newData.Outputs[inp.Id] = inp
                        }

                        newData.Nodes[data.Id].Inputs.sort(sortIo)
                        newData.Nodes[data.Id].Outputs.sort(sortIo)
                    }
                    for (let data of resp.data.Graph.Edges) {
                        newData.Edges[data.Id] = data
                        newData.EdgeKeys.push(data.Id)
                    }

                    for (let data of resp.data.Graph.Risks) {
                        newData.Risks[data.Id] = data
                        newData.RiskKeys.push(data.Id)
                    }

                    for (let data of resp.data.Graph.Controls) {
                        newData.Controls[data.Id] = data
                        newData.ControlKeys.push(data.Id)
                    }

                    for (let data of resp.data.Graph.NodeRisk) {
                        newData.NodeRiskRelationships.add(
                            newData.Nodes[data.NodeId],
                            newData.Risks[data.RiskId])
                    }

                    for (let data of resp.data.Graph.NodeControl) {
                        newData.NodeControlRelationships.add(
                            newData.Nodes[data.NodeId],
                            newData.Controls[data.ControlId])
                    }

                    for (let data of resp.data.Graph.RiskControl) {
                        newData.RiskControlRelationships.add(
                            newData.Risks[data.RiskId],
                            newData.Controls[data.ControlId])
                    }
                    
                    context.commit('setCurrentProcessFlowFullData', newData)
                    context.commit('setProcessFlowBasicData', resp.data.Basic)
                }
            ).catch(
                (err : any) => {
                    // TODO: Somehow display something went wrong??
                    console.log("Failed to obtain process flow.", err)
                    context.commit('setCurrentProcessFlowFullData', null)
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
            if (!context.state.currentProcessFlowFullData) {
                return
            }

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
            if (!context.state.currentProcessFlowFullData) {
                return
            }
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
            if (!context.state.currentProcessFlowFullData) {
                return
            }
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
                    riskId: riskIds[i]
                })

                if (global) {
                    context.commit('deleteControlGlobal', {controlId: controlIds[i]})
                }
            }
        }
    },
    getters: {
        glLinkedToNode(state): (nodeId : number) => GeneralLedger | null {
            if (!state.currentProcessFlowFullData) {
                return (_ : number) => null
            }

            let fullMap : Record<number, GeneralLedger | null> = state.currentProcessFlowFullData.NodeGLLinks
            return (nodeId : number) : GeneralLedger | null => {
                let ledger : GeneralLedger | null = null
                if (nodeId in fullMap) {
                    ledger = fullMap[nodeId]
                }

                if (ledger === null) {
                    allNodeGLLink({
                        nodeId: nodeId,
                        orgId: PageParamsStore.state.organization!.Id,
                    }).then((resp : TAllNodeGLLinkOutput) => {
                        let gl : GeneralLedger = new GeneralLedger()
                        gl.rebuildGL(resp.data.Categories!, resp.data.Accounts!)
                        Vue.set(fullMap, nodeId, gl)
                    })
                    return null
                }
                return ledger
            }
        },
        systemsLinkedToNode(state) : (nodeId : number) => System[] | null {
            if (!state.currentProcessFlowFullData) {
                return (_ : number) => null
            }

            let fullMap : Record<number, System[] | null> = state.currentProcessFlowFullData.NodeSystemLinks
            return (nodeId : number) : System[] | null => {
                let systems : System[] | null = null
                if (nodeId in fullMap) {
                    systems = fullMap[nodeId]
                }

                if (systems === null) {
                    allNodeSystemLink({
                        nodeId: nodeId,
                        orgId: PageParamsStore.state.organization!.Id,
                    }).then((resp : TAllNodeSystemLinkOutput) => {
                        Vue.set(fullMap, nodeId, <System[]>resp.data)
                    })

                    return null
                }

                return systems
            }
        },
        isNodeSelected: (state) => {
            return state.selectedNodeId != -1
        },
        nodeInfo: (state) => (nodeId : number) : ProcessFlowNode | null => {
            if (!state.currentProcessFlowFullData) {
                return null
            }
            return state.currentProcessFlowFullData.Nodes[nodeId]
        },
        currentNodeInfo: (state, getters) : ProcessFlowNode | null  => {
            if (!getters.isNodeSelected) {
                return null
            }
            return getters.nodeInfo(state.selectedNodeId)
        },
        risksForNode: (state) => (nodeId : number) : ProcessFlowRisk[] => {
            if (!state.currentProcessFlowFullData) {
                return []
            }
            return state.currentProcessFlowFullData.NodeRiskRelationships.changed &&
                state.currentProcessFlowFullData.NodeRiskRelationships.getB(
                    state.currentProcessFlowFullData.Nodes[nodeId]
                )
        },
        controlsForNode: (state) => (nodeId : number) : ProcessFlowControl[] => {
            if (!state.currentProcessFlowFullData) {
                return []
            }
            return state.currentProcessFlowFullData.NodeControlRelationships.changed &&
                state.currentProcessFlowFullData.NodeControlRelationships.getB(
                    state.currentProcessFlowFullData.Nodes[nodeId]
                )
        },
        controlsForRisk: (state) => (riskId: number) : ProcessFlowControl[] => {
            if (!state.currentProcessFlowFullData) {
                return []
            }
            return state.currentProcessFlowFullData.RiskControlRelationships.changed && 
                state.currentProcessFlowFullData.RiskControlRelationships.getB(
                    state.currentProcessFlowFullData.Risks[riskId]
                )
        },
        risksForControl: (state) => (controlId: number) : ProcessFlowRisk[] => {
            if (!state.currentProcessFlowFullData) {
                return []
            }
            return state.currentProcessFlowFullData.RiskControlRelationships.changed && 
                state.currentProcessFlowFullData.RiskControlRelationships.getA(
                    state.currentProcessFlowFullData.Controls[controlId]
                )
        },

        controlsForRiskNode: (state) => (riskId : number, nodeId : number) : RiskControl[] => {
            if (!state.currentProcessFlowFullData) {
                return []
            }
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
                risk: state.currentProcessFlowFullData!.Risks[riskId],
                control: ele
            }))
        },
        riskList(state) : ProcessFlowRisk[] {
            if (!state.currentProcessFlowFullData) {
                return []
            }
            return state.currentProcessFlowFullData.RiskKeys.map((ele : number) => state.currentProcessFlowFullData!.Risks[ele])
        },
        linkedRiskList(state) : ProcessFlowRisk[] {
            if (!state.currentProcessFlowFullData) {
                return []
            }

            let riskSet : Set<number> = new Set<number>()
            for (let nodeId of state.currentProcessFlowFullData.NodeKeys) {
                let risks : ProcessFlowRisk[] = state.currentProcessFlowFullData.NodeRiskRelationships.changed &&
                    state.currentProcessFlowFullData.NodeRiskRelationships.getB(
                        state.currentProcessFlowFullData.Nodes[nodeId]
                    )

                for (let r of risks) {
                    riskSet.add(r.Id)
                }
            }

            return [...riskSet].map((ele : number) => state.currentProcessFlowFullData!.Risks[ele]).sort(
                (a : ProcessFlowRisk, b : ProcessFlowRisk) : number => {
                    if (a.Name < b.Name) {
                        return -1
                    } else if (a.Name > b.Name) {
                        return 1
                    }
                    return 0
                })
        },
        controlList(state) : ProcessFlowControl[] {
            if (!state.currentProcessFlowFullData) {
                return []
            }
            return state.currentProcessFlowFullData.ControlKeys.map((ele : number) => state.currentProcessFlowFullData!.Controls[ele])
        }
    }
}

export default {
    vuetify: new Vuetify(opts),
    store: new Vuex.Store<VuexState>(store),
}
