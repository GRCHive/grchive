interface VuexState {
    miniMainNavBar : boolean,
    primaryNavBarWidth : number,
    allProcessFlowBasicData : ProcessFlowBasicData[],
    currentProcessFlowIndex: number,
    currentProcessFlowFullData: FullProcessFlowData,
    fullProcessFlowRequestedId: number // -1 for no process flow
    selectedNodeId : number // -1 for no selection
}
