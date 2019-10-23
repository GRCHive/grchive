import { FullProcessFlowData } from './processFlow'

interface VuexState {
    primaryNavBarWidth : number,
    allProcessFlowBasicData : ProcessFlowBasicData[],
    currentProcessFlowIndex: number,
    currentProcessFlowFullData: FullProcessFlowData,
    fullProcessFlowRequestedId: number // -1 for no process flow
    selectedNodeId : number // -1 for no selection
    selectedEdgeId: number // -1 for no selection
}

export default VuexState
