import Vue from 'vue'
import Vuex, { StoreOptions } from 'vuex'
import MetadataStore from '../metadata'
import { connectProcessFlowNodeDisplaySettingsWebsocket } from '../websocket/processFlowNodeDisplaySettings'
import { FullProcessFlowData } from '../processFlow'
import VuexState from '../processFlowState'
import LocalSettings from '../localSettings'

// A Vuex store to share the layout of the process flow (nodes, plugs, etc.)
// across the entire application.
const NodeMargins = {
    left: 5,
    right: 5,
    top: 5,
    bottom: 5
}

const NodeIOMargins = {
    betweenGroups: 5,
    betweenPlugs: 10
}

const titleHeight: number = 26
const subtitleHeight: number = 21
const bodyHeight: number = 19
const inputOutputGap : number = 200
const plugWidth : number = 20
const plugHeight: number = 20

let websocketConnection : WebSocket

function processIOGroupLayout(layout : IOGroupLayout, initialTransform: TransformData) {
    let groupStartTransform = {...initialTransform}
    layout.transform = {...groupStartTransform}
    layout.titleTransform = <TransformData>{
        tx: 0,
        ty: 0
    }

    let inputStartTransform = <TransformData>{
        tx: 0,
        ty: subtitleHeight + NodeIOMargins.betweenPlugs
    }
    for (let input of layout.relevantInputs) {
        layout.inputLayouts[input.Id] = <IOPlugLayout>{
            textTransform: <TransformData>{
                tx: inputStartTransform.tx,
                ty: inputStartTransform.ty
            },
            plugTransform: <TransformData>{
                tx: inputStartTransform.tx - plugWidth,
                ty: inputStartTransform.ty + 5
            }
        }
        inputStartTransform.ty += bodyHeight + NodeIOMargins.betweenPlugs
    }

    let outputStartTransform = <TransformData>{
        tx: 0,
        ty: subtitleHeight + NodeIOMargins.betweenPlugs
    }
    for (let output of layout.relevantOutputs) {
        let outputLayout = <IOPlugLayout>{
            textTransform: <TransformData>{
                tx: inputOutputGap,
                ty: outputStartTransform.ty
            },
            plugTransform: <TransformData>{
                tx: 0,
                ty: outputStartTransform.ty + 5
            }
        }

        layout.outputLayouts[output.Id] = outputLayout
        outputStartTransform.ty += bodyHeight + NodeIOMargins.betweenPlugs
    }

    initialTransform.ty += Math.max(inputStartTransform.ty, outputStartTransform.ty)
}

// This function merely makes sure the group exists in groupLayout/groupKeys.
// Actually creating the proper transforms comes later.
function addIOToGroupLayout(layout: NodeLayout, io : ProcessFlowInputOutput, isInput: boolean) {
    const typeId : number = io.TypeId
    const key : string = MetadataStore.state.idToIoTypes[typeId].Name
    if (!(key in layout.groupLayout)) {
        layout.groupKeys.push(key)
        layout.groupLayout[key] = <IOGroupLayout>{
            transform: <TransformData>{
                tx: 0,
                ty: 0
            },
            titleTransform: <TransformData>{
                tx: 0,
                ty: 0
            },
            relevantInputs: [] as ProcessFlowInputOutput[],
            relevantOutputs: [] as ProcessFlowInputOutput[],
            inputLayouts: Object() as Record<number, IOPlugLayout>,
            outputLayouts: Object() as Record<number, IOPlugLayout>
        }
    }

    if (isInput) {
        layout.groupLayout[key].relevantInputs.push(io)
    } else {
        layout.groupLayout[key].relevantOutputs.push(io)
    }
}

function createDefaultNodeLayout(node : ProcessFlowNode) : NodeLayout {
    let layout = <NodeLayout>{
        transform: {...LocalSettings.state.viewBoxTransform},
        titleTransform: <TransformData>{
            tx: NodeMargins.left,
            ty: NodeMargins.top
        },
        groupLayout: Object() as Record<string, IOGroupLayout>,
        groupKeys: [] as string[],
        associatedNode: Object(),
        textWidth: 200,
        textHeight: 200,
        boxWidth: 200,
        boxHeight: 200
    }

    // Group the input and outputs by their type first.
    for (let input of node.Inputs) {
        addIOToGroupLayout(layout, input, true)
    }

    for (let output of node.Outputs) {
        addIOToGroupLayout(layout, output, false)
    }

    // Sort groups alphabetically and pull 'Flow' to the front.
    layout.groupKeys.sort()
    const executionKey = 'Flow'
    const executionIdx = layout.groupKeys.findIndex((ele) => ele == executionKey)
    if (executionIdx != -1) {
        layout.groupKeys.splice(executionIdx, 1)
        layout.groupKeys.unshift(executionKey)
    }

    let currentGroupTransform: TransformData = <TransformData>{
        tx: 0,
        ty: titleHeight
    }

    // Finally, process each group to determine where their input/output elements should lie.
    for (let groupKey of layout.groupKeys) {
        currentGroupTransform.ty += NodeIOMargins.betweenGroups
        processIOGroupLayout(layout.groupLayout[groupKey], currentGroupTransform)
    }

    return layout
}

function mergeNodeLayout(node : ProcessFlowNode, existingLayout: NodeLayout) : NodeLayout {
    let defaultLayout : NodeLayout = createDefaultNodeLayout(node)
    if (!existingLayout) {
        return defaultLayout
    }
    defaultLayout.transform = existingLayout.transform
    return defaultLayout
}

interface ProcessFlowRenderLayoutStoreState {
    nodeLayouts: Record<number, NodeLayout>
    ready: boolean
}

function onUpdateAssociatedNode(layout : NodeLayout) {
    // TODO: I'm not sure how much we actually want the store to know about the UI...
    //@ts-ignore
    const textgroup : SVGGraphicsElement = layout.associatedNode.$refs.textgroup
    layout.textWidth = textgroup.getBBox().width
    layout.textHeight = textgroup.getBBox().height

    layout.boxWidth = layout.textWidth + NodeMargins.left + NodeMargins.left
    layout.boxHeight = layout.textHeight + NodeMargins.top + NodeMargins.bottom

    for (let groupKey of layout.groupKeys) {
        let groupLayout = layout.groupLayout[groupKey]
        for (let output of groupLayout.relevantOutputs) {
            groupLayout.outputLayouts[output.Id].plugTransform.tx = layout.boxWidth

            // Also update textTransform for outputs since the box might actually be
            // longer than the input output gap.
            groupLayout.outputLayouts[output.Id].textTransform.tx = 
                Math.max(layout.textWidth, groupLayout.outputLayouts[output.Id].textTransform.tx)
        }
    }
}

function sendWebsocketUpdate(nodeId: number, layout: NodeLayout) {
    let data : Object = {...layout}
    // Remove associatedNode since JSON can't serialize it.
    //@ts-ignore
    delete data.associatedNode

    websocketConnection.send(JSON.stringify({
        NodeId: nodeId,
        Settings: data
    }))
}

const renderLayoutStore: StoreOptions<ProcessFlowRenderLayoutStoreState> = {
    state: {
        nodeLayouts : Object() as Record<number, NodeLayout>,
        ready: false
    },
    mutations: {
        resetNodeLayout(state) {
            state.ready = false
            state.nodeLayouts = Object() as Record<number, NodeLayout>
        },
        setNodeLayout(state, {nodeId, layout, isDefault}) {
            Vue.set(state.nodeLayouts, nodeId, layout)
            if (!isDefault) {
                sendWebsocketUpdate(nodeId, layout)
            }
        },
        addNodeDisplayTranslation(state, {nodeId, tx, ty}) {
            // Don't send websocket update every time this happens
            // since you'll get some noticeable lag.
            // Have the UI send a separate update once the user finishes
            // moving the element around.
            state.nodeLayouts[nodeId].transform.tx += tx
            state.nodeLayouts[nodeId].transform.ty += ty
        },
        setNodeDisplayTranslation(state, {nodeId, tx, ty, sendUpdate}) {
            // Don't send websocket update every time this happens
            // since you'll get some noticeable lag.
            // Have the UI send a separate update once the user finishes
            // moving the element around.
            state.nodeLayouts[nodeId].transform.tx = tx
            state.nodeLayouts[nodeId].transform.ty = ty
            if (sendUpdate) {
                sendWebsocketUpdate(nodeId, state.nodeLayouts[nodeId])
            }
        },
        setReady(state) {
            state.ready = true
        },
        commitNodeLayoutWithComponent(state, {nodeId, component}) {
            state.nodeLayouts[nodeId].associatedNode = component
        },
        updateNodeLayoutWithComponent(state, nodeId) {
            onUpdateAssociatedNode(state.nodeLayouts[nodeId])
        }
    },
    actions: {
        // initialize should be called as late as possible to ensure that the
        // metadata datastore has already been fully initialized.
        initialize(context, {host, csrf, processFlowStore}) {
            processFlowStore.watch((state : VuexState) => {
                return state.currentProcessFlowFullData
            }, (newFlowData : FullProcessFlowData, oldFlowData: FullProcessFlowData) => {
                    let newFlow : boolean = oldFlowData.FlowId != newFlowData.FlowId
                    context.dispatch(
                         newFlow ?
                            'recomputeLayout' :
                            'mergeLayout',
                        processFlowStore.state.currentProcessFlowFullData)

                    if (newFlow) {
                        if (!!websocketConnection && websocketConnection.readyState == WebSocket.OPEN) {
                            websocketConnection.close()
                        }

                        websocketConnection = connectProcessFlowNodeDisplaySettingsWebsocket(host, csrf, processFlowStore.state.currentProcessFlowFullData.FlowId)
                        websocketConnection.onopen = () => {
                            this.commit('setReady')
                        }
                        websocketConnection.onclose = () => {
                            // TODO Need to notify user of the close and tell them to refresh when relevant?
                        }
                        websocketConnection.onmessage = (e : MessageEvent) => {
                            // For now only grab the node's transform since everything else
                            // can just be computed.
                            let data : { NodeId: number, Settings: NodeLayout }= JSON.parse(e.data)
                            context.commit('setNodeDisplayTranslation', {
                                nodeId: data.NodeId,
                                tx: data.Settings.transform.tx,
                                ty: data.Settings.transform.ty,
                                sendUpdate: false  
                            })
                        }
                    }

            }, {
                deep: true
            })
        },
        // Assumes that we're looking at a new process flow and want to re-render
        // everything from scratch.
        recomputeLayout(context, processFlow : FullProcessFlowData) {
            context.commit('resetNodeLayout')

            // Populate layout with default values.
            for (let nodeKey of processFlow.NodeKeys) {
                context.commit('setNodeLayout', {
                    nodeId: nodeKey,
                    layout: createDefaultNodeLayout(processFlow.Nodes[nodeKey]),
                    isDefault: true
                })
            }
        },
        // Assume that we already have display data for the input process flow and only
        // want to update where necessary.
        mergeLayout(context, processFlow : FullProcessFlowData) {
            for (let nodeKey of processFlow.NodeKeys) {
                context.commit('setNodeLayout', {
                    nodeId: nodeKey,
                    layout: mergeNodeLayout(
                                processFlow.Nodes[nodeKey],
                                context.state.nodeLayouts[nodeKey]),
                    isDefault: false
                })
            }
        },
        associateNodeLayoutWithComponent(context, {nodeId, component}) {
            context.commit('commitNodeLayoutWithComponent', {nodeId, component})

            // On the next frame go through the node layout and update values
            // that need to change based on the associated node.
            Vue.nextTick(() => {
                context.commit('updateNodeLayoutWithComponent', nodeId)
            })
        },
        syncNodeTransform(context, {nodeId}) {
            sendWebsocketUpdate(nodeId, context.state.nodeLayouts[nodeId])
        }
    },
    getters: {
        isReadyForNode: (state, getters) => (nodeId : number) : boolean => {
            return state.ready && !!getters.nodeLayout(nodeId)
        },
        nodeLayout: (state) => (nodeId : number) : NodeLayout => {
            return state.nodeLayouts[nodeId]
        },
        getPlugLocation: (state, getters) => (nodeId : number, io : ProcessFlowInputOutput, isInput: boolean) : Point2D  => {
            let retPoint = <Point2D>{
                x: 0,
                y: 0
            }
            let nodeLayout : NodeLayout = getters.nodeLayout(nodeId)
            retPoint.x += nodeLayout.transform.tx
            retPoint.y += nodeLayout.transform.ty

            const typeId : number = io.TypeId
            const key : string = MetadataStore.state.idToIoTypes[typeId].Name

            let groupLayout : IOGroupLayout = nodeLayout.groupLayout[key]
            retPoint.x += groupLayout.transform.tx
            retPoint.y += groupLayout.transform.ty

            let plugLayout : IOPlugLayout
            if (isInput) {
                plugLayout = groupLayout.inputLayouts[io.Id]
            } else {
                plugLayout = groupLayout.outputLayouts[io.Id]
            }

            retPoint.x += plugLayout.plugTransform.tx
            retPoint.y += plugLayout.plugTransform.ty

            // Because the SVG rect is drawn with the origin at the top left corner.
            retPoint.x += plugWidth / 2
            retPoint.y += plugHeight / 2
            return retPoint
        }
    }
}

let store = new Vuex.Store<ProcessFlowRenderLayoutStoreState>(renderLayoutStore)
export default {
    store,
    params : {
        plugHeight,
        plugWidth,
        NodeMargins
    },
}
