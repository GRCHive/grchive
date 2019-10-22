interface FrequencyData {

}

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

interface ProcessFlowRisk {
    Id : number
    Name : string
    Description: string
    RelevantNodeIds: number[]
}

interface ProcessFlowControlType {
    Id : number
    Name : string
}

interface ProcessFlowControl {
    Id : number
    Name : string
    Description : string
    Type: ProcessFlowControlType
    Frequency: FrequencyData
    ProcessOwner: string
    RelevantRiskIds: number[]
    RelevantNodeIds: number[]
}

interface ProcessFlowNode {
    Id: number,
    Name: string,
    Description: string,
    ProcessFlowId: number,
    NodeTypeId: number,
    Inputs: ProcessFlowInputOutput[],
    Outputs: ProcessFlowInputOutput[],
    RiskIds: number[]
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
    Edges: ProcessFlowEdge[],
    Risks: ProcessFlowRisk[]
}

interface FullProcessFlowData {
    FlowId: number
    Nodes: Record<number, ProcessFlowNode>,
    NodeKeys: number[],
    Edges: Record<number, ProcessFlowEdge>,
    EdgeKeys: number[],
    Inputs: Record<number, ProcessFlowInputOutput>,
    Outputs: Record<number, ProcessFlowInputOutput>,
    Risks: Record<number, ProcessFlowRisk>
    RiskKeys: number[],
    Controls: Record<number, ProcessFlowControl>,
    ControlKeys: number[]
}

interface FullProcessFlowDisplayData {
    Nodes: Record<number, ProcessFlowNodeDisplay>
}
