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

interface ProcessFlowInput {
    Id: number
    Name: string
}

interface ProcessFlowOutput {
    Id: number
    Name: string
}

interface ProcessFlowNode {
    Id: number,
    Name: string,
    Description: string,
    ProcessFlowId: number,
    NodeTypeId: number,
    Inputs: ProcessFlowInput[],
    Outputs: ProcessFlowOutput[]
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
