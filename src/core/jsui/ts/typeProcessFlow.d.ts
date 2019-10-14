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

interface ProcessFlowIOType {
    Id: number
    Name : string
}

interface ProcessFlowInputOutput {
    Id: number
    Name: string
    ParentNodeId: number
    TypeId: number
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

interface ProcessFlowNodeDisplay {
    Tx: number,
    Ty: number
}

interface FullProcessFlowResponseData {
    Nodes: ProcessFlowNode[]
}

interface FullProcessFlowData {
    Nodes: Record<number, ProcessFlowNode>,
    NodeKeys: number[],
    NumNodes: number
}

interface FullProcessFlowDisplayData {
    Nodes: Record<number, ProcessFlowNodeDisplay>
}
