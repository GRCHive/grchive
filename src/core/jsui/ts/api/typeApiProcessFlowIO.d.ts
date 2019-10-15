interface TGetProcessFlowIOTypesInput { 
    csrf : string
}

interface TGetProcessFlowIOTypesOutput { 
    data : ProcessFlowIOType[]
}

interface TDeleteProcessFlowIOInput { 
    csrf : string,
    ioId: number,
    isInput: boolean
}

interface TDeleteProcessFlowIOOutput { 
}
