interface TNewRiskInput {
    csrf: string
    name: string
    description : string
    nodeId: number
}

interface TNewRiskOutput {
    data: ProcessFlowRisk
}

interface TDeleteRiskInput {
    csrf: string
    nodeId: number
    riskIds: number[]
    global: boolean
}

interface TDeleteRiskOutput {
}
