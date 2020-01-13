import { 
    ControlDocumentationFile,
    ControlDocumentationFileHandle,
    extractControlDocumentationFileHandle
} from './controls'
import { Server, ServerHandle, extractServerHandle } from './infrastructure'
import { VendorProduct } from './vendors'

export const KSelfHosted : number = 0
export const KVendorHosted : number = 1
export const KNoHost : number = -1

export interface SelfDeployment {
    Servers: Server[]
}

export interface VendorDeployment {
    Product: VendorProduct | null
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
        Product: null
    }
}

export function createEmptySelfDeployment() : SelfDeployment {
    return {
        Servers: []
    }
}

export function deepCopyFullDeployment(f : FullDeployment) : FullDeployment {
    let copy = JSON.parse(JSON.stringify(f)) as FullDeployment
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
    VendorId : number
    ProductId: number
}

export function createStrippedVendorDeployment(f: VendorDeployment) : StrippedVendorDeployment {
    if (!f.Product) {
        return {
            VendorId: -1,
            ProductId: -1,
        }
    }
    return {
        VendorId: f.Product.VendorId,
        ProductId: f.Product.Id,
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
