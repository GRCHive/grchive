export interface ControlDocumentationCategory {
    Id: number
    Name: string
    Description: string
}

export interface ControlDocumentationFile {
    Id: number
    StorageName: string
    RelevantTime: Date
    UploadTime: Date
    CategoryId: number
}

export interface FullControlData {
    Control: ProcessFlowControl
    Nodes: ProcessFlowNode[]
    Risks: ProcessFlowRisk[]
}
