export interface ControlDocumentationCategory {
    Id: number
    Name: string
    Description: string
}

export interface FullControlData {
    Control: ProcessFlowControl
    Nodes: ProcessFlowNode[]
    Risks: ProcessFlowRisk[]
    DocumentCategories: ControlDocumentationCategory[]
}
