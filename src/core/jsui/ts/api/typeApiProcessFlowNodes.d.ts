interface TGetProcessFlowNodeTypesInput { 
    csrf : string
}

interface TGetProcessFlowNodeTypesOutput { 
    data : ProcessFlowNodeType[]
}

interface TEditProcessFlowNodeInput { 
    csrf : string,
    nodeId: number,
    name: string,
    description: string,
    type: number
}

interface TEditProcessFlowNodeOutput { 
    data: ProcessFlowNode
}
