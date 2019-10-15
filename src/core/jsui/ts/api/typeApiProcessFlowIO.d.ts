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

interface TEditProcessFlowIOInput { 
    csrf : string,
    ioId: number,
    isInput: boolean,
    name: string,
    type: number
}

interface TEditProcessFlowIOOutput { 
    data: ProcessFlowInputOutput
}
