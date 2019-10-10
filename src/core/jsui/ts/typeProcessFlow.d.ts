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

interface FullProcessFlowData {
    Nodes: ProcessFlowNode[]
}

interface FullProcessFlowDisplayData {
    Nodes: Record<number, ProcessFlowNodeDisplay>
}
