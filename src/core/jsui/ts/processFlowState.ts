import { FullProcessFlowData } from './processFlow'

interface VuexState {
    primaryNavBarWidth : number
    currentProcessFlowBasicData: ProcessFlowBasicData | null
    currentProcessFlowFullData: FullProcessFlowData | null
    selectedNodeId : number // -1 for no selection
    selectedEdgeId: number // -1 for no selection
}

export default VuexState
