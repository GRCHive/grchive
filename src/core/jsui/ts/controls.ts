export interface ControlDocumentationCategory {
    Id: number
    Name: string
    Description: string
}

export interface ControlDocumentationFile {
    Id: number
    StorageName: string
    RelevantTime: Date
    CategoryId: number
    AltName: string
    Description: string
}

export interface FileVersion {
	FileId        : number
	StorageId     : number
	OrgId         : number
	VersionNumber : number
}

export interface VersionedMetadata {
    File: ControlDocumentationFile
    Version: FileVersion
}

export interface FileStorageData {
	Id           : number
	MetadataId   : number
	OrgId        : number
	UploadTime   : Date
	UploadUserId : number
}

export interface ControlDocumentationFileHandle {
    Id: number
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

export function extractControlDocumentationFileHandle(f : ControlDocumentationFile) : ControlDocumentationFileHandle {
    return {
        Id: f.Id,
        CategoryId: f.CategoryId
    }
}

export function cleanJsonControlDocumentationFile(f : ControlDocumentationFile) {
    f.RelevantTime = new Date(f.RelevantTime)
}

export function cleanJsonFileStorageData(f : FileStorageData) {
    f.UploadTime = new Date(f.UploadTime)
}
