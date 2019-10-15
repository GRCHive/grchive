import Vue from 'vue'
import Vuex, { StoreOptions, ActionContext } from 'vuex'
import MetadataStore from '../metadata'

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
                tx: outputStartTransform.tx + inputOutputGap,
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
        transform: <TransformData>{
            tx: 0,
            ty: 0
        },
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

    // Sort groups alphabetically and pull 'Execution' to the front.
    layout.groupKeys.sort()
    const executionKey = 'Execution'
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
        }
    }
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
        setNodeLayout(state, {nodeId, layout}) {
            Vue.set(state.nodeLayouts, nodeId, layout)
        },
        addNodeDisplayTranslation(state, {nodeId, tx, ty}) {
            state.nodeLayouts[nodeId].transform.tx += tx
            state.nodeLayouts[nodeId].transform.ty += ty
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
        initialize(context, {processFlowStore}) {
            processFlowStore.watch((state : VuexState) => {
                return state.currentProcessFlowFullData
            }, () => {
                context.dispatch('recomputeLayout', processFlowStore.state.currentProcessFlowFullData)
            })
        },
        recomputeLayout(context, processFlow : FullProcessFlowData) {
            context.commit('resetNodeLayout')

            // Populate layout with default values.
            for (let nodeKey of processFlow.NodeKeys) {
                context.commit('setNodeLayout', {
                    nodeId: nodeKey,
                    layout: createDefaultNodeLayout(processFlow.Nodes[nodeKey])
                })
            }
            context.commit('setReady')

            // TODO: Query server for more up-to-date settings.
        },
        associateNodeLayoutWithComponent(context, {nodeId, component}) {
            context.commit('commitNodeLayoutWithComponent', {nodeId, component})

            // On the next frame go through the node layout and update values
            // that need to change based on the associated node.
            Vue.nextTick(() => {
                context.commit('updateNodeLayoutWithComponent', nodeId)
            })
        }
    },
    getters: {
        nodeLayout: (state) => (nodeId : number) : NodeLayout => {
            return state.nodeLayouts[nodeId]
        },
        getPlugLocation: (state, getters) => (nodeId : number, io : ProcessFlowInputOutput, isInput: boolean) : Point2D  => {
            let retPoint = <Point2D>{
                x: 0,
                y: 0
            }
            console.log(nodeId)
            let nodeLayout : NodeLayout = getters.nodeLayout(nodeId)
            console.log(nodeLayout)
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
