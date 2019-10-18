interface IOPlugLayout {
    textTransform: TransformData
    plugTransform: TransformData
}

interface IOGroupLayout {
    transform: TransformData
    titleTransform: TransformData

    relevantInputs: ProcessFlowInputOutput[],
    relevantOutputs: ProcessFlowInputOutput[]

    inputLayouts: Record<number, IOPlugLayout>,
    outputLayouts: Record<number, IOPlugLayout>
}

interface NodeLayout {
    // Transform of the entire node.
    transform: TransformData

    // Relative transform of the title.
    titleTransform: TransformData

    // IO Groups
    groupLayout: Record<string, IOGroupLayout>
    groupKeys: string[]

    // Associated node so we can query things about the DOM.
    associatedNode: Object

    textWidth: number
    textHeight: number

    boxWidth: number
    boxHeight: number
}
