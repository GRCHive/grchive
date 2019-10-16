interface ProcessFlowBasicData {
    Id : number
    Name : string
    Description: string
    CreationTime: Date
    LastUpdatedTime: Date
}

interface ProcessFlowNodeType {
    Id : number
    Name : string
    Description : string
}

interface ProcessFlowNode {
    Id: number,
    Name: string,
    Description: string,
    ProcessFlowId: number,
    NodeTypeId: number,
    Inputs: ProcessFlowInputOutput[],
    Outputs: ProcessFlowInputOutput[]
}

interface ProcessFlowEdge {
    Id: number,
    InputIoId: number,
    OutputIoId: number
}

interface ProcessFlowNodeDisplay {
    Tx: number,
    Ty: number
}

interface FullProcessFlowResponseData {
    Nodes: ProcessFlowNode[],
    Edges: ProcessFlowEdge[]
}

interface FullProcessFlowData {
    FlowId: number
    Nodes: Record<number, ProcessFlowNode>,
    NodeKeys: number[],
    Edges: Record<number, ProcessFlowEdge>,
    EdgeKeys: number[],
    Inputs: Record<number, ProcessFlowInputOutput>,
    Outputs: Record<number, ProcessFlowInputOutput>
}

interface FullProcessFlowDisplayData {
    Nodes: Record<number, ProcessFlowNodeDisplay>
}
