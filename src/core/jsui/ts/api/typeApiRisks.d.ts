interface TDeleteRiskInput {
    csrf: string
    nodeId: number
    riskIds: number[]
    global: boolean
}

interface TDeleteRiskOutput {
}

interface TAddExistingRiskInput {
    csrf: string
    nodeId: number
    riskIds: number[]
}

interface TAddExistingRiskOutput {
}
