import RelationshipMap from './relationship'

export interface FullProcessFlowData {
    FlowId: number
    Nodes: Record<number, ProcessFlowNode>
    NodeKeys: number[]
    Edges: Record<number, ProcessFlowEdge>
    EdgeKeys: number[]
    Inputs: Record<number, ProcessFlowInputOutput>
    Outputs: Record<number, ProcessFlowInputOutput>
    Risks: Record<number, ProcessFlowRisk>
    RiskKeys: number[]
    Controls: Record<number, ProcessFlowControl>
    ControlKeys: number[]
    NodeRiskRelationships: RelationshipMap<ProcessFlowNode, ProcessFlowRisk>
    NodeControlRelationships: RelationshipMap<ProcessFlowNode, ProcessFlowControl>
    RiskControlRelationships: RelationshipMap<ProcessFlowRisk, ProcessFlowControl>
}

export function isProcessFullDataEmpty(data : FullProcessFlowData) : boolean {
    if (!data) {
        return true
    }

    if (!data.Nodes || data.NodeKeys.length == 0) {
        return true
    }

    return false
}
