export enum IntegrationType {
    SapErp = 1
}

export interface GenericIntegration {
    Id: number
    OrgId: number
    Type : IntegrationType
    Name : string
    Description : string
}

export function emptyGenericIntegration() : GenericIntegration {
    return {
        Id: -1,
        OrgId: -1,
        Type: IntegrationType.SapErp,
        Name: "",
        Description: "",
    }
}

import { SapErpIntegrationSetup } from './sap'
export type UnionIntegrationSetup = SapErpIntegrationSetup | null
export type NNUnionIntegrationSetup = SapErpIntegrationSetup
