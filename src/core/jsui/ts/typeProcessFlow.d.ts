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
}

interface ProcessFlowControlType {
    Id : number
    Name : string
}

interface ProcessFlowControl {
    Id : number
    Name : string
    Description : string
    ControlTypeId: number
    FrequencyType: number
    FrequencyInterval: number 
    FrequencyOther: string
    OwnerId: number | null
    Manual: boolean
}

interface ProcessFlowNode {
    Id: number
    Name: string
    Description: string
    ProcessFlowId: number
    NodeTypeId: number
    Inputs: ProcessFlowInputOutput[]
    Outputs: ProcessFlowInputOutput[]
}

interface ProcessFlowEdge {
    Id: number
    InputIoId: number
    OutputIoId: number
}

interface ProcessFlowNodeDisplay {
    Tx: number
    Ty: number
}

interface NodeRiskRelationship {
    NodeId: number
    RiskId: number
}

interface NodeControlRelationship {
    NodeId: number
    ControlId: number
}

interface RiskControlRelationship {
    RiskId: number
    ControlId: number
}

interface FullProcessFlowResponseData {
    Nodes: ProcessFlowNode[]
    Edges: ProcessFlowEdge[]
    Risks: ProcessFlowRisk[]
    Controls: ProcessFlowControl[]
    NodeRisk: NodeRiskRelationship[]
    NodeControl: NodeControlRelationship[]
    RiskControl: RiskControlRelationship[]
}

interface FullProcessFlowDisplayData {
    Nodes: Record<number, ProcessFlowNodeDisplay>
}

interface RiskControl {
    risk: ProcessFlowRisk
    control: ProcessFlowControl
}
