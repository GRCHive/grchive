interface TGetControlTypesInput {
    csrf : string
}

interface TGetControlTypesOutput {
    data: ProcessFlowControlType[]
}

interface TNewControlInput {
    csrf : string
    name : string 
    description: string
    controlType : number
    frequencyType : number
    frequencyInterval : number
    ownerId : number
    nodeId:  number
    riskId : number
}

interface TNewControlOutput {
    data: ProcessFlowControl
}
