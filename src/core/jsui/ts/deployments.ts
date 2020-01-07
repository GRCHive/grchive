import { ControlDocumentationFile } from './controls'

export interface SelfDeployment {

}

export interface VendorDeployment {
    VendorName: string
    VendorProduct: string
    SocFiles: ControlDocumentationFile[]
}

export interface FullDeployment {
    Id: number
    OrgId: number
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
    return {}
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
