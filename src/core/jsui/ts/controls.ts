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
    InputDocCats: ControlDocumentationCategory[]
    OutputDocCats: ControlDocumentationCategory[]
}

export function compareDocumentationCategories(a : ControlDocumentationCategory | null, b : ControlDocumentationCategory | null) : boolean {
    if (!a || !b) {
        return false
    }
    return a.Id == b.Id
}

export function compareControls(a : ProcessFlowControl | null, b : ProcessFlowControl | null) : boolean {
    if (!a || !b) {
        return false
    }
    return a.Id == b.Id
}
