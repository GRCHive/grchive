import { 
    ControlDocumentationFile,
    ControlDocumentationFileHandle,
    extractControlDocumentationFileHandle
} from './controls'
import { Server, ServerHandle, extractServerHandle } from './infrastructure'

export const KSelfHosted : number = 0
export const KVendorHosted : number = 1
export const KNoHost : number = -1

export interface SelfDeployment {
    Servers: Server[]
}

export interface VendorDeployment {
    VendorName: string
    VendorProduct: string
    SocFiles: ControlDocumentationFile[]
}

export interface FullDeployment {
    Id: number
    OrgId: number
    DeploymentType: number
    VendorDeployment: VendorDeployment | null
    SelfDeployment: SelfDeployment | null
}

export function createEmptyVendorDeployment() : VendorDeployment {
    return {
        VendorName: "",
        VendorProduct: "",
        SocFiles: []
    }
}

export function createEmptySelfDeployment() : SelfDeployment {
    return {
        Servers: []
    }
}

export function deepCopyFullDeployment(f : FullDeployment) : FullDeployment {
    let copy = JSON.parse(JSON.stringify(f)) as FullDeployment
    
    if (!!copy.VendorDeployment) {
        for (let file of copy.VendorDeployment.SocFiles) {
            file.UploadTime = new Date(file.UploadTime)
            file.RelevantTime = new Date(file.RelevantTime)
        }
    }

    return copy
}

export interface StrippedSelfDeployment {
    Servers: ServerHandle[]
}

export function createStrippedSelfDeployment(f : SelfDeployment) : StrippedSelfDeployment {
    return {
        Servers: f.Servers.map(extractServerHandle)
    }
}

export interface StrippedVendorDeployment {
    VendorName: string
    VendorProduct: string
    SocFiles: ControlDocumentationFileHandle[]
}

export function createStrippedVendorDeployment(f: VendorDeployment) : StrippedVendorDeployment {
    return {
        VendorName: f.VendorName,
        VendorProduct: f.VendorProduct,
        SocFiles: f.SocFiles.map((ele: ControlDocumentationFile) => extractControlDocumentationFileHandle(ele))
    }
}

export interface StrippedFullDeployment {
    Id: number
    OrgId: number
    DeploymentType: number
    VendorDeployment? : StrippedVendorDeployment
    SelfDeployment? : StrippedSelfDeployment
}


export function createStrippedDeployment(f : FullDeployment) : StrippedFullDeployment {
    let ret = <StrippedFullDeployment>{
        Id: f.Id,
        OrgId: f.OrgId,
        DeploymentType: f.DeploymentType,
    }

    if (f.VendorDeployment) {
        ret.VendorDeployment = createStrippedVendorDeployment(f.VendorDeployment!)
    }

    if (f.SelfDeployment) {
        ret.SelfDeployment = createStrippedSelfDeployment(f.SelfDeployment!)
    }

    return ret
}
