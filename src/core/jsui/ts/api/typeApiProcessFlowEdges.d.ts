interface TNewProcessFlowEdgeInput {
    csrf: string
    inputIoId: number
    outputIoId: number
}

interface TNewProcessFlowEdgeOutput {
    data: ProcessFlowEdge
}

interface TDeleteProcessFlowEdgeInput {
    csrf: string
    edgeId: number
}

interface TDeleteProcessFlowEdgeOutput {
}
